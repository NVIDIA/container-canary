# Maintaining

This file contains notes and checklists for Container Canary maintainers.

For development guidelines see [CONTRIBUTING.md](CONTRIBUTING.md).

## Releasing

Draft releases are created automatically via the [release drafter GitHub Action](https://github.com/release-drafter/release-drafter) and kept up to date with each new PR.

The action will generate a changelog based on PR labels. For configuration options and the release notes template see [release-drafter.yml](.github/release-drafter.yml).

To publish a new release:

- Head to [releases](https://github.com/NVIDIA/container-canary/releases) and find the most recent draft.
- Edit the draft.
- Check the title and tag match the desired version (release drafter will attempt to bump this automatically based on PR tags).
- Check the body text makes sense.
- Hit publish release!

Once the release and tag have been created [a GitHub Actions workflow](.github/workflows/upload-release-assets.yaml) will run to build assets for the tagged commit and attach them to the release.

_Under the hood this runs `make package` to generate binaries for all supported platforms, then generates sha256 sums for each file and uploads them._
