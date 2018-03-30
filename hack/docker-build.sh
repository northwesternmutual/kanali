#!/bin/bash

set -e
source hack/vars.sh

docker build --build-arg VERSION=$IMAGE_TAG -t $IMAGE_NAME:$IMAGE_TAG .