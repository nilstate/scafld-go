# Release

Release gates include formatting, vet, unit tests, race tests, golden tests,
binary e2e tests, and architecture tests.

Import-boundary tests are release-blocking.

The primary repository and Go module path is `github.com/nilstate/scafld`.
Release tags are `vX.Y.Z`.

The release workflow publishes:

- GitHub release binaries for Linux, macOS, and Windows on amd64/arm64
- `checksums.txt`
- `manifest.json`
- npm package `scafld`
- PyPI package `scafld`

Run a local release artifact snapshot with:

```bash
make release-snapshot
```

See [distribution.md](distribution.md) for package-manager architecture.
