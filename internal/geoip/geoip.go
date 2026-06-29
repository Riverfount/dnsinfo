package geoip

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"time"
)

type Info struct {
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
	RequestTimeout = 5 * time.Second
	ipInfoBaseURL  = "https://ipinfo.io/%s/json"
)

func CurrencyForCountry(country string) string {
	currency, ok := currencyByCountry[country]
	if !ok {
		currency = "N/A"
	}
	return currency
}

func Fetch(ctx context.Context, ip4 net.IP) (Info, error) {
	urlAddr := fmt.Sprintf(ipInfoBaseURL, ip4.To4().String())
	client := &http.Client{Timeout: RequestTimeout}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, urlAddr, nil)
	if err != nil {
		return Info{}, err
	}
	resp, err := client.Do(req)
	if err != nil {
		return Info{}, err
	}
	defer func() { _ = resp.Body.Close() }()
	if resp.StatusCode != http.StatusOK {
		return Info{}, fmt.Errorf("request failed with status %s", resp.Status)
	}
	var info Info
	if err := json.NewDecoder(resp.Body).Decode(&info); err != nil {
		return Info{}, err
	}
	return info, nil
}
