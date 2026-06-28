package main

import (
	"fmt"
	"net/url"
	"os"
	"strings"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "Need to be informed the hostname!")
		os.Exit(1)
	}
	urlRaw := os.Args[1]
	if !strings.Contains(urlRaw, "//") {
		urlRaw = "//" + urlRaw
	}
	urlParsed, err := url.Parse(urlRaw)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error parse URL")
		os.Exit(1)
	}
	fmt.Println(urlParsed.Hostname())
}
