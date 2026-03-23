#!/usr/bin/env bash
set -euo pipefail

APP_NAME="xresources"
VERSION="${GITHUB_REF_NAME#v}"
TOPDIR="$PWD/.rpmbuild"
OUTDIR="$PWD/dist/rpm-source"

mkdir -p "$TOPDIR"/{BUILD,RPMS,SOURCES,SPECS,SRPMS} "$OUTDIR"

git archive --format=tar.gz --prefix="${APP_NAME}-${VERSION}/" -o "$TOPDIR/SOURCES/${APP_NAME}-${VERSION}.tar.gz" HEAD
cp packaging/rpm/app.spec "$TOPDIR/SPECS/"

rpmbuild \
  --define "_topdir $TOPDIR" \
  --define "version $VERSION" \
  -bs "$TOPDIR/SPECS/app.spec"

cp "$TOPDIR/SRPMS"/*.src.rpm "$OUTDIR/"
