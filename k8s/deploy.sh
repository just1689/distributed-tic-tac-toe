#!/usr/bin/env bash

#APP_NAME=$(cat config/APP_NAME)
#RUN_AS="kubectl apply -f build/${CLUSTER}-${APP_NAME}.yaml"
#echo Running: $RUN_AS
#sudo su - ${CLUSTER} bash -c "${RUN_AS}"

APP_NAME=$(cat config/APP_NAME)
kubectl apply -f build/${CLUSTER}-${APP_NAME}.yaml
