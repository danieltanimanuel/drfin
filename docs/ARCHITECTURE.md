# Architecture (current repository)

## Runtime shape

drfin is currently a single Go HTTP application. `main.go` opens the log file, configures `slog`, opens PostgreSQL through `app/db`, constructs the API server, and starts it on `localhost:8080`.

`app/api/api.go` creates a Gorilla Mux router and wires repositories, services, and handlers. The code is organized by business area, with a repeated handler → service → repository pattern.

## Main packages

| Package | Responsibility |
|---|---|
| `app/api` | Router, middleware, dependency wiring, server startup |
| `app/db` | PostgreSQL connection |
| `app/accounts` | Account data and account balance operations |
| `app/accountgroups` | Account-group CRUD |
| `app/categories` | Category CRUD |
| `app/transactiontypes` | Transaction-type CRUD |
| `app/transactions` | Transaction CRUD, aggregation, and balance effects |
| `app/dashboard` | Dashboard HTTP handlers and view models |
| `app/shared/renderer` | Server-side HTML template rendering and formatting |
| `app/shared/logdr` | Shared logging helpers |
| `app/gmail` | OAuth, Gmail message sync, HTML extraction, bank-specific parsing |
| `app/lastemail` | Last-sync checkpoint |
| `app/processedgmail` | Per-message success/failure tracking |

## Request and data flow

### Web requests

1. `main.go` initializes logging and the database.
2. `api.NewAPIServer` creates the HTTP server.
3. `APIServer.Run` constructs repositories, services, and handlers.
4. Handlers call services; services call repositories for PostgreSQL operations.
5. Dashboard handlers render templates from `app/dashboard/templates`.

### Gmail import

1. `POST /scan-email` calls `gmail.Service.syncMail`.
2. Gmail OAuth is loaded using the local token and client-secret file paths defined in `app/gmail/auth.go`.
3. The service queries Gmail using supported subject lines and a stored internal-date checkpoint.
4. Messages are fetched, sorted by Gmail internal date, and checked against the processed-message table.
5. `extractor.go` extracts relevant HTML fields; `parse.go` normalizes dates and amounts.
6. Subject-specific handlers in `bca.go`, `blu.go`, and `credit_card_bca.go` map extracted data into transaction DTOs.
7. The transaction service persists the transaction and applies account balance changes.
8. The Gmail service records success/failure and advances the checkpoint after successful processing.

## Current design constraints

- Keep this a simple monolith. Do not introduce microservices, queues, dependency-injection frameworks, or a new ORM without a concrete requirement.
- Follow the existing package boundaries unless a task specifically addresses architecture.
- Keep SQL in the existing repository layer.
- Prefer `context.Context` on request-scoped operations and use bounded timeouts for external/database work.
- Keep HTML parsing and bank-specific rules within the Gmail area; avoid mixing them into HTTP handlers.
- Be careful: some existing methods do not consistently pass context to database calls. Do not expand that inconsistency.

## Important implementation caveats

- `app/db/db.go` currently hard-codes a local PostgreSQL connection string.
- `main.go` opens `app.log`; `app.yaml` declares App Engine runtime `go122`, while `go.mod` declares Go 1.23.0 and toolchain 1.24.2. Reconcile deployment versions before deployment changes.
- `api.go` contains request-ID middleware, but it stores the value under a string context key and logs with the standard logger. Avoid changing this incidentally.
- The renderer parses templates on each render; do not introduce caching or a template framework as part of unrelated work.
- There is currently no evidence in the inspected tree of a frontend framework or a separate frontend build pipeline.
