package route53

import (
	"context"
	"errors"
	"net"
	"time"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/route53"
)

// Route53 represents notification of DNS records.
type Route53 struct {
	Domain       string
	HostedZoneID string
	TTL          int32
}

// Notify notifies DNS records to Route53.
func (r *Route53) Notify(ctx context.Context, ipv4Addr net.IP, ipv6Addr net.IP) error {
	if ipv4Addr == nil && ipv6Addr == nil {
		return errors.New("either or both of ipv4Addr and ipv6Addr required")
	}

	var resSets []*route53.ResourceRecordSet
	if ipv4Addr != nil {
		resSets = append(resSets, r.makeChange(ipv4Addr, route53.RRTypeA))
	}
	if ipv6Addr != nil {
		resSets = append(resSets, r.makeChange(ipv6Addr, route53.RRTypeAaaa))
	}

	var changes []*route53.Change
	action := string(route53.ChangeActionUpsert)
	for _, r := range resSets {
		changes = append(
			changes,
			&route53.Change{
				Action:            &action,
				ResourceRecordSet: r,
			},
		)
	}
	comment := time.Now().Format(time.RFC3339)
	input := &route53.ChangeResourceRecordSetsInput{
		ChangeBatch: &route53.ChangeBatch{
			Changes: changes,
			Comment: &comment,
		},
		HostedZoneId: &r.HostedZoneID,
	}

	svc := route53.New(session.New())
	_, err := svc.ChangeResourceRecordSets(input)

	return err
}

func (r *Route53) makeChange(ip net.IP, recordType string) *route53.ResourceRecordSet {
	ttl := int64(r.TTL)
	ipAddr := ip.String()
	return &route53.ResourceRecordSet{
		Name: &r.Domain,
		TTL:  &ttl,
		Type: &recordType,
		ResourceRecords: []*route53.ResourceRecord{
			{
				Value: &ipAddr,
			},
		},
	}
}
