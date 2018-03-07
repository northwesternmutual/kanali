#!/bin/bash

set -o errexit
set -o nounset
set -o pipefail

source hack/vars.sh

SCRIPT_ROOT=$(dirname ${BASH_SOURCE})/..
CODEGEN_PKG=${CODEGEN_PKG:-$(cd ${SCRIPT_ROOT}; ls -d -1 ./vendor/k8s.io/code-generator 2>/dev/null || echo ${GOPATH}/src/k8s.io/code-generator)}

vendor/k8s.io/code-generator/generate-groups.sh all \
  $PATH_IMPORT/pkg/client $PATH_IMPORT/pkg/apis \
  kanali.io:v2 \
  --go-header-file ${SCRIPT_ROOT}/hack/custom-boilerplate.go.txt \
  --v=1