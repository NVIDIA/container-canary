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

import v1 "k8s.io/api/core/v1"

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

type Probe struct {
	InitialDelaySeconds int `yaml:"initialDelaySeconds"`

	TimeoutSeconds int `yaml:"timeoutSeconds"`

	PeriodSeconds int `yaml:"periodSeconds"`

	SuccessThreshold int `yaml:"successThreshold"`

	FailureThreshold int `yaml:"failureThreshold"`

	TerminationGracePeriodSeconds *int `yaml:"terminationGracePeriodSeconds"`

	Exec *v1.ExecAction `yaml:"exec"`

	HTTPGet *HTTPGetAction `yaml:"httpGet"`
}

type HTTPGetAction struct {
	// Path to access on the HTTP server.
	// +optional
	Path string `json:"path,omitempty"`
	// Number of the port to access on the container.
	// Number must be in the range 1 to 65535.
	Port int `json:"port"`
	// Host name to connect to, defaults to the pod IP. You probably want to set
	// "Host" in httpHeaders instead.
	// +optional
	Host string `json:"host,omitempty"`
	// Scheme to use for connecting to the host.
	// Defaults to HTTP.
	// +optional
	Scheme v1.URIScheme `json:"scheme,omitempty"`
	// Custom headers to set in the request. HTTP allows repeated headers.
	// +optional
	HTTPHeaders []v1.HTTPHeader `json:"httpHeaders,omitempty"`
	// Headers expected in the response. Check will fail if any are missing.
	// +optional
	ResponseHTTPHeaders []v1.HTTPHeader `json:"responseHttpHeaders,omitempty"`
}

type Volume struct {
	// Path to mount in the container
	MountPath string `yaml:"mountPath,omitempty"`

	// Path to mount from host, will use empty volume if omitted
	// +optional
	Path string `yaml:"path,omitempty"`
}
