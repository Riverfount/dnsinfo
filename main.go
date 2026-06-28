package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
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

type IPInfo struct {
	IP       string `json:"ip"`
	City     string `json:"city"`
	Region   string `json:"region"`
	Country  string `json:"country"`
	Loc      string `json:"loc"`
	Org      string `json:"org"`
	Postal   string `json:"postal"`
	Timezone string `json:"timezone"`
	Anycast  bool   `json:"anycast"`
}

var currencyByCountry = map[string]string{
	"US": "USD",
	"BR": "BRL",
	"GB": "GBP",
	"CA": "CAD",
	"AU": "AUD",
	"NZ": "NZD",
	"JP": "JPY",
	"CN": "CNY",
	"IN": "INR",
	"MX": "MXN",
	"AR": "ARS",
	"CL": "CLP",
	"CO": "COP",
	"PE": "PEN",
	"UY": "UYU",
	"PY": "PYG",
	"DE": "EUR",
	"FR": "EUR",
	"IT": "EUR",
	"ES": "EUR",
	"PT": "EUR",
	"NL": "EUR",
	"CH": "CHF",
	"SE": "SEK",
	"NO": "NOK",
	"DK": "DKK",
	"PL": "PLN",
	"TR": "TRY",
	"RU": "RUB",
	"ZA": "ZAR",
	"KR": "KRW",
	"SG": "SGD",
	"HK": "HKD",
	"AE": "AED",
	"SA": "SAR",
	"IL": "ILS",
	"TH": "THB",
	"VN": "VND",
	"ID": "IDR",
	"PH": "PHP",
}

const (
	requestTimeout = 5 * time.Second
	ipInfoBaseURL  = "https://ipinfo.io/%s/json"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatalf("usage: %s <hostname>", os.Args[0])
	}
	host, err := normalizeHost(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	var resolver net.Resolver
	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	defer cancel()
	ipAddrs, err := resolver.LookupIPAddr(ctx, host)
	if err != nil {
		log.Fatal(err)
	}
	ip4, err := firstIPv4(ipAddrs)
	if err != nil {
		log.Fatal(err)
	}
	info, err := fetchIPInfo(ctx, ip4)
	if err != nil {
		log.Fatal(err)
	}

	currency := currencyForCountry(info.Country)

	fmt.Printf("IP: %s\n", info.IP)
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

func currencyForCountry(country string) string {
	currency, ok := currencyByCountry[country]
	if !ok {
		currency = "N/A"
	}
	return currency
}

func fetchIPInfo(ctx context.Context, ip4 net.IP) (IPInfo, error) {
	urlAddr := fmt.Sprintf(ipInfoBaseURL, ip4.To4().String())

	client := &http.Client{Timeout: requestTimeout}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, urlAddr, nil)
	if err != nil {
		return IPInfo{}, err
	}
	resp, err := client.Do(req)
	if err != nil {
		return IPInfo{}, err
	}
	defer func() { _ = resp.Body.Close() }()
	if resp.StatusCode != http.StatusOK {
		return IPInfo{}, fmt.Errorf("request failed with status %s", resp.Status)
	}
	var info IPInfo
	if err := json.NewDecoder(resp.Body).Decode(&info); err != nil {
		return IPInfo{}, err
	}
	return info, nil
}
