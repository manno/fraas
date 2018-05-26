package models

import (
	"time"
)

// Keep in sync with buffalo app, couldn't use pop model because of missing db config
type Deployment struct {
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`

	Name             string `json:"name" db:"name"`
	ContactEmail     string `json:"contact_email" db:"contact_email"`
	DatabaseID       string `json:"database_id" db:"database_id"`
	DatabasePassword string `json:"database_password" db:"database_password"`
	Domain           string `json:"domain" db:"domain"`
	ExternalDomain   string `json:"external_domain" db:"external_domain"`
	FromEmail        string `json:"from_email" db:"from_email"`
	SecretKeyBase    string `json:"secret_key_base" db:"secret_key_base"`
}

const prefix = "frab-"

func (d *Deployment) DeploymentID() string {
	return prefix + d.Name
}

func (d *Deployment) ConfigName() string {
	return prefix + d.Name + "-config"
}

func (d *Deployment) SecretName() string {
	return prefix + d.Name + "-secret"
}

func (d *Deployment) WebName() string {
	return prefix + d.Name + "-web"
}

func (d *Deployment) RailsContainerName() string {
	return prefix + d.Name + "-app"
}

func (d *Deployment) IPName() string {
	return prefix + d.Name + "-web-ip"
}

func (d *Deployment) BackendName() string {
	return prefix + d.Name + "-backend"
}

func (d *Deployment) TLSSecretName() string {
	return prefix + d.Name + "-app-tls"
}
