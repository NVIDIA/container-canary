package platforms

import (
	"fmt"

	"github.com/jacobtomlinson/containercanairy/internal/checks"
	"github.com/jacobtomlinson/containercanairy/internal/containertools/container"
	"github.com/jacobtomlinson/containercanairy/internal/terminal"
)

// Validate image against Kubeflow requirements
// https://www.kubeflow.org/docs/components/notebooks/container-images/#custom-images
func ValidateKubeflow(image string) bool {
	fmt.Printf("Validating %s for Kubeflow Notebooks\n", image)

	// Start image
	c := container.New(image)
	c.Start()
	defer c.Remove()

	// Run checks
	var allChecks []bool

	// Network checks
	// expose an HTTP interface on port 8888
	// web service honours NB_PREFIX
	// Responses contain Access-Control-Allow-Origin: * header

	// Container checks

	// username is jovyan
	check, err := checks.CheckUser(c, "jovyan")
	allChecks = append(allChecks, check)
	terminal.PrintCheckItem("👩", "User is jovyan", getStatus(check, err))

	// home directory is /home/jovyan
	check, err = checks.CheckPWD(c, "/home/jovyan")
	allChecks = append(allChecks, check)
	terminal.PrintCheckItem("🏠", "Home directory is /home/jovyan", getStatus(check, err))

	// /home/jovyan is empty
	// TODO

	// UID is 1000
	check, err = checks.CheckUID(c, "1000")
	allChecks = append(allChecks, check)
	terminal.PrintCheckItem("🆔", "UID is 1000", getStatus(check, err))

	// Clear up container
	return all(allChecks)

}
