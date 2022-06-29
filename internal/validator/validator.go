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
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	canaryv1 "github.com/nvidia/container-canary/internal/apis/v1"
	"github.com/nvidia/container-canary/internal/config"
	"github.com/nvidia/container-canary/internal/container"
	"github.com/spf13/cobra"
)

var helpStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("8")).Render
var passedStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("10")).Render
var failedStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("9")).Render
var highlightStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("6")).Render

type checkResult struct {
	Description string
	Passed      bool
	Error       error
}

type containerFailed struct {
	Error error
}
type containerStarted struct {
	Container container.ContainerInterface
}
type containerStopped struct{}
type configLoaded struct {
	Config *canaryv1.Validator
	Error  error
}

type probeCallable func(container.ContainerInterface, *canaryv1.Probe) (bool, error)

func Validate(image string, configPath string, cmd *cobra.Command, debug bool) (bool, error) {
	var tty io.Reader
	isTty := true
	tty, err := os.Open("/dev/tty")
	if err != nil {
		tty = bufio.NewReader(os.Stdin)
		isTty = false
	}
	m := model{
		sub:              make(chan checkResult),
		configPath:       configPath,
		containerStarted: false,
		spinner:          spinner.New(),
		progress:         progress.New(progress.WithSolidFill("#f2e63a")),
		allChecksPassed:  true,
		debug:            debug,
		image:            image,
		tty:              isTty,
	}
	p := tea.NewProgram(m, tea.WithInput(tty), tea.WithOutput(cmd.OutOrStderr()))
	out, err := p.StartReturningModel()
	if err != nil {
		return false, err
	}
	if out, ok := out.(model); ok {
		return out.allChecksPassed, out.err
	} else {
		return false, errors.New("program returned unknown model")
	}

}

func loadConfig(filePath string) tea.Cmd {
	return func() tea.Msg {
		var validatorConfig *canaryv1.Validator
		var err error

		if strings.Contains(filePath, "://") {
			validatorConfig, err = config.LoadValidatorFromURL(filePath)
		} else {
			validatorConfig, err = config.LoadValidatorFromFile(filePath)
		}
		if err != nil {
			return configLoaded{
				Config: nil,
				Error:  err,
			}
		}
		if len(validatorConfig.Checks) == 0 {
			return configLoaded{
				Config: nil,
				Error:  errors.New("no checks found"),
			}
		}

		return configLoaded{
			Config: validatorConfig,
			Error:  nil,
		}
	}
}

func startContainer(image string, validator *canaryv1.Validator) tea.Cmd {
	return func() tea.Msg {
		container := container.New(image, validator.Env, validator.Ports, validator.Volumes, validator.Command)
		err := container.Start()
		if err != nil {
			return containerFailed{Error: err}
		}
		return containerStarted{
			Container: container,
		}
	}
}

func shutdown(container container.ContainerInterface) tea.Cmd {
	return func() tea.Msg {
		err := container.Remove()
		if err != nil {
			return containerFailed{Error: err}
		}
		return containerStopped{}
	}
}

func waitForChecks(sub chan checkResult) tea.Cmd {
	return func() tea.Msg {
		return <-sub
	}
}

func runCheck(results chan<- checkResult, c container.ContainerInterface, check canaryv1.Check) tea.Cmd {
	return func() tea.Msg {
		var p bool
		var err error
		// TODO Make more SOLID (O)
		if check.Probe.Exec != nil {
			p, err = executeCheck(ExecCheck, c, &check.Probe)
		} else if check.Probe.HTTPGet != nil {
			p, err = executeCheck(HTTPGetCheck, c, &check.Probe)
		} else if check.Probe.TCPSocket != nil {
			p, err = executeCheck(TCPSocketCheck, c, &check.Probe)
		} else {
			results <- checkResult{check.Description, false, fmt.Errorf("check '%s' has no known probes", check.Name)}
			return nil
		}
		results <- checkResult{check.Description, p, err}
		return nil
	}
}

// Run a check method with appropriate delay, retries and retry interval
func executeCheck(method probeCallable, c container.ContainerInterface, probe *canaryv1.Probe) (bool, error) {
	time.Sleep(time.Duration(probe.InitialDelaySeconds) * time.Second)
	passes := 0
	fails := 0
	start := time.Now()
	for {
		passFail, err := method(c, probe)
		if err != nil {
			return false, err
		}
		if passFail {
			passes += 1
			fails = 0
		} else {
			fails += 1
			passes = 0
		}
		if passes >= probe.SuccessThreshold || fails >= probe.FailureThreshold {
			return passFail, err
		}
		if time.Since(start) > time.Duration(probe.TimeoutSeconds)*time.Second {
			return false, fmt.Errorf("check timed out after %d seconds", probe.TimeoutSeconds)
		}
		time.Sleep(time.Duration(probe.PeriodSeconds) * time.Second)
	}
}

func getStatus(check bool, err error) string {
	if err != nil {
		return failedStyle(fmt.Sprintf("error - %s", err.Error()))
	} else {
		if check {
			return passedStyle("passed")
		} else {
			return failedStyle("failed")
		}

	}
}
