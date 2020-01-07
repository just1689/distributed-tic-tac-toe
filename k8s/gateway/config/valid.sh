#!/usr/bin/bash

# Errors out when $2 is empty, providing $1 as part of message
errorOnEmpty() {
  if [ -z "$2" ]; then
    echo $1 cannot be empty
    exit 1
  fi
}

# Errors out when $2 contains a new line, providing $1 as part of message
errorOnNewLine() {
  if [ $(echo "$2" | wc -l) -gt 1 ]; then
    echo $1 cannot contain new line
    exit 1
  fi
}

# Provide the relavant registry by the provided cluster
getRegistryByCluster() {
  REGISTRY="registry.captainjustin.space/"
}
