# REQUIREMENTS

1. Create GKE project
1. Create service accounts

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

1. Install gcloud, kubectl and helm command line tools

# DEPLOYMENT

1. Export variables needed by `./setup-gke.sh`
1. Create `./fraas.yml` from `./fraas.yml.template`
1. Export configuration variables, i.e. via `.envrc`
```
export EMAIL=test@example.org
export PROJECT_ID=fraas-1234
export CLUSTER_NAME=fraas-1
export ZONE=eu-west3-a
export DNS_DOMAIN=example.com
export ACCOUNT_COMPUTE=compute-address
export ACCOUNT_SQL=sql
export DB_PASSWORD=123
export DB_ADMIN_PASSWORD=123
export APP_SECRET=123
```
1. Run `./setup-gke.sh`
1. Run `./install.sh`
