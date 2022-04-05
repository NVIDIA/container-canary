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
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidate(t *testing.T) {
	assert := assert.New(t)
	rootCmd.SetArgs([]string{"validate", "--file", "../examples/kubeflow.yaml", "container-canary/kubeflow:shouldpass"})
	err := rootCmd.Execute()
	assert.Nil(err, "should not error")
}

func TestFileDoesNotExist(t *testing.T) {
	assert := assert.New(t)
	rootCmd.SetArgs([]string{"validate", "--file", "foo.yaml", "container-canary/kubeflow:shouldpass"})
	err := rootCmd.Execute()
	assert.NotNil(err, "did not error")
}
