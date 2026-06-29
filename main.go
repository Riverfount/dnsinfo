package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"net/url"
	"os"
	"strings"

	"github.com/Riverfount/dnsinfo/internal/geoip"
)

func normalizeHost(urlRaw string) (string, error) {
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

func firstIPv4(ipAddrs []net.IPAddr) (net.IP, error) {
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

func main() {
	if len(os.Args) < 2 {
		log.Fatalf("usage: %s <hostname>", os.Args[0])
	}
	host, err := normalizeHost(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	resolver := &net.Resolver{}
	ctx, cancel := context.WithTimeout(context.Background(), geoip.RequestTimeout)
	defer cancel()
	ipAddrs, err := resolver.LookupIPAddr(ctx, host)
	if err != nil {
		log.Fatal(err)
	}
	ip4, err := firstIPv4(ipAddrs)
	if err != nil {
		log.Fatal(err)
	}
	info, err := geoip.Fetch(ctx, ip4)
	if err != nil {
		log.Fatal(err)
	}
	hostnames := reverseHostname(ctx, resolver, ip4)

	currency := geoip.CurrencyForCountry(info.Country)

	fmt.Printf("IP: %s\n", info.IP)
	fmt.Printf("Hostnames: %s\n", strings.Join(hostnames, ", "))
	fmt.Printf("Organization: %s\n", info.Org)
	fmt.Printf("Anycast: %t\n", info.Anycast)
	fmt.Printf("City: %s\n", info.City)
	fmt.Printf("Region: %s\n", info.Region)
	fmt.Printf("Postal: %s\n", info.Postal)
	fmt.Printf("Country: %s\n", info.Country)
	fmt.Printf("Currency: %s\n", currency)
	fmt.Printf("Location: %s\n", info.Loc)
	fmt.Printf("Timezone: %s\n", info.Timezone)
}

func reverseHostname(ctx context.Context, resolver *net.Resolver, ip net.IP) []string {
	hostnames, err := resolver.LookupAddr(ctx, ip.To4().String())
	if err != nil {
		return nil
	}
	for i := range hostnames {
		hostnames[i] = strings.TrimSuffix(hostnames[i], ".")
	}
	return hostnames
}
