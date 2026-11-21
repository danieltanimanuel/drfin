# drfin

A personal finance web application for tracking accounts and transactions, with Gmail-based import of supported bank transaction notifications.

## Current capabilities

- Manage account groups, accounts, transaction types, and categories.
- Create, list, update, and delete transactions through the existing handlers and services.
- Maintain account balances when transactions are created or updated.
- Render a web dashboard with account summaries, transaction history, and expense-category summaries.
- Trigger Gmail synchronization through `POST /scan-email`.
- Parse selected BCA, blu, and BCA credit-card notification formats.
- Post BCA QRIS payments to `CC BCA` or `Tahapan Gold` based on the email's `Source of Fund` field (empty falls back to `Tahapan Gold`).
- Import BCA MyBCA pocket transfers as `Transfer` transactions: `Tahapan Gold` → pocket, pocket → pocket, and pocket → `Tahapan Gold`. Pocket accounts must exist as `MyBCAPocket <pocket name>` (for example `MyBCAPocket Mam`).
- Record processed Gmail message IDs and failed processing attempts.
- Create accounts from the dashboard and adjust an account's balance (recorded as an `Adjustment` transaction).

See [docs/FEATURES.md](docs/FEATURES.md) for implementation status and limitations.

## Technology

- Go (module declares Go 1.23.0 and toolchain Go 1.24.2)
- PostgreSQL via `database/sql` and `github.com/lib/pq`
- HTTP routing with `github.com/gorilla/mux`
- Server-rendered HTML templates and static CSS
- Gmail API with Google OAuth 2.0 (`google.golang.org/api/gmail/v1`, `golang.org/x/oauth2`)
- HTML parsing with `github.com/PuerkitoBio/goquery`
- Deployment configuration currently includes Google App Engine (`app.yaml`) and dispatch configuration.

The exact dependencies are listed in `go.mod`.

## Repository layout

```text
.
├── app/
│   ├── accountgroups/   # Account group handlers, services, repository, model
│   ├── accounts/        # Account handlers, services, repository, model
│   ├── api/             # HTTP router and dependency wiring
│   ├── categories/      # Category CRUD
│   ├── dashboard/       # Dashboard handlers, view models, templates
│   ├── db/              # PostgreSQL connection setup
│   ├── gmail/           # Gmail OAuth, sync, extraction, and bank parsers
│   ├── lastemail/       # Gmail sync checkpoint persistence
│   ├── processedgmail/  # Processed/failed Gmail message tracking
│   ├── shared/          # Logging and HTML rendering helpers
│   ├── transactions/    # Transaction logic and persistence
│   ├── transactiontypes/# Transaction type CRUD
│   └── static/          # CSS
├── resources/
│   ├── SQL/             # Database creation/configuration and seed SQL
│   └── sample_html/     # Bank notification samples used by parsers
├── main.go              # Application entry point
├── go.mod
└── app.yaml
```

## Local development

The current database connection is hard-coded in `app/db/db.go` to a local PostgreSQL DSN. Configure a local PostgreSQL instance matching that DSN, or change the application configuration as a separate task before running elsewhere.

1. Install the Go version required by `go.mod`.
2. Create a PostgreSQL database and apply the SQL scripts in `resources/SQL/` in the appropriate order:
   - `initial_table.sql`
   - `initial_table_config.sql`
   - `populate_table.sql`
3. Provide Gmail OAuth client credentials locally (see the security note below).
4. Run:

```bash
go run .
```

The current entry point starts the HTTP server on `localhost:8080`.

## Security — read before running

**Do not commit OAuth client secrets or user OAuth tokens.** The repository has historically contained `token.json` and `resources/drfin-secret.json`. Treat any credentials that have ever been committed to this public repository as exposed: revoke/rotate them, remove the files, and consider cleaning Git history. Keep local credential files out of Git and deploy secrets through an appropriate secret store or environment-based configuration.

The current implementation expects:
- `token.json`
- `resources/drfin-secret.json`

Those paths are legacy implementation details, not a recommended secrets-management design. Never upload these files to an AI coding assistant.

## Known limitations

This is documentation of the current repository, not a claim that the future product roadmap is implemented. Gmail account-to-account configuration, reconciliation, robust transfer matching, and multi-Gmail-account support are not established by the current code inspected. See [docs/FEATURES.md](docs/FEATURES.md).

## Development checks

```bash
go test ./...
go vet ./...
```

Run a single package while iterating:

```bash
go test ./app/shared/logdr/ -v
go test ./app/gmail/ -v
```

### Current tests

| Package | File | Covers |
|---|---|---|
| `app/shared/logdr` | `logger_test.go` | `frames()` stack-trace helper: first frame is the caller, and `skip` moves it up one level |
| `app/gmail` | `parse_test.go` | `ParseDate` for blu emails: Indonesian month abbreviations (`Agu`, `Agt`, `Okt`, `Des`, `Mei`) and an English month |
| `app/gmail` | `credit_card_bca_test.go` | `isSkippedCurrency`: `MYR` detection (case, leading space), and that IDR, other currencies and empty input are not skipped |   
| `app/accounts` | `money_test.go` | `NormalizeAmount` (plain decimals, at most 2 decimals, no float conversion) and `validateCreateAccount` (name, group, initial balance, credit limit) |
| `app/gmail` | `bca_test.go` | `accountFromSourceOfFund`: `BCA KKI` → `CC BCA`, `TAHAPAN - 0840` → `Tahapan Gold`, empty/whitespace fallback to `Tahapan Gold`, and unrecognized values rejected |

Coverage is minimal. Only Gmail date parsing, the credit-card currency skip check, and the BCA QRIS source-of-fund account mapping are tested. Email extraction, bank parsers, amount parsing, transaction logic, balance updates and the `syncMail` skip handling have no tests yet, so add focused tests when changing them. The SQL paths for account creation and balance adjustment have no tests (no database test setup exists). The BCA pocket mapping (`buildPocketTransferDTO`) is also untested.

### Log output

`logdr.Info`, `Warn` and `Debug` add two fields:
- `caller`: `file:line` of the log call.
- `trace`: up to 5 frames, formatted as `file:line func`. `trace[0]` equals `caller`.

`logdr.Error` logs a full `debug.Stack()` instead of `trace`.