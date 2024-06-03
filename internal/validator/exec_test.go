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
	"testing"

	canaryv1 "github.com/nvidia/container-canary/internal/apis/v1"
	"github.com/stretchr/testify/assert"
	v1 "k8s.io/api/core/v1"
)

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
