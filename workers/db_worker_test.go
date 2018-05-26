package workers

import (
	"testing"

	fh "manno.name/mm/faas/faas-helpers"
)

func TestAdminDB(t *testing.T) {
	if err := fh.ConfigFromEnv(); err != nil {
		t.Error("load config", err)
	}

	db, err := NewAdminConnection()
	if err != nil {
		t.Error("new db connextion", err)
	}
	err = db.RawQuery("CREATE DATABASE testdb").Exec()
	if err != nil {
		t.Error("create db", err)
	}
	err = db.RawQuery("DROP DATABASE testdb").Exec()
	if err != nil {
		t.Error("drop db", err)
	}
}
