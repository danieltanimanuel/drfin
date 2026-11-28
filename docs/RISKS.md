# Known correctness and reliability risks

Moved from `docs/FEATURES.md`. These are observations for future review, not permission to refactor everything at once. Read only when changing the related area.

1. Monetary values are represented using `*big.Float` in Go and `numeric(24,2)` in PostgreSQL. Review precision, rounding, and float conversions before changing financial calculations.
2. Transaction types are resolved by name, but some logic relies on numeric IDs (for example, Expense/Income/Transfer IDs). Confirm seed data and constraints before modifying.
3. `transactions.Service` applies/reverses account balances. Any change to transaction type, amount, source account, or destination account must be traced through both reversal and reapplication paths.
4. Fixed in the multi-owner change: `syncMail` now returns the `listMessages` error. Still untested.
5. A Gmail import can create a transaction and then separately record the message as processed. Review failure windows and idempotency before changing sync behavior.
6. Resolved by multi-owner: the processing key is `(owner_id, gmail_message_id)`.
7. Some parser functions use string slicing and assumptions about exact email layouts. Test malformed and changed email formats. `ParseDate` strips the timezone by removing the last 4 characters (" WIB") for blu and BCA credit-card emails, and parsed times carry no timezone offset. Only the Indonesian month abbreviations seen in real emails are mapped, so a new variant fails that email's import.
8. Do not mark a message successful or advance the checkpoint if transaction persistence fails.
9. A failed email stops that owner's `syncMail` loop (`break`) and the checkpoint does not advance past it. Every later email is retried on each scan until the failure is fixed, so one unparsable email blocks all newer imports **for that owner** (other owners are unaffected). The exception is an intentional skip (`ErrSkipEmail`, currently only `MYR` credit-card emails): it does not stop the loop and the checkpoint advances.
10. **Possible duplicate for credit card QRIS payments.** A QRIS payment made with the credit card may also trigger a separate credit card notification handled in `credit_card_bca.go`. Not yet verified whether both are imported for the same purchase.
11. **Pocket withdrawals always post to `Tahapan Gold`.** For a BCA pocket email with `Source Pocket` + `Beneficiary Account`, the destination is not checked against the beneficiary. A pocket transfer to another person would be recorded as an internal transfer to `Tahapan Gold`. Accepted for now because only internal withdrawals are used.
12. **Pocket top-ups from `Source of Fund` assume `Tahapan Gold`.** The masked source value is not inspected.
13. **Pocket accounts depend on exact names.** Accounts are built as `"MyBCAPocket " + <pocket name in email>`. If an account with that name doesn't exist, the import fails and (per item 9) blocks later emails until fixed (per owner; names are looked up within the owner). The Transfer balance effect for pocket accounts has not been verified end to end.
14. **`/scan-email` responds with JSON** (one entry per owner) and returns HTTP 500 if any owner failed. Callers expecting the old plain-text response must be checked. All owners in one request share one 15-minute timeout.
15. **Adjust balance, cancel and category update work by id** and do not verify that the id belongs to the owner selected in the dashboard.
16. **Blu payments to another owner's account import as an Expense**, not a Transfer, because destination lookup is owner-scoped. Cross-owner transfers are not supported.
17. **A new owner needs a `users` row and a `last_email` row before its first sync.** A `last_email.internal_date` of `0` imports the entire mailbox history, and every import changes balances, so set real balances and a start date first.
18. **blu parsing depends on `users.bank_holder_name` matching the name printed in that owner's emails exactly.** A wrong value makes the parser miss the beneficiary fields, so the email imports without a destination.