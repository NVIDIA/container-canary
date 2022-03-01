// Copyright (c) 2022, NVIDIA CORPORATION.

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
