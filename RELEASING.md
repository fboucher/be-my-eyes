# Releasing be-my-eyes

This document describes how to create a new release of be-my-eyes.

## Prerequisites

- Push access to the repository
- The `HOMEBREW_TAP_TOKEN` secret configured (for Homebrew tap updates)

## Release Process

### 1. Create and Push a Tag

Create a semantic version tag following the format `vX.Y.Z`:

```bash
# Example: creating version 1.0.0
git checkout main
git pull origin main
git tag v1.0.0
git push origin v1.0.0
```

### 2. Monitor the Release Workflow

The GitHub Actions workflow will automatically:

- Build cross-platform binaries (Linux and macOS, both amd64 and arm64)
- Create Debian packages (.deb for amd64 and arm64)
- Generate checksums
- Create a GitHub Release with all artifacts
- Update the Homebrew tap (if token is available)

Watch the [Actions tab](https://github.com/fboucher/be-my-eyes/actions) to monitor progress.

### 3. Verify the Release

Once complete, check the [Releases page](https://github.com/fboucher/be-my-eyes/releases) for:

- `be-my-eyes_X.Y.Z_linux_amd64.tar.gz`
- `be-my-eyes_X.Y.Z_linux_arm64.tar.gz`
- `be-my-eyes_X.Y.Z_darwin_amd64.tar.gz`
- `be-my-eyes_X.Y.Z_darwin_arm64.tar.gz`
- `be-my-eyes_X.Y.Z_amd64.deb`
- `be-my-eyes_X.Y.Z_arm64.deb`
- `checksums.txt`

Test the Homebrew formula:

```bash
brew tap fboucher/tap
brew install be-my-eyes
be-my-eyes --version
```

## Version Numbering

Follow [Semantic Versioning](https://semver.org/):

- **MAJOR** (X.0.0): Incompatible API changes
- **MINOR** (0.X.0): New functionality in a backward-compatible manner
- **PATCH** (0.0.X): Backward-compatible bug fixes

Always prefix versions with `v` (e.g., `v1.0.0`).
