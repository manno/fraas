package grifts

import (
	"fmt"

	"github.com/gobuffalo/pop"
	"github.com/markbates/grift/grift"
	"github.com/pkg/errors"
	fh "manno.name/mm/faas/faas-helpers"
	"manno.name/mm/faas/models"
)

var _ = grift.Namespace("db", func() {
	grift.Desc("seed", "Seeds a database")
	grift.Add("seed", func(c *grift.Context) error {
		return models.DB.Transaction(func(tx *pop.Connection) error {
			if err := tx.TruncateAll(); err != nil {
				return errors.WithStack(err)
			}

			password := fh.GeneratePassword(12)
			u := &models.User{
				Email:                "admin@example.org",
				Password:             password,
				PasswordConfirmation: password,
			}

			if _, err := u.Create(tx); err != nil {
				return errors.WithStack(err)
			}
			fmt.Println("created user admin@example.org with password: " + password)
			return nil
		})
	})

})
