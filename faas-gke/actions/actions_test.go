package actions_test

import (
	"os"
	"testing"

	fg "manno.name/mm/faas/faas-gke"
)

func TestMain(m *testing.M) {
	fg.Setup()
	os.Exit(m.Run())
}
