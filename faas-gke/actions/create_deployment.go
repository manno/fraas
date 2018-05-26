package actions

import (
	"google.golang.org/api/googleapi"
	api "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"k8s.io/client-go/kubernetes"

	"manno.name/mm/faas/faas-gke/models"
	"manno.name/mm/faas/faas-gke/specs"
	fh "manno.name/mm/faas/faas-helpers"
)

type CreateDeployment struct {
	clientset  *kubernetes.Clientset
	deployment *models.Deployment
}

func NewCreateDeployment(clientset *kubernetes.Clientset, deployment *models.Deployment) *CreateDeployment {
	return &CreateDeployment{clientset: clientset, deployment: deployment}
}

func (a *CreateDeployment) Apply() error {
	config := fh.Config()

	deployment := specs.NewDeploymentSpec(a.deployment, config)
	client := a.clientset.AppsV1beta1().Deployments(api.NamespaceDefault)
	if _, err := client.Create(deployment); err != nil {
		if googleapi.IsNotModified(err) {
			ResourceNotModified("deployment", a.deployment.WebName(), err)
			return nil
		}
		if errors.IsAlreadyExists(err) {
			_, err = client.Update(deployment)
			if googleapi.IsNotModified(err) {
				ResourceNotModified("deployment", a.deployment.WebName(), err)
				return nil
			}
			return err
		}
		return err
	}

	return nil
}

func (a *CreateDeployment) Revert() error {
	deletePolicy := metav1.DeletePropagationBackground
	client := a.clientset.AppsV1beta1().Deployments(api.NamespaceDefault)
	if err := client.Delete(a.deployment.WebName(), &metav1.DeleteOptions{PropagationPolicy: &deletePolicy}); err != nil {
		return err
	}

	return nil
}

func (s *CreateDeployment) Check() error {
	// podsClient := s.clientset.CoreV1().Pods(api.NamespaceDefault)
	// if pods, err := podsClient.List(metav1.ListOptions{LabelSelector: "app=" + deploymentID}); err != nil {
	//         return err
	// } else {
	//         for _, v := range pods.Items {
	//                 fmt.Printf("delete %s\n", v.GetName())
	//                 if err := podsClient.Delete(v.GetName(), &metav1.DeleteOptions{}); err != nil {
	//                         return err
	//                 }
	//         }
	// }
	return nil
}
