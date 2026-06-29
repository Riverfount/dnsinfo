# dnsinfo

A command-line tool that, given a hostname, resolves its IPv4 address and enriches it with geolocation and network data.

Built as a learning project to explore Go's standard library, HTTP clients, and testing.

## Features

Given a hostname (e.g. `www.google.com`), `dnsinfo` reports:

- **IP** — the resolved IPv4 address
- **Hostnames** — reverse DNS (PTR) names for the IP, if any
- **Organization** — the network operator (ASN + name)
- **Anycast** — whether the IP is an anycast address
- **City**, **Region**, **Postal**, **Country** — geolocation
- **Currency** — derived from the country code
- **Location** — latitude and longitude
- **Timezone** — the IANA timezone

Geolocation and network data come from the [ipinfo.io](https://ipinfo.io) API; reverse hostnames are resolved locally using Go's standard library.

> **Note on accuracy:** IP geolocation is approximate, especially for CDN, cloud, and anycast addresses. For example, Google's front-end IPs report "Mountain View" (the registration location), not where your packets actually land. Treat the geographic data as indicative, not exact.

## Requirements

- Go 1.22 or later
- No API key required — `dnsinfo` uses ipinfo.io's anonymous tier and resolves reverse hostnames locally.

## Usage

Clone and run directly:

```bash
git clone https://github.com/Riverfount/dnsinfo.git
cd dnsinfo
go run . www.google.com
```

Or build a binary:

```bash
go build -o dnsinfo .
./dnsinfo www.google.com
```

The tool accepts a bare hostname, a hostname with a port, or a full URL — it extracts the host in all cases:

```bash
go run . www.google.com
go run . https://www.google.com/some/path
go run . dns.google
```

## Example output

```
$ go run . dns.google
IP: 8.8.8.8
Hostnames: dns.google
Organization: AS15169 Google LLC
Anycast: true
City: Mountain View
Region: California
Postal: 94043
Country: US
Currency: USD
Location: 37.4056,-122.0775
Timezone: America/Los_Angeles
```

Hosts whose IPs have no reverse DNS record simply show an empty `Hostnames` field.

## Project structure

```
dnsinfo/
├── main.go                     orchestration and output
└── internal/
    ├── dnsresolve/             hostname normalization and DNS resolution
    └── geoip/                  ipinfo.io lookup and currency mapping
```

## Tests

```bash
go test ./...
```

## License

MIT