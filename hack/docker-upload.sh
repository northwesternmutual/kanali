#!/bin/bash

set -e
source hack/vars.sh

docker login -u $DOCKER_USER -p $DOCKER_PASS $IMAGE_REGISTRY

if [ -z "$IMAGE_LATEST" ]; then
  docker tag $IMAGE_NAME:$IMAGE_TAG $IMAGE_NAME:latest
  docker push $IMAGE_NAME:latest
fi

docker push $IMAGE_NAME:$IMAGE_TAG
