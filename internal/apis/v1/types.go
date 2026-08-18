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
	"fmt"

	v1 "k8s.io/api/core/v1"
)

// Validator contains validator specification
type Validator struct {
	// The validator name.
	// +optional
	Name string

	// The validator description.
	// +optional
	Description string

	// A link to the documentation where the requirements are defined.
	// +optional
	Documentation string

	// A list of checks to perform validation against.
	Checks []Check

	// A list of environment variables to set on the container before starting.
	Env []v1.EnvVar

	// A list of ports to expose on the container.
	Ports []v1.ServicePort

	// A list of volumes to mount on the container.
	Volumes []Volume

	// A command to run in the container
	// +optional
	Command []string

	// Additional flags to pass to the docker CLI.
	// +optional
	DockerRunOptions []string `yaml:"dockerRunOptions"`
}

type Check struct {
	// The check name.
	// +optional
	Name string

	// The check description.
	// +optional
	Description string

	// A probe to run.
	Probe Probe
}

func (c *Check) UnmarshalYAML(unmarshal func(interface{}) error) error {
	raw := rawCheck{Probe: newRawProbe()}
	if err := unmarshal(&raw); err != nil {
		return err
	}
	probe, err := raw.Probe.probe()
	if err != nil {
		return fmt.Errorf("check %q: %w", raw.Name, err)
	}
	if err := probe.validate(); err != nil {
		return fmt.Errorf("check %q: %w", raw.Name, err)
	}

	*c = Check{Name: raw.Name, Description: raw.Description, Probe: probe}
	return nil
}

type Probe struct {
	InitialDelaySeconds int `yaml:"initialDelaySeconds"`

	TimeoutSeconds int `yaml:"timeoutSeconds"`

	PeriodSeconds int `yaml:"periodSeconds"`

	SuccessThreshold int `yaml:"successThreshold"`

	FailureThreshold int `yaml:"failureThreshold"`

	TerminationGracePeriodSeconds int `yaml:"terminationGracePeriodSeconds"`

	Exec *v1.ExecAction `yaml:"exec"`

	HTTPGet *HTTPGetAction `yaml:"httpGet"`

	TCPSocket *TCPSocketAction `yaml:"tcpSocket" `
}

func (p *Probe) UnmarshalYAML(unmarshal func(interface{}) error) error {
	raw := newRawProbe()
	if err := unmarshal(&raw); err != nil {
		return err
	}
	probe, err := raw.probe()
	if err != nil {
		return err
	}
	if err := probe.validate(); err != nil {
		return err
	}

	*p = probe
	return nil
}

func (p Probe) validate() error {
	actionCount := 0
	if p.Exec != nil {
		actionCount++
	}
	if p.HTTPGet != nil {
		actionCount++
	}
	if p.TCPSocket != nil {
		actionCount++
	}
	if actionCount != 1 {
		return fmt.Errorf("probe must define exactly one supported action, found %d", actionCount)
	}
	if p.Exec != nil && len(p.Exec.Command) == 0 {
		return fmt.Errorf("exec probe must define at least one command")
	}
	if p.HTTPGet != nil {
		if err := p.HTTPGet.validate(); err != nil {
			return err
		}
	}
	if p.TCPSocket != nil {
		if err := p.TCPSocket.validate(); err != nil {
			return err
		}
	}

	return nil
}

type HTTPGetAction struct {
	// Path to access on the HTTP server.
	// +optional
	Path string `yaml:"path,omitempty"`
	// Number of the port to access on the container.
	// Number must be in the range 1 to 65535.
	Port int `yaml:"port"`
	// Scheme to use for connecting to the host.
	// Defaults to HTTP.
	// +optional
	Scheme v1.URIScheme `yaml:"scheme,omitempty"`
	// Custom headers to set in the request. HTTP allows repeated headers.
	// +optional
	HTTPHeaders []v1.HTTPHeader `yaml:"httpHeaders,omitempty"`
	// Headers expected in the response. Check will fail if any are missing.
	// +optional
	ResponseHTTPHeaders []v1.HTTPHeader `yaml:"responseHttpHeaders,omitempty"`
}

func (a *HTTPGetAction) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var raw rawHTTPGetAction
	if err := unmarshal(&raw); err != nil {
		return err
	}
	action, err := raw.action()
	if err != nil {
		return err
	}

	*a = action
	return nil
}

func (a HTTPGetAction) validate() error {
	if a.Port < 1 || a.Port > 65535 {
		return fmt.Errorf("httpGet probe port must be between 1 and 65535")
	}
	return nil
}

