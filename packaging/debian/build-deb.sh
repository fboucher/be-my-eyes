#!/usr/bin/env bash
set -euo pipefail

VERSION="${VERSION:-dev}"
ARCH="${ARCH:-amd64}"
BIN_NAME="be-my-eyes"
REPO_ROOT="$(git rev-parse --show-toplevel 2>/dev/null || pwd)"
OUTDIR="${OUTDIR:-$REPO_ROOT/dist}"

mkdir -p "$OUTDIR"
STAGE="$OUTDIR/deb-stage-$ARCH"
rm -rf "$STAGE"
mkdir -p "$STAGE/DEBIAN" "$STAGE/usr/bin"

echo "Building binary ($ARCH)..."
GOARCH="$ARCH" GOOS=linux go build -ldflags "-s -w -X github.com/fboucher/be-my-eyes/internal/version.Version=$VERSION" -o "$STAGE/usr/bin/$BIN_NAME" ./cmd/be-my-eyes

CONTROL_TEMPLATE="$REPO_ROOT/packaging/debian/control.template"
CONTROL_FILE="$STAGE/DEBIAN/control"
sed "s/{{VERSION}}/$VERSION/; s/{{ARCH}}/$ARCH/" "$CONTROL_TEMPLATE" > "$CONTROL_FILE"

echo "Generating md5sums..."
pushd "$STAGE/usr" >/dev/null
find . -type f -exec md5sum {} + > ../DEBIAN/md5sums
popd >/dev/null

DEB_NAME="${BIN_NAME}_${VERSION}_${ARCH}.deb"
dpkg-deb --build "$STAGE" "$OUTDIR/$DEB_NAME"
echo "Created $OUTDIR/$DEB_NAME"
