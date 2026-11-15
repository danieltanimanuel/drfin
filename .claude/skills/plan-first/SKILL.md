---
name: plan-first
description: Investigate and plan non-trivial drfin changes before editing, especially database, Gmail, transaction, and balance work.
---

# Plan first

1. Restate the requested behavior and define what is out of scope.
2. Read `CLAUDE.md` and only the relevant sections of `docs/`.
3. Inspect the smallest relevant set of source files, including callers and tests where necessary.
4. Report the files inspected and the current execution/data flow.
5. Identify assumptions, edge cases, compatibility risks, and missing information.
6. Propose the smallest implementation plan, listing files to change and tests to add.
7. Do not edit files until the user approves the plan.

Do not keep exploring unrelated packages. Ask for another file only when a concrete dependency requires it.
