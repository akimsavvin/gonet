## Context

See `proposal.md` — Why. The structural facts that shape the approach:

- `Container.accessors` is `map[serviceIdentifier]*serviceAccessorsList`. An identifier is `{Type, Key, HasKey}`; a list holds one or more `*serviceAccessor` in registration order.
- An accessor holds either a `*serviceFactory` or a pre-built instance value. An instance-backed accessor has no dependencies.
- `serviceFactory` already captures everything an edge needs: `Type` (the `reflect.Type` of the function), `DepsCount`, `ReturnType`. Parameter types are `factory.Type.In(i)`.
- Runtime dependency lookup lives in `Container.resolveFactoryDeps`, which builds `serviceIdentifier{Type: depType}` — unkeyed — and calls `getService`. `getService` strips a slice type to its element type before the map lookup.
- All options are applied inside `NewContainer` before it returns, and `NewContainer` appends `WithValue(c)` so the container registers itself. The full graph is therefore known at the end of construction, and nothing can register afterwards.
- v2 is still pre-release (`v2.0.0-rc.2`), so a breaking constructor signature can land before the stable tag.

## Goals / Non-Goals

**Goals:**

- Verification and runtime resolution agree by construction, not by parallel maintenance.
- Report the whole broken graph in one pass, with enough type information to locate each fault.
- Zero effect on resolution semantics and no factory invocation during verification.
- Linear cost in registrations plus dependency edges.

**Non-Goals:**

- Changing what `getService` does at runtime. Verification mirrors it; it does not correct it.
- Converting the existing `log.Panicf` registration-shape checks into errors. Those fire inside `Option.apply`, before a graph exists.
- Detecting faults that are not decidable from the registration set (a factory that returns nil, a factory that panics, a factory that returns an error).

## Decisions

### Graph nodes are accessors, not identifiers

An identifier can map to several accessors, and runtime treats them differently: single resolution takes `Last()`, slice resolution instantiates all of them. If nodes were identifiers, a broken registration reachable only through `[]T` would never be checked, and two registrations of the same type could not be distinguished in a cycle report.

So each `*serviceAccessor` is a node, keyed by pointer identity. An edge from an accessor is: for each factory parameter, resolve the parameter to an identifier, then fan out to **every** accessor in that identifier's list.

*Alternative considered:* nodes as identifiers. Smaller graph, but it silently skips non-last registrations — exactly the case `WithService[T]` registered twice is designed to support.

### One shared function builds a dependency identifier

The requirement that verification and runtime never disagree is enforced by extracting the parameter-to-identifier mapping — the unkeyed identifier plus the slice-element stripping currently split across `resolveFactoryDeps` and `getService` — into a single helper used by both paths.

This is the load-bearing decision. If verification re-derived lookup rules, the two would drift on the next change to keying or slice handling, and verification would start passing containers that fail at runtime (or vice versa). It also means the known keyed-dependency gap is mirrored rather than accidentally papered over: verification reports an unkeyed parameter to a keyed-only registration as missing, which is precisely what runtime does.

*Alternative considered:* a separate lookup written for verification. Rejected — drift risk is the whole point of the feature.

### Depth-first traversal with three-state marking

Walk every accessor with an iterative DFS carrying an explicit path stack, marking nodes unvisited / in-progress / done. An edge to an in-progress node is a back-edge and yields a cycle: the path stack from that node to the current node is the cycle, in order. A node marked done is skipped, so a shared dependency is walked once and its faults are reported once.

*Alternative considered:* Kahn's algorithm / topological sort. It detects that a cycle exists but does not hand back the participating path, and the spec requires naming the ordered sequence. DFS produces the path as a by-product of the stack.

Iterative over recursive: a deep chain should not risk the goroutine stack, and an explicit stack is what makes the cycle path available anyway.

### Dedicated fault types that stay compatible with `ErrServiceNotFound`

Two fault types:

- a missing-dependency fault carrying the requesting type and the missing dependency type, wrapping the existing `ErrServiceNotFound` sentinel so `errors.Is(err, di.ErrServiceNotFound)` keeps working;
- a cycle fault carrying the ordered `[]reflect.Type` of the cycle.

*Alternative considered:* reusing the existing `DependencyError` for missing dependencies. It already carries both types and unwraps to a cause, and callers may know it — but its message reads "failed to create dependency … for service …", which is false for a static check that creates nothing. Keeping `DependencyError` as the runtime-only error and adding a construction-time type keeps both messages honest. Wrapping the shared sentinel preserves the `errors.Is` path that matters most.

### Aggregate type with `Unwrap() []error`

A verification error holds `[]error` of faults, implements `Unwrap() []error` so `errors.As` and `errors.Is` traverse into every fault, and renders a header plus one line per fault.

*Alternative considered:* plain `errors.Join`. Same unwrapping behavior for free, but the message is a bare newline-joined list with no header and no way to attach a count or ordering. A thin named type costs little and gives a stable, greppable message shape.

### Cycle reporting is canonicalized and deduplicated

A cycle is discovered once per back-edge, but the same cycle can be entered from different starting nodes across the outer loop. Cycles are rotated to a canonical starting point and deduplicated before being added to the fault list, so a three-service cycle is reported once, not three times.

### Verification runs at the end of `NewContainer`, and failure yields a nil container

Verification is the last step, after every option is applied — including the self-registration of `*Container`. On failure `NewContainer` returns `(nil, err)`.

*Alternative considered:* returning the half-valid container alongside the error, letting a caller resolve the sound parts. Rejected — it invites ignoring the error, and the whole change exists to make that hard.

## Risks / Trade-offs

- **Every call site breaks, including four examples, the README quickstart, and the container test suite.** → v2 is pre-release, so this lands before the stable tag rather than forcing a `/v3` module path. The compiler flags every site; the migration is mechanical.
- **Verification is only as correct as the shared lookup helper.** A bug there now produces both a wrong runtime result and a wrong verification verdict, instead of one of them catching the other. → Tests assert the agreement property directly: for each fixture graph, verification's verdict and an actual resolution attempt must match.
- **Extracting the shared helper touches the hot resolution path.** → It is a pure refactor of identifier construction with no allocation change; existing resolution tests must pass unmodified apart from the constructor signature.
- **Reporting a keyed-only registration as a missing dependency will read as a bug to users** who registered the service and expect it to be injected. → The message names the dependency as unkeyed, and the limitation is stated in the proposal's non-goals; the underlying gap is tracked as separate work.
- **Slice dependencies with zero registrations now fail at construction** rather than at first resolution, which is a behavior change in timing for anyone relying on lazy discovery. → It matches the existing runtime verdict, so no previously working program starts failing; only the moment of failure moves earlier.
- **A container built by hand rather than through `NewContainer` skips verification.** → All construction paths go through `NewContainer`; `Container` has no exported fields and no exported way to register after construction.

## Open Questions

None blocking. The keyed-dependency resolution gap and the panic-versus-error treatment of registration-shape validation are both recorded as non-goals and sized as separate changes.
