package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/Riverfount/dnsinfo/internal/dnsresolve"
	"github.com/Riverfount/dnsinfo/internal/geoip"
)

const requestTimeout = 5 * time.Second

func main() {
	if len(os.Args) < 2 {
		log.Fatalf("usage: %s <hostname>", os.Args[0])
	}

	host, err := dnsresolve.NormalizeHost(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	defer cancel()

	ip4, err := dnsresolve.ResolveIPv4(ctx, host)
	if err != nil {
		log.Fatal(err)
	}

	info, err := geoip.Fetch(ctx, ip4)
	if err != nil {
		log.Fatal(err)
	}

	hostnames := dnsresolve.ReverseHostname(ctx, ip4)
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
