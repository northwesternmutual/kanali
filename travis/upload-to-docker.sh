#!/bin/bash

set -e

if [[ "$DOCKER_REPO" == "" ]]; then
  echo "skip docker upload, DOCKER_REPO=$DOCKER_REPO"
  exit 0
fi

if [[ "$DOCKER_TAG" == "" ]]; then
  echo "skip docker upload, DOCKER_TAG=$DOCKER_TAG"
  exit 0
fi

docker login -u $DOCKER_USER -p $DOCKER_PASS
docker push $DOCKER_REPO:$DOCKER_TAG
