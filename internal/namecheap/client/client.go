package client

import (
	"context"
	"errors"
	"fmt"
	"net/netip"

	sdk "github.com/namecheap/go-namecheap-sdk/v2/namecheap"
)

type NamecheapDNSClient interface {
	ListZones(ctx context.Context) ([]Zone, error)
	ListRecordSets(ctx context.Context, zoneID string) ([]RecordSet, error)
	UpsertRecordSets(ctx context.Context, req UpsertRequest) error
}

type RecordSet struct {
	Name string
	Type string
	TTL  int64
	Data []string
}

type UpsertRequest struct {
	DnsZoneID    string
	Deletions    []RecordSet
	Replacements []RecordSet
	Merges       []RecordSet
}

type Zone struct {
	ID        string
	Name      string
	IsPrivate bool
}

type NamecheapClient struct {
	sdk *sdk.Client
}

func NewNamecheapClient(username, apiKey, clientIp string, useSandbox bool) (*NamecheapClient, error) {
	if username == "" || len(username) > 20 {
		return nil, fmt.Errorf("Username %s malformed", username)
	}
	if apiKey == "" || len(apiKey) > 50 {
		return nil, fmt.Errorf("APIKey malformed")
	}
	if clientIp == "" || len(clientIp) > 15 {
		return nil, fmt.Errorf("Client IP %s malformed", username)
	}

	_, err := netip.ParseAddr(clientIp)
	if err != nil {
		return nil, err
	}

	client := sdk.NewClient(&sdk.ClientOptions{
		UserName:   username,
		ApiUser:    username,
		ApiKey:     apiKey,
		ClientIp:   clientIp,
		UseSandbox: useSandbox,
	})
	return &NamecheapClient{
		sdk: client,
	}, nil
}
func (c *NamecheapClient) Validate(ctx context.Context) error {
	// Get something with current client config, return error if it fails
	_, err := c.sdk.Domains.GetList(&sdk.DomainsGetListArgs{
		PageSize: sdk.Int(1),
	})
	return err
}

func (c *NamecheapClient) ListZones(ctx context.Context) ([]Zone, error) {
	zones := make([]Zone, 0)

	//depaginate
	pageSize := 100
	page := 1
	for {
		res, err := c.sdk.Domains.GetList(&sdk.DomainsGetListArgs{
			Page:     sdk.Int(page),
			PageSize: sdk.Int(pageSize),
		})
		if err != nil {
			return nil, err
		}
		for _, domain := range *res.Domains {
			zones = append(zones, Zone{
				ID:        *domain.ID,
				Name:      *domain.Name,
				IsPrivate: *domain.IsLocked,
			})
		}
		// If ceiling is larger than total items, break
		if *res.Paging.CurrentPage**res.Paging.PageSize > *res.Paging.TotalItems {
			break
		}
	}

	return zones, nil
}

func (c *NamecheapClient) ListRecordSets(ctx context.Context, zoneID string) ([]RecordSet, error) {
	return nil, errors.New("not implemented")
}

func (c *NamecheapClient) UpsertRecordSets(ctx context.Context, req UpsertRequest) error {
	return errors.New("not implemented")
}
