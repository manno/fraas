package actions_test

import (
	"testing"

	fg "manno.name/mm/faas/faas-gke"
	"manno.name/mm/faas/faas-gke/actions"
	"manno.name/mm/faas/faas-gke/models"

	. "github.com/onsi/gomega"
)

func TestCreateDeployment(t *testing.T) {
	RegisterTestingT(t)

	d := &models.Deployment{
		Name:             "frab5",
		ContactEmail:     "mail@example.org",
		DatabaseID:       "ffrab5",
		DatabasePassword: "JUM(K!B@7j",
		Domain:           "frab5.frab.app",
		ExternalDomain:   "frab5.frab.app",
		FromEmail:        "mail@example.org",
		SecretKeyBase:    "uVt17MwiXMriYoBgap8C632V2zits4ZQDxQ2og5eLEDXei6upd5DIpDKJKtMV1hCtipD41mUCCSBIuYXPF1g8ix2zyIhKEhYHdkyWQAvhpxOxLOcyTZrF6WdzwyNlVfb",
	}
	s := actions.NewCreateDeployment(fg.KubeClientset, d)

	t.Run("Apply", func(t *testing.T) { Expect(s.Apply()).ToNot(HaveOccurred()) })
	t.Run("Revert", func(t *testing.T) { Expect(s.Revert()).ToNot(HaveOccurred()) })
}
