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

- Always use `context.Context` with DB operations.
- Don't assume consistent cancellation behavior across drivers.
- For predictable results, prefer production-ready drivers that are known to handle context correctly.


