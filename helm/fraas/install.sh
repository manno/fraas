#!/bin/bash

set -ev

: ${EMAIL:?}
: ${PROJECT_ID:?}
: ${CLUSTER_NAME:?}
: ${DNS_DOMAIN:?}
: ${DB_PASSWORD:?}
: ${APP_SECRET:?}

echo "Installing nginx loadbalancer..."
helm install stable/nginx-ingress --name nginx-controller \
  #--set controller.service.enableHttp=false \
  --set controller.publishService.enabled=true

echo "Installing external DNS updater..."
helm install --name dns --namespace kube-system stable/external-dns \
  --set provider=google \
  --set google.project="$PROJECT_ID" \
  --set google.serviceAccountSecret="external-dns-credentials" \
  --set txtOwnerId="$CLUSTER_NAME" \
  --set rbac.create=true \
  --set policy=sync

echo "Installing cert manager..."
helm upgrade --install cert-manager stable/cert-manager \
  --namespace kube-system \
  --set ingressShim.extraArgs='{--default-issuer-name=letsencrypt-prod,--default-issuer-kind=ClusterIssuer}' \
  --set ingressShim.enabled=false

echo "Installing FRAAS..."
helm upgrade --install fraas \
  --set ingress.dnsName="admin.$DNS_DOMAIN" \
  --set ingress.email="$EMAIL" \
  --set ingress.issuer="letsencrypt-prod" \
  --set image.source="gcr.io/$PROJECT_ID/fraas:latest" \
  --set project.id="$PROJECT_ID" \
  --set app.database_url="postgres://fraas:${DB_PASSWORD}@127.0.0.1:5432/fraas?sslmode=disable" \
  --set app.secret="$APP_SECRET" \
  $EXTRA_ARGS \
  .

# TODO post install db seed: kubectl exec -it faas-web-6dd99d6d7f-f65ds  -c faas-app  -- /bin/app task db:seed
