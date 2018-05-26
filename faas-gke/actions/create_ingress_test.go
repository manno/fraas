package actions_test

import (
	"testing"

	fg "manno.name/mm/faas/faas-gke"
	"manno.name/mm/faas/faas-gke/actions"
	"manno.name/mm/faas/faas-gke/models"
	"manno.name/mm/faas/faas-gke/specs"

	. "github.com/onsi/gomega"
)

func TestCreateIngress(t *testing.T) {
	RegisterTestingT(t)

	d := &models.Deployment{Name: "frab5", Domain: "test.frab.app", ExternalDomain: "test-frab.example.org"}
	s := actions.NewCreateIngress(fg.KubeClientset, specs.NewIngressSpec(d))

	t.Run("Apply", func(t *testing.T) { Expect(s.Apply()).ToNot(HaveOccurred()) })
	t.Run("Revert", func(t *testing.T) { Expect(s.Revert()).ToNot(HaveOccurred()) })
}

func TestCreateIngressTLS(t *testing.T) {
	RegisterTestingT(t)

	d := &models.Deployment{Name: "frab5", Domain: "test.frab.app", ExternalDomain: "test-frab.example.org"}
	s := actions.NewCreateIngress(fg.KubeClientset, specs.NewIngressTLSSpec(d))

	t.Run("Apply", func(t *testing.T) { Expect(s.Apply()).ToNot(HaveOccurred()) })
	t.Run("Revert", func(t *testing.T) { Expect(s.Revert()).ToNot(HaveOccurred()) })
}
