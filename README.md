# Welcome to FRAAS

FRAAS is intended to be a sign-up and deployment manager for frab installations in the cloud.

FRAAS is built on buffalo and currently uses GKE.

## Requirements

* Download the buffalo cli.
* Working PostgreSQL

## Database Setup

Look at the `database.yml` file and set the `DATABASE_URL` to point to your PostgreSQL for a production setup.

### Create Your Databases

Ok, so you've edited the "database.yml" file and started postgres, now Buffalo can create the databases in that file for you:

	$ buffalo db create -a

## Starting the Application

Buffalo ships with a command that will watch your application and automatically rebuild the Go binary and any assets for you. To do that run the "buffalo dev" command:

	$ buffalo dev

Point your browser to [http://127.0.0.1:3000](http://127.0.0.1:3000).

## Configuration

### Site-Wide Configuration

The environment variable `FRAAS_CONFIG` contains the site-wide configuration as YAML.

```
---
google:
  project_id: fraas-1234
  zone: europe-west3-a
  region: europe-west3
  cluster_id: fraas-9
  dnszone: frabapp

domain: frab.app.
docker_image: gcr.io/fraas-1234/frab:current

mail:
  smtp_server: 'localhost'
  smtp_server_port: '25'
  smtp_notls: 'true'
  smtp_user_name: 'frab'
  smtp_password: 'barf'
  exception_email: 'admin@example.org'

database:
  instance: frab-pq
  admin_user: postgres
  admin_password: password123
```

When running `export FRAAS_CONFIG=$(cat fraas.yml)` watch out for newlines.

## Google Kubernetes Engine Configuration

### Cluster

To be able to modify DNS:

```
% gcloud container clusters create "fraas-1" \
    --num-nodes 1 \
    --scopes "https://www.googleapis.com/auth/ndev.clouddns.readwrite"

or use nodeSelector on deployment

% gcloud container node-pools create fraas-pool --cluster fraas-1 --scopes https://www.googleapis.com/auth/ndev.clouddns.readwrite
```

### IAM

1. service account for deployments, dns.admin, compute.admin, used by FRAAS
1. service account with sql role, later uploaded to be used by frab deployments

```
% gcloud projects get-iam-policy fraas-1234
bindings:
- members:
  - serviceAccount:sql@fraas-1234.iam.gserviceaccount.com
  role: roles/cloudsql.client
- members:
  - serviceAccount:compute-address@fraas-1234.iam.gserviceaccount.com
  role: roles/compute.admin
- members:
  - serviceAccount:compute-address@fraas-1234.iam.gserviceaccount.com
  role: roles/compute.networkAdmin
- members:
  - serviceAccount:compute-address@fraas-1234.iam.gserviceaccount.com
  role: roles/dns.admin
```

#### Inside GKE

```
kubectl create serviceaccount fraas
kubectl create rolebinding fraas-admin \
  --clusterrole=cluster-admin \
  --serviceaccount=default:fraas \
  --namespace=default
```

### Docker

Upload docker image: `gcloud docker -- push eu.gcr.io/fraas-1234/frab`

### Database

Create GKE Cloud SQL instance:

```
gcloud sql databases create fraas \
  --instance frab-pq --project fraas-1234
```

Set admin user password on instance, so FRAAS can create databases:

```
gcloud sql users set-password postgres no-host \
  --instance frab-pq --project fraas-1234 \
  --password password123
```

Push credentials with permissions on Cloud SQL (`sql@fraas-1234.iam.gserviceaccount.com`) to GKE:

```
kubectl create secret generic cloudsql-frab-credentials \
  --from-file=credentials.json=FRAAS-123456789012.json
```

### DNS

The configured zone is managed by Google DNS servers: `gcloud dns managed-zones`.

<<<<<<< HEAD
```
gcloud dns --project=fraas-1234 managed-zones create frabapp --description='frab zone' --dns-name=frab.app.
```

### TLS - LetsEncrypt

Follow https://github.com/ahmetb/gke-letsencrypt/blob/master/30-setup-letsencrypt.md

Upload issuer with your email:

```
curl -sSL https://rawgit.com/ahmetb/gke-letsencrypt/master/yaml/letsencrypt-issuer.yaml | \
    sed -e "s/email: ''/email: $EMAIL/g" | \
    kubectl apply -f-
```
