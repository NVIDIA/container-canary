package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidator(t *testing.T) {
	assert := assert.New(t)

	validator, err := LoadValidatorFromFile("../../examples/kubeflow.yaml")

	assert.Nil(err)
	assert.Equal("kubeflow", validator.Name)
	assert.Equal("Kubeflow notebooks", validator.Description)

	assert.GreaterOrEqual(len(validator.Checks), 1)

	check := validator.Checks[0]

	assert.Equal("user", check.Name)
	assert.Equal("👩 User is jovyan", check.Description)

	assert.Equal(0, check.Probe.InitialDelaySeconds)

}
