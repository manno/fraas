package actions

import (
	apiv1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"manno.name/mm/fraas/fraas-gke/models"
	"manno.name/mm/fraas/fraas-gke/specs"
	fh "manno.name/mm/fraas/fraas-helpers"
)

type CreateConfigMap struct {
	clientset  *kubernetes.Clientset
	deployment *models.Deployment
}

func NewCreateConfigMap(clientset *kubernetes.Clientset, deployment *models.Deployment) *CreateConfigMap {
	return &CreateConfigMap{clientset: clientset, deployment: deployment}
}

func (a *CreateConfigMap) Apply() error {
	client := a.clientset.CoreV1().ConfigMaps(apiv1.NamespaceDefault)

	c := fh.Config()
	if _, err := client.Create(specs.NewRailsConfigMap(a.deployment, c)); err != nil {
		if errors.IsAlreadyExists(err) {
			ResourceAlreadyExists("configmap", a.deployment.ConfigName(), err)
			return nil
		}
		return err
	}
	return nil
}

func (a *CreateConfigMap) Revert() error {
	client := a.clientset.CoreV1().ConfigMaps(apiv1.NamespaceDefault)
	return client.Delete(a.deployment.ConfigName(), &metav1.DeleteOptions{})
}
