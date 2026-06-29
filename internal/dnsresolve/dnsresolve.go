package dnsresolve

import (
	"context"
	"errors"
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

func ReverseHostname(ctx context.Context, resolver *net.Resolver, ip net.IP) []string {
	hostnames, err := resolver.LookupAddr(ctx, ip.To4().String())
	if err != nil {
		return nil
	}
	for i := range hostnames {
		hostnames[i] = strings.TrimSuffix(hostnames[i], ".")
	}
	return hostnames
}
