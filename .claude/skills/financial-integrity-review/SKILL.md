---
name: financial-integrity-review
description: Review drfin changes that can affect money amounts, transaction classification, account balances, or Gmail idempotency.
---

# Financial integrity review

Review the full path from input to persisted transaction and account balance.

Check:
- Expense decreases the source account.
- Income increases the receiving account.
- Transfer decreases the source and increases the destination by the same amount; it is not income or expense.
- Transfer source and destination must be valid and different.
- Transaction row changes and balance changes are atomic.
- Editing a transaction reverses the old balance effect before applying the new one.
- Deletion reverses the balance effect exactly once.
- Reprocessing the same Gmail message cannot create another transaction.
- Failed imports remain retryable and do not move the safe checkpoint forward.
- Amount parsing is exact for supported formats and never silently rounds or truncates.
- Nil destination, unknown type, missing account, and malformed amount cases fail safely.

Do not approve a financial change based only on code inspection. Request or add tests for the relevant cases. State what could not be verified.
