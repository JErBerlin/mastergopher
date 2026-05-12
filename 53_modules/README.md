# Notes

## Key points

- Packages organize Go code inside a module. One directory normally contains one package.
- Exported and unexported names define the package API boundary. This is controlled by capitalization.
- Executable programs use package `main` and a `main` function. Larger modules often put executable entry points under `cmd/...`.
- Modules are the dependency and versioning unit. A module is defined by `go.mod` and can contain one or more packages.
- `go.mod` records the module path, Go version, and required module versions. `go.sum` records checksums used to verify downloaded module content.
- Dependencies are selected from `go.mod` and the module graph. Go does not automatically upgrade everything to the newest version during normal builds.
- `go get` changes dependency versions. `go mod tidy` synchronizes module metadata with the imports used by the module and its tests.
- The Go command can use the module cache and module proxies during normal builds. Private modules usually need `GOPRIVATE` configuration.
- Go does not allow circular package imports, so package boundaries must follow dependency direction.
- Major versions `v2+` use semantic import versioning, so the major version appears in the module path.
- Workspaces (`go.work`) help when developing several local modules together.

## Package structure and visibility

- Go code is organized in packages.
- A package is made from the `.go` files in one directory.
- One directory should contain only one package, except for special test package cases.
- Every `.go` source file must declare its package name at the top.
- Package names should be short, lower-case, and normally match the folder name.
- Export only what other packages really need.
- Exported names form part of the package API and are harder to change later.

## File names

- Go does not force file names inside a package.
- Common practice: use a file named after the package as the main entry file, for example `store/store.go` for package `store`.
- Package-level documentation usually belongs near the main package file.

## Modules and dependencies

- A module is a collection of related Go packages versioned together.
- A module is defined by a `go.mod` file at the module root.
- The module path is the import prefix for packages inside the module.
- The `go` line defines the minimum Go version required by the module.
- `require` lines define minimum required dependency versions.
- Dependencies are resolved from imports and the module graph.
- Downloaded modules are stored in the module cache. Check the exact location with `go env GOMODCACHE`.
- The module cache avoids downloading the same module version again and makes repeated builds faster.
- `go.sum` stores checksums for downloaded module content and related `go.mod` files.
- A `// indirect` requirement is a module needed through another dependency, or a module that provides packages used only indirectly by the current module or its tests.

Create a module:

```sh
go mod init example.com/myapp
```

Use this after adding, removing, or changing imports:

```sh
go mod tidy
```

`go mod tidy` adds missing module requirements and removes unused ones based on the packages imported by the module and its tests. It also updates `go.sum`.

## Daily workflow with modules

Normal workflow for adding a dependency:

1. Import the package in the Go file.
2. Run `go test ./...`, `go build ./...`, or `go mod tidy`.
3. Go resolves the missing module and updates `go.mod` / `go.sum` as needed.

Use `go get` when you want to control the dependency version explicitly. For executable tools, use `go install package@version` instead.

After changing dependencies, run:

```sh
go mod tidy
go test ./...
```

Do not normally edit `go.mod` by hand for dependency changes. Let the Go command update it, then review the result.

## Dependency command reference

Go does not automatically move dependencies to the newest versions during normal builds. The selected versions come from `go.mod` and the module graph. This is intentional: builds should be repeatable. If you want newer versions, make that explicit with `go get`, then clean and test the module.

Start by checking the current dependency graph:

```sh
go list -m all
```

Check which modules have newer versions available:

```sh
go list -m -u all
```

Add a dependency or upgrade one dependency using Go's default upgrade query:

```sh
go get example.com/pkg
```

Use the latest available version explicitly:

```sh
go get example.com/pkg@latest
```

Use a specific version:

```sh
go get example.com/pkg@v1.2.3
```

Upgrade all modules that provide packages imported by the current module:

```sh
go get -u ./...
```

Upgrade all of those modules only to latest patch versions:

```sh
go get -u=patch ./...
```

