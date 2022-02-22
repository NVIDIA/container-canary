package validator

import (
	"fmt"
	"time"

	"github.com/jacobtomlinson/containercanairy/internal/apis/config"
	"github.com/jacobtomlinson/containercanairy/internal/checks"
	"github.com/jacobtomlinson/containercanairy/internal/containertools/container"
	"github.com/jacobtomlinson/containercanairy/internal/terminal"
)

func Validate(image string, validator *config.Validator) (bool, error) { // Start image
	c := container.New(image, validator.Env, validator.Ports, validator.Volumes)
	c.Start()
	defer c.Remove()

	// Run checks
	var allChecks []bool

	if len(validator.Checks) == 0 {
		return false, fmt.Errorf("error no checks found")
	}

	// TODO Make checks async
	for _, check := range validator.Checks {
		time.Sleep(time.Duration(check.Probe.InitialDelaySeconds) * time.Second)
		if check.Probe.Exec != nil {
			passFail, err := checks.ExecCheck(c, check.Probe.Exec)
			allChecks = append(allChecks, passFail)
			terminal.PrintCheckItem("", check.Description, getStatus(passFail, err))
			continue
		}
		if check.Probe.HTTPGet != nil {
			passFail, err := checks.HTTPGetCheck(c, check.Probe.HTTPGet)
			allChecks = append(allChecks, passFail)
			terminal.PrintCheckItem("", check.Description, getStatus(passFail, err))
			continue
		}
		return false, fmt.Errorf("check '%s' has no known probes", check.Name)
	}
	// Clear up container
	return all(allChecks), nil
}
