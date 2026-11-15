# Feature inventory

Status is based on source code and SQL currently present in the public `main` branch when this snapshot was prepared. “Implemented” means code paths exist; it does not mean they are production-ready or fully tested.

## Implemented code paths

| Area | Observed implementation | Notes |
|---|---|---|
| Accounts | Account listing; lookup by name/account number; account balance updates; receivable-account creation | Account CRUD is not fully exposed in the inspected handlers |
| Account groups | Create, list, read, update, delete | Standard CRUD handlers/service/repository |
| Categories | Create, list, read, update, delete | Standard CRUD handlers/service/repository |
| Transaction types | Create, list, read, update, delete | Standard CRUD handlers/service/repository |
| Transactions | List, create, update; delete is exposed through dashboard handler | Transaction service also groups daily totals and computes expense category aggregates |
| Dashboard | Dashboard page, account grid, transaction grid, expense category summary, create/edit transaction views | Server-rendered templates |
| Gmail OAuth | Gmail read-only scope and local token/client-secret file loading | Single local credential path; not a multi-account setup |
| Gmail sync | `POST /scan-email`; subject/date-based search; message sorting; checkpoint; processed-message records | Sync is initiated by an HTTP request |
| BCA parsing | Internet Transaction Journal cases include transfer, QRIS, credit-card/paylater transfer, Flazz top-up | Current mapping uses fixed account/category names |
| blu parsing | blu transaction notification; BI-FAST, QRIS, Flazz top-up cases | Destination account-number lookup is used to distinguish known accounts from outside recipients |
| BCA credit-card parsing | Credit-card transaction notification | Maps as an expense against the configured `CC BCA` account |
| Gmail processing records | Success/failure records keyed by Gmail message ID | Failed records are stored for diagnosis |

## Data model concepts currently present

- Accounts have a name, account number, balance, account group, status, timestamps, and optional credit limit.
- Transactions have type, category, source account, optional destination account, amount, description, status, trade date, and timestamps.
- Transaction types seeded in SQL include Expense, Income, Transfer, Lend, Borrow, and Adjustment.
- Gmail processing has a last-email checkpoint and a processed-message table.

## Not established by the current code

Do not tell users or assume that these are implemented just because they are part of the desired product direction:

- Multiple Gmail accounts with independent OAuth credentials.
- User-configurable mapping of Gmail sources to accounts/pockets.
- A configurable parser registry or UI for managing bank parsers.
- Generalized transfer matching across multiple notification emails with robust cross-message deduplication.
- Manual income workflow beyond the generic transaction creation path.
- Balance reconciliation against externally entered actual balances.
- A transaction review queue/UI for failed or uncertain imports.
- Automatic expense categorization beyond the hard-coded mappings currently present.
- Full audit trail or ledger-based immutable accounting model.

## Known correctness and reliability risks to investigate

These are observations for future review, not permission to refactor everything at once:

1. Monetary values are represented using `*big.Float` in Go and `numeric(24,2)` in PostgreSQL. Review precision, rounding, and float conversions before changing financial calculations.
2. Transaction types are resolved by name, but some logic relies on numeric IDs (for example, Expense/Income/Transfer IDs). Confirm seed data and constraints before modifying.
3. `transactions.Service` applies/reverses account balances. Any change to transaction type, amount, source account, or destination account must be traced through both reversal and reapplication paths.
4. Gmail `syncMail` currently ignores the error returned by `listMessages` (`msgs, _ := ...`) before checking a different `err` variable. This deserves a focused bug fix with tests.
5. A Gmail import can create a transaction and then separately record the message as processed. Review failure windows and idempotency before changing sync behavior.
6. The Gmail processing repository uses Gmail message ID as the unique key; current schema does not include Gmail account identity in that key.
7. Some parser functions use string slicing and assumptions about exact email layouts. Test malformed and changed email formats.
8. Do not mark a message successful or advance the checkpoint if transaction persistence fails.

## TODO list in the existing README

The previous README mentions request/correlation IDs and additional bank/email parsing work. Treat those notes as a historical backlog, not proof that a feature is complete. Re-check code and sample fixtures before taking on each item.