Downgrade one dependency:

```sh
go get example.com/pkg@v1.1.0
```

Remove a dependency requirement:

```sh
go get example.com/pkg@none
```

After any dependency change, clean and verify:

```sh
go mod tidy
go test ./...
```

Install an executable tool:

```sh
go install example.com/tool/cmd/tool@latest
```

## Module proxies and private modules

- A module proxy is a server used by the Go command to download module versions, `go.mod` files, and module zip files.
- This is part of normal module builds, not an ad hoc workaround. By default, Go uses `https://proxy.golang.org,direct` through `GOPROXY`.
- Proxies improve repeatability and speed because a requested public module version can stay available even if the original repository changes later.
- In companies, the usual proxy-related issue is private modules. A private module path should not be sent to the public proxy or public checksum database.
- Use `GOPRIVATE` for private module path patterns:

```sh
go env -w GOPRIVATE=github.com/mycompany/*
```

- Check current proxy/private settings:

```sh
go env GOPROXY GOPRIVATE GONOPROXY GONOSUMDB
```

- For normal public dependencies, the default proxy setup is usually fine.
- For private dependencies, configure `GOPRIVATE` and make sure Git authentication works in local development and CI.

## Inspecting dependencies

Useful commands when the dependency graph is unclear:

```sh
go list -m all              # show selected module versions
go list -m -u all           # show available updates
go mod why example.com/pkg  # explain why a package is needed
go mod graph                # print module dependency edges
```

Use `go mod why` when you do not know why a package is present. Use `go mod graph` when you need to see which module requires which other module. For large output, filter it with tools like `grep`.

The Go command is integrated with normal build/test commands. If a needed dependency is missing, commands such as `go test` can download it and update module metadata.

## Executable packages

- An executable Go program must use package `main`.
- It must declare exactly one `main` function.
- `main` receives no arguments and returns no values.
- Use `go install package@version` to install a command at a specific version. Do not use `go get` for this in modern Go.

## Package organization

- Go does not allow circular imports.
- Package boundaries should follow dependency direction: entry packages import application packages, and application packages import lower-level packages. The dependency should not point back upward. For example, `cmd/server` can import `httpapi`, and `httpapi` can import `store`, but `store` should not import `httpapi` just to call a handler or reuse a helper.
- If two packages need the same helper, move that helper to a lower-level package that both can import. This avoids sibling packages importing each other only to share code.
- A common layout is to put executable programs under `cmd/...` and keep reusable code outside `cmd`, for example `cmd/server/main.go` for startup code and packages like `store`, `httpapi`, or `cache` for reusable code.
- Do not over-design the layout too early. Write working code first, then split or merge packages when the shape becomes clear.

## Versioning and semantic import versioning

- Go modules use semantic versioning: `v1.2.3` means major version `1`, minor version `2`, patch version `3`.
- A breaking change requires a new major version.
- For major version `v2` or higher, the module path must include the major version suffix.

Example:

```go
module example.com/payment/v2
```

And imports must also use the `v2` path:

```go
import "example.com/payment/v2/client"
```

This allows code to use `v1` and `v2` of the same module at the same time, because they are treated as different import paths.

## Go workspaces

- A workspace is useful when several local modules are developed together.
- It is defined by a `go.work` file.
- A `replace` directive in `go.mod` can point one module dependency to a local directory.
- This is useful for temporary local development, but it changes the module file itself.
- A workspace solves the same local-development problem at a higher level: it tells the Go command to use several local modules together without adding local `replace` directives to each module.
- Keep `go.work` mainly for local development. Do not use it as a substitute for proper module versions in released code.

Example:

```sh
go work init ./app ./lib
```

This creates a workspace that includes the `app` and `lib` modules.

## Takeaways

