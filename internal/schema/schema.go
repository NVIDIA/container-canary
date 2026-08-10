// Package schema validates Validator manifests against the published schema.
package schema

import (
	"encoding/json"
	"fmt"
	"os"

	jsonschema "github.com/santhosh-tekuri/jsonschema/v6"
	"gopkg.in/yaml.v3"
)

const Path = "schema/container-canary.nvidia.com/v1/validator.schema.json"

// ValidateFile validates a YAML or JSON Validator manifest without loading it
// through the Container Canary runtime.
func ValidateFile(schemaPath, manifestPath string) error {
	sch, err := Compile(schemaPath)
	if err != nil {
		return err
	}

	contents, err := os.ReadFile(manifestPath)
	if err != nil {
		return err
	}

	var manifest any
	if err := yaml.Unmarshal(contents, &manifest); err != nil {
		return fmt.Errorf("parse manifest: %w", err)
	}

	if err := copiedKubernetesField(manifest); err != nil {
		return err
	}

	if err := sch.Validate(manifest); err != nil {
		return fmt.Errorf("%s: %w", manifestPath, err)
	}
	return nil
}

// Compile compiles the schema at schemaPath and validates the schema itself.
func Compile(schemaPath string) (*jsonschema.Schema, error) {
	contents, err := os.ReadFile(schemaPath)
	if err != nil {
		return nil, err
	}

	var document any
	if err := json.Unmarshal(contents, &document); err != nil {
		return nil, fmt.Errorf("parse schema: %w", err)
	}

	compiler := jsonschema.NewCompiler()
	if err := compiler.AddResource("file://"+schemaPath, document); err != nil {
		return nil, fmt.Errorf("add schema: %w", err)
	}
	sch, err := compiler.Compile("file://" + schemaPath)
	if err != nil {
		return nil, fmt.Errorf("compile schema: %w", err)
	}
	return sch, nil
}
