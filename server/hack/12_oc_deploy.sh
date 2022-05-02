#!/usr/bin/env bash
set -e
echo 'Deploying scm server with spi operator and oauth service'

kustomize build server/config/os | oc apply -f -

oc rollout status deployment/spi-system-file-retriever-server  -n spi-system
oc rollout status deployment/spi-controller-manager  -n spi-system
oc rollout status deployment/spi-oauth-service  -n spi-system
