#!/usr/bin/env bash
set -e
echo 'Deploying scm server with spi operator and oauth service'

oc create secret generic oauth-config \
    --save-config --dry-run=client \
    --from-file=server/config/os/config.yaml \
    -n spi-system \
    -o yaml |
oc apply -f -



kustomize build server/config/os | oc apply -f -

oc rollout status deployment/spi-scm-file-retriever-server  -n spi-scm
oc rollout status deployment/spi-controller-manager  -n spi-system
oc rollout status deployment/spi-oauth-service  -n spi-system
