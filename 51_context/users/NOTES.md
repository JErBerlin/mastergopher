# Notes

## Key points

- `db.Query` internally calls `db.QueryContext` — same for `Exec`, `Prepare`, etc.
- Passing a `context.Context` allows:
  - Cancelling long-running or stuck queries
  - Enforcing deadlines or timeouts
  - Reducing resource usage on aborted requests
- Contexts start from `context.Background()` and are extended using `WithTimeout`, `WithCancel`, etc.
- A common function signature in concurrent Go code has a context as the first argument:

    `func GetUsers(ctx context.Context, db *sql.DB) ([]User, error)`

## Context behavior per driver

- Drivers must observe context but differ in how they react to cancellation.
  - PostgreSQL supports clean aborts on cancellation.
  - SQLite in-memory may not terminate the query immediately.
- Don't assume identical cancellation behavior across drivers.

## Takeaways

- Use `QueryContext`, `ExecContext`, etc., with a proper `context.Context` when:
  - Handling incoming HTTP requests
  - Running background jobs with time limits
  - Managing shutdowns or user-initiated cancellations

- Build a context with timeout in top level functions like this:  
    ```go  
    ctx := context.Background()  
    ctx, cancel := context.WithTimeout(ctx, 50*time.Millisecond)  
    defer cancel()
    ```

- Use `context.TODO()` as a placeholder when the actual context is not available yet. This is useful during refactoring or early stages of integration.

## Try it out

The code shows how to use `context.Context` when executing database queries. It demonstrates how to control execution time and handle cases where a query takes too long or the context is cancelled.

The timeout is defined as a constant:  
    ```go
    const queryTimeout = 5 * time.Millisecond`
    ```

You can modify this value to observe different behaviours:

- Change the value to a very low number (e.g. 1 microsecond) to force the deadline to expire before the query finishes. This will cause the operation to be aborted and print a timeout message.
- Change the value to a higher number (e.g. 100 millisecond) to allow the query to complete and print the user data normally.

