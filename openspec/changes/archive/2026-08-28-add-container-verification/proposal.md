## Why

Today a `di.Container` can be structurally broken and stay silent. A missing dependency is discovered only when something resolves that specific branch, so a container assembled wrong at startup can run for hours before failing in a request path. A circular dependency is worse: `serviceAccessor.Instance` guards creation with `sync.Once`, which is not re-entrant, so `A -> B -> A` produces `fatal error: all goroutines are asleep - deadlock!` — an unrecoverable process kill, not an error a caller can handle.

Both faults are fully determinable from the registration set alone. Verifying the graph once at construction converts a late, non-local runtime failure into an immediate, actionable startup error.

## What Changes

- **BREAKING**: `di.NewContainer(opts ...Option) (*Container, error)` — the constructor now returns an error alongside the container. Every existing call site must be updated.
- `NewContainer` builds a dependency graph from all registered accessors after applying options, and walks it before returning.
- **Missing dependency detection**: every factory parameter must resolve to a registered service identifier, transitively. An unresolvable parameter is reported with both the requesting service type and the missing dependency type.
- **Circular dependency detection**: a cycle among factory dependencies is reported as an error naming the full cycle path, replacing today's deadlock.
- **Aggregated reporting**: one walk collects every problem in the graph and returns them joined, so a single run surfaces the whole broken graph rather than one fault per fix-and-rerun cycle.
- Slice dependencies (`[]T`) with zero registrations are reported as missing, mirroring current runtime behavior in `Container.getService`, which returns `ErrServiceNotFound` for an unregistered `[]T`.
- Verification is construction-time only and creates no service instances — factories are not called, so registration order and lazy-singleton semantics are unchanged.
- `examples/` (4 programs) and `README.md` updated for the new constructor signature.

## Capabilities

### New Capabilities
- `container-verification`: construction-time validation of the container dependency graph — missing dependency detection, circular dependency detection, and aggregated error reporting.

### Modified Capabilities
<!-- None. The project has no existing specs under openspec/specs/. -->

## Impact

**Public API (breaking)**
- `di.NewContainer` return type changes from `*Container` to `(*Container, error)`.
- New exported error types for verification failures so callers can inspect problems with `errors.As` / `errors.Is`.

**Affected code**
- `di/container.go` — `NewContainer` signature and the new post-registration verification pass.
- New file for graph construction and traversal.
- `di/container_test.go` — every `NewContainer` call site.
- `examples/simple`, `examples/factory_error`, `examples/multiple_implementations`, `examples/value_providers` — `main.go` in each.
- `README.md` — quickstart snippet.

**Not changed (explicit non-goals)**
- Keyed dependency resolution. `Container.resolveFactoryDeps` builds dependency identifiers with `HasKey` false, so keyed registrations are unreachable as factory parameters. Verification mirrors runtime lookup exactly and therefore reports such a parameter as missing; it does not fix the underlying gap. Tracked separately.
- Existing registration-shape validation. `newServiceFactory` and the option `apply` methods use `log.Panicf` for malformed factories and non-assignable types. These remain panics — they are option-construction errors detected before a graph exists, distinct from graph-level faults.
- Service lifetimes, scopes, and disposal.
