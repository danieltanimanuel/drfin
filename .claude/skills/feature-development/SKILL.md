---
name: feature-development
description: Implement a focused drfin feature safely using the existing architecture.
---

# Feature development

1. Use `plan-first` for any task spanning multiple files or affecting data/financial behavior.
2. Identify the existing route, handler, service, repository, model, SQL, and template paths relevant to the feature.
3. Confirm the current behavior from code; do not treat TODOs or docs as proof of implementation.
4. Define acceptance criteria and edge cases before coding.
5. Implement only the approved scope using existing conventions.
6. Add focused tests. For money and transactions, cover positive/negative direction, rounding/precision assumptions, source/destination accounts, and rollback behavior as relevant.
7. Run `gofmt`, focused tests, `go test ./...`, and `go vet ./...` where feasible.
8. Review the diff for unrelated edits, secrets, logging of sensitive data, and unintended schema changes.
9. Summarize changed files, behavior, tests run, and remaining risks.
