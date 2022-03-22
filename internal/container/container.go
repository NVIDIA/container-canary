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
	"encoding/json"
	"errors"
	"fmt"
	"os/exec"
	"time"

	"github.com/google/uuid"
	canaryv1 "github.com/nvidia/container-canary/internal/apis/v1"
	v1 "k8s.io/api/core/v1"
)

type ContainerState struct {
	Status     string
	Running    bool
	Paused     bool
	Restarting bool
	OOMKilled  bool
	Dead       bool
	Pid        int
	ExitCode   int
	Error      string
	StartedAt  string
	FinishedAt string
}

type ContainerInfo struct {
	Id    string
	State ContainerState
}

type Container struct {
	Name    string
	Id      string
	Image   string
	Runtime string
	Command string
	Env     []v1.EnvVar
	Ports   []v1.ServicePort
	Volumes []canaryv1.Volume
}

// Start a container
func (c Container) Start() error {

	commandArgs := []string{"run", "-d"}

	commandArgs = append(commandArgs, "--name", c.Name)

	for _, e := range c.Env {
		commandArgs = append(commandArgs, "-e", fmt.Sprintf("%s=%s", e.Name, e.Value))
	}

	for _, p := range c.Ports {
		commandArgs = append(commandArgs, "-p", fmt.Sprintf("%d:%d/%s", p.Port, p.Port, p.Protocol))
	}

	for _, v := range c.Volumes {
		if v.Path != "" {
			commandArgs = append(commandArgs, "-v", fmt.Sprintf("%s:%s", v.Path, v.MountPath))
		} else {
			commandArgs = append(commandArgs, "-v", v.MountPath)
		}
	}

	commandArgs = append(commandArgs, c.Image)

	if c.Command != "" {
		commandArgs = append(commandArgs, c.Command)
	}

	_, err := exec.Command(c.Runtime, commandArgs...).Output()

	for {
		info, err := c.Status()
		if err != nil {
			return err
		}
		if info.State.Status == "exited" {
			return errors.New("container failed to start")
		}
		if info.State.Running {
			break
		}
		time.Sleep(time.Second)
	}

	if err != nil {
		return err
	}

	return nil
}

// Remove a container
func (c Container) Remove() error {
	_, err := exec.Command(c.Runtime, "rm", "-f", c.Name).Output()
	return err
}

// Get container status
func (c Container) Status() (*ContainerInfo, error) {

	output, err := exec.Command(c.Runtime, "inspect", c.Name).Output()

	if err != nil {
		return nil, err
	}

	var info []ContainerInfo

	err = json.Unmarshal(output, &info)

	if err != nil {
		return nil, err
	}

	if len(info) != 1 {
		return nil, fmt.Errorf("expected 1 container, got %d", len(info))
	}

	return &info[0], nil
}

// Exec a command inside a container
func (c Container) Exec(command ...string) (string, error) {

	args := append([]string{"exec", c.Name}, command...)
	out, err := exec.Command(c.Runtime, args...).Output()
	return string(out), err
}

// Get container logs
func (c Container) Logs() (string, error) {
	out, err := exec.Command(c.Runtime, "logs", c.Name).Output()
	return string(out), err
}

func New(image string, env []v1.EnvVar, ports []v1.ServicePort, volumes []canaryv1.Volume, command string) Container {
	name := fmt.Sprintf("%s%s", "canary-runner-", uuid.New().String()[:8])
	return Container{Name: name, Image: image, Runtime: "docker", Command: command, Env: env, Ports: ports, Volumes: volumes}
}
