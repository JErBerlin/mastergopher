# Notes

## Key points

* Go’s error handling is explicit and based on return values. Functions return `error` as the last return value.
* The `error` type is an interface with a single method: `Error() string`.
* Errors are values. They can be assigned, passed, and compared like any other value.
* Panic stops the normal execution of a function. It can be recovered using `recover()` inside a `defer` block.
    * `recover()` only works in a deferred function in the **same** goroutine as the panic.
    * Panic+recover is sometimes used internally to abort deep recursion and convert expected panics back into `error` returns.
* `errors.Is` and `errors.As` are used to inspect wrapped errors returned with `fmt.Errorf(... %w ...)`.
    * `errors.Is(err, targetErr)` walks the wrapped chain and does an identity (`==`) check against a **sentinel variable**.
    * `errors.As(err, &targetType)` finds the first error in the chain assignable to **that type**, then lets you access its fields.


## Sentinel errors

* Sentinel errors are predefined values used to indicate a specific condition.

  * Example: `io.EOF`, `sql.ErrNoRows`.
  * You can define your own: `var ErrX = errors.New("message")`.
* Equality comparison only works if the same sentinel instance is reused.
* Constants with identical content are **fungible** — interchangeable and equal.

  * `const a = "msg"`, `const b = "msg"` → `a == b` is true.
  * `errors.New("msg")` creates a new value each time — not equal.

* Comparison only works if the **same** sentinel instance is reused:
        
```go
  if errors.Is(err, fs.ErrNotExist) {
    // checks package-level var fs.ErrNotExist (in io/fs/fs.go)
  }
```

* Pitfall: using `errors.New("msg")` inline or recreating the var at runtime produces a distinct value—`errors.Is` will never match.
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

* Use type assertions or `errors.As` to access fields in wrapped errors:
       
```go 
  var pe *os.PathError
  if errors.As(err, &pe) {
    // pe.Op, pe.Path, pe.Err available
  }
```

* Useful when you want to pass additional data along with the error.

## Takeaways

* Avoid ignoring errors. Always check `if err != nil` after function calls.
* Use `errors.New` or `fmt.Errorf` (with `%w`) to generate and wrap errors.
* Use `errors.Is` to compare against **sentinel variables**, and `errors.As` to extract **typed** errors.
* Do not recreate sentinel errors at runtime; declare them once at package scope.
* Use panic+recover internally for complex control flow (e.g., deep recursion), but always re‑panic unexpected values.
* Keep public APIs clean by exposing only `error` returns, not panic details.
* Avoid using `panic` in application code—reserve it for unrecoverable or internal-only failures.

## Try it out

```go
package main

import (
    "errors"
    "fmt"
    "log"
    "runtime/debug"
)

var errEmptyInput = errors.New("input cannot be empty")

type missingError struct {
    key string
}

func (e missingError) Error() string {
    return fmt.Sprintf("no entry found for %q", e.key)
}

func fetch(key string) (string, error) {
    db := map[string]string{"lang": "Go", "tool": "gopher", "pkg": "errors"}

    if key == "" {
        return "", errEmptyInput
    }
    if val, ok := db[key]; ok {
        return val, nil
    }
    return "", missingError{key}
}

// TODO: refactor this to use errors.Is
func checkEmpty() {
    key := ""
    v, err := fetch(key)
    if err != nil {
        if err != errEmptyInput {
            log.Fatal(err)
        }
        fmt.Println("please provide a non-empty key")
        return
    }
    fmt.Printf("value: %q\n", v)
}

// TODO: refactor this to use errors.As
func checkMissing() {
    key := "unknown"
    v, err := fetch(key)
    if err != nil {
        if _, ok := err.(missingError); !ok {
            log.Fatal(err)
        }
        fmt.Printf("missing: %s\n", err)
        return
    }
    fmt.Printf("value: %q\n", v)
}

// TODO: modify this to recover when depth == 5 and print a stack trace
func nestedWalk(depth int) {
    if depth == 0 {
        return
    }
    fmt.Println("depth", depth)
    nestedWalk(depth - 1)
}

func main() {
    checkEmpty()
    checkMissing()
    // after refactoring, this will panic at depth 5 and recover
    nestedWalk(7)
}
```
