package container

import (
	"strings"
	"testing"
)

func TestContainer(t *testing.T) {
	c := New("nginx")

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
	c := New("nginx")

	err := c.Start()
	if err != nil {
		t.Errorf("Failed to start container: %s", err.Error())
	}

	err = c.Remove()
	if err != nil {
		t.Errorf("Failed to remove container: %s", err.Error())
	}
}
