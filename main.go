package main

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/url"
	"os"
	"strings"
	"time"
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

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "usage: %s <hostname>\n", os.Args[0])
		os.Exit(1)
	}
	host, err := normalizeHost(os.Args[1])
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	var resolver net.Resolver
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	ipAddrs, err := resolver.LookupIPAddr(ctx, host)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	var ip4 net.IP
	for _, ip := range ipAddrs {
		if ip.IP.To4() != nil {
			ip4 = ip.IP
			break
		}
	}
	if ip4 == nil {
		fmt.Fprintln(os.Stderr, "ip4 not found")
		os.Exit(1)
	}
	fmt.Println(ip4)
}
