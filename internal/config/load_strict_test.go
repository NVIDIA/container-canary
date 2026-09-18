package config

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoadValidatorRejectsUnknownFields(t *testing.T) {
	tests := []struct {
		name     string
		manifest string
		field    string
	}{
		{"validator", "name: example\ndescripton: misspelled\nchecks:\n  - name: command\n    probe:\n      exec:\n        command: [\"true\"]\n", "descripton"},
		{"check", "checks:\n  - name: command\n    descripton: misspelled\n", "descripton"},
		{"probe", "checks:\n  - probe:\n      timeoutSecond: 10\n", "timeoutSecond"},
		{"exec", "checks:\n  - probe:\n      exec:\n        commmand: [true]\n", "commmand"},
		{"HTTP action", "checks:\n  - probe:\n      httpGet:\n        paht: /health\n", "paht"},
		{"request header", "checks:\n  - probe:\n      httpGet:\n        httpHeaders:\n          - name: Accept\n            vale: application/json\n", "vale"},
		{"response header", "checks:\n  - probe:\n      httpGet:\n        responseHttpHeaders:\n          - name: Content-Type\n            vale: application/json\n", "vale"},
		{"TCP action", "checks:\n  - probe:\n      tcpSocket:\n        prot: 8080\n", "prot"},
		{"environment", "env:\n  - name: GREETING\n    vale: hello\n", "vale"},
		{"port", "ports:\n  - port: 8080\n    protcol: TCP\n", "protcol"},
		{"volume", "volumes:\n  - mountPaht: /data\n", "mountPaht"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			validator, err := LoadValidatorFromBytes([]byte(tt.manifest))
			require.Error(t, err)
			assert.Nil(t, validator)
			assert.ErrorContains(t, err, "field "+tt.field+" not found")
			assert.ErrorContains(t, err, "line ")
		})
	}
}

func TestLoadValidatorRejectsDuplicateFields(t *testing.T) {
	validator, err := LoadValidatorFromBytes([]byte("name: first\nname: second\n"))
	require.Error(t, err)
	assert.Nil(t, validator)
	assert.ErrorContains(t, err, "field name already set")
}

func TestLoadValidatorSourcesRejectUnknownFields(t *testing.T) {
	const manifest = "name: example\ndescripton: misspelled\n"
	path := filepath.Join(t.TempDir(), "validator.yaml")
	require.NoError(t, os.WriteFile(path, []byte(manifest), 0600))
	t.Run("file", func(t *testing.T) {
		validator, err := LoadValidatorFromFile(path)
		require.Error(t, err)
		assert.Nil(t, validator)
		assert.ErrorContains(t, err, "line 2: field descripton not found")
	})
	t.Run("URL", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			_, _ = w.Write([]byte(manifest))
		}))
		defer server.Close()

		validator, err := LoadValidatorFromURL(server.URL)
		require.Error(t, err)
		assert.Nil(t, validator)
		assert.ErrorContains(t, err, "line 2: field descripton not found")
	})
}

func TestLoadValidatorExamples(t *testing.T) {
	paths, err := filepath.Glob("../../examples/*.yaml")
	require.NoError(t, err)
	require.NotEmpty(t, paths)
	for _, path := range paths {
		t.Run(filepath.Base(path), func(t *testing.T) {
			validator, err := LoadValidatorFromFile(path)
			require.NoError(t, err)
			assert.NotNil(t, validator)
		})
	}
}

func TestLoadValidatorPreservesProbeDefaultsAndValues(t *testing.T) {
	const manifest = `name: example
env:
  - name: GREETING
    value: "hello: world"
checks:
  - name: defaults
    probe:
      exec:
        command: ["printf", "hello: world"]
  - name: overrides
    probe:
      timeoutSeconds: 0
      periodSeconds: 5
      httpGet:
        port: 8080
        httpHeaders:
          - name: X-Greeting
            value: "hello: world"
`
	validator, err := LoadValidatorFromBytes([]byte(manifest))
	require.NoError(t, err)
	require.Len(t, validator.Checks, 2)
	defaults := validator.Checks[0].Probe
	assert.Equal(t, 30, defaults.TimeoutSeconds)
	assert.Equal(t, 1, defaults.PeriodSeconds)
	assert.Equal(t, 1, defaults.SuccessThreshold)
	assert.Equal(t, 1, defaults.FailureThreshold)
	assert.Equal(t, 30, defaults.TerminationGracePeriodSeconds)
	require.NotNil(t, defaults.Exec)
	assert.Equal(t, []string{"printf", "hello: world"}, defaults.Exec.Command)
	overrides := validator.Checks[1].Probe
	assert.Zero(t, overrides.TimeoutSeconds)
	assert.Equal(t, 5, overrides.PeriodSeconds)
	require.NotNil(t, overrides.HTTPGet)
	require.Len(t, overrides.HTTPGet.HTTPHeaders, 1)
	assert.Equal(t, "hello: world", overrides.HTTPGet.HTTPHeaders[0].Value)
	require.Len(t, validator.Env, 1)
	assert.Equal(t, "hello: world", validator.Env[0].Value)
}
