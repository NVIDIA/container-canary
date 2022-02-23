package validator

import (
	"fmt"
	"time"

	canaryv1 "github.com/jacobtomlinson/containercanary/internal/apis/v1"
	"github.com/jacobtomlinson/containercanary/internal/container"
	"github.com/jacobtomlinson/containercanary/internal/terminal"
)

func Validate(image string, validator *canaryv1.Validator) (bool, error) { // Start image
	c := container.New(image, validator.Env, validator.Ports, validator.Volumes)
	c.Start()
	defer c.Remove()

	// Run checks
	var allChecks []bool

	if len(validator.Checks) == 0 {
		return false, fmt.Errorf("error no checks found")
	}

	// TODO Make checks async
	// TODO Retry each check on fail "failureThreshold" times with "periodSeconds" sleep between
	// TODO Retry each check on success "successThreshold" times with "periodSeconds" sleep between
	for _, check := range validator.Checks {
		time.Sleep(time.Duration(check.Probe.InitialDelaySeconds) * time.Second)
		if check.Probe.Exec != nil {
			passFail, err := ExecCheck(c, check.Probe.Exec)
			allChecks = append(allChecks, passFail)
			terminal.PrintCheckItem("", check.Description, getStatus(passFail, err))
			continue
		}
		if check.Probe.HTTPGet != nil {
			passFail, err := HTTPGetCheck(c, check.Probe.HTTPGet)
			allChecks = append(allChecks, passFail)
			terminal.PrintCheckItem("", check.Description, getStatus(passFail, err))
			continue
		}
		return false, fmt.Errorf("check '%s' has no known probes", check.Name)
	}
	// Clear up container
	return all(allChecks), nil
}
