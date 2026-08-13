#!/bin/bash

CR_NAME="payments-api"
NAMESPACE="default"

echo "Starting simple GatewayAPI controller..."

while true
do
  echo "Reconciling GatewayAPI: $CR_NAME"

  TEAM_NAME=$(kubectl get gapi "$CR_NAME" \
    -n "$NAMESPACE" \
    -o jsonpath='{.spec.teamName}')

  API_NAME=$(kubectl get gapi "$CR_NAME" \
    -n "$NAMESPACE" \
    -o jsonpath='{.spec.apiName}')

  AUTH_TYPE=$(kubectl get gapi "$CR_NAME" \
    -n "$NAMESPACE" \
    -o jsonpath='{.spec.authentication.type}')

  kubectl create configmap "${CR_NAME}-config" \
    -n "$NAMESPACE" \
    --from-literal=teamName="$TEAM_NAME" \
    --from-literal=apiName="$API_NAME" \
    --from-literal=authentication="$AUTH_TYPE" \
    --dry-run=client \
    -o yaml | kubectl apply -f -

  sleep 5
done