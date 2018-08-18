package actions_test

import (
	"testing"

	fg "manno.name/mm/fraas/fraas-gke"
	"manno.name/mm/fraas/fraas-gke/actions"
	"manno.name/mm/fraas/fraas-gke/models"
	"manno.name/mm/fraas/fraas-gke/specs"

	. "github.com/onsi/gomega"
)

func TestCreateCertificate(t *testing.T) {
	RegisterTestingT(t)

	d := &models.Deployment{Name: "test"}
	a := actions.NewCreateCertificate(fg.KubeConfig, d, specs.NewCertificate(d))

	t.Run("Apply", func(t *testing.T) { Expect(a.Apply()).NotTo(HaveOccurred()) })
	t.Run("Revert", func(t *testing.T) { Expect(a.Revert()).ToNot(HaveOccurred()) })
}
