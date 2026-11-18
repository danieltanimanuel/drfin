# Development notes

## Prerequisites

- Go version compatible with `go.mod` (Go 1.23.0; toolchain directive 1.24.2).
- PostgreSQL.
- Gmail OAuth client credentials only if testing Gmail synchronization (see `docs/GMAIL_SETUP.md`).

## Run locally

The current `app/db/db.go` uses this hard-coded DSN:

```text
postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable
```

Create the database/schema expected by the SQL scripts, then run:

```bash
go run .
```

The server address is currently `localhost:8080`.

## Useful checks

```bash
gofmt -w <changed-go-files>
go test ./...
go vet ./...
```

Tests exist but coverage is minimal: `app/shared/logdr` (stack-frame helper) and `app/gmail` (`ParseDate`, `isSkippedCurrency`). Extraction, bank parsers, amount parsing, transaction logic, balance updates and `syncMail` have no tests yet. See the test table in `README.md`.

## Secrets and local files

Never commit:
- `token.json`
- `resources/drfin-secret.json`
- `.env` or `.env.*`
- local database dumps
- real bank email samples or exports containing personal data

The existing code reads OAuth files from fixed paths. For local development, keep those files untracked and provision them manually. A later task can move configuration to environment variables or a secret manager.

## Task workflow

1. Inspect `CLAUDE.md` and the relevant docs.
2. Inspect only relevant source files and their callers/tests.
3. For complex or financial changes, write a plan before editing.
4. Implement the smallest coherent change.
5. Add focused tests, especially for amount parsing, date parsing, transaction types, transfer balance effects, duplicate Gmail messages, and failure/retry behavior.
6. Run tests and vet; report anything not run.
7. Review `git diff` and ensure no secrets or unrelated changes are included.