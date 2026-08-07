package schema

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSchemaCompilesAndExamplesValidate(t *testing.T) {
	root := filepath.Join("..", "..")
	sch, err := Compile(filepath.Join(root, Path))
	require.NoError(t, err)

	examples, err := filepath.Glob(filepath.Join(root, "examples", "*.yaml"))
	require.NoError(t, err)
	require.NotEmpty(t, examples)
	for _, example := range examples {
		err := ValidateFile(filepath.Join(root, Path), example)
		require.NoError(t, err, example)
	}

	require.NotNil(t, sch)
}
