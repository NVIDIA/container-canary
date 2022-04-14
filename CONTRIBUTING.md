# Container Canary Contribution Rules

Contributions to Container Canary are very welcome!

## Developing

This project is written in Go and uses `make` for task running. Code should be formatted with `gofmt`.

### Building

You can build Canary into `./bin/canary` with:

```shell
$ make build
go build -o bin/canary .
```

### Testing

Tests require some example containers to use. Before you run the tests you must build them.

```console
$ make testprep
docker build -t container-canary/kubeflow:shouldpass -f internal/testdata/containers/kubeflow.Dockerfile .
[+] Building 1.7s (10/10) FINISHED
...
 => => naming to docker.io/container-canary/kubeflow:shouldpass
 ```

You can then invoke tests with:

```shell
$ make test
go test -v ./...
...
PASS
```

## Signing Your Work

* We require that all contributors "sign-off" on their commits. This certifies that the contribution is your original work, or you have rights to submit it under the same license, or a compatible license.

  * Any contribution which contains commits that are not Signed-Off will not be accepted.

* To sign off on a commit you simply use the `--signoff` (or `-s`) option when committing your changes:

  ```bash
  git commit -s -m "Add cool feature."
  ```

  This will append the following to your commit message:

  ```
  Signed-off-by: Your Name <your@email.com>
  ```

* Full text of the DCO:

  ```
    Developer Certificate of Origin
    Version 1.1

    Copyright (C) 2004, 2006 The Linux Foundation and its contributors.
    1 Letterman Drive
    Suite D4700
    San Francisco, CA, 94129

    Everyone is permitted to copy and distribute verbatim copies of this license document, but changing it is not allowed.
  ```

  ```
    Developer's Certificate of Origin 1.1

    By making a contribution to this project, I certify that:

    (a) The contribution was created in whole or in part by me and I have the right to submit it under the open source license indicated in the file; or

    (b) The contribution is based upon previous work that, to the best of my knowledge, is covered under an appropriate open source license and I have the right under that license to submit that work with modifications, whether created in whole or in part by me, under the same open source license (unless I am permitted to submit under a different license), as indicated in the file; or

    (c) The contribution was provided directly to me by some other person who certified (a), (b) or (c) and I have not modified it.

    (d) I understand and agree that this project and the contribution are public and that a record of the contribution (including all personal information I submit with it, including my sign-off) is maintained indefinitely and may be redistributed consistent with this project or the open source license(s) involved.
  ```
