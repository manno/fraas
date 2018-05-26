package actions

import (
	"google.golang.org/api/googleapi"
	api "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/rest"
	"manno.name/mm/faas/faas-gke/models"
)

type CreateCertificate struct {
	client     dynamic.ResourceInterface
	deployment *models.Deployment
	spec       *unstructured.Unstructured
}

func NewCreateCertificate(config *rest.Config, deployment *models.Deployment, spec *unstructured.Unstructured) *CreateCertificate {
	clientset := NewDynamicClient(config)
	client := clientset.Resource(schema.GroupVersionResource{
		Group:    "certmanager.k8s.io",
		Version:  "v1alpha1",
		Resource: "certificates",
	}).Namespace(api.NamespaceDefault)
	return &CreateCertificate{deployment: deployment, client: client, spec: spec}
}

func (a *CreateCertificate) Apply() error {
	_, err := a.client.Create(a.spec)
	if err != nil {
		if googleapi.IsNotModified(err) {
			ResourceNotModified("certificate", a.deployment.TLSSecretName(), err)
			return nil
		}
		if errors.IsAlreadyExists(err) {
			ResourceAlreadyExists("certificate", a.deployment.TLSSecretName(), err)
			return nil
		}
	}

	return err
}

func (a *CreateCertificate) Revert() error {
	if err := a.client.Delete(a.deployment.TLSSecretName(), &metav1.DeleteOptions{}); err != nil {
		return err
	}
	return nil
}
