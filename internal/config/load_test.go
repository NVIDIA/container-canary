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

func TestLoadValidatorFromBytesValidatesProbeActions(t *testing.T) {
	tests := []struct {
		name    string
		probe   string
		wantErr string
	}{
		{
			name:    "no action",
			probe:   ` {}`,
			wantErr: `check "test-check": probe must define exactly one supported action, found 0`,
		},
		{
			name: "multiple actions",
			probe: `
      exec:
        command: ["true"]
      tcpSocket:
        port: 8080`,
			wantErr: `check "test-check": probe must define exactly one supported action, found 2`,
		},
		{
			name: "empty exec action",
			probe: `
      exec: {}`,
			wantErr: `check "test-check": exec probe must define at least one command`,
		},
		{
			name: "empty HTTP GET action",
			probe: `
      httpGet: {}`,
			wantErr: `check "test-check": httpGet probe port must be between 1 and 65535`,
		},
		{
			name: "empty TCP socket action",
			probe: `
      tcpSocket: {}`,
			wantErr: `check "test-check": tcpSocket probe port must be between 1 and 65535`,
		},
		{
			name: "HTTP GET port above range",
			probe: `
      httpGet:
        port: 65536`,
			wantErr: `check "test-check": httpGet probe port must be between 1 and 65535`,
		},
		{
			name: "TCP socket port above range",
			probe: `
      tcpSocket:
        port: 65536`,
			wantErr: `check "test-check": tcpSocket probe port must be between 1 and 65535`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			manifest := []byte(`name: example
checks:
  - name: test-check
    probe:` + tt.probe + "\n")

			validator, err := LoadValidatorFromBytes(manifest)

			assert.Nil(t, validator)
			assert.EqualError(t, err, tt.wantErr)
		})
	}
}

func TestLoadValidatorFromBytesAcceptsUsableProbeActions(t *testing.T) {
	manifest := []byte(`name: example
checks:
  - name: exec
    probe:
      exec:
        command: ["true"]
  - name: http
    probe:
      httpGet:
        port: 8080
  - name: tcp
    probe:
      tcpSocket:
        port: 8080
`)

	validator, err := LoadValidatorFromBytes(manifest)

	assert.NoError(t, err)
	assert.Len(t, validator.Checks, 3)
}

func TestLoadValidatorFromBytesRejectsMissingProbe(t *testing.T) {
	validator, err := LoadValidatorFromBytes([]byte("name: example\nchecks:\n  - name: test-check\n"))

	assert.Nil(t, validator)
	assert.EqualError(t, err, "check \"test-check\": probe must define exactly one supported action, found 0")
}

func TestLoadValidatorFromBytesIdentifiesInvalidCheck(t *testing.T) {
	validator, err := LoadValidatorFromBytes([]byte(`name: example
checks:
  - name: valid-check
    probe:
      exec:
        command: ["true"]
  - name: invalid-check
    probe: {}`))

	assert.Nil(t, validator)
	assert.EqualError(t, err, "check \"invalid-check\": probe must define exactly one supported action, found 0")
}

func TestLoadValidatorFromBytesIdentifiesMalformedPortCheck(t *testing.T) {
	validator, err := LoadValidatorFromBytes([]byte("name: example\nchecks:\n  - name: invalid-check\n    probe:\n      httpGet:\n        port: not-a-number\n"))

	assert.Nil(t, validator)
	assert.EqualError(t, err, "check \"invalid-check\": httpGet probe port must be an integer between 1 and 65535")
}
