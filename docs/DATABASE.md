# Database (current schema)

The application uses PostgreSQL through `database/sql` and `github.com/lib/pq`. SQL definitions and seed data are under `resources/SQL/`.

## Core tables

| Table | Purpose | Key fields observed |
|---|---|---|
| `accounts` | Tracked accounts/wallets/cash/credit accounts | `id`, `name`, `balance numeric(24,2)`, `account_group_id`, `status`, `account_number`, `credit_limit` |
| `account_groups` | Account classification | `id`, `name` |
| `categories` | Transaction categories | `id`, `name` |
| `transaction_types` | Expense/Income/Transfer/etc. | `id`, `name` |
| `transactions` | Financial activity | `id`, `parent_id`, `trade_date`, `type_id`, `category_id`, `account_id`, `account_dest_id`, `amount numeric(24,2)`, `status`, `source`, timestamps |
| `last_email` | Gmail sync checkpoint | `message_id`, `internal_date` |
| `processed_gmail_messages` | Per-message processing state | Gmail message ID, internal date, transaction ID, status, error, processed timestamp |
| `owned_account` | Account labels/numbers used by older logic/data | `account_label`, `account_number` |
| `users`, `users_token` | User/token-related tables in config SQL | Present in SQL; current app wiring does not establish multi-user support |

## Relationships

- `transactions.type_id` → `transaction_types.id`
- `transactions.category_id` → `categories.id`
- `transactions.account_id` → `accounts.id`
- `transactions.account_dest_id` → `accounts.id` (nullable)
- `accounts.account_group_id` → `account_groups.id`
- `processed_gmail_messages.gmail_message_id` has a unique constraint.

## Seeded transaction types

The seed script inserts IDs for Expense, Income, Transfer, Lend, Borrow, and Adjustment. Existing service logic contains ID-specific branches, so changing seed IDs or names can break behavior.

## Adjustment transactions

Balance adjustments are stored in `transactions` with:
- `type_id` of the `Adjustment` type (looked up by name, seeded as ID 6)
- `category_id` of the category named `Unknown` (looked up by name; `category_id` is NOT NULL). If your database lacks this exact name, adjusting fails with "category ... not found".
- `amount` = signed difference (`new balance - old balance`), negative when the balance decreases. `transactions.amount` has no check constraint against negatives.
- `account_dest_id` NULL, `parent_id` 0, `status` `Confirmed`, `source` `Manual`.

Daily income/expense totals and the expense-category summary only count types named `Income` and `Expense`, so adjustments are excluded from them.

New accounts store `initial_balance` directly in `accounts.balance` with no transaction row. Account names must be unique (case-insensitively checked in the service, exact match enforced by the DB).

## Database safety rules for Claude

- Never run destructive SQL against a real database.
- Do not alter schema or seed IDs without explaining migration and compatibility impact.
- Prefer additive, reversible migrations for schema changes.
- Preserve transaction atomicity when updating a transaction and account balances.
- Add constraints/indexes only after checking existing data and query behavior.
- Do not assume `parent_id` semantics; inspect its usage before relying on it.
- Do not change money representation as an incidental cleanup. Propose a migration and tests separately.
- Never include real account numbers, tokens, or personal financial data in docs, fixtures, logs, or prompts.
- `populate_table.sql` inserts accounts before `account_groups` and uses explicit IDs for some tables. On a fresh database this can violate the account group foreign key, and sequences may not advance. Check with `SELECT setval('<table>_id_seq', (SELECT MAX(id) FROM <table>));` if inserts fail with a duplicate key.
