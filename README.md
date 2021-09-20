# git-volume-reloader

Synchronise a directory's contents with a git repository. Synchronisation is
triggered by a webhook sent by the git service provider.

## Usage

To package the git-volume-reloader into a container image and push the image to
a container registry:

```bash
make build push IMAGE=ghcr.io/padok-team/git-volume-reloader:latest
```

An example of how to deploy the git-volume-reloader as a sidecar in a Kubernetes
Pod is available in the [examples/git-mkdocs](./examples/git-mkdocs/README.md) directory.

## LICENSE

© 2021 [Padok](https://www.padok.fr/).

Licensed under the [Apache License](https://www.apache.org/licenses/LICENSE-2.0), Version 2.0 ([LICENSE](./LICENSE))
