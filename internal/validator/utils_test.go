/*
* SPDX-FileCopyrightText: Copyright (c) <2024> NVIDIA CORPORATION & AFFILIATES. All rights reserved.
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
	"bytes"
	"errors"

	"github.com/nvidia/container-canary/internal/container"
	"github.com/rs/zerolog"
)

type dummyContainer struct{}

func (c *dummyContainer) Start() error { return nil }

func (c *dummyContainer) Remove() error { return nil }

func (c *dummyContainer) Status() (*container.ContainerInfo, error) { return nil, nil }

func (c *dummyContainer) Exec(command ...string) (exitCode int, stdout string, stderr string, err error) {
	switch command[0] {
	case "success":
		return 0, "Success stdout", "Success stderr", nil
	case "failure":
		return 1, "Failure stdout", "Failure stderr", nil
	case "error":
		return 0, "", "", errors.New("This command is an error")
	default:
		return 1, "", "", nil
	}
}

func (c *dummyContainer) Logs() (string, error) { return "", nil }

func logger() (zerolog.Logger, *bytes.Buffer) {
	buf := new(bytes.Buffer)
	logger := zerolog.New(buf)
	return logger, buf
}
