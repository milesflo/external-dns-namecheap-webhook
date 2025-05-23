package provider

import (
	"context"
	"errors"

	"github.com/milesflo/external-dns-namecheap-webhook/internal/namecheap/client"
	"sigs.k8s.io/external-dns/endpoint"
	"sigs.k8s.io/external-dns/plan"
	"sigs.k8s.io/external-dns/provider"
)

// namecheap provider type
type namecheapProvider struct {
	provider.BaseProvider
	client client.NamecheapDNSClient

	// only consider hosted zones managing domains ending in this suffix
	domainFilter endpoint.DomainFilter
	dryRun       bool
}

func NewNamecheapProvider(domainFilter endpoint.DomainFilter, dryRun bool, client client.NamecheapDNSClient) provider.Provider {
	return &namecheapProvider{
		client:       client,
		domainFilter: domainFilter,
		dryRun:       dryRun,
	}
}

func (p namecheapProvider) ApplyChanges(ctx context.Context, changes *plan.Changes) error {
	return errors.New("not implemented")
}

func (p namecheapProvider) Records(ctx context.Context) ([]*endpoint.Endpoint, error) {
	return nil, errors.New("not implemented")
}
