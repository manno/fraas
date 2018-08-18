package models

import (
	"encoding/json"
	"time"

	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/uuid"
	"github.com/gobuffalo/validate"
	"github.com/gobuffalo/validate/validators"

	fh "manno.name/mm/fraas/fraas-helpers"
)

type Deployment struct {
	ID        uuid.UUID `json:"id" db:"id"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`

	//DeploymentRequest / SignUp  / Request -> Deployment
	//ContactEmail     string `json:"contact_email" db:"contact_email"`
	//Reason     string `json:"reason" db:"reason"`

	Name             string `json:"name" db:"name"`
	ContactEmail     string `json:"contact_email" db:"contact_email"`
	DatabaseID       string `json:"database_id" db:"database_id"`
	DatabasePassword string `json:"database_password" db:"database_password"`
	Domain           string `json:"domain" db:"domain"`
	ExternalDomain   string `json:"external_domain" db:"external_domain"`
	FromEmail        string `json:"from_email" db:"from_email"`
	SecretKeyBase    string `json:"secret_key_base" db:"secret_key_base"`
}

func NewDeployment() *Deployment {
	return &Deployment{}
}

func (d *Deployment) SetDefaults() {
	d.DatabaseID = "f" + d.Name
	d.DatabasePassword = fh.GeneratePassword(16)
	d.SecretKeyBase = fh.GenerateSecret(128)
	d.Domain = d.Name + "." + fh.Config().Domain
}

// String is not required by pop and may be deleted
func (d Deployment) String() string {
	jd, _ := json.Marshal(d)
	return string(jd)
}

// Deployments is not required by pop and may be deleted
type Deployments []Deployment

// String is not required by pop and may be deleted
func (d Deployments) String() string {
	jd, _ := json.Marshal(d)
	return string(jd)
}

// from compute [a-z]([-a-z0-9]*[a-z0-9])?
const ACRONYM = "^[a-zA-Z][a-zA-Z0-9]{0,62}$"
const DNS_LABEL = "^[a-zA-Z0-9]|[a-zA-Z0-9][a-zA-Z0-9-]{0,61}[a-zA-Z0-9]$"

func (d *Deployment) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(
		&validators.EmailIsPresent{Field: d.FromEmail, Name: "FromEmail"},
		&validators.EmailLike{Field: d.ContactEmail, Name: "ContactEmail"},
		&validators.StringIsPresent{Field: d.Name, Name: "Name"},
		&validators.StringLengthInRange{Field: d.Name, Name: "Name", Min: 3, Max: 12},
		&validators.RegexMatch{Field: d.Name, Name: "Name", Expr: ACRONYM},
		// &validators.StringIsPresent{Field: d.ExternalDomain, Name: "ExternalDomain"},
		// &validators.StringLengthInRange{Field: d.ExternalDomain, Name: "ExternalDomain", Min: 3, Max: 63},
		// &validators.RegexMatch{Field: d.ExternalDomain, Name: "ExternalDomain", Expr: DNS_LABEL},
	), nil
}

// ValidateCreate gets run every time you call "pop.ValidateAndCreate" method.
// This method is not required and may be deleted.
// TODO does it use validate instead
// func (d *Deployment) ValidateCreate(tx *pop.Connection) (*validate.Errors, error) {
//         return validate.NewErrors(), nil
// }

// // ValidateUpdate gets run every time you call "pop.ValidateAndUpdate" method.
// // This method is not required and may be deleted.
// func (d *Deployment) ValidateUpdate(tx *pop.Connection) (*validate.Errors, error) {
//         return validate.NewErrors(), nil
// }
