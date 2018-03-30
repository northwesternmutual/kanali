#!/bin/bash

set -e
source hack/vars.sh

E2E_PREFIX=e2e
E2E_SOURCE=./test/e2e
E2E_CLOUD_PROVIDER=${E2E_CLOUD_PROVIDER:-minikube}
E2E_KUBE_CONFIG=${E2E_KUBE_CONFIG:-$HOME/.kube/config}
E2E_KUBE_VERSION=${E2E_KUBE_VERSION:-v1.9.0}
E2E_KANALI_HOST=${E2E_KANALI_HOST:-$(minikube ip)}
E2E_KANALI_IMAGE_NAME=${E2E_KANALI_IMAGE_NAME:-$IMAGE_NAME}
E2E_KANALI_IMAGE_TAG=${E2E_KANALI_IMAGE_TAG:-local}
E2E_PROJECT_COMMIT_SHA=${E2E_PROJECT_COMMIT_SHA:-local}

GINKGO_PREFIX=ginkgo

go test -v -race $E2E_SOURCE \
  -$E2E_PREFIX.kube.config $E2E_KUBE_CONFIG \
  -$E2E_PREFIX.cloud.provider $E2E_CLOUD_PROVIDER \
  -$E2E_PREFIX.kanali.host $E2E_KANALI_HOST \
  -$E2E_PREFIX.kanali.image.name $E2E_KANALI_IMAGE_NAME \
  -$E2E_PREFIX.kanali.image.tag $E2E_KANALI_IMAGE_TAG \
  -$E2E_PREFIX.project.commit_sha $E2E_PROJECT_COMMIT_SHA \
  -$GINKGO_PREFIX.slowSpecThreshold=60 \
  -$GINKGO_PREFIX.flakeAttempts=2 \
  -$GINKGO_PREFIX.progress=true \
  -$GINKGO_PREFIX.randomizeAllSpecs=true
