package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/jessevdk/go-flags"
	"gopkg.in/go-playground/validator.v9"

	"github.com/xruins/route53-updater/ipfetcher"
	"github.com/xruins/route53-updater/route53"
)

type cmdOptions struct {
	Domain             string `short:"d" long:"domain" description:"domain name to notify." validator:"required,fqdn"`
	IPV4Address        string `short:"4" long:"ipv4" description:"IPV4 address to notify. specify 'omit' to skip notification of IPv4 address." validator:"ip"`
	IPV6Address        string `short:"6" long:"ipv6" description:"IPV6 address to notify. specify 'omit' to skip notification of IPv4 address." validator:"ip"`
	DisableAutoFetchV4 bool   `long:"disable-auto-fetch-v4" description:"If true, notification for IPv4 address will not be execute without --ipv4 flag."`
	DisableAutoFetchV6 bool   `long:"disable-auto-fetch-v6" description:"If true, notification for IPv4 address will not be execute without --ipv6 flag."`
	HostedZoneID       string `short:"h" long:"hosted-zone-id" description:"HostedZoneID of Route53." validator:"required,alphanum"`
	TTL                int32  `short:"t" long:"ttl" description:"time to live in second for DNS records." validator:"min=0"`
}

func main() {
	var options cmdOptions
	var parser = flags.NewParser(&options, flags.Default)
	if _, err := parser.Parse(); err != nil {
		if flagsErr, ok := err.(*flags.Error); ok && flagsErr.Type == flags.ErrHelp {
			os.Exit(0)
		} else {
			os.Exit(1)
		}
	}

	validate := validator.New()
	err := validate.Struct(options)
	if err != nil {
		fmt.Fprintf(os.Stderr, "validation of command-line options failed. err: %s", err)
		os.Exit(1)
	}

	if options.IPV4Address == "" && options.DisableAutoFetchV4 {
		fmt.Fprintf(os.Stderr, "--ipv4 option is required when --disable-auto-fetch-v4 option is true.")
		os.Exit(1)
	}

	if options.IPV6Address == "" && options.DisableAutoFetchV6 {
		fmt.Fprintf(os.Stderr, "--ipv6 option is required when --disable-auto-fetch-v6 option is true.")
		os.Exit(1)
	}

	ctx := context.Background()

	var ipf ipfetcher.IPFetcher
	ipf = &ipfetcher.IfconfigIoFetcher{}

	ipv4, err := ipf.FetchIPv4(ctx)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to fetch IPv4 address. err: %s", err)
	}
	if ipv4 == nil {
		fmt.Fprintf(os.Stderr, "fetched malformed IPv4 address. address: %s", ipv4)
	}

	ipv6, err := ipf.FetchIPv6(ctx)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to fetch IPv6 address. err: %s", err)
	}
	if ipv6 == nil {
		fmt.Fprintf(os.Stderr, "fetched malformed IPv6 address. address: %s", ipv6)
	}

	r53 := &route53.Route53{
		Domain:       options.Domain,
		HostedZoneID: options.HostedZoneID,
		TTL:          time.Duration(options.TTL) * time.Second,
	}
	err = r53.Notify(ctx, &ipv4, &ipv6)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to notify DNS record to Route53. err: %s", err)
		os.Exit(1)
	}
	fmt.Printf("succeed to notify DNS record to Route53. domain: %s, hostedZoneID: %s", options.Domain, options.HostedZoneID)
}
