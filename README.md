# Container Canary

A little bird to validate your container images.

```console
$ canary validate --file somespec.yaml foo/bar:latest
Validating foo/bar:latest against somespec
 📦 Required packages are installed                  [true]
 🤖 Expected services are running                    [true]
 🎉 Your container is awesome                        [true]
PASSED
```

Many modern compute platforms support bring-your-own-container models where the user can provide container images with their custom software environment. However platforms commonly have a set of requirements that the container must conform to, such as using a non-root user, having the home directory in a specific location, having certain packages installed or running web applications on specific ports.

Container Canary is a tool for recording those requirements in a YAML manifest and then validating containers against that manifest. This is particularly useful in CI environments to avoid regressions in containers.

## Installation

You can find binaries and instructions on [our releases page](https://github.com/NVIDIA/container-canary/releases).

## Example (Kubeflow)

The [Kubeflow](https://www.kubeflow.org/) documentation has a [list of requirements](https://www.kubeflow.org/docs/components/notebooks/container-images/#custom-images) for container images that can be used in the Kubeflow Notebooks service.

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
PASSED
```

For more examples [see the examples directory](examples/).
