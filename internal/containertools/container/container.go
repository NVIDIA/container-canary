package container

import (
	"fmt"
	"os/exec"

	"github.com/google/uuid"
	"github.com/jacobtomlinson/containercanary/internal/apis/config"
	v1 "k8s.io/api/core/v1"
)

type Container struct {
	Name    string
	Id      string
	Image   string
	Runtime string
	Command string
	Env     []v1.EnvVar
	Ports   []v1.ServicePort
	Volumes []config.Volume
}

func (c Container) Start() error {
	// Start a container
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

	id, err := exec.Command(c.Runtime, commandArgs...).Output()

	if err != nil {
		return err
	}

	c.Id = string(id)
	return nil
}

func (c Container) Remove() error {
	// Remove a container
	_, err := exec.Command(c.Runtime, "rm", "-f", c.Name).Output()
	return err
}

func (c Container) Status() (string, error) {
	// Get container status
	_, err := exec.Command(c.Runtime, "inspect", c.Name).Output()

	if err != nil {
		return "Error", err
	}
	return "Created", nil
}

func (c Container) Exec(command ...string) (string, error) {
	// Exec a command inside a container
	args := append([]string{"exec", c.Name}, command...)
	out, err := exec.Command(c.Runtime, args...).Output()
	return string(out), err
}

func New(image string, env []v1.EnvVar, ports []v1.ServicePort, volumes []config.Volume) Container {
	name := fmt.Sprintf("%s%s", "canary-runner-", uuid.New().String()[:8])
	return Container{Name: name, Image: image, Runtime: "docker", Command: "", Env: env, Ports: ports, Volumes: volumes}
}
