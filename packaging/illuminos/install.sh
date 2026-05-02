#!/bin/bash

set -e

VERSION="0.1.0"
PREFIX="/opt/local"
BUILD_DIR="/tmp/freecode-build"

mkdir -p "$BUILD_DIR"
cd "$BUILD_DIR"

git clone https://github.com/cloudbsdorg/freecode.git
cd freecode
git checkout v${VERSION}

go build -ldflags="-s -w" -o freecode ./cmd/freecode

mkdir -p ${PREFIX}/bin
mkdir -p ${PREFIX}/share/man/man1
mkdir -p ${PREFIX}/share/freecode

cp freecode ${PREFIX}/bin/
gzip -c freecode.1 > ${PREFIX}/share/man/man1/freecode.1.gz

rm -rf "$BUILD_DIR"

echo "Freecode ${VERSION} installed to ${PREFIX}"
