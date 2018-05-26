package actions

import (
	"google.golang.org/api/googleapi"
	api "k8s.io/api/core/v1"
	ext "k8s.io/api/extensions/v1beta1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type CreateIngress struct {
	clientset *kubernetes.Clientset
	spec      *ext.Ingress
}

func NewCreateIngress(clientset *kubernetes.Clientset, spec *ext.Ingress) *CreateIngress {
	return &CreateIngress{clientset: clientset, spec: spec}
}

func (a *CreateIngress) Apply() error {

	client := a.clientset.ExtensionsV1beta1().Ingresses(api.NamespaceDefault)
	if _, err := client.Create(a.spec); err != nil {
		if googleapi.IsNotModified(err) {
			ResourceNotModified("ingress", a.spec.GetName(), err)
			return nil
		}
		if errors.IsAlreadyExists(err) {
			_, err := client.Update(a.spec)
			if googleapi.IsNotModified(err) {
				ResourceNotModified("ingress", a.spec.GetName(), err)
				return nil
			}
			return err
		}
		return err
	}

	return nil
}

func (a *CreateIngress) Revert() error {
	client := a.clientset.ExtensionsV1beta1().Ingresses(api.NamespaceDefault)
	if err := client.Delete(a.spec.GetName(), &metav1.DeleteOptions{}); err != nil {
		return err
	}
	return nil
}
