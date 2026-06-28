package main

import (
	"errors"
	"fmt"
	"net/url"
	"os"
	"strings"
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
		fmt.Fprintln(os.Stderr, "Need to be informed the hostname!")
		os.Exit(1)
	}
	host, err := normalizeHost(os.Args[1])
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	fmt.Println(host)
}
