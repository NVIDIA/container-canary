/*
* SPDX-FileCopyrightText: Copyright (c) <2026> NVIDIA CORPORATION & AFFILIATES.
* SPDX-License-Identifier: Apache-2.0
 */

package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
	yaml "gopkg.in/yaml.v2"
)

func TestProbeUnmarshalYAMLValidatesActions(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr string
	}{
		{name: "no action", input: `{}`, wantErr: "probe must define exactly one supported action, found 0"},
		{name: "null action", input: "exec: null", wantErr: "probe must define exactly one supported action, found 0"},
		{name: "unknown action", input: "unknown: value", wantErr: "probe must define exactly one supported action, found 0"},
		{name: "multiple actions", input: "exec:\n  command: [true]\ntcpSocket:\n  port: 8080", wantErr: "probe must define exactly one supported action, found 2"},
		{name: "empty exec", input: "exec: {}", wantErr: "exec probe must define at least one command"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var probe Probe
			assert.EqualError(t, yaml.Unmarshal([]byte(tt.input), &probe), tt.wantErr)
		})
	}
}

func TestCheckUnmarshalYAMLRejectsMissingProbe(t *testing.T) {
	var check Check
	assert.EqualError(t, yaml.Unmarshal([]byte("name: test"), &check), "check \"test\": probe must define exactly one supported action, found 0")
}

func TestProbeUnmarshalYAMLAcceptsValidActionAndAppliesDefaults(t *testing.T) {
	var probe Probe
	err := yaml.Unmarshal([]byte("exec:\n  command: [true]"), &probe)

	assert.NoError(t, err)
	assert.NotNil(t, probe.Exec)
	assert.Equal(t, 30, probe.TimeoutSeconds)
	assert.Equal(t, 1, probe.PeriodSeconds)
}

func TestProbeUnmarshalYAMLAcceptsAliasedAction(t *testing.T) {
	var probe Probe
	err := yaml.Unmarshal([]byte("defaults: &http\n  port: 8080\nhttpGet: *http"), &probe)

	assert.NoError(t, err)
	assert.Equal(t, 8080, probe.HTTPGet.Port)
}

func TestHTTPGetActionUnmarshalYAMLValidatesPort(t *testing.T) {
	testActionPort(t, "httpGet probe port must be between 1 and 65535", func(input string) error {
		var action HTTPGetAction
		return yaml.Unmarshal([]byte(input), &action)
	})
}

func TestTCPSocketActionUnmarshalYAMLValidatesPort(t *testing.T) {
	testActionPort(t, "tcpSocket probe port must be between 1 and 65535", func(input string) error {
		var action TCPSocketAction
		return yaml.Unmarshal([]byte(input), &action)
	})
}

func testActionPort(t *testing.T, rangeError string, unmarshal func(string) error) {
	t.Helper()
	for _, tt := range []struct {
		name    string
		input   string
		wantErr string
		invalid bool
	}{
		{name: "missing", input: "{}", wantErr: rangeError},
		{name: "zero", input: "port: 0", wantErr: rangeError},
		{name: "negative", input: "port: -1", wantErr: rangeError},
		{name: "too large", input: "port: 65536", wantErr: rangeError},
		{name: "non-numeric", input: "port: not-a-number", invalid: true},
		{name: "minimum", input: "port: 1"},
		{name: "maximum", input: "port: 65535"},
	} {
		t.Run(tt.name, func(t *testing.T) {
			err := unmarshal(tt.input)
			if tt.wantErr != "" {
				assert.EqualError(t, err, tt.wantErr)
			} else if tt.invalid {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
