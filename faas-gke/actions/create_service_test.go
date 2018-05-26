package actions_test

import (
	"testing"

	fg "manno.name/mm/faas/faas-gke"
	"manno.name/mm/faas/faas-gke/actions"
	"manno.name/mm/faas/faas-gke/models"

	. "github.com/onsi/gomega"
)

func TestCreateService(t *testing.T) {
	RegisterTestingT(t)

	d := &models.Deployment{Name: "frab5"}
	s := actions.NewCreateService(fg.KubeClientset, d)

	t.Run("Apply", func(t *testing.T) { Expect(s.Apply()).ToNot(HaveOccurred()) })
	t.Run("Revert", func(t *testing.T) { Expect(s.Revert()).ToNot(HaveOccurred()) })
}
