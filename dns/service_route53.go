package dns

import (
	"context"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/route53"
	"github.com/aws/aws-sdk-go-v2/service/route53/types"
	"github.com/pkg/errors"
)

// Route53 are the methods we need exposed to manage DNS entries.
type Route53 interface {
	ListHostedZones(ctx context.Context, params *route53.ListHostedZonesInput, optFns ...func(*route53.Options)) (*route53.ListHostedZonesOutput, error)
	ChangeResourceRecordSets(ctx context.Context, params *route53.ChangeResourceRecordSetsInput, optFns ...func(*route53.Options)) (*route53.ChangeResourceRecordSetsOutput, error)
}

// ServiceRoute53 handles updating DNS records using route53 as a backend.
type ServiceRoute53 struct {
	client Route53
}

// NewServiceRoute53 iniitalizes a new ServiceRoute53.
func NewServiceRoute53(client Route53) *ServiceRoute53 {
	s := &ServiceRoute53{
		client: client,
	}

	return s
}

// DNS handles dns requests.
func (s *ServiceRoute53) DNS(ctx context.Context, req Request) error {
	zoneID, err := getZoneID(ctx, s.client, req.Name)
	if err != nil {
		return err
	}

	params := route53.ChangeResourceRecordSetsInput{
		HostedZoneId: &zoneID,
		ChangeBatch: &types.ChangeBatch{
			Changes: []types.Change{
				{
					Action: types.ChangeActionUpsert,
					ResourceRecordSet: &types.ResourceRecordSet{
						Name: &req.Name,
						Type: types.RRTypeA,
						TTL:  aws.Int64(300),
						ResourceRecords: []types.ResourceRecord{
							{Value: &req.Value},
						},
					},
				},
			},
		},
	}

	_, err = s.client.ChangeResourceRecordSets(ctx, &params)
	return err
}

func listAllHostedZones(ctx context.Context, client Route53) ([]types.HostedZone, error) {
	var zones []types.HostedZone
	var marker string
	for {
		params := route53.ListHostedZonesInput{}
		if marker != "" {
			params.Marker = &marker
		}
		resp, err := client.ListHostedZones(ctx, &params)
		if err != nil {
			return nil, err
		}

		zones = append(zones, resp.HostedZones...)

		if !resp.IsTruncated {
			break
		}
		marker = *resp.NextMarker
	}

	return zones, nil
}

func getZoneID(ctx context.Context, client Route53, record string) (string, error) {
	zones, err := listAllHostedZones(ctx, client)
	if err != nil {
		return "", errors.Wrap(err, "retrieving zones from route53")
	}

	return determineZoneID(zones, record)
}

// This function tries to determine the most appropriate zone to place the record into. If a zone cannot be determined, an error is returned.
func determineZoneID(zones []types.HostedZone, record string) (string, error) {
	if !strings.HasSuffix(record, ".") {
		record += "."
	}

	var zone *types.HostedZone
	for i := range zones {
		if !strings.HasSuffix(record, "."+*zones[i].Name) {
			continue
		}

		if zone == nil || len(*zones[i].Name) > len(*zone.Name) {
			zone = &zones[i]
		}
	}

	if zone == nil {
		return "", errZoneNotFound
	}

	return *zone.Id, nil
}

var errZoneNotFound = errors.New("suitable zone not found")
