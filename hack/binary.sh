#!/usr/bin/env bash

set -e

if [ -z "$1" ]; then
    VERSION=dev-$(date +%FT%T%z)
else
    VERSION=$1
fi

PROJECTPATH=$GOPATH/src/github.com/northwesternmutual/kanali
BINARIES=($(for i in $(ls -d $PROJECTPATH/cmd/*/); do echo ${i%%/} | awk -F "/" '{print $NF}'; done))
GITCOMMIT=$(git rev-parse --short HEAD)
BUILD_FLAGS=(-ldflags="-X main.version=${VERSION} -X main.commit=${GITCOMMIT}")

for i in "${BINARIES[@]}"
do
  rm -f ${GOPATH%%:*}/bin/$i

  go install \
      "${BUILD_FLAGS[@]}" \
      github.com/northwesternmutual/kanali/cmd/$i

  if [ $? -eq 0 ]; then
    echo "Build successful. Binary located at ${GOPATH%%:*}/bin/$i"
  else
    echo "Build failed."
  fi
done
