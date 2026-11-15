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

## Database safety rules for Claude

- Never run destructive SQL against a real database.
- Do not alter schema or seed IDs without explaining migration and compatibility impact.
- Prefer additive, reversible migrations for schema changes.
- Preserve transaction atomicity when updating a transaction and account balances.
- Add constraints/indexes only after checking existing data and query behavior.
- Do not assume `parent_id` semantics; inspect its usage before relying on it.
- Do not change money representation as an incidental cleanup. Propose a migration and tests separately.
- Never include real account numbers, tokens, or personal financial data in docs, fixtures, logs, or prompts.
