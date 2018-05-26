package actions

import (
	"errors"

	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

const cloudsql_creds = "cloudsql-frab-credentials"

type CheckDBCreds struct {
	clientset *kubernetes.Clientset
}

func NewCheckDBCreds(clientset *kubernetes.Clientset) *CheckDBCreds {
	return &CheckDBCreds{clientset: clientset}
}

func (a *CheckDBCreds) Apply() error {
	client := a.clientset.CoreV1().Secrets(apiv1.NamespaceDefault)

	if secrets, err := client.List(metav1.ListOptions{}); err != nil {
		return err
	} else {
		for _, v := range secrets.Items {
			if v.GetName() == cloudsql_creds {
				return nil
			}
		}
	}

	return errors.New("failed to find database creds: " + cloudsql_creds)
}

func (a *CheckDBCreds) Revert() error {
	return nil
}
