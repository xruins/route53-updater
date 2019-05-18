# Overview

`route53-updater` is a CLI tool to notify your IP addressed to Amazon Route53, implemented by golang.

# Prerequisites

You have to AWS credential has previleges to access Amazon Route53.
1) place your credential to default location; on `~/.aws/credential`
2) set environment variable. both of  `AWS_ACCESS_KEY_ID` and `AWS_SECRET_ACCESS_KEY`, or `AWS_SESSION_TOKEN`

cf. [Configuring the AWS CLI \- AWS Command Line Interface](https://docs.aws.amazon.com/cli/latest/userguide/cli-chap-configure.html)

# Usage

```
Usage:
  route53-updater [OPTIONS]

Application Options:
  -d, --domain=         domain name to notify.
  -4, --ipv4=           IPV4 address to notify. specify 'omit' to skip notification of IPv4 address.
  -6, --ipv6=           IPV6 address to notify. specify 'omit' to skip notification of IPv4 address.
  -z, --hosted-zone-id= HostedZoneID of Route53.
  -t, --ttl=            time to live in second for DNS records.

Help Options:
  -h, --help            Show this help message
```
