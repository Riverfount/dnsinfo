package geoip

import "testing"

func TestCurrencyForCountry(t *testing.T) {
	tests := []struct {
		name    string
		country string
		want    string
	}{
		{
			name:    "returns the currency for the country",
			country: "US",
			want:    "USD",
		},
		{
			name:    "returns the currency for the country that share the currency",
			country: "DE",
			want:    "EUR",
		},
		{
			name:    "returns N/A when the country does not exist",
			country: "XX",
			want:    "N/A",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CurrencyForCountry(tt.country); got != tt.want {
				t.Errorf("currencyForCountry() = %v, want %v", got, tt.want)
			}
		})
	}
}
