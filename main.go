package main

import (
	"context"
	"flag"
	"fmt"
	"io"
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
	outFile := flag.String("o", "", "write output to a file instead of stdout")
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "usage: %s [--json] [-o file] <hostname>\n", os.Args[0])
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

	var dst io.Writer = os.Stdout
	if *outFile != "" {
		f, err := os.Create(*outFile)
		if err != nil {
			log.Fatal(err)
		}
		defer func() { _ = f.Close() }()
		dst = f
	}
	if *jsonFlag {
		if err := output.PrintJSON(dst, result); err != nil {
			log.Fatal(err)
		}
	} else {
		output.Print(dst, result)
	}

}
