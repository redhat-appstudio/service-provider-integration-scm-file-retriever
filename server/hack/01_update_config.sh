#!/usr/bin/env bash
set -e
echo 'Updating host variables'
SCM_HOST_VALUE='scm.'$(minikube ip)'.nip.io'
OAUTH_URL='oauth.'$(minikube ip)'.nip.io'
echo "scm="$SCM_HOST_VALUE
echo "oauth="$OAUTH_URL

SCRIPT_DIR=$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )

yq -i e '.spec.rules[0].host = "'$SCM_HOST_VALUE'"' $SCRIPT_DIR'/../config/k8s/ingress.yaml'
jq 'map(select(.op == "replace").value |= "'$OAUTH_URL'")' $SCRIPT_DIR'/../config/k8s/ingress-patch.json' > tmp.$$.json && mv tmp.$$.json $SCRIPT_DIR'/../config/k8s/ingress-patch.json'
