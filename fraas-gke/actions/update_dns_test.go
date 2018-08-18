package actions_test

import (
	"testing"

	"manno.name/mm/fraas/fraas-gke/actions"
	"manno.name/mm/fraas/fraas-gke/models"

	. "github.com/onsi/gomega"
)

// need to export GOOGLE_APPLICATION_CREDENTIALS=/Users/mm/Downloads/PROJECT_COMPUTE_NETWORK_ADMIN.json service account
func TestUpdateDNS(t *testing.T) {
	RegisterTestingT(t)

	d := &models.Deployment{Name: "test"}
	a := actions.NewUpdateDNS(d)

	t.Run("Apply", func(t *testing.T) { Expect(a.Apply()).NotTo(HaveOccurred()) })
	t.Run("Revert", func(t *testing.T) { Expect(a.Revert()).ToNot(HaveOccurred()) })
}
