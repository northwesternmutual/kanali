#!/usr/bin/env bash

set -e

source hack/vars.sh

if [ -z "$1" ]; then
    VERSION=dev-$(date +%FT%T%z)
else
    VERSION=$1
fi

BINARIES=($(for i in $(ls -d $PATH_PROJECT/cmd/*/); do echo ${i%%/} | awk -F "/" '{print $NF}'; done))
BUILD_FLAGS=(-ldflags="-X $PATH_IMPORT/pkg/version.version=${VERSION} -X $PATH_IMPORT/pkg/version.commit=${GIT_COMMIT}")

for i in "${BINARIES[@]}"
do
  rm -f ${GOPATH%%:*}/bin/$i

  go install \
      "${BUILD_FLAGS[@]}" \
      $PATH_IMPORT/cmd/$i

  if [ $? -eq 0 ]; then
    echo "Build successful. Binary located at ${GOPATH%%:*}/bin/$i"
  else
    echo "Build failed."
  fi
done
