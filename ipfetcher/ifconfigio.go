package ipfetcher

import (
	"context"
	"fmt"
	"net"
)

// IfconfigIoFetcher represents to fetch remote IPs from ifconfig.io
type IfconfigIoFetcher struct{}

const (
	ifconfigIoURL = "https://ifconfig.io"
)

// FetchIPv4 fetches IPv4 address from ifconfig.io
func (i *IfconfigIoFetcher) FetchIPv4(ctx context.Context) (net.IP, error) {
	ipBytes, err := getByIPv4(ctx, ifconfigIoURL)
	if err != nil {
		return nil, err
	}

	ip := net.ParseIP(string(ipBytes))
	if ip == nil {
		return nil, fmt.Errorf("fetched string cannot be parsed as IP address. string: %s")
	}
	return ip, nil
}

// FetchIPv4 fetches IPv6 address from ifconfig.io
func (i *IfconfigIoFetcher) FetchIPv6(ctx context.Context) (net.IP, error) {
	ipBytes, err := getByIPv6(ctx, ifconfigIoURL)
	if err != nil {
		return nil, err
	}

	ip := net.ParseIP(string(ipBytes))
	if ip == nil {
		return nil, fmt.Errorf("fetched string cannot be parsed as IP address. string: %s")
	}
	return ip, nil
}
