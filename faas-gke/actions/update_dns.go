package actions

import (
	"fmt"
	"log"
	"strings"
	"time"

	"golang.org/x/net/context"

	"github.com/Rican7/retry"
	"github.com/Rican7/retry/backoff"
	"github.com/Rican7/retry/strategy"

	"google.golang.org/api/compute/v1"
	"google.golang.org/api/dns/v1"
	"google.golang.org/api/googleapi"

	"manno.name/mm/faas/faas-gke/models"
	fh "manno.name/mm/faas/faas-helpers"
)

type UpdateDNS struct {
	ctx           context.Context
	computeClient *compute.Service
	dnsClient     *dns.Service
	deployment    *models.Deployment
}

func NewUpdateDNS(deployment *models.Deployment) *UpdateDNS {
	ctx := context.Background()
	computeClient := NewComputeClient(ctx)
	dnsClient := NewDNSClient(ctx)
	return &UpdateDNS{deployment: deployment, computeClient: computeClient, dnsClient: dnsClient, ctx: ctx}
}

func (a *UpdateDNS) Apply() error {
	cfg := fh.Config()

	var externalIP string
	err := retry.Retry(
		func(attempt uint) error {
			var err error
			externalIP, err = a.getExternalIP(cfg.Google.ProjectID, a.deployment.IPName())
			if err != nil || externalIP == "" {
				log.Printf("waiting for external ip: %s\n", err)
				return fmt.Errorf("failed to fetch (attempt #%d) with reason %s", attempt, err)
			}
			log.Printf("DNSUPDATE: acquired external ip: %s\n", externalIP)
			return nil
		},
		strategy.Limit(15),
		strategy.Backoff(backoff.Fibonacci(10*time.Millisecond)),
	)
	if err != nil || externalIP == "" {
		return fmt.Errorf("giving up on fetching external ip %s: %s", a.deployment.IPName(), err)
	}

	change := &dns.Change{Additions: recordSet(a.deployment.Domain, externalIP)}
	_, err = a.dnsClient.Changes.Create(cfg.Google.ProjectID, cfg.Google.DNSZone, change).Context(a.ctx).Do()
	if err != nil {
		serr, ok := err.(*googleapi.Error)
		if ok && strings.Contains(serr.Message, "already exists") {
			ResourceAlreadyExists("dns record-set", a.deployment.Domain, err)
			return nil
		}

		return err
	}
	return err
}

func (a *UpdateDNS) Revert() error {
	cfg := fh.Config()
	externalIP, err := a.getExternalIP(cfg.Google.ProjectID, a.deployment.IPName())
	if err != nil {
		return err
	}

	change := &dns.Change{Deletions: recordSet(a.deployment.Domain, externalIP)}
	_, err = a.dnsClient.Changes.Create(cfg.Google.ProjectID, cfg.Google.DNSZone, change).Context(a.ctx).Do()
	return err
}

func (a *UpdateDNS) getExternalIP(projectID string, ipName string) (string, error) {
	address, err := a.computeClient.GlobalAddresses.Get(projectID, ipName).Context(a.ctx).Do()
	if err != nil {
		return "", err
	}
	if address == nil {
		return "", fmt.Errorf("failed to fetch address")
	}

	return address.Address, nil
}

func recordSet(domain, externalIP string) []*dns.ResourceRecordSet {
	return []*dns.ResourceRecordSet{
		&dns.ResourceRecordSet{
			Name:    domain + ".",
			Type:    "A",
			Ttl:     300,
			Rrdatas: []string{externalIP},
		},
	}
}
