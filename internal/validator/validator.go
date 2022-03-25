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
	"os"
	"os/signal"
	"syscall"
	"time"

	canaryv1 "github.com/nvidia/container-canary/internal/apis/v1"
	"github.com/nvidia/container-canary/internal/container"
	"github.com/nvidia/container-canary/internal/terminal"
	"github.com/spf13/cobra"
)

type checkResult struct {
	Passed bool
	Error  error
}

func Validate(image string, validator *canaryv1.Validator, cmd *cobra.Command, debug bool) (bool, error) { // Start image
	c := container.New(image, validator.Env, validator.Ports, validator.Volumes, validator.Command)
	err := c.Start()
	defer c.Remove()
	if debug {
		status, _ := c.Status()
		cmd.Printf("Running container with command '%s'\n", status.RunCommand)
	}
	if err != nil {
		logs, logsErr := c.Logs()
		if logsErr == nil {
			cmd.Print(logs)
		} else {
			cmd.Println("Unable to get logs")
			cmd.Println(logsErr)
		}
		return false, err
	}

	if len(validator.Checks) == 0 {
		return false, fmt.Errorf("no checks found")
	}

	allChecksPassed := true
	results := make(chan checkResult)

	// Start checks
	for _, check := range validator.Checks {
		go runCheck(results, c, check)
	}

	// Wait for checks
	for j := 1; j <= len(validator.Checks); j++ {
		cr := <-results
		if cr.Error != nil {
			return false, cr.Error
		}
		if !cr.Passed {
			allChecksPassed = false
		}
	}
	if debug && !allChecksPassed {
		done := make(chan os.Signal, 1)
		signal.Notify(done, syscall.SIGINT, syscall.SIGTERM)
		cmd.Println("Leaving container running for debugging, press ctrl+c to exit...")
		<-done
	}
	return allChecksPassed, nil
}

func runCheck(results chan<- checkResult, c container.ContainerInterface, check canaryv1.Check) {
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
		results <- checkResult{false, fmt.Errorf("check '%s' has no known probes", check.Name)}
		return
	}
	results <- checkResult{p, err}
	terminal.PrintCheckItem("", check.Description, getStatus(p, err))
}

type probeCallable func(container.ContainerInterface, *canaryv1.Probe) (bool, error)

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
		if passes > probe.SuccessThreshold || fails > probe.FailureThreshold {
			return passFail, err
		}
		if time.Since(start) > time.Duration(probe.TimeoutSeconds)*time.Second {
			return false, fmt.Errorf("check timed out after %d seconds", probe.TimeoutSeconds)
		}
		time.Sleep(time.Duration(probe.PeriodSeconds) * time.Second)
	}
}
