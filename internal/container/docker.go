package container

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os/exec"
	"strings"
	"time"

	canaryv1 "github.com/nvidia/container-canary/internal/apis/v1"
	v1 "k8s.io/api/core/v1"
)

type DockerContainer struct {
	Name       string
	Id         string
	Image      string
	Command    []string
	Env        []v1.EnvVar
	Ports      []v1.ServicePort
	Volumes    []canaryv1.Volume
	runCommand string
}

// Start a container
func (c *DockerContainer) Start() error {

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

	if len(c.Command) > 0 {
		commandArgs = append(commandArgs, c.Command...)
	}

	if err := CheckForDocker(); err != nil {
		return err
	}
	_, err := exec.Command("docker", commandArgs...).Output()
	c.runCommand = fmt.Sprintf("docker %s", strings.Join(commandArgs, " "))

	for startTime := time.Now(); ; {
		info, err := c.Status()
		if err != nil {
			return err
		}
		if info.State.Status == "exited" {
			err := c.Remove()
			if err != nil {
				return err
			}
			return errors.New("container failed to start")
		}
		if info.State.Running {
			break
		}
		if time.Since(startTime) > (time.Second * 10) {
			err := c.Remove()
			if err != nil {
				return err
			}
			return errors.New("container failed to start after 10 seconds")
		}
		time.Sleep(time.Second)
	}

	if err != nil {
		return err
	}

	return nil
}

func CheckForDocker() error {
	if _, err := exec.LookPath("docker"); err != nil {
		return errors.New("Docker is missing")
	}

	if _, err := exec.Command("docker", "ps").Output(); err != nil {
		return errors.New("Docker requires root privileges to run")
	}
	return nil
}

// Remove a container
func (c DockerContainer) Remove() error {
	_, err := exec.Command("docker", "rm", "-f", c.Name).Output()
	return err
}

// Get container status
func (c DockerContainer) Status() (*ContainerInfo, error) {

	output, err := exec.Command("docker", "inspect", c.Name).Output()

	if err != nil {
		return nil, err
	}

	var infoList []ContainerInfo

	err = json.Unmarshal(output, &infoList)

	if err != nil {
		return nil, err
	}

	if len(infoList) != 1 {
		return nil, fmt.Errorf("expected 1 container, got %d", len(infoList))
	}

	info := infoList[0]

	info.RunCommand = c.runCommand

	return &info, nil
}

// Exec a command inside a container
func (c DockerContainer) Exec(command ...string) (exitCode int, stdout string, stderr string, err error) {

	args := append([]string{"exec", c.Name}, command...)
	cmd := exec.Command("docker", args...)
	stdoutPipe, err := cmd.StdoutPipe()
	if err != nil {
		return
	}
	stderrPipe, err := cmd.StderrPipe()
	if err != nil {
		return
	}

	if err = cmd.Start(); err != nil {
		return
	}

	stdoutBytes, err := io.ReadAll(stdoutPipe)
	if err != nil {
		return
	}
	stderrBytes, err := io.ReadAll(stderrPipe)
	if err != nil {
		return
	}

	if err = cmd.Wait(); err != nil {
		switch t := err.(type) {
		case *exec.ExitError:
			exitCode = t.ExitCode()
		default:
			return
		}
	}

	return exitCode, string(stdoutBytes), string(stderrBytes), nil
}

// Get container logs
func (c DockerContainer) Logs() (string, error) {
	out, err := exec.Command("docker", "logs", c.Name).Output()
	return string(out), err
}
