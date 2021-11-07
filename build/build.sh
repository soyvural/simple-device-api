#!/usr/bin/env bash

# exit immediately when an error occurs.
set -e

BUILD_VERSION=$1

GIT_ROOT=$(git rev-parse --show-toplevel)

echo "contanier build with version $BUILD_VERSION"

docker build -f $GIT_ROOT/build/Dockerfile -t soyvural/simple-device-api:$BUILD_VERSION .
echo $DOCKER_PASSWORD | docker login registry-1.docker.io --username $DOCKER_USERNAME --password-stdin
docker push docker.io/soyvural/simple-device-api:$BUILD_VERSION