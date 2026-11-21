# Feature inventory

Status is based on source code and SQL currently present in the public `main` branch when this snapshot was prepared. “Implemented” means code paths exist; it does not mean they are production-ready or fully tested.

Related (read only when relevant):
- `docs/DASHBOARD.md` — dashboard fragments, calculation rules, and account-group semantics.
- `docs/RISKS.md` — known correctness and reliability risks.
- `docs/todo/foreign-currency-cc.md` — open design for foreign-currency credit-card imports.

## Implemented code paths

| Area | Observed implementation | Notes |
|---|---|---|
| Accounts | Account listing; lookup by name/account number; account balance updates; receivable-account creation | Account CRUD is not fully exposed in the inspected handlers |
| Account groups | Create, list, read, update, delete | Standard CRUD handlers/service/repository |
| Categories | Create, list, read, update, delete | Standard CRUD handlers/service/repository |
| Transaction types | Create, list, read, update, delete | Standard CRUD handlers/service/repository |
| Transactions | List, create, update; delete is exposed through dashboard handler | Transaction service also groups daily totals and computes expense and income category aggregates |
| Dashboard | Overview cards (month income, expenses, net cash flow, balance blocks), account cards grouped by account group with adjust-balance and available credit, 6-month income vs expenses bars, expense breakdown with percentages and date filter, recent transactions, full transaction grid, create/edit/cancel transaction views, add-account modal | Server-rendered fragments loaded by `dashboard.html`. See `docs/DASHBOARD.md` |
| Gmail OAuth | Gmail read-only scope and local token/client-secret file loading | Single local credential path; not a multi-account setup |
| Gmail sync | `POST /scan-email`; subject/date-based search; message sorting; checkpoint; processed-message records | Sync is initiated by an HTTP request |
| Gmail date parsing | `ParseDate` in `app/gmail/parse.go` maps Indonesian month abbreviations (`Agt`/`Agu`, `Okt`, `Des`, `Mei`) to English before `time.Parse` | Covered by `app/gmail/parse_test.go`. Full month names are not handled |
| BCA parsing | Internet Transaction Journal cases include transfer, QRIS, credit-card/paylater transfer, Flazz top-up, pocket transfers, | QRIS payments pick the account from `Source of Fund` (`BCA KKI` → `CC BCA`, `TAHAPAN - 0840` → `Tahapan Gold`, missing → `Tahapan Gold`); a non-empty unrecognized value fails the message. Other cases still use fixed account/category names. Mapping covered by `app/gmail/bca_test.go`. Pocket emails (`Transaction Type` = `Pocket Transaction`) map to `Transfer` by layout: `Beneficiary Pocket` + `Source of Fund` → `Tahapan Gold` → `MyBCAPocket <name>`; `Beneficiary Pocket` + `Source Pocket` → pocket → pocket; `Source Pocket` + `Beneficiary Account` → `MyBCAPocket <name>` → `Tahapan Gold` (always, regardless of beneficiary). An unrecognized layout fails the message. Logic is in `app/gmail/bca_pocket.go`, not yet covered by tests. |
| blu parsing | blu transaction notification; BI-FAST, QRIS, Flazz top-up cases | Destination account-number lookup is used to distinguish known accounts from outside recipients |
| BCA credit-card parsing | Credit-card transaction notification | Maps as an expense against the configured `CC BCA` account. Emails whose `Sejumlah` currency is `MYR` are skipped (not imported); see `docs/todo/foreign-currency-cc.md` |
| Gmail processing records | Success/failure records keyed by Gmail message ID | Failed records are stored for diagnosis. Intentionally skipped emails (`ErrSkipEmail`) are stored as successes with transaction ID `0` |

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