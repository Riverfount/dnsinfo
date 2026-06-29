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

Output is available as human-readable text (default) or JSON, and can be written to the terminal or to a file.

Geolocation and network data come from the [ipinfo.io](https://ipinfo.io) API; reverse hostnames are resolved locally using Go's standard library.

> **Note on accuracy:** IP geolocation is approximate, especially for CDN, cloud, and anycast addresses. For example, Google's front-end IPs report "Mountain View" (the registration location), not where your packets actually land. Treat the geographic data as indicative, not exact.

## Requirements

- Go 1.22 or later
- No API key required — `dnsinfo` uses ipinfo.io's anonymous tier and resolves reverse hostnames locally.

## Installation

```bash
git clone https://github.com/Riverfount/dnsinfo.git
cd dnsinfo
go build -o dnsinfo .
```

Or run directly without building:

```bash
go run . www.google.com
```

## Usage

```
usage: dnsinfo [--json] [-o file] <hostname>
  -json
        output as JSON
  -o string
        write output to a file instead of stdout
```

The hostname can be a bare host, a host with a port, or a full URL — the host is extracted in all cases. Flags must come before the hostname.

Examples:

```bash
dnsinfo www.google.com                    # text output to the terminal
dnsinfo --json www.google.com             # JSON output to the terminal
dnsinfo -o result.txt www.google.com      # text output to a file
dnsinfo --json -o result.json dns.google  # JSON output to a file
```

## Example output

Text (default):

```
$ dnsinfo dns.google
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

JSON (`--json`):

```json
{
  "ip": "8.8.8.8",
  "city": "Mountain View",
  "region": "California",
  "country": "US",
  "loc": "37.4056,-122.0775",
  "org": "AS15169 Google LLC",
  "postal": "94043",
  "timezone": "America/Los_Angeles",
  "anycast": true,
  "hostnames": [
    "dns.google"
  ],
  "currency": "USD"
}
```

Hosts whose IPs have no reverse DNS record show an empty `Hostnames` field (or an empty `hostnames` array in JSON).

## Project structure

```
dnsinfo/
├── main.go                     flag parsing and orchestration
└── internal/
    ├── dnsresolve/             hostname normalization and DNS resolution
    ├── geoip/                  ipinfo.io lookup and currency mapping
    └── output/                 text and JSON formatting
```

## Tests

```bash
go test ./...
```

## License

MIT