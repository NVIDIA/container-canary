package config

import v1 "k8s.io/api/core/v1"

// Validator contains validator specification
type Validator struct {
	// The validator name.
	// Optional
	Name string

	// The validator description.
	// Optional
	Description string

	// A list of checks to perform validation against.
	Checks []Check
}

type Check struct {
	// The check name.
	// Optional
	Name string

	// The check description.
	// Optional
	Description string

	// A probe to use to test the container.
	Probe v1.Probe
}
