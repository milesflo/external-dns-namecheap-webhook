package provider

import (
	"context"

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

	recordsToCreate := client.UpsertRequest{}

	for _, ep := range changes.UpdateNew {
		for _, target := range ep.Targets {
			recordsToCreate.Creates = append(recordsToCreate.Creates, client.RecordSet{
				Name:    ep.DNSName,
				Type:    ep.RecordType,
				TTL:     int(ep.RecordTTL),
				Address: target,
			})
		}
	}

	for _, ep := range changes.Create {
		for _, target := range ep.Targets {
			recordsToCreate.Creates = append(recordsToCreate.Creates, client.RecordSet{
				Name:    ep.DNSName,
				Type:    ep.RecordType,
				TTL:     int(ep.RecordTTL),
				Address: target,
			})
		}
	}

	return p.client.UpsertRecordSets(ctx, recordsToCreate)
}

func (p namecheapProvider) Records(ctx context.Context) ([]*endpoint.Endpoint, error) {
	zones, err := p.client.ListZones(ctx)
	if err != nil {
		return nil, err
	}

	var endpoints []*endpoint.Endpoint

	for _, zone := range zones {
		records, err := p.client.ListRecordSets(ctx, zone.Name)
		if err != nil {
			return nil, err
		}

		for _, record := range records {
			ep := endpoint.NewEndpointWithTTL(zone.Name, record.Type, endpoint.TTL(record.TTL), record.Data...)

			endpoints = append(endpoints, ep)
		}
	}

	return endpoints, nil
}
