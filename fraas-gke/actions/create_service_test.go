package actions_test

import (
	"testing"

	fg "manno.name/mm/fraas/fraas-gke"
	"manno.name/mm/fraas/fraas-gke/actions"
	"manno.name/mm/fraas/fraas-gke/models"

	. "github.com/onsi/gomega"
)

func TestCreateService(t *testing.T) {
	RegisterTestingT(t)

	d := &models.Deployment{Name: "frab5"}
	s := actions.NewCreateService(fg.KubeClientset, d)

	t.Run("Apply", func(t *testing.T) { Expect(s.Apply()).ToNot(HaveOccurred()) })
	t.Run("Revert", func(t *testing.T) { Expect(s.Revert()).ToNot(HaveOccurred()) })
}
