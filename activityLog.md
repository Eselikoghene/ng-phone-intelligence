# Activity Log

## 2026-10-11

Feature or change: Stages 2–4 portfolio expansion

Short summary: Expanded phone package table-driven tests; introduced structured JSON API error responses; verified and documented NCC mobile prefix allocations, including careful partial handling for MAFAB 0801.

Reason for change: Continue portfolio development with stronger tests, a stable client error contract, and allocation coverage grounded in the official NCC table.

## 2026-10-12

Feature or change: README accuracy review

Short summary: Reworked README.md to match the current API: Go 1.27.2 requirement, install/run steps, all endpoints, curl examples, error codes with HTTP 400, input policy, supported/unsupported prefixes, and explicit non-claims around portability and subscriber status.

Reason for change: Keep external documentation accurate, complete, and aligned with the implemented service without inventing features.

## 2026-10-12

Feature or change: Reliability audit and hardening

Short summary: Audited validation, prefix lookup, and HTTP contract; added GET-only method checks, safer JSON encoding, HTTP server timeouts with graceful shutdown, unicode/edge-case tests, concurrency coverage, and local benchmarks. README updated for METHOD_NOT_ALLOWED and operational limits.

Reason for change: Improve real-world correctness and operability for a small portfolio API without expanding scope or inventing allocations.
