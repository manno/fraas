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

## Running FRAAS on Google Kubernetes Engine

Use the `setup-gke.sh` and helm charts in `helm/fraas` as described by the (README)[helm/fraas/README.md].

### Notes

#### DNS

The configured DNS zone is managed by Google DNS servers: `gcloud dns managed-zones list`.
To be able to modify DNS the node-pool used by FRAAS needs access to the scope `"https://www.googleapis.com/auth/ndev.clouddns.readwrite"`.

#### IAM

FRAAS uses two pre-existing service accounts:

1. service account for deployments, dns.admin, compute.admin, used by FRAAS itself
1. service account with sql role, used by frab deployments to access the DB

FRAAS uses a separate admin connection to create databases for new frab deployments.
