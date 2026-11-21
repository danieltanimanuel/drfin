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
- Do not create, edit or hand over test files (`*_test.go`) for now. This pause overrides any "add tests" guidance in this file or in `docs/`. Say which new behavior is untested instead.
- Never claim tests passed unless they were executed successfully.
- Run focused tests first, followed by `go test ./...` and `go vet ./...` when feasible.
- Report executed commands, results, and unverified behavior.

## Project References
Read these only when relevant to the task:

- `docs/ARCHITECTURE.md` — project structure and request flow.
- `docs/FEATURES.md` — implemented features and what is not implemented.
- `docs/RISKS.md` — known correctness/reliability risks (read when touching Gmail sync, parsers, or balance logic).
- `docs/DATABASE.md` — PostgreSQL schema.
- `docs/DEVELOPMENT.md` — setup, commands, and safety notes.
- `docs/GMAIL_SETUP.md` — Gmail OAuth credential setup (only for auth/credential problems).
- `docs/todo/` — one file per unimplemented feature design. Read only the file for the task at hand.

## Token Efficiency
- Never scan or read the entire repository unless explicitly requested.
- Start with relevant entry points, targeted searches, and task-specific files.
- Follow direct dependencies only as needed to understand the implementation.
- Read documentation only when it answers a task-specific question.
- Do not reread files already inspected unless they may have changed.
- Keep command output targeted. Avoid dumping entire files, directory trees, or large logs.
- Do not repeat established project context or generate unsolicited documentation and summaries.
- Keep progress updates concise. Report findings, decisions, and blockers.
- Docs layout: keep `CLAUDE.md` short. Put new topics in their own file under `docs/` (designs in `docs/todo/`) and add a one-line pointer here. Do not append large sections to existing docs.

## Definition of Done
A change is complete when it:
- Implements only the requested scope.
- Follows existing architecture and patterns.
- Includes no new test files while the test pause in "Testing" is active; lists new untested behavior under remaining risks.
- Has no unrelated modifications.
- Reports changed files, executed tests, and remaining risks or assumptions.