package ipfetcher

import (
	"context"
	"io/ioutil"
	"net"
	"net/http"
	"time"

	"github.com/hashicorp/go-retryablehttp"
)

// IPFetcher is the interface to define fetch remote ipv4/v6 addresses
type IPFetcher interface {
	FetchIPv4(ctx context.Context) (net.IP, error)
	FetchIPv6(ctx context.Context) (net.IP, error)
}

// getByIPv4 returns response of GET request with IPV4.
// it will be retried several times when request fails.
func getByIPv4(ctx context.Context, url string) ([]byte, error) {
	return getWithExponentialBackoff(ctx, url, "tcp4")
}

// getByIPv6 returns response of GET request with IPV6.
// it will be retried several times when request fails.
func getByIPv6(ctx context.Context, url string) ([]byte, error) {
	return getWithExponentialBackoff(ctx, url, "tcp6")
}

func getWithExponentialBackoff(ctx context.Context, url string, proto string) ([]byte, error) {
	client := &http.Client{}
	client.Transport = &http.Transport{
		Dial: (func(network, addr string) (net.Conn, error) {
			return (&net.Dialer{
				Timeout:   3 * time.Second,
				LocalAddr: nil,
				DualStack: false,
			}).Dial(proto, addr)
		}),
	}

	rtClient := retryablehttp.NewClient()
	rtClient.HTTPClient = client

	resp, err := retryablehttp.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}
