#!/bin/bash

PLATFORMS="darwin/386 darwin/amd64 freebsd/386 freebsd/amd64 freebsd/arm linux/386 linux/amd64 linux/arm windows/386 windows/amd64 openbsd/386 openbsd/amd64"

function go-alias {
	local GOOS=${1%/*} 
  local GOARCH=${1#*/}
  local OUT_DIR="release/$1"
  echo "Building: $OUT_DIR"
  mkdir -p $OUT_DIR
  GOOS=$GOOS GOARCH=$GOARCH go build -o "$OUT_DIR/rip"
}

echo "Creating release/ directory"
mkdir release

for PLATFORM in $PLATFORMS; do
	go-alias $PLATFORM
done
unset -f go-alias

echo "Copying README"
cp README.md release/

echo "Creating Archives"
tar -czf release.tar.gz release
zip -r release.zip release

