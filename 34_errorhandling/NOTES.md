# Notes

## Key points

* Go’s error handling is explicit and based on return values. Functions return `error` as the last return value.
* The `error` type is an interface with a single method: `Error() string`.
* Errors are values. They can be assigned, passed, and compared like any other value.
* Panic stops the normal execution of a function. It can be recovered using `recover()` inside a `defer` block.
* `errors.Is` and `errors.As` are used to inspect wrapped errors returned with `fmt.Errorf(... %w ...)`.

## Sentinel errors

* Sentinel errors are predefined values used to indicate a specific condition.

  * Example: `io.EOF`, `sql.ErrNoRows`.
  * You can define your own: `var ErrX = errors.New("message")`.
* Equality comparison only works if the same sentinel instance is reused.
* Constants with identical content are **fungible** — interchangeable and equal.

  * `const a = "msg"`, `const b = "msg"` → `a == b` is true.
  * `errors.New("msg")` creates a new value each time — not equal.
* Sentinel errors expose internal logic and can couple packages.
* Prefer exposing helper functions like `IsX(err)` instead of exporting the error variable.

## Custom error types and inspection

* A custom error type is any type that implements `Error() string`.

  * Example:

    ```go
    type MyError struct {
        Msg  string
        Code int
    }

    func (e MyError) Error() string {
        return fmt.Sprintf("%d: %s", e.Code, e.Msg)
    }
    ```

* Use type assertions to access fields:

  ```go
  if err, ok := err.(MyError); ok {
      fmt.Println(err.Code)
  }
  ```

* `errors.As` can be used for wrapped errors.

* Useful when you want to pass additional data along with the error.

## Takeaways

* Avoid ignoring errors. Always check `if err != nil` after function calls.
* Use `errors.New` or `fmt.Errorf` to generate errors.
* When returning wrapped errors, prefer `fmt.Errorf(... %w ...)` for context.
* Use `errors.Is` to compare, and `errors.As` to extract specific error types.
* Avoid using `panic`, except for unrecoverable internal failures.
* Do not expose sentinel errors as public variables; instead, expose helper match functions.

## Try it out

*to be filled*

