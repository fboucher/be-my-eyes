# Releasing be-my-eyes

This document describes how to create a new release of be-my-eyes.

## Prerequisites

- Push access to the repository
- The `HOMEBREW_TAP_TOKEN` secret configured (for Homebrew tap updates)

## Release Process

### 1. Prepare the Release

1. Ensure all changes are merged to `main`
2. Verify that the build passes on `main`
3. Update any necessary documentation

### 2. Create and Push a Tag

Create a semantic version tag following the format `vX.Y.Z`:

```bash
# Example: creating version 1.0.0
git checkout main
git pull origin main
git tag v1.0.0
git push origin v1.0.0
```

### 3. Monitor the Release Workflow

1. Go to the [Actions tab](https://github.com/fboucher/be-my-eyes/actions)
2. Watch the "Release" workflow run
3. The workflow will:
   - Build cross-platform binaries (Linux and macOS, both amd64 and arm64)
   - Create Debian packages (.deb for amd64 and arm64)
   - Generate checksums
   - Create a GitHub Release with all artifacts
   - Update the Homebrew tap (fboucher/homebrew-tap) if the token is available

### 4. Verify the Release

Once the workflow completes successfully:

1. Check the [Releases page](https://github.com/fboucher/be-my-eyes/releases)
2. Verify all expected artifacts are present:
   - `be-my-eyes_X.Y.Z_linux_amd64.tar.gz`
   - `be-my-eyes_X.Y.Z_linux_arm64.tar.gz`
   - `be-my-eyes_X.Y.Z_darwin_amd64.tar.gz`
   - `be-my-eyes_X.Y.Z_darwin_arm64.tar.gz`
   - `be-my-eyes_X.Y.Z_amd64.deb`
   - `be-my-eyes_X.Y.Z_arm64.deb`
   - `checksums.txt`
3. Test the release by downloading and running a binary for your platform
4. If Homebrew tap was updated, test the formula:
   ```bash
   brew tap fboucher/homebrew-tap
   brew install be-my-eyes
   be-my-eyes --version
   ```

## Manual Release (Emergency)

If you need to manually trigger a release (e.g., to fix a failed release):

1. Go to the [Release workflow](https://github.com/fboucher/be-my-eyes/actions/workflows/release.yml)
2. Click "Run workflow"
3. Enter the tag name (e.g., `v1.0.0`)
4. Click "Run workflow"

## Troubleshooting

### Release workflow fails

1. Check the workflow logs in the Actions tab
2. Common issues:
   - **Missing ARM64 cross-compiler**: The workflow should install it automatically. If it fails, check the apt-get logs.
   - **GoReleaser errors**: Check the `.goreleaser.yml` configuration
   - **Homebrew token issues**: Verify the `HOMEBREW_TAP_TOKEN` secret is set correctly

### Artifacts missing from release

If the release was created but some artifacts are missing:

1. Re-run the failed workflow jobs
2. Or delete the tag and release, fix any issues, and create a new tag

### Need to update an existing release

If you need to add or update artifacts on an existing release:

1. Use the manual workflow dispatch to re-run for that tag
2. Or manually build and upload using `gh release upload`

## Version Numbering

Follow [Semantic Versioning](https://semver.org/):

- **MAJOR** (X.0.0): Incompatible API changes
- **MINOR** (0.X.0): New functionality in a backward-compatible manner
- **PATCH** (0.0.X): Backward-compatible bug fixes

Always prefix versions with `v` (e.g., `v1.0.0`).
