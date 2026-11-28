# Multi-owner support

drfin supports several owners, each with one Gmail mailbox. An owner is a row in the `users` table. Each owner has their own accounts, transactions and Gmail sync state. Account groups, categories and transaction types are shared.

## What is scoped and what is shared

| Scope | Tables |
|---|---|
| Per owner | `accounts` (`owner_id`), `last_email` (`owner_id`), `processed_gmail_messages` (`owner_id`) |
| Per owner, indirectly | `transactions` (no `owner_id` column; ownership is `transactions.account_id` -> `accounts.owner_id`) |
| Shared | `account_groups`, `categories`, `transaction_types` |
| Unused, unchanged | `owned_account` (not referenced by the code), `users_token` |

Constraints that changed:

- `accounts`: `UNIQUE (owner_id, name)` replaces the global `accounts_name_key`. Two owners may both have `Tahapan Gold`.
- `last_email`: `UNIQUE (owner_id)`, one checkpoint row per owner.
- `processed_gmail_messages`: `UNIQUE (owner_id, gmail_message_id)` replaces `uq_gmail_message`, because Gmail message IDs are unique per mailbox, not globally.
- `users`: `name`, `bank_holder_name` added; `email` is unique.

`bank_holder_name` is the account holder's name exactly as printed in that owner's blu emails. `parseBLUXML` uses it as a landmark to find the beneficiary fields. Do not commit real names to the repository; set them locally with `UPDATE users SET bank_holder_name = '...' WHERE id = N;`.

## Rules

- A transaction belongs to the owner of its source account. Transfers across owners are not supported.
- `CreateTransaction` requires `owner_id`. Source and destination accounts are looked up by name **within that owner**.
- `UpdateTransaction` takes the owner from the existing transaction's source account, never from the request, so a transaction cannot move to another owner's accounts.
- `accounts.CreateAccount` and `CreateReceivable` require an owner.
- Account lookups by name or account number are always owner-scoped.
- Queries that list or aggregate take an optional owner (`*int`) and use `($1::int IS NULL OR a.owner_id = $1)`. `nil` means all owners; this is what a future combined view will use.

## Adding an owner (manual)

1. Insert a `users` row (`name`, `email`, `bank_holder_name`).
2. Insert a `last_email` row for that owner. Its `internal_date` (Unix milliseconds) is where the first sync starts; `0` imports the whole mailbox history.
3. Create the owner's accounts from the dashboard. Names the parsers use must exist: `Tahapan Gold`, `CC BCA`, `Blu`, and `MyBCAPocket <name>` for pocket transfers. Set current balances before importing, because every import changes balances.
4. Add the Gmail address as a test user on the Google OAuth consent screen.
5. Authorize the mailbox: `go run . -authorize-owner <id>` and sign in with that owner's Gmail. The token is stored in `tokens/<id>.json`.

## Gmail sync

- `POST /scan-email?owner=<id>` syncs one owner. Without `owner`, every owner in `users` is synced, one after another and independently: one owner's failure does not stop the others.
- The response is JSON, one entry per owner (`owner_id`, `name`, `ok`, `error`). If any owner failed, the HTTP status is 500.
- A request shares one 15-minute timeout across all owners it syncs.
- Tokens are per owner (`tokens/<id>.json`, git-ignored). The shared OAuth client is `resources/drfin-secret.json`. A sync never prompts for authorization; a missing token is an error that names the `-authorize-owner` command.
- A blu email for an owner whose `bank_holder_name` is empty fails with an error instead of importing as an Expense with no destination.
- Within one owner, a failing email stops that owner's sync at that email and is retried on the next scan.
- Blu transfers are detected by looking up the beneficiary account number among the **same owner's** accounts, so a blu payment to another owner's account imports as an Expense, not a Transfer.

## Dashboard

- `/dashboard?owner=<id>` selects the owner; no parameter means owner 1. An unknown owner id returns 404.
- The sidebar lists every row of `users` by name. The shell page renders the selected id into `<body data-owner-id>`, and the page JavaScript adds `owner=<id>` to every fragment request and to the create requests.
- Fragment routes (`/dashboard/accounts`, `/dashboard/overview`, `/dashboard/cashflow`, `/dashboard/categories`, `/dashboard/transactions`, `/dashboard/recent`, `/dashboard/pie`, `/dashboard/transactions/new`) accept `?owner=<id>`.
- The transaction edit form takes its owner from the transaction being edited.
- Adjust balance, cancel, and category update work by id and do not verify that the id belongs to the selected owner.

## Not implemented

- Combined view across owners (planned: period filter defaulting to the current month; combined income, expense and net with a per-owner table; combined assets by shared account group; combined expense breakdown by category).
- Cross-owner transfers.
- A UI to manage owners.
- Owner checks on adjust balance, cancel and category update.

## Verification status

- No automated tests cover multi-owner behavior (tests are paused, see `CLAUDE.md`). Untested: owner-scoped SQL, per-owner checkpoint, token and processed-message handling, `SyncOwners` isolation, the `/scan-email` response, `bank_holder_name` matching, and the `UpdateTransaction` owner derivation.
- Owner switching in the dashboard has been checked manually. Gmail sync for a second owner had not been run when this was written.