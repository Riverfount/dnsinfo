package output

import (
	"fmt"
	"strings"

	"github.com/Riverfount/dnsinfo/internal/geoip"
)

type Result struct {
	geoip.Info
	Hostnames []string `json:"hostnames"`
	Currency  string   `json:"currency"`
}

func Print(result Result) {
	fmt.Printf("IP: %s\n", result.IP)
	fmt.Printf("Hostnames: %s\n", strings.Join(result.Hostnames, ", "))
	fmt.Printf("Organization: %s\n", result.Org)
	fmt.Printf("Anycast: %t\n", result.Anycast)
	fmt.Printf("City: %s\n", result.City)
	fmt.Printf("Region: %s\n", result.Region)
	fmt.Printf("Postal: %s\n", result.Postal)
	fmt.Printf("Country: %s\n", result.Country)
	fmt.Printf("Currency: %s\n", result.Currency)
	fmt.Printf("Location: %s\n", result.Loc)
	fmt.Printf("Timezone: %s\n", result.Timezone)
}
