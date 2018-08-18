package actions

import (
	"google.golang.org/api/googleapi"
	api "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"k8s.io/client-go/kubernetes"

	"manno.name/mm/fraas/fraas-gke/models"
	"manno.name/mm/fraas/fraas-gke/specs"
)

type CreateService struct {
	clientset  *kubernetes.Clientset
	deployment *models.Deployment
}

func NewCreateService(clientset *kubernetes.Clientset, deployment *models.Deployment) *CreateService {
	return &CreateService{clientset: clientset, deployment: deployment}
}

func (a *CreateService) Apply() error {
	serviceDeployment := specs.NewServiceSpec(a.deployment)
	servicesClient := a.clientset.CoreV1().Services(api.NamespaceDefault)
	if _, err := servicesClient.Create(serviceDeployment); err != nil {
		if googleapi.IsNotModified(err) {
			ResourceNotModified("service", a.deployment.BackendName(), err)
			return nil
		}
		if errors.IsAlreadyExists(err) {
			ResourceAlreadyExists("service", a.deployment.BackendName(), err)
			return nil
		}
		return err
	}

	return nil
}

func (a *CreateService) Revert() error {
	servicesClient := a.clientset.CoreV1().Services(api.NamespaceDefault)
	if err := servicesClient.Delete(a.deployment.BackendName(), &metav1.DeleteOptions{}); err != nil {
		return err
	}

	return nil
}
