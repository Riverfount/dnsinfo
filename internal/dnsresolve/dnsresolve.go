package dnsresolve

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/url"
	"strings"
)

func NormalizeHost(urlRaw string) (string, error) {
	if !strings.Contains(urlRaw, "://") {
		urlRaw = "//" + urlRaw
	}
	urlParsed, err := url.Parse(urlRaw)
	if err != nil {
		return "", errors.New("error parse url")
	}
	hostName := urlParsed.Hostname()
	if hostName == "" {
		return "", errors.New("hostname is empty")
	}
	return hostName, nil
}

func FirstIPv4(ipAddrs []net.IPAddr) (net.IP, error) {
	var ip4 net.IP
	for _, ip := range ipAddrs {
		if ip.IP.To4() != nil {
			ip4 = ip.IP
			break
		}
	}
	if ip4 == nil {
		return nil, errors.New("no IPv4 address found")
	}
	return ip4, nil
}

func ReverseHostname(ctx context.Context, ip net.IP) []string {
	resolver := &net.Resolver{}
	hostnames, err := resolver.LookupAddr(ctx, ip.To4().String())
	if err != nil {
		return nil
	}
	for i := range hostnames {
		hostnames[i] = strings.TrimSuffix(hostnames[i], ".")
	}
	return hostnames
}

func ResolveIPv4(ctx context.Context, host string) (net.IP, error) {
	resolver := &net.Resolver{}
	ipAddrs, err := resolver.LookupIPAddr(ctx, host)
	if err != nil {
		return nil, fmt.Errorf("resolving %s: %w", host, err)
	}

	ip4, err := FirstIPv4(ipAddrs)
	if err != nil {
		return nil, fmt.Errorf("%s has no IPv4 address: %w", host, err)
	}

	return ip4, nil
}
