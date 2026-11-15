---
name: gmail-parser
description: Work on Gmail ingestion and bank notification parsing in drfin without breaking idempotency or financial data.
---

# Gmail parser work

- Inspect `docs/FEATURES.md`, `docs/ARCHITECTURE.md`, and relevant files under `app/gmail/`, `app/processedgmail/`, `app/lastemail/`, `app/transactions/`, and `resources/sample_html/`.
- Never read OAuth secret/token files or real mailbox contents.
- Treat email HTML as untrusted and variable input.
- Preserve the distinction between extraction, normalization, bank-specific mapping, and transaction persistence.
- Test supported subject lines, missing fields, malformed HTML, unexpected amount formats, date/time variants, and unknown transaction types.
- Check parser behavior for BCA, blu, and BCA credit-card formats independently.
- Verify that a failed parse or failed transaction write does not mark the message successful or advance the checkpoint.
- Verify duplicate-message behavior and consider that message IDs should be scoped by Gmail account if multi-account support is introduced.
- Do not invent bank parsing behavior from sample names alone; inspect fixtures and current parser code.
- Do not change balance logic as part of a parser-only task unless the approved scope requires it.
