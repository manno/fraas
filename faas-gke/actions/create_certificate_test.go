package actions_test

import (
	"testing"

	fg "manno.name/mm/faas/faas-gke"
	"manno.name/mm/faas/faas-gke/actions"
	"manno.name/mm/faas/faas-gke/models"
	"manno.name/mm/faas/faas-gke/specs"

	. "github.com/onsi/gomega"
)

func TestCreateCertificate(t *testing.T) {
	RegisterTestingT(t)

	d := &models.Deployment{Name: "test"}
	a := actions.NewCreateCertificate(fg.KubeConfig, d, specs.NewCertificate(d))

	t.Run("Apply", func(t *testing.T) { Expect(a.Apply()).NotTo(HaveOccurred()) })
	t.Run("Revert", func(t *testing.T) { Expect(a.Revert()).ToNot(HaveOccurred()) })
}
