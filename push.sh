#!/bin/bash

HASH="$(git rev-parse HEAD)"
REGISTRY=$(cat k8s/config/REGISTRY)
GROUP_NAME=$(cat k8s/config/GROUP_NAME)
APP_NAME=$(cat k8s/config/APP_NAME)
LOCAL_URL=${GROUP_NAME}/${APP_NAME}:${HASH}
source ./k8s/config/valid.sh
tagAndPush() {
  FULL_URL=${REGISTRY}${GROUP_NAME}/${APP_NAME}:${HASH}
  docker tag $LOCAL_URL $FULL_URL
  docker push ${FULL_URL}
}

tagAndPush prod

# Insist that the version file is tracked by git
git add --force VERSION

# Commit and push the version change. This will preserve the state.
git commit -am "Bumped version of docker image to $HASH"
#git push

