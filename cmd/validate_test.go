package cmd

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidate(t *testing.T) {
	assert := assert.New(t)

	b := new(bytes.Buffer)
	rootCmd.SetOut(b)
	rootCmd.SetErr(b)
	rootCmd.SetArgs([]string{"validate", "--file", "../examples/kubeflow.yaml", "daskdev/dask-notebook:latest"})
	rootCmd.Execute()

	assert.Contains(b.String(), "Validating daskdev/dask-notebook:latest against kubeflow", "did not validate")
}

func TestFileDoesNotExist(t *testing.T) {
	assert := assert.New(t)

	b := new(bytes.Buffer)
	rootCmd.SetOut(b)
	rootCmd.SetErr(b)
	rootCmd.SetArgs([]string{"validate", "--file", "foo.yaml", "nginx"})
	rootCmd.Execute()

	assert.Contains(b.String(), "Error", "did not error")
}
