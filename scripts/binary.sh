#!/usr/bin/env bash

set -e

if [ -z "$1" ]; then
    VERSION=dev-$(date +%FT%T%z)
else
    VERSION=$1
fi

GITCOMMIT=$(git rev-parse --short HEAD)
BUILD_FLAGS=(-ldflags="-X main.version=${VERSION} -X main.commit=${GITCOMMIT}")

rm -f ${GOPATH%%:*}/bin/kanali*

go install \
    "${BUILD_FLAGS[@]}" \
    github.com/northwesternmutual/kanali/cmd/kanali

if [ $? -eq 0 ]; then
  echo "Build successful. Binary located at ${GOPATH%%:*}/bin"
else
  echo "Build failed."
fi
