package actions

import (
	"strings"

	"golang.org/x/net/context"

	compute "google.golang.org/api/compute/v1"
	"google.golang.org/api/googleapi"

	"manno.name/mm/fraas/fraas-gke/models"
	fh "manno.name/mm/fraas/fraas-helpers"
)

type ComputeIP struct {
	ctx        context.Context
	client     *compute.Service
	deployment *models.Deployment
}

func NewComputeIP(deployment *models.Deployment) *ComputeIP {
	ctx := context.Background()
	client := NewComputeClient(ctx)
	return &ComputeIP{deployment: deployment, client: client, ctx: ctx}
}

func (a *ComputeIP) Apply() error {
	cfg := fh.Config().Google
	_, err := a.client.GlobalAddresses.Insert(cfg.ProjectID, NewExternalIP(a)).Context(a.ctx).Do()

	if err != nil {
		serr, ok := err.(*googleapi.Error)
		if ok && strings.Contains(serr.Message, "already exists") {
			ResourceAlreadyExists("address", a.deployment.IPName(), err)
			return nil
		}

		return err
	}
	return nil
}

func NewExternalIP(a *ComputeIP) *compute.Address {
	return &compute.Address{
		Name:        a.deployment.IPName(),
		Description: a.deployment.DeploymentID(),
		IpVersion:   "IPV4",
	}
}

func (a *ComputeIP) Revert() error {
	cfg := fh.Config().Google
	_, err := a.client.GlobalAddresses.Delete(cfg.ProjectID, a.deployment.IPName()).Context(a.ctx).Do()
	return err
}
