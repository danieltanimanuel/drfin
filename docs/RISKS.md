# Known correctness and reliability risks

Moved from `docs/FEATURES.md`. These are observations for future review, not permission to refactor everything at once. Read only when changing the related area.

1. Monetary values are represented using `*big.Float` in Go and `numeric(24,2)` in PostgreSQL. Review precision, rounding, and float conversions before changing financial calculations.
2. Transaction types are resolved by name, but some logic relies on numeric IDs (for example, Expense/Income/Transfer IDs). Confirm seed data and constraints before modifying.
3. `transactions.Service` applies/reverses account balances. Any change to transaction type, amount, source account, or destination account must be traced through both reversal and reapplication paths.
4. Gmail `syncMail` currently ignores the error returned by `listMessages` (`msgs, _ := ...`) before checking a different `err` variable. This deserves a focused bug fix with tests.
5. A Gmail import can create a transaction and then separately record the message as processed. Review failure windows and idempotency before changing sync behavior.
6. The Gmail processing repository uses Gmail message ID as the unique key; current schema does not include Gmail account identity in that key.
7. Some parser functions use string slicing and assumptions about exact email layouts. Test malformed and changed email formats. `ParseDate` strips the timezone by removing the last 4 characters (" WIB") for blu and BCA credit-card emails, and parsed times carry no timezone offset. Only the Indonesian month abbreviations seen in real emails are mapped, so a new variant fails that email's import.
8. Do not mark a message successful or advance the checkpoint if transaction persistence fails.
9. A failed email stops the `syncMail` loop (`break`) and the checkpoint does not advance past it. Every later email is retried on each scan until the failure is fixed, so one unparsable email blocks all newer imports. The exception is an intentional skip (`ErrSkipEmail`, currently only `MYR` credit-card emails): it does not stop the loop and the checkpoint advances.
10. **Possible duplicate for credit card QRIS payments.** A QRIS payment made with the credit card may also trigger a separate credit card notification handled in `credit_card_bca.go`. Not yet verified whether both are imported for the same purchase.
11. **Pocket withdrawals always post to `Tahapan Gold`.** For a BCA pocket email with `Source Pocket` + `Beneficiary Account`, the destination is not checked against the beneficiary. A pocket transfer to another person would be recorded as an internal transfer to `Tahapan Gold`. Accepted for now because only internal withdrawals are used.
12. **Pocket top-ups from `Source of Fund` assume `Tahapan Gold`.** The masked source value is not inspected.
13. **Pocket accounts depend on exact names.** Accounts are built as `"MyBCAPocket " + <pocket name in email>`. If an account with that name doesn't exist, the import fails and (per item 9) blocks later emails until fixed. The Transfer balance effect for pocket accounts has not been verified end to end.
