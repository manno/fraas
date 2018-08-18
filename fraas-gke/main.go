package fraas_gke

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"reflect"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	// Uncomment the following line to load the gcp plugin (only required to authenticate against GKE clusters).
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"

	"k8s.io/client-go/rest"
	"manno.name/mm/fraas/fraas-gke/actions"
	"manno.name/mm/fraas/fraas-gke/models"
	"manno.name/mm/fraas/fraas-gke/specs"
	fh "manno.name/mm/fraas/fraas-helpers"
)

var KubeClientset *kubernetes.Clientset
var KubeConfig *rest.Config

func SetupKubeClient() error {
	clientset, err := kubernetes.NewForConfig(KubeConfig)
	if err != nil {
		return err
	}

	KubeClientset = clientset
	return nil
}

func SetupKubeConfig() error {
	kubeconfig := os.Getenv("KUBE_CONFIG")
	if kubeconfig == "" {
		kubeconfig = filepath.Join(os.Getenv("HOME"), ".kube", "config")
	}

	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		log.Printf("INFO: cannot use kube config: %s\n", err)
		config, err = rest.InClusterConfig()
		if err != nil {
			return err
		}
	}

	KubeConfig = config
	return nil
}

func Setup() {
	if err := fh.ConfigFromEnv(); err != nil {
		log.Fatal("Worker failed to load site config from ENV", err)
	}
	if err := SetupKubeConfig(); err != nil {
		log.Fatal(err)
	}
	if err := SetupKubeClient(); err != nil {
		log.Fatal(err)
	}
}

func Apply(d string) error {
	deployment := &models.Deployment{}
	err := json.Unmarshal([]byte(d), deployment)
	if err != nil {
		fmt.Printf("%#v\n", d)
	}

	actions := defaultActions(deployment)
	for _, action := range actions {
		log.Printf("INFO: applying %s for %s", reflect.TypeOf(action).Elem().Name(), deployment.Name)
		if err := action.Apply(); err != nil {
			// TODO revert actions up to this one
			log.Println(err)
			return err
		}
	}

	return nil
}

func Revert(d string) error {
	deployment := &models.Deployment{}
	err := json.Unmarshal([]byte(d), deployment)
	if err != nil {
		fmt.Printf("%#v\n", d)
	}

	actions := defaultActions(deployment)
	for i := len(actions) - 1; i >= 0; i-- {
		action := actions[i]
		log.Printf("INFO: reverting %s for %s", reflect.TypeOf(action).Elem().Name(), deployment.Name)
		if err := action.Revert(); err != nil {
			log.Println(err)
		}
	}
	return nil
}

func defaultActions(deployment *models.Deployment) []actions.Action {
	return []actions.Action{
		actions.NewCreateSecrets(KubeClientset, specs.NewRailsSecretsSpec(deployment)),
		actions.NewCreateConfigMap(KubeClientset, deployment),
		actions.NewCreateDeployment(KubeClientset, deployment),
		actions.NewCreateService(KubeClientset, deployment),
		actions.NewComputeIP(deployment),
		//actions.NewCreateIngress(KubeClientset, specs.NewIngressSpec(deployment)),
		actions.NewCheckDBCreds(KubeClientset),
		actions.NewUpdateDNS(deployment),
		actions.NewCreateCertificate(KubeConfig, deployment, specs.NewCertificate(deployment)),
		actions.NewCreateIngress(KubeClientset, specs.NewIngressTLSSpec(deployment)),
	}
}
