package main

import (
	"context"
	"fmt"
	"net"
	"os"
	"time"

	"github.com/jessevdk/go-flags"
	"gopkg.in/go-playground/validator.v9"

	"github.com/xruins/route53-updater/route53"
)

type cmdOptions struct {
	Domain       string `short:"d" long:"domain" description:"domain name to notify." validator:"required,fqdn"`
	IPV4Address  string `short:"4" long:"ipv4" description:"IPV4 address to notify." validator:"ipv4"`
	IPV6Address  string `short:"6" long:"ipv6" description:"IPV6 address to notify." validator:"ipv6"`
	HostedZoneID string `short:"z" long:"hosted-zone-id" description:"HostedZoneID of Route53." validator:"required,alphanum"`
	TTL          int32  `short:"t" long:"ttl" description:"time to live in second for DNS records." validator:"min=0"`
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
		fmt.Fprintf(os.Stderr, "validation of command-line options failed. err: %s\n", err)
		os.Exit(1)
	}

	if options.IPV4Address == "" && options.IPV6Address == "" {
		fmt.Fprintln(os.Stderr, "either --ipv4 or --ipv6 required.")
		os.Exit(1)
	}

	ipv4 := net.ParseIP(options.IPV4Address)
	ipv6 := net.ParseIP(options.IPV6Address)

	ctx := context.Background()

	r53 := &route53.Route53{
		Domain:       options.Domain,
		HostedZoneID: options.HostedZoneID,
		TTL:          time.Duration(options.TTL) * time.Second,
	}
	err = r53.Notify(ctx, ipv4, ipv6)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to notify DNS record to Route53. err: %s\n", err)
		os.Exit(1)
	}
	fmt.Printf("succeed to notify DNS record to Route53. domain: %s, hostedZoneID: %s\n", options.Domain, options.HostedZoneID)
}
