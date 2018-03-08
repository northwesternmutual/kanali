#!/bin/bash

set -e
source hack/vars.sh

docker login -u $DOCKER_USER -p $DOCKER_PASS $IMAGE_REGISTRY
docker push $IMAGE_NAME:$IMAGE_TAG
