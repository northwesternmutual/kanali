#!/bin/bash

set -e
source hack/vars.sh

DISTRIBUTIONS=(
  "darwin/386"
  "darwin/amd64"
  "linux/386"
  "linux/amd64"
  "linux/arm"
  "freebsd/386"
  "freebsd/amd64"
  "freebsd/arm"
  "openbsd/386"
  "openbsd/amd64"
  "netbsd/386"
  "netbsd/amd64"
  "netbsd/arm"
  "windows/amd64"
)
S3_BASE_PATH="s3://kanalictl/release"

echo $RELEASE >> latest.txt
aws s3 mv latest.txt $S3_BASE_PATH/latest.txt

BUILD_FLAGS=(-ldflags="-X $PATH_IMPORT/pkg/version.version=${RELEASE} -X $PATH_IMPORT/pkg/version.commit=${GIT_COMMIT}")

for distro in ${DISTRIBUTIONS[@]}
do
    remote_path="${S3_BASE_PATH}/${RELEASE}/${distro}/${PROJECT_NAME}ctl"
    raw_binary_name="${PROJECT_NAME}ctl_${distro}"
    actual_delimiter=\\/
    desired_delimiter=_
    gox_binary_name=$(echo "${raw_binary_name}" | sed s/$actual_delimiter/$desired_delimiter/)

    gox "${BUILD_FLAGS[@]}" -osarch=$distro ${PATH_IMPORT}/cmd/kanalictl

    if [[ $distro = *"windows"* ]]; then
      aws s3 mv ${gox_binary_name}.exe ${remote_path}.exe
    else
      aws s3 mv ${gox_binary_name} ${remote_path}
    fi
done