type TCPSocketAction struct {
	// Number or name of the port to access on the container.
	// Number must be in the range 1 to 65535.
	Port int `yaml:"port"`
}

func (a *TCPSocketAction) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var raw rawTCPSocketAction
	if err := unmarshal(&raw); err != nil {
		return err
	}
	action, err := raw.action()
	if err != nil {
		return err
	}

	*a = action
	return nil
}

func (a TCPSocketAction) validate() error {
	if a.Port < 1 || a.Port > 65535 {
		return fmt.Errorf("tcpSocket probe port must be between 1 and 65535")
	}
	return nil
}

type rawCheck struct {
	Name        string
	Description string
	Probe       rawProbe
}

type rawProbe struct {
	InitialDelaySeconds           int                 `yaml:"initialDelaySeconds"`
	TimeoutSeconds                int                 `yaml:"timeoutSeconds"`
	PeriodSeconds                 int                 `yaml:"periodSeconds"`
	SuccessThreshold              int                 `yaml:"successThreshold"`
	FailureThreshold              int                 `yaml:"failureThreshold"`
	TerminationGracePeriodSeconds int                 `yaml:"terminationGracePeriodSeconds"`
	Exec                          *v1.ExecAction      `yaml:"exec"`
	HTTPGet                       *rawHTTPGetAction   `yaml:"httpGet"`
	TCPSocket                     *rawTCPSocketAction `yaml:"tcpSocket"`
}

type rawHTTPGetAction struct {
	Path                string          `yaml:"path,omitempty"`
	Port                interface{}     `yaml:"port"`
	Scheme              v1.URIScheme    `yaml:"scheme,omitempty"`
	HTTPHeaders         []v1.HTTPHeader `yaml:"httpHeaders,omitempty"`
	ResponseHTTPHeaders []v1.HTTPHeader `yaml:"responseHttpHeaders,omitempty"`
}

type rawTCPSocketAction struct {
	Port interface{} `yaml:"port"`
}

func newRawProbe() rawProbe {
	return rawProbe{
		TimeoutSeconds:                30,
		PeriodSeconds:                 1,
		SuccessThreshold:              1,
		FailureThreshold:              1,
		TerminationGracePeriodSeconds: 30,
	}
}

func (p rawProbe) probe() (Probe, error) {
	probe := Probe{
		InitialDelaySeconds:           p.InitialDelaySeconds,
		TimeoutSeconds:                p.TimeoutSeconds,
		PeriodSeconds:                 p.PeriodSeconds,
		SuccessThreshold:              p.SuccessThreshold,
		FailureThreshold:              p.FailureThreshold,
		TerminationGracePeriodSeconds: p.TerminationGracePeriodSeconds,
		Exec:                          p.Exec,
	}
	if p.HTTPGet != nil {
		action, err := p.HTTPGet.action()
		if err != nil {
			return Probe{}, err
		}
		probe.HTTPGet = &action
	}
	if p.TCPSocket != nil {
		action, err := p.TCPSocket.action()
		if err != nil {
			return Probe{}, err
		}
		probe.TCPSocket = &action
	}
	return probe, nil
}

func (a rawHTTPGetAction) action() (HTTPGetAction, error) {
	port, err := decodePort(a.Port, "httpGet")
	if err != nil {
		return HTTPGetAction{}, err
	}
	return HTTPGetAction{Path: a.Path, Port: port, Scheme: a.Scheme, HTTPHeaders: a.HTTPHeaders, ResponseHTTPHeaders: a.ResponseHTTPHeaders}, nil
}

func (a rawTCPSocketAction) action() (TCPSocketAction, error) {
	port, err := decodePort(a.Port, "tcpSocket")
	if err != nil {
		return TCPSocketAction{}, err
	}
	return TCPSocketAction{Port: port}, nil
}

func decodePort(raw interface{}, action string) (int, error) {
	if raw == nil {
		return 0, fmt.Errorf("%s probe port must be between 1 and 65535", action)
	}
	port, ok := raw.(int)
	if !ok {
		return 0, fmt.Errorf("%s probe port must be an integer between 1 and 65535", action)
	}
	if port < 1 || port > 65535 {
		return 0, fmt.Errorf("%s probe port must be between 1 and 65535", action)
	}
	return port, nil
}

type Volume struct {
	// Path to mount in the container
	MountPath string `yaml:"mountPath,omitempty"`

	// Path to mount from host, will use empty volume if omitted
	// +optional
	Path string `yaml:"path,omitempty"`
}
