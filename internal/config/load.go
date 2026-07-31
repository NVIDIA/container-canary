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

package config

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	canaryv1 "github.com/nvidia/container-canary/internal/apis/v1"
	yaml "gopkg.in/yaml.v2"
)

func LoadValidatorFromURL(url string) (*canaryv1.Validator, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return LoadValidatorFromBytes(body)
}

func LoadValidatorFromFile(path string) (*canaryv1.Validator, error) {
	filename, err := filepath.Abs(path)
	if err != nil {
		return nil, err
	}

	if _, err := os.Stat(filename); errors.Is(err, os.ErrNotExist) {
		return nil, fmt.Errorf("no such file %s", filename)
	}

	yamlFile, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	return LoadValidatorFromBytes(yamlFile)
}

func LoadValidatorFromBytes(b []byte) (*canaryv1.Validator, error) {
	var validator canaryv1.Validator

	err := yaml.Unmarshal(b, &validator)
	if err != nil {
		return nil, err
	}
	if err := validateChecks(validator.Checks); err != nil {
		return nil, err
	}

	return &validator, nil
}

func validateChecks(checks []canaryv1.Check) error {
	for i, check := range checks {
		actionCount := 0
		if check.Probe.Exec != nil {
			actionCount++
		}
		if check.Probe.HTTPGet != nil {
			actionCount++
		}
		if check.Probe.TCPSocket != nil {
			actionCount++
		}

		checkIdentifier := fmt.Sprintf("checks[%d] %q", i, check.Name)
		if actionCount != 1 {
			return fmt.Errorf("%s: probe must define exactly one supported action, found %d", checkIdentifier, actionCount)
		}

		switch {
		case check.Probe.Exec != nil:
			if len(check.Probe.Exec.Command) == 0 {
				return fmt.Errorf("%s: exec probe must define at least one command", checkIdentifier)
			}
		case check.Probe.HTTPGet != nil:
			if check.Probe.HTTPGet.Port < 1 || check.Probe.HTTPGet.Port > 65535 {
				return fmt.Errorf("%s: httpGet probe port must be between 1 and 65535", checkIdentifier)
			}
		case check.Probe.TCPSocket != nil:
			if check.Probe.TCPSocket.Port < 1 || check.Probe.TCPSocket.Port > 65535 {
				return fmt.Errorf("%s: tcpSocket probe port must be between 1 and 65535", checkIdentifier)
			}
		}
	}

	return nil
}
