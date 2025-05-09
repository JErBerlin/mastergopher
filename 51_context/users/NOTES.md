# Notes

## Key points

- `db.Query` internally calls `db.QueryContext` — same for `Exec`, `Prepare`, etc.
- Passing a `context.Context` allows:
  - Cancelling long-running or stuck queries
  - Enforcing deadlines or timeouts
  - Reducing resource usage on aborted requests
- All contexts start with `context.Background()`, from which child contexts are created using `WithTimeout`, `WithCancel`, etc.

## Context behavior per driver

- Database drivers must observe context but do not behave the same.
  - PostgreSQL cancels queries cleanly on context cancellation.
  - SQLite in-memory mode may not abort immediately.
- Don't rely on consistent cancellation across drivers, but still pass context to allow cancellation when supported.

## Takeaways

- Use `QueryContext`, `ExecContext`, etc., with a proper `context.Context` when:
  - Handling incoming HTTP requests
  - Running background jobs with time limits
  - Managing shutdowns or user-initiated cancellations

- Example of building a context with timeout:  
    ```go  
    ctx := context.Background()`  
    ctx, cancel := context.WithTimeout(ctx, 50*time.Millisecond)`  
    defer cancel()`
    ```
- For CLI tools, simple scripts, or short queries in isolated code, using context is optional.

## Try it out

The code demonstrates how to use `context.Context` when executing database queries. It shows how to control execution time and handle cases where a query takes too long or the context is cancelled.

The timeout is defined as a constant:  
(Go) `const queryTimeout = 5 * time.Millisecond`

You can modify this value to observe different behaviours:

- Change the value to a very low number (e.g. 1 microsecond) to force the deadline to expire before the query finishes. This will cause the operation to be aborted and print a timeout message.
- Change the value to a higher number (e.g. 100 millisecond) to allow the query to complete and print the user data normally.

