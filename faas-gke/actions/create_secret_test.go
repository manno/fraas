package actions_test

import (
	"testing"

	fg "manno.name/mm/faas/faas-gke"
	"manno.name/mm/faas/faas-gke/actions"
	"manno.name/mm/faas/faas-gke/models"
	"manno.name/mm/faas/faas-gke/specs"

	. "github.com/onsi/gomega"
)

func TestCreateSecrets(t *testing.T) {
	RegisterTestingT(t)

	d := &models.Deployment{}
	spec := specs.NewRailsSecretsSpec(d)
	s := actions.NewCreateSecrets(fg.KubeClientset, spec)

	t.Run("Apply", func(t *testing.T) { Expect(s.Apply()).ToNot(HaveOccurred()) })
	t.Run("Revert", func(t *testing.T) { Expect(s.Revert()).ToNot(HaveOccurred()) })
}
