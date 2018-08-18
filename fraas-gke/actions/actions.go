package actions

import (
	"log"

	"golang.org/x/net/context"
	"golang.org/x/oauth2/google"
	compute "google.golang.org/api/compute/v1"
	dns "google.golang.org/api/dns/v1"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/rest"
)

type Action interface {
	Apply() error
	Revert() error
}

func ResourceAlreadyExists(origin string, name string, err error) {
	log.Printf("INFO: not updating %s '%s'. Already exists: %v", origin, name, err)
}

func ResourceNotModified(origin string, name string, err error) {
	log.Printf("INFO: %s resource '%s' is not modified: %v", origin, name, err)
}

func NewDNSClient(ctx context.Context) *dns.Service {
	c, err := google.DefaultClient(ctx, dns.CloudPlatformScope)
	if err != nil {
		log.Fatal(err)
	}
	client, _ := dns.New(c)
	if err != nil {
		log.Fatal(err)
	}
	return client
}

func NewComputeClient(ctx context.Context) *compute.Service {
	c, err := google.DefaultClient(ctx, compute.ComputeScope)
	if err != nil {
		log.Fatal(err)
	}
	client, err := compute.New(c)
	if err != nil {
		log.Fatal(err)
	}
	return client
}

func NewDynamicClient(config *rest.Config) dynamic.Interface {
	clientset, err := dynamic.NewForConfig(config)
	if err != nil {
		log.Fatal(err)
	}
	return clientset
}
