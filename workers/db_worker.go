package workers

import (
	"errors"
	"fmt"
	"log"

	"github.com/gobuffalo/buffalo/worker"
	"github.com/gobuffalo/pop"
	fh "manno.name/mm/fraas/fraas-helpers"
	"manno.name/mm/fraas/models"
)

func SetDB(args worker.Args) error {
	d := &models.Deployment{}
	if err := models.DB.Find(d, args["deployment_id"]); err != nil {
		return errors.New("Worker failed to load deployment from DB")
	}

	db, err := NewAdminConnection()
	if err != nil {
		return err
	}

	err = db.RawQuery(fmt.Sprintf("CREATE ROLE %s WITH CREATEDB LOGIN password '%s'", d.DatabaseID, d.DatabasePassword)).Exec()
	if err != nil {
		return err
	}

	err = db.RawQuery(fmt.Sprintf("GRANT %s TO %s;", d.DatabaseID, fh.Config().Database.AdminUser)).Exec()
	if err != nil {
		return err
	}

	err = db.RawQuery(fmt.Sprintf("CREATE DATABASE %s ENCODING 'UTF8' OWNER %s", d.DatabaseID, d.DatabaseID)).Exec()
	if err != nil {
		return err
	}

	return nil
}

func UnsetDB(args worker.Args) error {
	d := &models.Deployment{}
	if err := models.DB.Find(d, args["deployment_id"]); err != nil {
		return errors.New("Worker failed to load deployment from DB")
	}

	db, err := NewAdminConnection()
	if err != nil {
		return err
	}

	err = db.RawQuery("DROP DATABASE " + d.DatabaseID).Exec()
	if err != nil {
		log.Println(err)
	}

	err = db.RawQuery("DROP ROLE " + d.DatabaseID).Exec()
	if err != nil {
		log.Println(err)
	}

	return nil
}

func NewAdminConnection() (*pop.Connection, error) {
	url := fmt.Sprintf(
		"postgres://%s:%s@127.0.0.1:5432/?sslmode=disable",
		fh.Config().Database.AdminUser,
		fh.Config().Database.AdminPassword,
	)
	options := &pop.ConnectionDetails{URL: url}
	options.Finalize()

	db, err := pop.NewConnection(options)
	if err != nil {
		log.Panic(err)
	}

	err = db.Open()
	if err != nil {
		log.Panic(err)
	}
	return db, err
}
