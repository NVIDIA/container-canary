package schema

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCopiedKubernetesFieldsExplainHowToAdapt(t *testing.T) {
	root := filepath.Join("..", "..")
	schemaPath := filepath.Join(root, Path)

	tests := []struct {
		name    string
		fixture string
		message string
	}{
		{"EnvVar valueFrom", "kubernetes-env-value-from.yaml", "env[0].valueFrom is a Kubernetes core/v1.EnvVar field, but Container Canary supports only name and value"},
		{"gRPC probe", "kubernetes-grpc.yaml", "checks[0].probe.grpc is a Kubernetes probe action, but it is unsupported in container-canary.nvidia.com/v1; use exec, httpGet, or tcpSocket"},
		{"Pod probe", "kubernetes-liveness-probe.yaml", "livenessProbe is a Kubernetes Pod field; Container Canary uses checks[].probe instead"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateFile(schemaPath, filepath.Join("testdata", tt.fixture))
			assert.EqualError(t, err, tt.message)
		})
	}
}
