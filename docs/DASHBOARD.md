# Dashboard

`/dashboard` renders only the shell (`dashboard.html`). JavaScript then loads each section as an HTML fragment into a placeholder element. All JavaScript stays in `dashboard.html` (see `docs/ARCHITECTURE.md`).

## Sections and fragments

| Section | Element ID | Route | Template | Data source |
|---|---|---|---|---|
| Financial Overview | `overview-grid` | `/dashboard/overview` | `overview_grid.html` | Month Income/Expense aggregates + `GetAllAccounts` |
| My Accounts | `accounts-grid` | `/dashboard/accounts` | `accounts_grid.html` | `GetAllAccounts`, grouped by account group |
| Income vs Expenses (6 months) | `cashflow-grid` | `/dashboard/cashflow` | `cashflow_grid.html` | Income/Expense aggregates, one pair of queries per month |
| Expense Breakdown | `categories-grid` | `/dashboard/categories?from=&to=` | `categories_grid.html` | `GetExpenseCategoryAggregates` |
| Recent Transactions | `recent-grid` | `/dashboard/recent` | `recent_grid.html` | `GetAllTransactions`, flattened (latest 8, cancelled excluded) |
| All Transactions | `transactions-grid` | `/dashboard/transactions` | `transactions_grid.html` | `GetAllTransactions` |
| Add Account modal | `modal-content` | `/dashboard/accounts/new` | `account_create.html` | Static form |

`/dashboard/pie` still exists but the dashboard no longer uses it.

## Owner selection

- `/dashboard?owner=<id>` selects the owner (default 1; unknown id → 404). The sidebar lists every row of `users`. The shell renders the id into `<body data-owner-id>` and the JavaScript adds `owner=<id>` to every fragment and create request.
- Every fragment route above accepts `?owner=<id>`. The transaction edit form takes its owner from the transaction being edited, not from the URL.
- Account groups are shared across owners. Overview blocks, month totals and breakdowns are per owner.

Handlers are in `app/dashboard/handlers.go` (existing) and `app/dashboard/overview_handlers.go` (new fragments). Owner parsing and the shell view are in `app/dashboard/owner.go`. View models and pure builders are in `app/dashboard/overview_view.go`. Styles: `app/static/dashboard.css` (loaded after `style.css` and `transaction-grid.css`).

## Calculation rules

- Month income/expenses reuse `GetCategoryAggregatesByType`: only types named `Income` / `Expense`, `status <> 'Cancelled'`. Transfers, Adjustments, Lend and Borrow are excluded. `Waiting` (email-imported) rows are counted, consistent with balances, which are applied at creation.
- Net cash flow = Income − Expense.
- The month's expense total equals the sum of the breakdown rows for the same date range.
- Sums use `big.Rat` built from 2-decimal strings (`moneyRat`). `big.Float` is used only to hand values to `formatNumber` for display. Percentages and bar heights are display-only.
- "This month" uses the server's local time zone.
- The Overview always shows the current month. The breakdown's date filter affects only the breakdown.

## Account groups

Group IDs are the seeded `account_groups` IDs, as already hard-coded in the create-account form and `CreateReceivableAccount` (Loan = 7).

| ID | Group | Overview block |
|---|---|---|
| 1, 2, 4, 5, 8 | Cash, Debit, E-Money, E-Wallet, Saving | Available funds |
| 7 | Loan (money owed to you) | Owed to you |
| 3, 6 | Credit, Debt (stored negative) | You owe (shown as a positive amount owed) |
| other | — | Other accounts |

Blocks are never added together. Available credit on a Credit account = `credit_limit + balance`, shown only when a limit above zero is set.

## Recent Transactions display

- Income is green with `+`, Expense is red with `-`. Transfer, Adjustment, Lend and Borrow are neutral.
- Adjustment amounts are stored as signed deltas and displayed as sign plus magnitude.
- A destination account is shown only for Transfers (other types may store one that balance logic ignores).
- Category is hidden for Transfers and Adjustments.

## Known limitations

- The accounts grid, Add Account form, and group classification depend on fixed group IDs 1–8.
- Adjust balance still uses `prompt()`; after create/adjust/cancel the page reloads or reloads every section.
- Closed accounts (`status`) are not filtered out.
- No combined view across owners yet (see `docs/MULTI_OWNER.md`).

## Untested behavior

No `*_test.go` files exist for the dashboard (test pause in `CLAUDE.md`). Untested: owner scoping and the owner sidebar, month totals, the 6-month series, group subtotals and block classification, available credit, category percentages, recent-transaction flattening and the `formatNumber` negative-number fix.