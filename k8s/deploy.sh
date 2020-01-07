#!/usr/bin/env bash

APP_NAME=$(cat config/APP_NAME)
RUN_AS="kubectl apply -f build/${CLUSTER}-${APP_NAME}.yaml"
echo Running: $RUN_AS
su ${CLUSTER} bash -c "${RUN_AS}"
