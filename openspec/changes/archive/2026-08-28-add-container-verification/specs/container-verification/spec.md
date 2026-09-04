## Purpose

Validates a dependency injection container's registration graph at construction time, so that a structurally broken container is rejected immediately with an actionable error instead of failing later during service resolution or deadlocking the process.

## ADDED Requirements

### Requirement: Container construction reports graph faults

Container construction SHALL return an error value alongside the container. When the registration graph is sound, construction SHALL return a usable container and a nil error. When the graph contains one or more faults, construction SHALL return a non-nil error.

This is a breaking change to the container constructor's signature.

#### Scenario: Sound graph constructs successfully

- **WHEN** a container is constructed with registrations whose every dependency is satisfiable
- **THEN** construction returns a usable container and a nil error
- **AND** service resolution from that container behaves exactly as before this change

#### Scenario: Faulty graph is rejected at construction

- **WHEN** a container is constructed with registrations containing at least one graph fault
- **THEN** construction returns a non-nil error describing the fault

#### Scenario: Empty container is sound

- **WHEN** a container is constructed with no registrations
- **THEN** construction returns a usable container and a nil error

### Requirement: Missing dependency detection

Construction SHALL verify that every dependency required by every registered factory resolves to a registered service, transitively across the whole graph. A dependency that does not resolve SHALL be reported as a fault.

A reported missing-dependency fault SHALL identify both the type of the service that requires the dependency and the type of the dependency that could not be resolved.

Dependency lookup during verification SHALL use the same matching rules as runtime resolution, so that any dependency verification accepts is resolvable at runtime and any dependency verification rejects would fail at runtime.

#### Scenario: Direct missing dependency

- **WHEN** a registered factory requires a type that is not registered
- **THEN** construction returns an error identifying the requiring service type and the missing dependency type

#### Scenario: Transitive missing dependency

- **WHEN** a registered factory's dependency is itself satisfiable, but that dependency's own factory requires an unregistered type
- **THEN** construction returns an error identifying the fault at the level where the dependency is actually missing

#### Scenario: Registered instances satisfy dependencies

- **WHEN** a factory requires a type that was registered as a value or instance rather than a factory
- **THEN** that dependency is treated as satisfied and produces no fault

#### Scenario: The container itself is always available

- **WHEN** a registered factory requires the container type
- **THEN** that dependency is treated as satisfied and produces no fault

#### Scenario: Keyed registrations do not satisfy unkeyed dependencies

- **WHEN** a type is registered only under a key, and a factory requires that type without a key
- **THEN** construction reports the dependency as missing, matching runtime resolution behavior

### Requirement: Slice dependency verification

A dependency declared as a slice of a type SHALL be treated as satisfied when at least one registration exists for the element type, and SHALL be reported as missing when no registration exists for the element type. This mirrors runtime resolution, which reports an unregistered slice element type as a not-found service rather than producing an empty slice.

#### Scenario: Slice dependency with registrations

- **WHEN** a factory requires a slice of a type and one or more implementations of that type are registered
- **THEN** the dependency is treated as satisfied and produces no fault

#### Scenario: Slice dependency with no registrations

- **WHEN** a factory requires a slice of a type and no implementation of that type is registered
- **THEN** construction reports the slice element type as a missing dependency

### Requirement: Circular dependency detection

Construction SHALL detect dependency cycles among registered factories and report each cycle as a fault. Construction SHALL NOT block, hang, or terminate the process when a cycle is present.

A reported cycle fault SHALL name the ordered sequence of service types forming the cycle, so the caller can locate every participant.

#### Scenario: Two-service cycle

- **WHEN** two registered factories each require the other's produced type
- **THEN** construction returns an error naming both types as a cycle
- **AND** construction returns rather than deadlocking

#### Scenario: Longer cycle

- **WHEN** three or more registered factories form a dependency cycle
- **THEN** construction returns an error naming the ordered sequence of types in that cycle

#### Scenario: Self-dependency

- **WHEN** a registered factory requires the type it produces
- **THEN** construction reports it as a cycle

#### Scenario: Shared dependency is not a cycle

- **WHEN** two or more registered factories require the same dependency, and that dependency does not depend back on either
- **THEN** no cycle is reported

### Requirement: Aggregated fault reporting

A single construction attempt SHALL report every fault found in the graph, not only the first. The returned error SHALL expose the individual faults so a caller can inspect them programmatically, and SHALL render a message listing each fault when formatted as text.

#### Scenario: Multiple missing dependencies

- **WHEN** a container is constructed with three unrelated missing dependencies
- **THEN** the single returned error reports all three faults

#### Scenario: Mixed fault kinds

- **WHEN** a container is constructed with both a missing dependency and a dependency cycle
- **THEN** the single returned error reports both faults

#### Scenario: Faults are programmatically inspectable

- **WHEN** a caller receives a construction error and inspects it with the standard error-unwrapping helpers
- **THEN** the caller can distinguish a missing-dependency fault from a cycle fault and read the types involved in each

### Requirement: Verification does not instantiate services

Verification SHALL be a static analysis of the registration graph. It SHALL NOT invoke any registered factory, and SHALL NOT create, cache, or memoize any service instance.

Lazy singleton semantics SHALL be preserved: after successful construction, a service instance is still created on first resolution and reused thereafter.

#### Scenario: Factories are not called during construction

- **WHEN** a container is constructed with factories that record whether they were invoked
- **THEN** no factory has been invoked once construction returns

#### Scenario: Factory errors are not surfaced at construction

- **WHEN** a registered factory would return an error, but the graph is otherwise sound
- **THEN** construction succeeds with a nil error
- **AND** the factory's error surfaces on first resolution of that service, as before this change

#### Scenario: First resolution still creates the instance

- **WHEN** a service is resolved for the first time from a successfully constructed container
- **THEN** its factory is invoked at that moment
- **AND** a second resolution of the same service returns the same instance without invoking the factory again
