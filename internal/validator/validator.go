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
	"time"

	canaryv1 "github.com/jacobtomlinson/containercanary/internal/apis/v1"
	"github.com/jacobtomlinson/containercanary/internal/container"
	"github.com/jacobtomlinson/containercanary/internal/terminal"
)

type checkResult struct {
	Passed bool
	Error  error
}

func Validate(image string, validator *canaryv1.Validator) (bool, error) { // Start image
	c := container.New(image, validator.Env, validator.Ports, validator.Volumes)
	c.Start()
	defer c.Remove()

	if len(validator.Checks) == 0 {
		return false, fmt.Errorf("error no checks found")
	}

	allChecksPassed := true
	results := make(chan checkResult)

	// Start checks
	for _, check := range validator.Checks {
		go runCheck(results, &c, check)
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
	return allChecksPassed, nil
}

func runCheck(results chan<- checkResult, c *container.Container, check canaryv1.Check) {
	// TODO Retry each check on fail "failureThreshold" times with "periodSeconds" sleep between
	// TODO Retry each check on success "successThreshold" times with "periodSeconds" sleep between
	time.Sleep(time.Duration(check.Probe.InitialDelaySeconds) * time.Second)
	if check.Probe.Exec != nil {
		p, err := ExecCheck(c, check.Probe.Exec)
		results <- checkResult{p, nil}
		terminal.PrintCheckItem("", check.Description, getStatus(p, err))
		return
	}
	if check.Probe.HTTPGet != nil {
		p, err := HTTPGetCheck(c, check.Probe.HTTPGet)
		results <- checkResult{p, nil}
		terminal.PrintCheckItem("", check.Description, getStatus(p, err))
		return
	}
	results <- checkResult{false, fmt.Errorf("check '%s' has no known probes", check.Name)}
}
