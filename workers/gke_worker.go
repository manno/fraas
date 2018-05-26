package workers

import (
	"errors"

	"github.com/gobuffalo/buffalo/worker"
	gke "manno.name/mm/faas/faas-gke"
	"manno.name/mm/faas/models"
)

func SetGKE(args worker.Args) error {
	deployment := &models.Deployment{}
	if err := models.DB.Find(deployment, args["deployment_id"]); err != nil {
		return errors.New("Worker failed to load deployment from DB")
	}

	gke.Setup()
	return gke.Apply(deployment.String())
}

func UnsetGKE(args worker.Args) error {
	deployment := &models.Deployment{}
	if err := models.DB.Find(deployment, args["deployment_id"]); err != nil {
		return errors.New("Worker failed to load deployment from DB")
	}

	gke.Setup()
	return gke.Revert(deployment.String())
}