- Start with a simple layout. Split packages when dependencies and responsibilities become clear.
- Put startup code in `cmd/...` when the module has more than one executable or when it keeps reusable code cleaner.
- Keep lower-level packages independent from higher-level packages. If two packages need the same helper, move it to a lower-level package that both can import.
- Export names intentionally. Once another package uses an exported name, changing it becomes harder.
- Do not treat unexported fields as secret data. They are an API boundary, not a security boundary.
- For a new library dependency, usually import it in code, then run `go mod tidy` and `go test ./...`.
- Use `go get module@version` when you intentionally add, upgrade, downgrade, or remove a dependency version.
- Check dependency updates explicitly with `go list -m -u all`. Then update deliberately, tidy, and test.
- Use `go install package@version` for executable tools, not `go get`.
- Use `go mod why`, `go mod graph`, and `go list -m all` when a dependency appears unexpectedly.
- Configure `GOPRIVATE` for private module paths before working with private dependencies locally or in CI.
- Use `go.work` for active local development across several modules. Prefer real module versions for released code.
- When importing a `v2+` module, include the major version suffix in the import path.

## Try it out

### Exercise 1: package organization inside one module

Create this structure:

```txt
modules/
├── cmd/
│   └── server/
│       └── main.go
├── store/
│   └── store.go
└── httpapi/
    └── handler.go
```

1. Put the executable startup code in `cmd/server/main.go`.
2. Put a reusable type in `store/store.go`.
3. Import `store` from `cmd/server` using the module path.
4. Add one exported function and one unexported helper in `store`.
5. Try to call both from `cmd/server` and observe which one is visible.
6. Make `httpapi` import `store`.
7. Then try to make `store` import `httpapi`. Observe the circular import error.
8. Fix it by moving shared code into a lower-level package, for example `internal/ids` or `internal/validate`.

### Exercise 2: daily dependency workflow with a third-party module

Use `github.com/google/uuid`. It is a small common package used to create and parse UUIDs.

1. From the chapter directory/module root, check the current module state:

```sh
go list -m all
```

2. In `store/store.go`, import the package and use it:

```go
package store

import "github.com/google/uuid"

type User struct {
    ID    string
    Email string
}

func NewUser(email string) User {
    return User{
        ID:    uuid.NewString(),
        Email: email,
    }
}
```

3. Run tests or build all packages:

```sh
go test ./...
```

If the dependency is missing, Go resolves it and updates module metadata.

4. Clean the module files:

```sh
go mod tidy
```

5. Inspect the result:

```sh
git diff go.mod go.sum
go list -m all
```

6. Upgrade the dependency explicitly:

```sh
go get github.com/google/uuid@latest
go mod tidy
```

7. Downgrade to a specific older version:

```sh
go get github.com/google/uuid@v1.3.0
go mod tidy
```

8. Return to the latest version:

```sh
go get github.com/google/uuid@latest
go mod tidy
```

9. Remove the dependency from the code by deleting the import and replacing the ID generation with a simple placeholder:

```go
func NewUser(email string) User {
    return User{
        ID:    "temporary-id",
        Email: email,
    }
}
```

10. Run tidy again and inspect the diff:

```sh
go mod tidy
git diff go.mod go.sum
```

11. Optional comparison: add the dependency again, then remove it with `go get @none`:

```sh
go get github.com/google/uuid@latest
go get github.com/google/uuid@none
go mod tidy
```

Compare this with removing the import and running `go mod tidy`.

### Exercise 3: tool dependency vs library dependency

Use this only as a comparison point.

- `github.com/google/uuid` is a library dependency. It is imported by your code and belongs in `go.mod` when used.
- A command-line tool is different. It is installed with `go install package@version`.

Example:

```sh
go install golang.org/x/tools/cmd/stringer@latest
```

This installs an executable tool. It is not the same as importing a library package in application code.

### Optional workspace exercise

1. Create two small modules side by side: `app` and `lib`.
2. Make `app` import `lib`.
3. First connect `app` to local `lib` with a `replace` directive.
4. Then remove the `replace` directive and use:

```sh
go work init ./app ./lib
```

5. Compare how the local dependency is resolved in both cases.

