package main

import (
	"fmt"
	"net/url"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprint(os.Stderr, "Need to be informed the hostname!")
		os.Exit(1)
	}
	urlParsed, err := url.Parse(os.Args[1])
	if err != nil {
		fmt.Fprint(os.Stderr, "Error parse URL")
		os.Exit(1)
	}
	fmt.Println(urlParsed.Hostname())
}
