# Nigerian Phone Number Intelligence API

Small Go HTTP API for Nigerian phone number normalization and original mobile
network allocation lookup.

## Purpose

Given a Nigerian phone number, the service:

1. Validates and normalizes it to E.164 form (`+234...`).
2. Identifies the mobile operator originally allocated that number prefix by the
   Nigerian Communications Commission (NCC).

It is intended as a clear, testable portfolio API, not a live telecom lookup
service.

### What this service does not do

- It does not determine the current serving operator. Mobile number portability
  means a number may have moved away from its originally allocated network.
- It does not verify subscriber identity.
- It does not check whether a number is assigned, active, or able to receive
  calls or messages.
- Format validation only confirms that an input matches the accepted Nigerian
  number shapes. It is not proof that a subscriber exists.

## Requirements

- Go `1.27.2` (from `go.mod`)

## Install and run

```bash
git clone https://github.com/eselikoghene/ng-phone-intelligence.git
cd ng-phone-intelligence
go run .
```

The server listens on `http://localhost:8080`.

Stop the process with `Ctrl+C` (SIGINT) or SIGTERM. The server shuts down
gracefully with a short timeout.

There are no third-party runtime dependencies beyond the Go standard library.

## Endpoints

| Method | Path | Description |
| --- | --- | --- |
| `GET` | `/` | Plain-text service banner |
| `GET` | `/health` | Health check JSON |
| `GET` | `/api/v1/phone` | Normalize a number and identify its original network |

Unsupported methods on `/api/v1/phone` and `/health` return HTTP `405` with
`Allow: GET` and the structured error envelope. Extra query parameters on
`/api/v1/phone` are ignored; only `number` is used.

### `GET /`

Returns plain text confirming the API is running.

### `GET /health`

```json
{
  "status": "ok"
}
```

### `GET /api/v1/phone`

Query parameter:

| Parameter | Required | Description |
| --- | --- | --- |
| `number` | Yes | Nigerian phone number to normalize and look up |

#### Accepted formats

- Local: `08031234567` (11 digits, leading `0`)
- International with plus: `+2348031234567` (14 characters)
- Leading and trailing whitespace is trimmed (spaces, tabs, newlines)

#### Rejected formats

- Empty or whitespace-only input
- Bare international form without `+` (for example `2348031234567`)
- Letters, hyphens, internal spaces, or other non-digit characters
- Unicode digits (only ASCII `0`-`9` are accepted)
- Wrong length (too short or too long)
- Non-Nigerian E.164 country codes (for example `+1...`)
- Incorrect combinations such as `+234` followed by a retained leading `0`

The API does not strip punctuation or rewrite arbitrary characters to force a
match. Inputs must already be in an accepted shape after trimming.

#### Successful response

HTTP `200`, `Content-Type: application/json`:

```json
{
  "phone_number": "+2348031234567",
  "network": "MTN",
  "valid": true
}
```

If the number format is valid but the prefix is not in the supported allocation
set, the API still returns HTTP `200` with `"network": "Unknown"`.

#### Error response

Client and method errors use this JSON envelope:

```json
{
  "error": {
    "code": "INVALID_PHONE_NUMBER",
    "message": "The phone number format is invalid."
  }
}
```

| Error code | HTTP status | When |
| --- | --- | --- |
| `MISSING_PHONE_NUMBER` | `400` | `number` query parameter is absent |
| `EMPTY_PHONE_NUMBER` | `400` | `number` is empty or only whitespace |
| `INVALID_PHONE_NUMBER` | `400` | format is not a supported Nigerian number |
| `UNSUPPORTED_COUNTRY_CODE` | `400` | E.164 country code is not `+234` |
| `METHOD_NOT_ALLOWED` | `405` | HTTP method is not `GET` |

#### Examples

Valid Nigerian number:

```bash
curl "http://localhost:8080/api/v1/phone?number=08031234567"
```

```json
{
  "phone_number": "+2348031234567",
  "network": "MTN",
  "valid": true
}
```

Invalid number:

```bash
curl "http://localhost:8080/api/v1/phone?number=12345"
```

```json
{
  "error": {
    "code": "INVALID_PHONE_NUMBER",
    "message": "The phone number format is invalid."
  }
}
```

Valid format, unknown network prefix:

```bash
curl "http://localhost:8080/api/v1/phone?number=08001234567"
```

```json
{
  "phone_number": "+2348001234567",
  "network": "Unknown",
  "valid": true
}
```

## Supported network allocations

Prefix mappings come from the NCC Mobile Number Allocation Table (allocations as
at December 2023; page last updated 25 November 2024), verified for this project
on 11 October 2026:

https://www.ncc.gov.ng/operators/mobile-number-allocation-table

Supported originally allocated operators and prefixes:

| Network | Prefixes |
| --- | --- |
| MTN | `0703`, `0704`, `0706`, `0707`, `0803`, `0806`, `0810`, `0813`, `0814`, `0816`, `0903`, `0906`, `0913`, `0916` |
| Airtel | `0701`, `0708`, `0802`, `0808`, `0812`, `0901`, `0902`, `0904`, `0907`, `0911`, `0912` |
| Glo | `0705`, `0805`, `0807`, `0811`, `0815`, `0905`, `0915` |
| 9mobile | `0809`, `0817`, `0818`, `0908`, `0909` |
| M-Tel | `0804` |
| Mafab | `08010xxxxxxx` only (partial NCC allocation) |

Unknown or unsupported prefixes return `"network": "Unknown"` without treating
the number as malformed.

Intentionally unsupported until a clear verified mapping exists:

- `0700`, `0800` special-service / VAS shared ranges
- `0900` vanity reservation
- `0702` mixed partial subranges
- `0709` withdrawn
- `0819` not verified as a current full mobile allocation
- `0801` outside the Mafab subrange `08010xxxxxxx`

## Development

```bash
go test ./...
go vet ./...
go test -race ./...
go test -bench=. -benchmem ./...
```

Local benchmarks measure single-process throughput on developer hardware. They
do not establish production capacity.
