# Examples

This directory contains example validation manifests. You can download and adapt these manifests or use them directly from GitHub.

For example the Kubeflow example from the README can be used directly like this.

```console
$ canary validate --file https://raw.githubusercontent.com/NVIDIA/container-canary/main/examples/kubeflow.yaml public.ecr.aws/j1r0q0g6/notebooks/notebook-servers/jupyter-scipy:v1.5.0-rc.1
Validating public.ecr.aws/j1r0q0g6/notebooks/notebook-servers/jupyter-scipy:v1.5.0-rc.1 against kubeflow
 👩 User is jovyan                                   [true]
 🆔 User ID is 1000                                  [true]
 🏠 Home directory is /home/jovyan                   [true]
 🌏 Exposes an HTTP interface on port 8888           [true]
 🧭 Correctly routes the NB_PREFIX                   [true]
 🔓 Sets 'Access-Control-Allow-Origin: *' header     [true]
validation passed
```

[Contributing](../CONTRIBUTING.md) more manifests here is highly encouraged!
