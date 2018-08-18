package actions

import (
	"google.golang.org/api/googleapi"
	apiv1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type CreateSecrets struct {
	clientset *kubernetes.Clientset
	spec      *apiv1.Secret
}

func NewCreateSecrets(clientset *kubernetes.Clientset, spec *apiv1.Secret) *CreateSecrets {
	return &CreateSecrets{clientset: clientset, spec: spec}
}

func (a *CreateSecrets) Apply() error {
	client := a.clientset.CoreV1().Secrets(apiv1.NamespaceDefault)

	if _, err := client.Create(a.spec); err != nil {
		if errors.IsAlreadyExists(err) {
			_, err := client.Update(a.spec)
			if googleapi.IsNotModified(err) {
				ResourceNotModified("secret", a.spec.GetName(), err)
				return nil
			}
			return err
		}
		return err
	}

	return nil
}

func (a *CreateSecrets) Revert() error {
	secretsClient := a.clientset.CoreV1().Secrets(apiv1.NamespaceDefault)
	return secretsClient.Delete(a.spec.GetName(), &metav1.DeleteOptions{})
}
