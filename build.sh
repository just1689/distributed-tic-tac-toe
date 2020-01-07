#!/bin/bash

# Use the current git commit hash as the version
HASH="$(git rev-parse HEAD)"

GROUP_NAME=$(cat k8s/config/GROUP_NAME)
APP_NAME=$(cat k8s/config/APP_NAME)
LOCAL_URL=${GROUP_NAME}/${APP_NAME}:${HASH}

# Create a docker image with that hash
docker build -t ${LOCAL_URL} .

# Overwrite the version in the VERSION file
echo ${HASH} >VERSION
