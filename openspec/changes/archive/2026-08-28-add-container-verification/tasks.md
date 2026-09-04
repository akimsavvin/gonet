## 1. Shared dependency lookup

- [x] 1.1 Extract the parameter-to-identifier mapping into one helper in `di` — building the unkeyed `serviceIdentifier` from a factory parameter type, plus the slice-element stripping currently inside `getService` — returning the lookup identifier and whether the parameter was a slice; verify by `go build ./...` succeeding
- [x] 1.2 Rewrite `Container.resolveFactoryDeps` and `Container.getService` to call that helper instead of constructing identifiers inline; verify the existing `di` test suite passes unchanged with `go test ./di/...`

## 2. Verification fault types

- [x] 2.1 Add a missing-dependency fault type carrying the requesting `reflect.Type` and the missing dependency `reflect.Type`, wrapping `ErrServiceNotFound`; verify a unit test asserts `errors.Is(fault, ErrServiceNotFound)` and that both types are readable from the value
- [x] 2.2 Add a circular-dependency fault type carrying the ordered `[]reflect.Type` of the cycle; verify a unit test asserts the ordered types round-trip and the message names every participant
- [x] 2.3 Add the aggregate verification error holding `[]error` with `Unwrap() []error` and a header-plus-one-line-per-fault message; verify a unit test asserts `errors.As` finds both fault kinds inside it and that the message lists every fault

## 3. Graph traversal

- [x] 3.1 Build the graph over `*serviceAccessor` nodes keyed by pointer identity, fanning each factory parameter out through the task 1.1 helper to every accessor in the target identifier's list; verify a unit test on a container with two registrations of one type shows both registrations present as nodes
- [x] 3.2 Implement the iterative depth-first walk with unvisited / in-progress / done marking and an explicit path stack; verify a unit test on a sound multi-level graph visits every accessor exactly once
- [x] 3.3 Emit a missing-dependency fault when a parameter's identifier has no registered accessors, covering direct, transitive, slice-element, and keyed-only-registration cases; verify unit tests for each of those four cases
- [x] 3.4 Emit a circular-dependency fault from the path stack on a back-edge to an in-progress node, covering two-service, three-or-more-service, and self-dependency cycles; verify unit tests for each and that a shared non-cyclic dependency produces no fault
- [x] 3.5 Canonicalize cycles to a fixed starting point and deduplicate them before collecting; verify a unit test that a three-service cycle is reported exactly once regardless of which node the outer loop starts from
- [x] 3.6 Collect all faults from the full walk into one aggregate error and return nil when the graph is sound; verify unit tests covering three unrelated missing dependencies in one error, a mixed missing-plus-cycle graph, and an empty container returning nil

## 4. Constructor integration

- [x] 4.1 Change `NewContainer` to `NewContainer(opts ...Option) (*Container, error)`, running verification after every option is applied and after the container's `WithValue(c)` self-registration, returning `(nil, err)` on failure; verify `go vet ./di/...` passes and a test asserts a faulty graph yields a nil container with a non-nil error
- [x] 4.2 Verify no factory is invoked during construction — add a test whose factories record invocation and assert none ran after `NewContainer` returns
- [x] 4.3 Verify lazy singleton semantics survive — add a test asserting a factory runs on first resolution, returns the same instance on second resolution, and that a factory which returns an error still constructs successfully and surfaces that error only on resolution
- [x] 4.4 Add an agreement test asserting verification's verdict matches an actual resolution attempt for each fixture graph, so the two paths cannot drift; verify it passes for sound, missing-dependency, and keyed-only fixtures

## 5. Call site migration

- [x] 5.1 Update every `NewContainer` call site in `di/container_test.go` and `di/accessor_test.go` for the new signature; verify `go test ./di/...` passes
- [x] 5.2 Update `examples/simple/main.go`, `examples/factory_error/main.go`, `examples/multiple_implementations/main.go`, and `examples/value_providers/main.go` to handle the returned error; verify each runs to completion with `go run ./examples/<name>` and prints its documented output
- [x] 5.3 Update the `README.md` quickstart to the new constructor signature; verify the snippet compiles when pasted into a scratch main package

## 6. Full verification

- [x] 6.1 Run `go build ./... && go vet ./... && go test ./... -race` and confirm all pass
- [x] 6.2 Confirm the previously deadlocking case is fixed — construct a container with an `A -> B -> A` cycle and verify it returns a cycle error rather than producing `fatal error: all goroutines are asleep - deadlock!`
