package main

import (
	"net"
	"testing"
)

func TestFirstIPv4(t *testing.T) {
	tests := []struct {
		name    string
		ips     []net.IPAddr
		want    net.IP
		wantErr bool
	}{
		{
			name: "finds first ipv4",
			ips: []net.IPAddr{
				{IP: net.ParseIP("2606:4700:4700::1111")},
				{IP: net.ParseIP("2001:4860:4860::8888")},
				{IP: net.ParseIP("::ffff:192.0.2.128")},
				{IP: net.ParseIP("192.168.0.1")},
				{IP: net.ParseIP("10.0.0.5")},
				{IP: net.ParseIP("172.16.254.3")},
				{IP: net.ParseIP("8.8.8.8")},
				{IP: net.ParseIP("203.0.113.42")},
				{IP: net.ParseIP("2001:db8::1")},
				{IP: net.ParseIP("fe80::1")},
			},
			want:    net.ParseIP("192.0.2.128"),
			wantErr: false,
		},
		{
			name: "only ipv6",
			ips: []net.IPAddr{
				{IP: net.ParseIP("2001:db8::1")},
				{IP: net.ParseIP("2001:db8:1::abcd")},
				{IP: net.ParseIP("2001:db8:2::dead:beef")},
				{IP: net.ParseIP("fe80::1")},
				{IP: net.ParseIP("fe80::a6db:30ff:fe98:e946")},
				{IP: net.ParseIP("2001:4860:4860::8888")},
				{IP: net.ParseIP("2606:4700:4700::1111")},
				{IP: net.ParseIP("2001:db8:85a3::8a2e:370:7334")},
				{IP: net.ParseIP("::1")},
				{IP: net.ParseIP("fd12:3456:789a::1")},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "empty ip slice",
			ips:     []net.IPAddr{},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := firstIPv4(tt.ips)
			if (err != nil) != tt.wantErr {
				t.Fatalf("firstIPv4() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}
			if !got.Equal(tt.want) {
				t.Errorf("firstIPv4() = %v, want %v", got, tt.want)
			}
		})
	}
}

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
			if got := currencyForCountry(tt.country); got != tt.want {
				t.Errorf("currencyForCountry() = %v, want %v", got, tt.want)
			}
		})
	}
}
