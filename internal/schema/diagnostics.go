package schema

import "fmt"

func copiedKubernetesField(manifest any) error {
	root, ok := manifest.(map[string]any)
	if !ok {
		return nil
	}

	for _, field := range []string{"livenessProbe", "readinessProbe", "startupProbe"} {
		if _, found := root[field]; found {
			return fmt.Errorf("%s is a Kubernetes Pod field; Container Canary uses checks[].probe instead", field)
		}
	}

	if env, ok := root["env"].([]any); ok {
		for i, item := range env {
			if value, ok := item.(map[string]any); ok {
				if _, found := value["valueFrom"]; found {
					return fmt.Errorf("env[%d].valueFrom is a Kubernetes core/v1.EnvVar field, but Container Canary supports only name and value", i)
				}
			}
		}
	}

	if checks, ok := root["checks"].([]any); ok {
		for i, item := range checks {
			check, ok := item.(map[string]any)
			if !ok {
				continue
			}
			probe, ok := check["probe"].(map[string]any)
			if !ok {
				continue
			}
			if _, found := probe["grpc"]; found {
				return fmt.Errorf("checks[%d].probe.grpc is a Kubernetes probe action, but it is unsupported in container-canary.nvidia.com/v1; use exec, httpGet, or tcpSocket", i)
			}
		}
	}

	return nil
}
