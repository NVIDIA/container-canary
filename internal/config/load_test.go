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
	"github.com/stretchr/testify/require"
)

func TestLoadValidatorPreservesIdentity(t *testing.T) {
	tests := []struct {
		name       string
		identity   string
		apiVersion string
		kind       string
	}{
		{
			name:       "documented identity",
			identity:   "apiVersion: container-canary.nvidia.com/v1\nkind: Validator\n",
			apiVersion: "container-canary.nvidia.com/v1",
			kind:       "Validator",
		},
		{
			name:       "custom identity",
			identity:   "apiVersion: example.invalid/v99\nkind: SomethingElse\n",
			apiVersion: "example.invalid/v99",
			kind:       "SomethingElse",
		},
		{
			name:     "missing API version",
			identity: "kind: Validator\n",
			kind:     "Validator",
		},
		{
			name:       "missing kind",
			identity:   "apiVersion: container-canary.nvidia.com/v1\n",
			apiVersion: "container-canary.nvidia.com/v1",
		},
		{
			name: "missing identity",
		},
	}

	const manifest = `name: example
checks:
  - name: command
    probe:
      exec:
        command: ["true"]
`

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			validator, err := LoadValidatorFromBytes([]byte(tt.identity + manifest))
			require.NoError(t, err)
			assert.Equal(t, tt.apiVersion, validator.APIVersion)
			assert.Equal(t, tt.kind, validator.Kind)
		})
	}
}

func TestValidator(t *testing.T) {
	assert := assert.New(t)

	validator, err := LoadValidatorFromFile("../../examples/kubeflow.yaml")

	assert.Nil(err)
	assert.Equal("container-canary.nvidia.com/v1", validator.APIVersion)
	assert.Equal("Validator", validator.Kind)
	assert.Equal("kubeflow", validator.Name)
	assert.Equal("Kubeflow notebooks", validator.Description)

	assert.GreaterOrEqual(len(validator.Checks), 1)

	check := validator.Checks[0]

	assert.Equal("user", check.Name)
	assert.Equal("👩 User is jovyan", check.Description)

	assert.Equal(0, check.Probe.InitialDelaySeconds)

	check = validator.Checks[5]

	assert.Equal("allow-origin-all", check.Name)
	assert.Equal("🔓 Sets 'Access-Control-Allow-Origin: *' header", check.Description)

	assert.Equal("/hub/jovyan/lab", check.Probe.HTTPGet.Path)
	assert.Equal(8888, check.Probe.HTTPGet.Port)

	header := check.Probe.HTTPGet.HTTPHeaders[0]
	assert.Equal("User-Agent", header.Name)
	assert.Equal("container-canary/0.2.1", header.Value)

	header = check.Probe.HTTPGet.ResponseHTTPHeaders[0]
	assert.Equal("Access-Control-Allow-Origin", header.Name)
	assert.Equal("*", header.Value)
}
