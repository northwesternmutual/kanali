#!/bin/bash

PATH_IMPORT=github.com/northwesternmutual/kanali
PATH_PROJECT=$GOPATH/src/$PATH_IMPORT

GIT_COMMIT=$(git rev-parse --short HEAD)

E2E_PREFIX=e2e
E2E_SOURCE=./test/e2e
E2E_CLOUD_PROVIDER=${E2E_CLOUD_PROVIDER:-minikube}
E2E_KUBE_CONFIG=${E2E_KUBE_CONFIG:-$HOME/.kube/config}
E2E_KUBE_VERSION=${E2E_KUBE_VERSION:-v1.9.0}
E2E_KANALI_HOST=${E2E_KANALI_HOST:-$(minikube ip)}
E2E_KANALI_IMAGE_NAME=${E2E_KANALI_IMAGE_NAME:-kanali}
E2E_KANALI_IMAGE_TAG=${E2E_KANALI_IMAGE_TAG:-local}
E2E_PROJECT_COMMIT_SHA=${E2E_PROJECT_COMMIT_SHA:-local}

GINKGO_PREFIX=ginkgo