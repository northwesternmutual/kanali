#!/bin/bash

ORG_NAME=northwesternmutual
PROJECT_NAME=kanali

PATH_IMPORT=github.com/northwesternmutual/kanali
PATH_PROJECT=$GOPATH/src/$PATH_IMPORT

GIT_COMMIT=$(git rev-parse --short HEAD)

IMAGE_NAME=$ORG_NAME/$PROJECT_NAME
IMAGE_TAG=${IMAGE_TAG:-local}
