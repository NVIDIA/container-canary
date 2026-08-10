package schema

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPublishedContractHasNoExternalReferences(t *testing.T) {
	contents, err := os.ReadFile(filepath.Join("..", "..", Path))
	require.NoError(t, err)

	assert.NotContains(t, string(contents), "\"$ref\": \"http")
	assert.NotContains(t, string(contents), "openapi")
	assert.NotContains(t, string(contents), "kubernetes.io")
	assert.Contains(t, string(contents), "https://json-schema.org/draft/2020-12/schema")
	assert.Contains(t, string(contents), "https://raw.githubusercontent.com/NVIDIA/container-canary/<release-tag>/schema/container-canary.nvidia.com/v1/validator.schema.json")
	assert.Contains(t, string(contents), "container-canary.nvidia.com/v1")
	assert.Contains(t, string(contents), "\"Validator\"")
}

func TestDocumentationPublishesReleasePinnedSchemaPattern(t *testing.T) {
	contents, err := os.ReadFile(filepath.Join("..", "..", "README.md"))
	require.NoError(t, err)

	docs := string(contents)
	assert.Contains(t, docs, "https://raw.githubusercontent.com/NVIDIA/container-canary/<release-tag>/schema/container-canary.nvidia.com/v1/validator.schema.json")
	assert.Contains(t, docs, "Do not use `main`")
	assert.Contains(t, docs, "core/v1.EnvVar")
	assert.Contains(t, docs, "core/v1.ServicePort")
	assert.Contains(t, docs, "core/v1.ExecAction")
	assert.Contains(t, docs, "core/v1.HTTPHeader")
	assert.Contains(t, docs, "env[0].valueFrom")
	assert.Contains(t, docs, "checks[0].probe.grpc")
	assert.Contains(t, docs, "livenessProbe")
	assert.Contains(t, docs, "Supported subset")

}
