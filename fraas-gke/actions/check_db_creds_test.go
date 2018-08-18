package actions_test

import (
	"testing"

	fg "manno.name/mm/fraas/fraas-gke"
	"manno.name/mm/fraas/fraas-gke/actions"

	. "github.com/onsi/gomega"
)

// need to export GOOGLE_APPLICATION_CREDENTIALS=/Users/mm/Downloads/PROJECT_COMPUTE_NETWORK_ADMIN.json service account
func TestCheckDBCreds(t *testing.T) {
	RegisterTestingT(t)

	a := actions.NewCheckDBCreds(fg.KubeClientset)

	t.Run("Apply", func(t *testing.T) { Expect(a.Apply()).NotTo(HaveOccurred()) })
}
