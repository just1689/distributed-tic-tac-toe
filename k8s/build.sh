#!/bin/bash

# The docker image to use
HASH=$(cat ../VERSION)

## REGISTRY must include trailing "/" unless empty
REGISTRY=$(cat config/REGISTRY)
GROUP_NAME=$(cat config/GROUP_NAME)
APP_NAME=$(cat config/APP_NAME)
FULL_URL=${REGISTRY}${GROUP_NAME}/${APP_NAME}:${HASH}

# Check that required config is not empty
source ./config/valid.sh
errorOnEmpty "GROUP_NAME" $GROUP_NAME
errorOnEmpty "APP_NAME" $APP_NAME
errorOnNewLine "GROUP_NAME" $GROUP_NAME
errorOnNewLine "APP_NAME" $APP_NAME
errorOnNewLine "REGISTRY" $REGISTRY

echo Building K8s yaml for image ${GROUP_NAME}/${APP_NAME} against tag ${HASH}

# Refresh the build area
rm -rf build/
mkdir -p build/tmp

cd overlays
# For each D in current directory
for D in *; do
  # If D is a directory...
  if [ -d "${D}" ]; then
    printf "resources:\n - out.yaml" >../build/tmp/kustomization.yaml

    # Kustomize using the overlay
    kustomize build ${D} >../build/tmp/out.yaml

    # Edit the image per overlay from the tmp directory
    FULL_URL=${REGISTRY}${GROUP_NAME}/${APP_NAME}:${HASH}
    (cd ../build/tmp && kustomize edit set image $GROUP_NAME\$APP_NAME=$FULL_URL)

    # Build the resulting artefact to the build directory
    kustomize build ../build/tmp >../build/${D}-${APP_NAME}.yaml

    echo Now building ${D}-${APP_NAME} for image ${FULL_URL}

  fi
done

# Clean up tmp area
cd ../
rm -rf build/tmp
