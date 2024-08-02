/*
* SPDX-FileCopyrightText: Copyright (c) <2022> NVIDIA CORPORATION & AFFILIATES. All rights reserved.
* SPDX-License-Identifier: Apache-2.0
*
* Licensed under the Apache License, Version 2.0 (the "License");
* you may not use this file except in compliance with the License.
* You may obtain a copy of the License at
*
* http://www.apache.org/licenses/LICENSE-2.0
*
* Unless required by applicable law or agreed to in writing, software
* distributed under the License is distributed on an "AS IS" BASIS,
* WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
* See the License for the specific language governing permissions and
* limitations under the License.
 */

package validator

import (
	"fmt"

	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	canaryv1 "github.com/nvidia/container-canary/internal/apis/v1"
	"github.com/nvidia/container-canary/internal/container"
)

type model struct {
	sub                     chan checkResult
	container               container.ContainerInterface
	containerStarted        bool
	containerStartupTimeout int
	results                 []checkResult
	allChecksPassed         bool
	spinner                 spinner.Model
	progress                progress.Model
	debug                   bool
	image                   string
	validator               *canaryv1.Validator
	configPath              string
	err                     error
	tty                     bool
}

func (m model) Init() tea.Cmd {
	return tea.Batch(
		spinner.Tick,
		waitForChecks(m.sub),
		loadConfig(m.configPath),
	)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	// Application events
	case tea.KeyMsg:
		return handleKeypress(m, msg)
	case configLoaded:
		return handleConfigLoaded(m, msg)
	case containerStarted:
		return handleContainerStarted(m, msg)
	case containerStopped:
		return m, tea.Quit
	case containerFailed:
		return handleContainerFailed(m, msg)
	case checkResult:
		return handleCheckResult(m, msg)

	// UI events
	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	case progress.FrameMsg:
		progressModel, cmd := m.progress.Update(msg)
		m.progress = progressModel.(progress.Model)
		return m, cmd

	default:
		return m, nil
	}
}

func (m model) View() string {
	help := helpStyle("Press q to quit...")

	if m.tty {
		if m.validator == nil {
			return fmt.Sprintf("%s Loading config\n", m.spinner.View()) + help
		}

		if !m.containerStarted {
			return fmt.Sprintf("%s Starting container\n", m.spinner.View()) + help
		}

		if len(m.results) < len(m.validator.Checks) {
			return m.progress.View() + "\n" + help
		}
	}

	return help
}

func (m model) Passed() bool {
	return m.allChecksPassed
}

func handleKeypress(m model, keypress tea.KeyMsg) (model, tea.Cmd) {
	switch keypress.String() {
	// These keys should exit the program.
	case "ctrl+c", "q":
		m.allChecksPassed = false
		if m.containerStarted {
			return m, shutdown(m.container)
		} else {
			return m, tea.Quit
		}
	default:
		return m, nil
	}

}

func handleConfigLoaded(m model, msg configLoaded) (model, tea.Cmd) {
	var commands []tea.Cmd
	if msg.Error != nil {
		m.err = msg.Error
		return m, tea.Batch(tea.Printf("Error: %s\n", msg.Error.Error()), tea.Quit)
	}
	m.validator = msg.Config
	if !m.tty {
		commands = append(commands, tea.Printf("Starting container"))
	}
	commands = append(commands, startContainer(m.image, m.validator, m.containerStartupTimeout))
	return m, tea.Batch(commands...)
}

func handleContainerStarted(m model, msg containerStarted) (model, tea.Cmd) {
	var commands []tea.Cmd
	m.containerStarted = true
	m.container = msg.Container

	if m.debug {
		status, _ := m.container.Status()
		commands = append(commands, tea.Printf("Running container with command '%s'", status.RunCommand))
	}
	commands = append(commands, tea.Printf("Validating %s against %s", highlightStyle(m.image), highlightStyle(m.validator.Name)))
	for _, check := range m.validator.Checks {
		commands = append(commands, runCheck(m.sub, m.container, check))
	}
	return m, tea.Batch(commands...)
}

func handleContainerFailed(m model, msg containerFailed) (model, tea.Cmd) {
	m.err = msg.Error
	return m, tea.Batch(tea.Printf("Error: %s\n", m.err.Error()), tea.Quit)
}

func handleCheckResult(m model, msg checkResult) (model, tea.Cmd) {
	var commands []tea.Cmd
	var printCommands []tea.Cmd
	m.results = append(m.results, msg)
	if !msg.Passed {
		m.allChecksPassed = false
	}
	printCommands = append(printCommands,
		tea.Printf(" %-50s [%s]", msg.Description, getStatus(msg.Passed, msg.Error)))
	if len(m.results) == len(m.validator.Checks) {
		if m.allChecksPassed {
			printCommands = append(printCommands, tea.Println(passedStyle("validation passed")))
		} else {
			printCommands = append(printCommands, tea.Println(failedStyle("validation failed")))
		}
		if !m.allChecksPassed && m.debug {
			printCommands = append(printCommands, tea.Println("Leaving container running for debugging..."))
		} else {
			commands = append(commands, shutdown(m.container))
		}
	} else {
		if m.allChecksPassed {
			m.progress.FullColor = "10"
		} else {
			m.progress.FullColor = "9"
		}
		commands = append(commands,
			waitForChecks(m.sub),
			m.progress.SetPercent(float64(len(m.results))/float64(len(m.validator.Checks))))
	}
	return m, tea.Sequence(tea.Sequence(printCommands...), tea.Batch(commands...))
}
