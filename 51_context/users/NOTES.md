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

- Pass an explicit `context.Context` to DB methods (`QueryContext`, `ExecContext`, etc.) in server code to support request cancellation and timeouts.
- Context is useful when a query may take long or the request may be aborted (e.g. user closes browser, system sends shutdown signal).
- The database driver may abort the running query if the context is cancelled, though support varies across drivers.
- Use `http.Request.Context()` or `context.WithTimeout` to manage the context lifecycle.
- For short-lived CLI tools or scripts where the entire program runs synchronously, passing context to db calls is unnecessary.
