package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
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
	ip4, err := firstIPv4(ipAddrs)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	urlAddr := fmt.Sprintf("https://ipinfo.io/%s/json", ip4.To4().String())
	client := &http.Client{Timeout: 5 * time.Second}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, urlAddr, nil)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer func() { _ = resp.Body.Close() }()
	if resp.StatusCode != http.StatusOK {
		fmt.Fprintf(os.Stderr, "request error: %s\n", resp.Status)
		os.Exit(1)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	fmt.Println(string(body))
}
