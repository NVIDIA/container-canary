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
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidator(t *testing.T) {
	assert := assert.New(t)

	validator, err := LoadValidatorFromFile("../../examples/kubeflow.yaml")

	assert.Nil(err)
	assert.Equal("kubeflow", validator.Name)
	assert.Equal("Kubeflow notebooks", validator.Description)

	assert.GreaterOrEqual(len(validator.Checks), 1)

	check := validator.Checks[0]

	assert.Equal("user", check.Name)
	assert.Equal("👩 User is jovyan", check.Description)

	assert.Equal(0, check.Probe.InitialDelaySeconds)

}
