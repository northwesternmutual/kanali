#!/bin/bash

set -e

if [[ "$TRAVIS_SECURE_ENV_VARS" == "false" ]]; then
  echo "skip docker upload, TRAVIS_SECURE_ENV_VARS=$TRAVIS_SECURE_ENV_VARS"
  exit 0
fi

export DOCKER_REPO="northwesternmutual/kanali"

if [[ $TRAVIS_TAG =~ ^v[0-9]+\.[0-9]+\.[0-9]+$ ]]; then
  echo "uploading docker image TAG=$TRAVIS_TAG"
  export DOCKER_TAG=$TRAVIS_TAG
  docker build --build-arg VERSION=$DOCKER_TAG -t $DOCKER_REPO:$DOCKER_TAG .
  bash ./travis/upload-to-docker.sh
elif [ $TRAVIS_PULL_REQUEST == "false" ] && [ $TRAVIS_BRANCH == "master" ]; then
  echo "uploading docker image BRANCH=$TRAVIS_BRANCH"
  for component in latest ${TRAVIS_COMMIT::8}
  do
    export DOCKER_TAG=${component}
    docker build --build-arg VERSION=$DOCKER_TAG -t $DOCKER_REPO:$DOCKER_TAG .
    bash ./travis/upload-to-docker.sh
  done
else
  echo "skip docker upload - neither master branch nor tag"
  exit 0
fi

echo "docker upload successfull"
exit 0