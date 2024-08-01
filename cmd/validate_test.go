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

package cmd

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidate(t *testing.T) {
	assert := assert.New(t)
	b := new(bytes.Buffer)
	rootCmd.SetOut(b)
	rootCmd.SetErr(b)
	rootCmd.SetArgs([]string{"validate", "--file", "../examples/kubeflow.yaml", "container-canary/kubeflow:shouldpass"})
	err := rootCmd.Execute()

	assert.Nil(err, "should not error")
	assert.Contains(b.String(), "Validating container-canary/kubeflow:shouldpass against kubeflow", "did not validate")
	assert.Contains(b.String(), "validation passed", "did not pass")
}

func TestValidateFails(t *testing.T) {
	assert := assert.New(t)
	b := new(bytes.Buffer)
	rootCmd.SetOut(b)
	rootCmd.SetErr(b)
	rootCmd.SetArgs([]string{"validate", "--file", "../examples/kubeflow.yaml", "container-canary/kubeflow:shouldfail"})
	err := rootCmd.Execute()

	assert.NotNil(err, "should fail")
	assert.Contains(b.String(), "validation failed", "did not fail")
}

func TestFileDoesNotExist(t *testing.T) {
	assert := assert.New(t)
	b := new(bytes.Buffer)
	rootCmd.SetOut(b)
	rootCmd.SetErr(b)
	rootCmd.SetArgs([]string{"validate", "--file", "foo.yaml", "container-canary/kubeflow:doesnotexist"})
	err := rootCmd.Execute()

	assert.NotNil(err, "did not error")
	assert.Contains(b.String(), "Cannot find container-canary/kubeflow:doesnotexist", "did not fail")
}

func TestValidateRespectsStartupTimeout(t *testing.T) {
	assert := assert.New(t)
	b := new(bytes.Buffer)
	rootCmd.SetOut(b)
	rootCmd.SetErr(b)
	rootCmd.SetArgs([]string{"validate", "--file", "../examples/kubeflow.yaml", "container-canary/long-sleep:local", "--startup-timeout", "3"})
	err := rootCmd.Execute()

	assert.NotNil(err, "should fail")
	assert.Contains(b.String(), "validation failed", "container failed to start after 3 seconds")
}
