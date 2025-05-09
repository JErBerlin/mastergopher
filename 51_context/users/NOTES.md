# Notes

## Key points

- `db.Query` internally calls `db.QueryContext` — same for other methods like `Exec`, `Prepare`, etc.
- Using `db` methods with `context.Context` allows:
  - Request cancellation
  - Deadline/timeout propagation
  - Driver-specific optimisations

## Context behavior per driver

- Drivers behave differently when context is cancelled.
  - Example:
    - PostgreSQL cancels query properly.
    - SQLite in-memory may not react the same way.
- Important: All drivers **must observe** the context and allow cancellation, but the way they do it is not unified.

## Takeaways

- Prefer using context-aware DB methods (`QueryContext`, `ExecContext`, etc.) with an explicit `context.Context`, especially in request-handling code (e.g. HTTP handlers, background workers with deadlines).
- Manage the `context.Context` lifecycle yourself (e.g. via `http.Request.Context()` or `context.WithTimeout`) to allow proper cancellation or timeout control.
- Avoid relying on default wrappers like `db.Query` unless context is not relevant to your use case (e.g., CLI scripts, short-lived tools).

