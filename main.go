package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/Riverfount/dnsinfo/internal/dnsresolve"
	"github.com/Riverfount/dnsinfo/internal/geoip"
	"github.com/Riverfount/dnsinfo/internal/output"
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

	result := output.Result{Info: info, Hostnames: hostnames, Currency: currency}

	output.Print(result)
}
