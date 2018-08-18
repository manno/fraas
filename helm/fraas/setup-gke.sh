#!/bin/bash

set -ev

: ${PROJECT_ID:?}
: ${CLUSTER_NAME:?}
: ${ZONE:?}
: ${DNS_DOMAIN:?}
: ${ACCOUNT_COMPUTE:?}
: ${ACCOUNT_SQL:?}
: ${DB_PASSWORD:?}
: ${DB_ADMIN_PASSWORD:?}

FRAAS_DIR=$(dirname "$BASH_SOURCE")/../..

echo "Creating cluster..."
gcloud container clusters create "$CLUSTER_NAME" \
  --project "$PROJECT_ID" \
  --num-nodes 1 \
  --scopes "https://www.googleapis.com/auth/ndev.clouddns.readwrite" \
  --disable-addons=HttpLoadBalancing

gcloud container clusters get-credentials "$CLUSTER_NAME" \
  --project "$PROJECT_ID" \
  --zone "$ZONE"

kubectl patch limitranges limits -p '{"spec":{"limits":[{"defaultRequest":{"cpu":"10m"},"type":"Container"}]}}'

echo "Creating service accounts..."
kubectl create serviceaccount fraas
kubectl create rolebinding fraas-admin \
  --clusterrole=cluster-admin \
  --serviceaccount=default:fraas \
  --namespace=default


echo "Creating database instance..."
gcloud sql instances create fraas-db \
  --gce-zone $ZONE \
  --database-version=POSTGRES_9_6 \
  --tier db-f1-micro
gcloud sql users set-password postgres no-host \
  --instance fraas-db --project "$PROJECT_ID" \
  --password "$DB_ADMIN_PASSWORD"

echo "Create FRAAS database..."
gcloud sql databases create fraas \
  --instance fraas-db --project "$PROJECT_ID"
gcloud sql users create fraas no-host \
  --instance fraas-db --project "$PROJECT_ID" \
  --password "$DB_PASSWORD"

echo "Create DB service account credentials..."
gcloud iam service-accounts keys create key-sql.json \
  --iam-account "${ACCOUNT_SQL}@${PROJECT_ID}.iam.gserviceaccount.com"
kubectl create secret generic cloudsql-fraas-credentials \
  --from-file=credentials.json="key-sql.json"

echo "Create DNS service account credentials..."
gcloud iam service-accounts keys create key-compute.json \
  --iam-account "${ACCOUNT_COMPUTE}@${PROJECT_ID}.iam.gserviceaccount.com"
kubectl create secret generic external-dns-credentials \
  --namespace kube-system \
  --from-file=credentials.json=key-compute.json

echo "Prepare DNS zone for frab deployments..."
gcloud dns --project="$PROJECT_ID" managed-zones create frabapp --description='frab zone' --dns-name="${DNS_DOMAIN}."

echo "Building and uploading FRAAS application..."
pushd $FRAAS_DIR
  buffalo build
  docker build -t "gcr.io/$PROJECT_ID/fraas:latest" .
  docker push "gcr.io/$PROJECT_ID/fraas:latest"
popd

echo "Installing helm..."
kubectl create serviceaccount -n kube-system tiller
kubectl create clusterrolebinding tiller-binding --clusterrole=cluster-admin --serviceaccount kube-system:tiller
helm init --wait --service-account tiller
helm repo update
