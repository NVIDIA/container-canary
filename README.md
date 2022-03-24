# Container Canary

[![Test](https://github.com/NVIDIA/container-canary/actions/workflows/test.yaml/badge.svg)](https://github.com/NVIDIA/container-canary/actions/workflows/test.yaml)
![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/nvidia/container-canary)
![GitHub tag (latest SemVer)](https://img.shields.io/github/v/tag/nvidia/container-canary?label=version)

A little bird to validate your container images.

```console
$ canary validate --file somespec.yaml foo/bar:latest
Validating foo/bar:latest against somespec
 📦 Required packages are installed                  [true]
 🤖 Expected services are running                    [true]
 🎉 Your container is awesome                        [true]
validation passed
```

Many modern compute platforms support bring-your-own-container models where the user can provide container images with their custom software environment. However platforms commonly have a set of requirements that the container must conform to, such as using a non-root user, having the home directory in a specific location, having certain packages installed or running web applications on specific ports.

Container Canary is a tool for recording those requirements as a manifest that can be versioned and then validating containers against that manifest. This is particularly useful in CI environments to avoid regressions in containers.

- [Container Canary](#container-canary)
  - [Installation](#installation)
  - [Example (Kubeflow)](#example-kubeflow)
  - [Validator reference](#validator-reference)
    - [Metadata](#metadata)
    - [Runtime options](#runtime-options)
      - [Environment variables](#environment-variables)
      - [Ports](#ports)
      - [Volumes](#volumes)
      - [Command](#command)
    - [Checks](#checks)
      - [Exec](#exec)
      - [HTTPGet](#httpget)
      - [Delays, timeouts, periods and thresholds](#delays-timeouts-periods-and-thresholds)
  - [Contributing](#contributing)
  - [Maintaining](#maintaining)
  - [License](#license)

## Installation

You can find binaries and instructions on [our releases page](https://github.com/NVIDIA/container-canary/releases).

## Example (Kubeflow)

The [Kubeflow](https://www.kubeflow.org/) documentation has a [list of requirements](https://www.kubeflow.org/docs/components/notebooks/container-images/#custom-images) for container images that can be used in the [Kubeflow Notebooks](https://www.kubeflow.org/docs/components/notebooks/) service.

That list looks like this:

- expose an HTTP interface on port `8888`:
  - kubeflow sets an environment variable `NB_PREFIX` at runtime with the URL path we expect the container be listening under
  - kubeflow uses IFrames, so ensure your application sets `Access-Control-Allow-Origin: *` in HTTP response headers
- run as a user called `jovyan`:
  - the home directory of `jovyan` should be `/home/jovyan`
  - the UID of `jovyan` should be `1000`
- start successfully with an empty PVC mounted at `/home/jovyan`:
  - kubeflow mounts a PVC at `/home/jovyan` to keep state across Pod restarts

With Container Canary we could write this list as the following YAML spec.

```yaml
# examples/kubeflow.yaml
apiVersion: container-canary.nvidia.com/v1
kind: Validator
name: kubeflow
description: Kubeflow notebooks
env:
  - name: NB_PREFIX
    value: /hub/jovyan/
ports:
  - port: 8888
    protocol: TCP
volumes:
  - mountPath: /home/jovyan
checks:
  - name: user
    description: 👩 User is jovyan
    probe:
      exec:
        command:
          - /bin/sh
          - -c
          - "[ $(whoami) = jovyan ]"
  - name: uid
    description: 🆔 User ID is 1000
    probe:
      exec:
        command:
          - /bin/sh
          - -c
          - "id | grep uid=1000"
  - name: home
    description: 🏠 Home directory is /home/jovyan
    probe:
      exec:
        command:
          - /bin/sh
          - -c
          - "[ $HOME = /home/jovyan ]"
  - name: http
    description: 🌏 Exposes an HTTP interface on port 8888
    probe:
      httpGet:
        path: /
        port: 8888
      initialDelaySeconds: 10
  - name: NB_PREFIX
    description: 🧭 Correctly routes the NB_PREFIX
    probe:
      httpGet:
        path: /hub/jovyan/lab
        port: 8888
      initialDelaySeconds: 10
  - name: allow-origin-all
    description: "🔓 Sets 'Access-Control-Allow-Origin: *' header"
    probe:
      httpGet:
        path: /
        port: 8888
        responseHttpHeaders:
          - name: Access-Control-Allow-Origin
            value: "*"
      initialDelaySeconds: 10
```

The Canary Validator spec reuses parts of the [Kubernetes](https://kubernetes.io/) configuration API including [probes](https://kubernetes.io/docs/tasks/configure-pod-container/configure-liveness-readiness-startup-probes/). In Kubernetes probes are used to check on the health of a pod, but in Container Canary we use them to validate if the container meets our specification.

We can then run our specification against any desired container image to see a pass/fail breakdown of requirements. We can test one of the default images that ships with Kubeflow as that should pass.

```console
$ canary validate --file examples/kubeflow.yaml public.ecr.aws/j1r0q0g6/notebooks/notebook-servers/jupyter-scipy:v1.5.0-rc.1
Validating public.ecr.aws/j1r0q0g6/notebooks/notebook-servers/jupyter-scipy:v1.5.0-rc.1 against kubeflow
 👩 User is jovyan                                   [true]
 🆔 User ID is 1000                                  [true]
 🏠 Home directory is /home/jovyan                   [true]
 🌏 Exposes an HTTP interface on port 8888           [true]
 🧭 Correctly routes the NB_PREFIX                   [true]
 🔓 Sets 'Access-Control-Allow-Origin: *' header     [true]
validation passed
```

For more examples [see the examples directory](examples/).

## Validator reference

Validator manifests are YAML files that describe how to validate a container image. Check out the [examples](examples/) directory for real world applications.

### Metadata

Each manifests starts with some metadata.

```yaml
# Manifest versioning
apiVersion: container-canary.nvidia.com/v1
kind: Validator

# Metadata
name: foo  # The name of the platform that this manifest validates for
description: Foo runs containers for you  # A description of that platform
documentation: https://example.com  # A link to the documentation that defines the container requirements in prose
```

### Runtime options

Next you can set runtime configuration for the container you are validating. You should set these to mimic the environment that the compute platform will create. When you validate a container it will be run locally using [Docker](https://www.docker.com/).

#### Environment variables

A list of environment variables that should be set on the container.

```yaml
env:
  - name: HELLO
    value: world
  - name: FOO
    value: bar
```

#### Ports

Ports that need to be exposed on the container. These need to be configured in order for Container Canary to perform connectivity tests.

```yaml
ports:
  - port: 8888
    protocol: TCP
```

#### Volumes

Volumes to be mounted to the container. This is useful if the compute platform will always mount an empty volume to a specific location.

```yaml
volumes:
  - mountPath: /home/jovyan
```

#### Command

You can specify a custom command to be run inside the container.

```yaml
command:
 - foo
 - --bar=true
```

### Checks

Checks are the tests that we want to run against the container to ensure it is compliant. Each check contains a probe, and those probes are superset of the Kubernetes [probes](https://kubernetes.io/docs/tasks/configure-pod-container/configure-liveness-readiness-startup-probes/) API and so any valid Kubernetes probe can be used in a check.

```yaml
checks:
  - name: mycheck  # Name of the check
    description: Ensuring a thing  # Descrption of what is being checked (will be used in output)
    probe:
      ...  # A probe to run
```

#### Exec

An exec check runs a command inside the running container. If the command exits with `0` the check will pass.

```yaml
checks:
  - name: uid
    description: User ID is 1234
    probe:
      exec:
        command:
          - /bin/sh
          - -c
          - "id | grep uid=1234"
```

#### HTTPGet

An HTTP Get check will perform an HTTP GET request against your container. If the response code is `<300` and the optional response headers match the check will pass.

```yaml
checks:
  - name: http
    description: Exposes an HTTP interface on port 80
    probe:
      httpGet:
        path: /
        port: 80
        httpHeaders:  # Optional, headers to set in the request
          - name: Foo-Header
            value: "myheader"
        responseHttpHeaders:  # Optional, headers that you expect to see in the response
          - name: Access-Control-Allow-Origin
            value: "*"
```

#### Delays, timeouts, periods and thresholds

Checks also support the same delays, timeouts, periods and thresholds that Kubernetes probes do.

```yaml
checks:
  - name: uid
    description: User ID is 1234
    probe:
      exec:
        command: [...]
      initialDelaySeconds: 0  # Delay after starting the container before the check should be run
      timeoutSeconds: 30  # Overall timeout for the check
      successThreshold: 1  # Number of times the check must pass before moving on
      failureThreshold: 1  # Number of times the check is allowed to fail before giving up
      periodSeconds: 1  # Interval between runs if threasholds are >1
```

## Contributing

Contributions are very welcome, be sure to review the [contribution guidelines](./CONTRIBUTING.md).

## Maintaining

Maintenance steps [can be found here](./MAINTAINING.md).

## License

Apache License Version 2.0, see [LICENSE](./LICENSE).
