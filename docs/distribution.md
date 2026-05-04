# Distribution

scafld is a Go CLI distributed through several package ecosystems. The Go
binary is authoritative; npm and PyPI packages are thin launchers that download
the matching GitHub release asset.

## Primary channels

- Go modules: `go install github.com/nilstate/scafld/cmd/scafld@latest`
- GitHub Releases: raw platform binaries plus `checksums.txt` and `manifest.json`
- npm: `npm install -g scafld`
- PyPI: `pipx install scafld`

## Secondary channels

These should be generated from the GitHub release assets, not rebuilt from
source in separate registry flows:

- Homebrew tap: `nilstate/homebrew-tap`
- Scoop bucket: `nilstate/scoop-bucket`
- WinGet manifest: submitted after each stable release
- Docker/OCI image: useful for CI runners, published from the same tag
- Debian/RPM/AUR/Nix: package from the release binary and checksum manifest

Templates for Homebrew, Scoop, WinGet, and OCI live under `package/`.
They are intentionally templates, because those registries either require
separate repositories or human/registry review. They must consume the GitHub
release artifacts and `checksums.txt`.

## Release contract

1. A tag `vX.Y.Z` is pushed in `github.com/nilstate/scafld`.
2. CI runs `make check`.
3. `scripts/build-release-artifacts.sh X.Y.Z` builds raw native binaries for
   Linux, macOS, and Windows on amd64/arm64.
4. GitHub release assets are published before npm/PyPI, because those wrappers
   download from the release by version.
5. npm and PyPI publish through OIDC trusted publishing from
   `.github/workflows/release.yml`.

Package wrappers must never reimplement scafld behavior. They may only locate,
download, cache, and execute the native binary.
