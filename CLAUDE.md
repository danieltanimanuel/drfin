# drfin — Claude Code Instructions

## Project
drfin is a personal finance tracker. It is an existing Go web application that tracks accounts and transactions and imports supported bank notifications from Gmail.

## Working rules
- Read this file first, then read only the documentation and source files relevant to the current task.
- Before editing, inspect the existing implementation and explain the smallest safe change. For multi-file, database, Gmail, or financial-logic work, plan first and wait for approval.
- Keep changes minimal. Prefer the existing patterns and Go standard library. Add a dependency only when it clearly reduces complexity or risk.
- Do not refactor unrelated code, rename packages, introduce layers/frameworks, or redesign the architecture unless asked.
- Preserve existing behavior unless the task explicitly changes it. Do not invent features that are only mentioned in future requirements.
- Never read, print, commit, or expose secrets, OAuth tokens, real email contents, or personal financial data. Treat `token.json` and `resources/drfin-secret.json` as sensitive.
- Do not change database schema or migration/SQL scripts without explaining compatibility and data-loss risks first.
- Financial logic is high risk: trace the full transaction and balance flow; verify Expense, Income, and Transfer separately; use database transactions where multiple records/balances must change.
- Do not assume `math/big.Float` is safe for monetary arithmetic. Understand existing numeric types and database precision before changing them. Avoid introducing floating-point conversions.
- Add or update focused tests for behavior changes. Never claim tests passed unless you ran them.
- After changes, run the narrowest relevant tests, then `go test ./...` and `go vet ./...` when feasible. Report commands and results.
- Keep comments and documentation concise and update docs only when behavior or architecture actually changes.

## Project references
- `docs/ARCHITECTURE.md` — current structure and request flow
- `docs/FEATURES.md` — implemented features versus known gaps
- `docs/DATABASE.md` — current PostgreSQL schema
- `docs/DEVELOPMENT.md` — local setup, commands, and safety notes

## Definition of done
A change is complete when it is scoped, understandable, tested where practical, has no unrelated edits, and its remaining risks or unverified assumptions are stated.
