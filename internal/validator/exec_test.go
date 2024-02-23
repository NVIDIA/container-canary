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
	"testing"

	canaryv1 "github.com/nvidia/container-canary/internal/apis/v1"
	"github.com/nvidia/container-canary/internal/container"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	v1 "k8s.io/api/core/v1"
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

func TestExecSuccess(t *testing.T) {
	assert := assert.New(t)

	c := &dummyContainer{}
	probe := &canaryv1.Probe{
		Exec: &v1.ExecAction{
			Command: []string{"success"},
		},
	}

	logger, buf := logger()
	e := logger.Info()
	result, err := ExecCheck(c, probe, e)
	e.Send()
	assert.True(result)
	assert.NoError(err)
	assert.Equal("{\"level\":\"info\",\"exitCode\":0,\"stdout\":\"Success stdout\",\"stderr\":\"Success stderr\"}\n", buf.String())
}

func TestExecFailure(t *testing.T) {
	assert := assert.New(t)

	c := &dummyContainer{}
	probe := &canaryv1.Probe{
		Exec: &v1.ExecAction{
			Command: []string{"failure"},
		},
	}

	logger, buf := logger()
	e := logger.Info()
	result, err := ExecCheck(c, probe, e)
	e.Send()
	assert.False(result)
	assert.NoError(err)
	assert.Equal("{\"level\":\"info\",\"exitCode\":1,\"stdout\":\"Failure stdout\",\"stderr\":\"Failure stderr\"}\n", buf.String())
}

func TestExecError(t *testing.T) {
	assert := assert.New(t)

	c := &dummyContainer{}
	probe := &canaryv1.Probe{
		Exec: &v1.ExecAction{
			Command: []string{"error"},
		},
	}

	logger, buf := logger()
	e := logger.Info()
	result, err := ExecCheck(c, probe, e)
	e.Send()
	assert.False(result)
	assert.Error(err, "This command is an error")
	assert.Equal("{\"level\":\"info\"}\n", buf.String())
}
