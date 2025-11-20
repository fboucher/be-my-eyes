## Packaging Overview

This directory contains metadata and helper scripts for creating distribution packages for Debian/Ubuntu (APT), Arch (AUR), and Homebrew.

### Debian / Ubuntu
Use `packaging/debian/build-deb.sh` to create a simple `.deb` from the current source.
Example:
```bash
VERSION=0.1.0 ARCH=amd64 bash packaging/debian/build-deb.sh
```
Artifacts land in `dist/`.

### Arch (AUR)
The `PKGBUILD.template` is parameterized. Replace `{{VERSION}}` and submit to AUR as `PKGBUILD`. Generate `.SRCINFO` via `makepkg --printsrcinfo > .SRCINFO`.

### Homebrew
`be-my-eyes.rb.template` is a formula template. Replace `{{VERSION}}` and SHA256 sums with release asset checksums, then commit into a tap repository (`fboucher/homebrew-reka-tap`). Users install via:
```bash
brew install fboucher/tap/be-my-eyes
```

### Version Injection
All build scripts pass `-ldflags "-X github.com/fboucher/be-my-eyes/internal/version.Version=<version>"` so `be-my-eyes --version` reports the packaged version.

### Next Steps
1. Create a tagged release (e.g., `v0.1.0`).
2. Generate binaries for each target OS/ARCH.
3. Produce `.deb` packages.
4. Compute SHA256 for each binary; update Homebrew formula.
5. Publish AUR PKGBUILD.
6. Host apt repository or use a service (Cloudsmith / PackageCloud).
