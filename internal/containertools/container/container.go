package container

import (
	"fmt"
	"os/exec"

	"github.com/google/uuid"
)

type Container struct {
	Name    string
	Id      string
	Image   string
	Runtime string
	Command string
}

func (c Container) Start() error {
	// Start a container
	id, err := exec.Command(c.Runtime, "run", "-d", "--name", c.Name, c.Image).Output()

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

func New(image string) Container {
	name := fmt.Sprintf("%s%s", "canairy-runner-", uuid.New().String()[:8])
	return Container{Name: name, Image: image, Runtime: "docker", Command: ""}
}
