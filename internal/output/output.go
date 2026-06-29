package output

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/Riverfount/dnsinfo/internal/geoip"
)

type Result struct {
	geoip.Info
	Hostnames []string `json:"hostnames"`
	Currency  string   `json:"currency"`
}

func Print(dst io.Writer, result Result) {
	fmt.Fprintf(dst, "IP: %s\n", result.IP)
	fmt.Fprintf(dst, "Hostnames: %s\n", strings.Join(result.Hostnames, ", "))
	fmt.Fprintf(dst, "Organization: %s\n", result.Org)
	fmt.Fprintf(dst, "Anycast: %t\n", result.Anycast)
	fmt.Fprintf(dst, "City: %s\n", result.City)
	fmt.Fprintf(dst, "Region: %s\n", result.Region)
	fmt.Fprintf(dst, "Postal: %s\n", result.Postal)
	fmt.Fprintf(dst, "Country: %s\n", result.Country)
	fmt.Fprintf(dst, "Currency: %s\n", result.Currency)
	fmt.Fprintf(dst, "Location: %s\n", result.Loc)
	fmt.Fprintf(dst, "Timezone: %s\n", result.Timezone)
}

func PrintJSON(dst io.Writer, result Result) error {
	out, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		return err
	}
	fmt.Fprintln(dst, string(out))
	return nil
}
