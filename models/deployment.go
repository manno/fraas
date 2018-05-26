package models

import (
	"encoding/json"
	"time"

	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/uuid"
	"github.com/gobuffalo/validate"
	"github.com/gobuffalo/validate/validators"
)

type Deployment struct {
	ID        uuid.UUID `json:"id" db:"id"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`

	Name             string `json:"name" db:"name"`
	DatabaseID       string `json:"database_id" db:"database_id"`
	DatabasePassword string `json:"database_password" db:"database_password"`
	Domain           string `json:"domain" db:"domain"`
	ExternalDomain   string `json:"external_domain" db:"external_domain"`
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

const ACRONYM = "^[a-zA-Z][a-zA-Z0-9]{0,62}$"
const DNS_LABEL = "^[a-zA-Z0-9]|[a-zA-Z0-9][a-zA-Z0-9-]{0,61}[a-zA-Z0-9]$"

func (d *Deployment) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(
		&validators.StringIsPresent{Field: d.Name, Name: "Name"},
		&validators.StringLengthInRange{Field: d.Name, Name: "Name", Min: 3, Max: 12},
		&validators.RegexMatch{Field: d.Name, Name: "Name", Expr: ACRONYM},
		&validators.StringIsPresent{Field: d.ExternalDomain, Name: "ExternalDomain"},
		&validators.StringLengthInRange{Field: d.ExternalDomain, Name: "ExternalDomain", Min: 3, Max: 63},
		&validators.RegexMatch{Field: d.ExternalDomain, Name: "ExternalDomain", Expr: DNS_LABEL},
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
