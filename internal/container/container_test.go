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
	"strings"
	"testing"

	canaryv1 "github.com/jacobtomlinson/containercanary/internal/apis/v1"
	v1 "k8s.io/api/core/v1"
)

func TestContainer(t *testing.T) {
	env := []v1.EnvVar{
		{Name: "FOO", Value: "BAR"},
	}
	ports := []v1.ServicePort{
		{Port: 80, Protocol: "TCP"},
	}
	volumes := []canaryv1.Volume{
		{MountPath: "/foo"},
	}
	c := New("nginx", env, ports, volumes)

	err := c.Start()
	defer c.Remove()
	if err != nil {
		t.Errorf("Failed to start container: %s", err.Error())
	}

	_, err = c.Status()
	if err != nil {
		t.Errorf("Failed to inspect container: %s", err.Error())
	}

	uname, err := c.Exec("uname", "-a")
	if err != nil {
		t.Errorf("Failed to exec command in container: %s", err.Error())
	}
	if !strings.Contains(uname, "Linux") {
		t.Error("Output for command 'uname' did not contain expected string 'Linux'")
	}
}
func TestContainerRemoves(t *testing.T) {
	c := New("nginx", nil, nil, nil)

	err := c.Start()
	if err != nil {
		t.Errorf("Failed to start container: %s", err.Error())
	}

	err = c.Remove()
	if err != nil {
		t.Errorf("Failed to remove container: %s", err.Error())
	}
}
