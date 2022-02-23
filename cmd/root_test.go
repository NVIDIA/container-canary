package cmd

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRoot(t *testing.T) {
	assert := assert.New(t)

	b := new(bytes.Buffer)
	rootCmd.SetOut(b)
	rootCmd.SetErr(b)
	rootCmd.Execute()

	assert.Contains(b.String(), "Usage:", "did not how usage")
}
