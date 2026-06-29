package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/Riverfount/dnsinfo/internal/dnsresolve"
	"github.com/Riverfount/dnsinfo/internal/geoip"
	"github.com/Riverfount/dnsinfo/internal/output"
)

const requestTimeout = 5 * time.Second

func main() {
	jsonFlag := flag.Bool("json", false, "output as JSON")
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "usage: %s [--json] <hostname>\n", os.Args[0])
		flag.PrintDefaults()
	}
	flag.Parse()
	if flag.NArg() < 1 {
		flag.Usage()
		os.Exit(1)
	}
	hostname := flag.Arg(0)
	host, err := dnsresolve.NormalizeHost(hostname)
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

	if *jsonFlag {
		if err := output.PrintJSON(result); err != nil {
			log.Fatal(err)
		}
	} else {
		output.Print(result)
	}

}
