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

package container

import (
	"fmt"

	"github.com/google/uuid"
	canaryv1 "github.com/nvidia/container-canary/internal/apis/v1"
	v1 "k8s.io/api/core/v1"
)

type ContainerState struct {
	Status  string
	Running bool
}

type ContainerInfo struct {
	Id         string
	State      ContainerState
	RunCommand string
}

type ContainerInterface interface {
	Start() error
	Remove() error
	Status() (*ContainerInfo, error)
	Exec(command ...string) (string, error)
	Logs() (string, error)
	GetStartupTimeout() int
}

func New(image string, env []v1.EnvVar, ports []v1.ServicePort, volumes []canaryv1.Volume, command []string, dockerRunOptions []string, startupTimeout int) ContainerInterface {
	name := fmt.Sprintf("%s%s", "canary-runner-", uuid.New().String()[:8])
	return &DockerContainer{Name: name, Image: image, Command: command, Env: env, Ports: ports, Volumes: volumes, RunOptions: dockerRunOptions, StartupTimeout: startupTimeout}
}
