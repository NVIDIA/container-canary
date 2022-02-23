package cmd

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestVersion(t *testing.T) {
	assert := assert.New(t)

	b := new(bytes.Buffer)
	rootCmd.SetOut(b)
	rootCmd.SetErr(b)
	rootCmd.SetArgs([]string{"version"})
	rootCmd.Execute()

	assert.Contains(b.String(), "Version:", "mission version")
	assert.Contains(b.String(), "Go Version:", "missing go version")
	assert.Contains(b.String(), "Commit:", "missing commit hash")
	assert.Contains(b.String(), "OS/Arch:", "missing OS and CPU arch")
	assert.Contains(b.String(), "Built:", "mission build time")
}
