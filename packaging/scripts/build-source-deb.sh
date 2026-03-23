#!/usr/bin/env bash
set -euo pipefail

APP_NAME="xresources"
VERSION="${GITHUB_REF_NAME#v}"
WORKDIR="/tmp/${APP_NAME}-${VERSION}"
OUTDIR="$PWD/dist/deb-source"

rm -rf "$WORKDIR"
mkdir -p "$WORKDIR" "$OUTDIR"

git archive --format=tar.gz --prefix="${APP_NAME}-${VERSION}/" -o "$OUTDIR/${APP_NAME}_${VERSION}.orig.tar.gz" HEAD

tar -xzf "$OUTDIR/${APP_NAME}_${VERSION}.orig.tar.gz" -C /tmp
cp -r packaging/debian "/tmp/${APP_NAME}-${VERSION}/debian"

(
  cd "/tmp/${APP_NAME}-${VERSION}"
  dch --create -v "${VERSION}-1" --package "$APP_NAME" "Automated source release"
  dpkg-buildpackage -S -sa
)

mv /tmp/${APP_NAME}_${VERSION}-1* "$OUTDIR/" || true
