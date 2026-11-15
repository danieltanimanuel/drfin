# drfin — Claude Code Instructions

## Project
drfin is a personal finance tracker built with Go. It tracks accounts and transactions and imports supported bank notifications from Gmail.

## Working Rules
- Read this file first. Inspect only documentation and source files relevant to the task.
- Investigate the existing implementation before proposing changes.
- Before editing, explain the smallest safe change and identify affected files.
- For multi-file, database, Gmail, or financial-logic changes, present a plan and wait for approval.
- Keep changes minimal. Follow existing patterns and prefer the Go standard library.
- Avoid unrelated refactoring, new abstractions, dependencies, or architectural changes unless necessary or requested.
- Preserve existing behavior. Do not implement features mentioned only in future requirements.
- Never expose, print, or commit secrets, OAuth tokens, real email contents, or personal financial data. Treat `token.json` and `resources/drfin-secret.json` as sensitive. Do not read their contents.
- Explain compatibility and data-loss risks before changing database schemas or SQL migrations.

## Financial Logic
- Trace the complete transaction and balance flow before making changes.
- Verify Expense, Income, and Transfer behavior separately.
- Use database transactions when multiple records or balances must change together.
- Understand existing numeric types and database precision before modifying monetary calculations.
- Avoid floating-point conversions. Do not assume `math/big.Float` is suitable for monetary arithmetic.

## Testing
- Add or update focused tests for behavior changes.
- Never claim tests passed unless they were executed successfully.
- Run focused tests first, followed by `go test ./...` and `go vet ./...` when feasible.
- Report executed commands, results, and unverified behavior.

## Project References
Read these only when relevant to the task:

- `docs/ARCHITECTURE.md` — project structure and request flow.
- `docs/FEATURES.md` — implemented features and known gaps.
- `docs/DATABASE.md` — PostgreSQL schema.
- `docs/DEVELOPMENT.md` — setup, commands, and safety notes.

## Token Efficiency
- Never scan or read the entire repository unless explicitly requested.
- Start with relevant entry points, targeted searches, and task-specific files.
- Follow direct dependencies only as needed to understand the implementation.
- Read documentation only when it answers a task-specific question.
- Do not reread files already inspected unless they may have changed.
- Keep command output targeted. Avoid dumping entire files, directory trees, or large logs.
- Do not repeat established project context or generate unsolicited documentation and summaries.
- Keep progress updates concise. Report findings, decisions, and blockers.

## Definition of Done
A change is complete when it:
- Implements only the requested scope.
- Follows existing architecture and patterns.
- Includes appropriate tests where practical.
- Has no unrelated modifications.
- Reports changed files, executed tests, and remaining risks or assumptions.