---
name: minimal-code
description: Keep drfin changes small, direct, and consistent with the existing Go codebase.
---

# Minimal code

- Prefer the smallest change that correctly solves the requested problem.
- Follow the existing package structure and handler/service/repository pattern.
- Prefer the Go standard library and existing dependencies. Do not add a library for convenience alone.
- Avoid speculative abstractions, generic frameworks, interfaces with only one implementation, helper layers, and premature configurability.
- Do not refactor neighboring code unless required for correctness.
- Reuse existing functions when their behavior is correct and clear.
- Keep functions focused; do not split code into many tiny helpers without a clear benefit.
- Do not duplicate business rules across handlers, services, and repositories.
- Add comments only for non-obvious reasoning or important invariants.
- Explain any proposed dependency, abstraction, or broader refactor before implementing it.
- Minimal does not mean skipping validation, error handling, transactions, or tests.
