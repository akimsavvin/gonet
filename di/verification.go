// 🔥 GoNet is the first full-fledged framework made for Golang!
// ⚡️ GoNet is inspired by .NET, NestJS and other languages frameworks
// 🤖 GitHub Repository: https://github.com/akimsavvin/gonet

package di

import (
	"fmt"
	"reflect"
	"slices"
	"strings"
)

// MissingDependencyError is a container verification fault reporting a factory
// dependency that is not registered in the Container.
//
// It wraps ErrServiceNotFound, so errors.Is(err, ErrServiceNotFound) reports
// true for a verification error containing this fault.
type MissingDependencyError struct {
	// RequestingType is the type of the service requiring the dependency
	RequestingType reflect.Type

	// DependencyType is the type of the dependency that is not registered
	DependencyType reflect.Type
}

// Error implements the error interface for MissingDependencyError.
func (e *MissingDependencyError) Error() string {
	return fmt.Sprintf(
		"service %q requires %q, which is not registered",
		e.RequestingType,
		e.DependencyType,
	)
}

// Unwrap returns ErrServiceNotFound, the sentinel a failed runtime resolution
// of the same dependency would return.
func (e *MissingDependencyError) Unwrap() error {
	return ErrServiceNotFound
}

// CircularDependencyError is a container verification fault reporting a cycle
// among the factory dependencies of the Container.
type CircularDependencyError struct {
	// Cycle is the ordered sequence of service types forming the cycle.
	// Each type depends on the next one, and the last one depends on the first
	Cycle []reflect.Type
}

// Error implements the error interface for CircularDependencyError.
func (e *CircularDependencyError) Error() string {
	names := make([]string, len(e.Cycle)+1)
	for i, typ := range e.Cycle {
		names[i] = fmt.Sprintf("%v", typ)
	}

	// Close the cycle so the message shows the dependency returning to its start
	if len(e.Cycle) > 0 {
		names[len(e.Cycle)] = fmt.Sprintf("%v", e.Cycle[0])
	} else {
		names = names[:0]
	}

	return fmt.Sprintf("circular dependency: %s", strings.Join(names, " -> "))
}

// VerificationError is the error returned by NewContainer when the registration
// graph contains one or more faults.
//
// It aggregates every fault found in a single verification pass. Faults are
// reachable with errors.As and errors.Is, which traverse all of them.
type VerificationError struct {
	// Faults is every fault found while verifying the Container.
	// Each entry is a *MissingDependencyError or a *CircularDependencyError
	Faults []error
}

// Error implements the error interface for VerificationError.
func (e *VerificationError) Error() string {
	var sb strings.Builder
	sb.WriteString("di: container verification failed:")

	for _, fault := range e.Faults {
		sb.WriteString("\n  - ")
		sb.WriteString(fault.Error())
	}

	return sb.String()
}

// Unwrap returns every aggregated fault, so errors.As and errors.Is inspect all
// of them.
func (e *VerificationError) Unwrap() []error {
	return e.Faults
}

// nodeType returns the type used to name a dependencyGraph node in a fault.
//
// A node created from a factory is named by the factory return type, the same
// type the resolution path reports in a DependencyError. A node created from an
// instance is named by the type it was registered for.
func nodeType(accessor *serviceAccessor) reflect.Type {
	if accessor.factory != nil {
		return accessor.factory.ReturnType
	}

	return accessor.id.Type
}

// dependencyEdge is a single factory parameter of a dependencyGraph node,
// resolved to the accessors serving it
type dependencyEdge struct {
	// DependencyType is the declared parameter type, before the slice reduction
	DependencyType reflect.Type

	// Targets are the accessors serving the dependency.
	// Empty if the dependency is not registered
	Targets []*serviceAccessor
}

// dependencyGraph is the dependency graph of the Container registrations.
//
// Nodes are service accessors rather than service identifiers, because a single
// identifier can hold several registrations and a slice dependency instantiates
// every one of them.
type dependencyGraph struct {
	// nodes is every registered accessor, in a deterministic order
	nodes []*serviceAccessor

	// edges are the outgoing dependency edges of every node with a factory
	edges map[*serviceAccessor][]dependencyEdge
}

// compareServiceIdentifiers orders two service identifiers, so that a graph
// built from the Container accessors map does not depend on the map iteration
// order
func compareServiceIdentifiers(a, b serviceIdentifier) int {
	if res := strings.Compare(a.Type.String(), b.Type.String()); res != 0 {
		return res
	}

	if res := strings.Compare(a.Type.PkgPath(), b.Type.PkgPath()); res != 0 {
		return res
	}

	if a.HasKey != b.HasKey {
		if a.HasKey {
			return 1
		}

		return -1
	}

	return strings.Compare(a.Key, b.Key)
}

// newDependencyGraph builds the dependencyGraph of the provided Container
func newDependencyGraph(c *Container) *dependencyGraph {
	ids := make([]serviceIdentifier, 0, len(c.accessors))
	for id := range c.accessors {
		ids = append(ids, id)
	}

	slices.SortFunc(ids, compareServiceIdentifiers)

	g := &dependencyGraph{
		nodes: make([]*serviceAccessor, 0, len(ids)),
		edges: make(map[*serviceAccessor][]dependencyEdge, len(ids)),
	}

	for _, id := range ids {
		for _, accessor := range c.accessors[id].Iter() {
			g.nodes = append(g.nodes, accessor)
		}
	}

	for _, node := range g.nodes {
		// An accessor registered with an instance has no dependencies
		if node.factory == nil {
			continue
		}

		edges := make([]dependencyEdge, node.factory.DepsCount)
		for i := range edges {
			depType := node.factory.Type.In(i)
			edge := dependencyEdge{
				DependencyType: depType,
			}

			// Resolve the parameter exactly the way the resolution path does,
			// so the graph cannot disagree with runtime resolution
			lookupID, _ := newDependencyLookup(depType)
			if accessors, ok := c.accessors[lookupID]; ok {
				edge.Targets = make([]*serviceAccessor, 0, accessors.Len())
				for _, target := range accessors.Iter() {
					edge.Targets = append(edge.Targets, target)
				}
			}

			edges[i] = edge
		}

		g.edges[node] = edges
	}

	return g
}

// nodeState is the traversal state of a dependencyGraph node
type nodeState uint8

const (
	// nodeUnvisited marks a node the walk has not reached yet
	nodeUnvisited nodeState = iota

	// nodeInProgress marks a node on the current traversal path.
	// An edge reaching such a node closes a cycle
	nodeInProgress

	// nodeDone marks a fully expanded node, which the walk never revisits
	nodeDone
)

// walkFrame is a dependencyGraph node being expanded by the walk
type walkFrame struct {
	// node is the accessor being expanded
	node *serviceAccessor

	// edgeIdx is the index of the edge being expanded
	edgeIdx int

	// targetIdx is the index of the next target of the current edge
	targetIdx int
}

// newCycleFault creates a CircularDependencyError for the provided traversal
// path segment.
//
// The cycle is rotated to start at its lowest named type, so the same cycle
// yields the same fault no matter where the walk entered it. Returns nil if an
// equal cycle was already recorded in seen.
func newCycleFault(cycle []*serviceAccessor, seen map[string]struct{}) *CircularDependencyError {
	if len(cycle) == 0 {
		return nil
	}

	start := 0
	names := make([]string, len(cycle))
	for i, accessor := range cycle {
		names[i] = nodeType(accessor).String()
		if names[i] < names[start] {
			start = i
		}
	}

	types := make([]reflect.Type, 0, len(cycle))
	for _, accessor := range cycle[start:] {
		types = append(types, nodeType(accessor))
	}
	for _, accessor := range cycle[:start] {
		types = append(types, nodeType(accessor))
	}

	rotated := make([]string, 0, len(types))
	for _, typ := range types {
		rotated = append(rotated, typ.String())
	}

	key := strings.Join(rotated, " -> ")
	if _, ok := seen[key]; ok {
		return nil
	}

	seen[key] = struct{}{}

	return &CircularDependencyError{
		Cycle: types,
	}
}

// walk traverses the dependencyGraph depth first and returns every fault found.
//
// Nodes are marked unvisited, in progress and done, so each node is expanded
// exactly once and a shared dependency reports its faults once. An edge to an
// in-progress node closes a cycle, which is read off the current path.
func (g *dependencyGraph) walk() []error {
	var faults []error

	states := make(map[*serviceAccessor]nodeState, len(g.nodes))
	pathIdx := make(map[*serviceAccessor]int, len(g.nodes))
	seenCycles := make(map[string]struct{})

	path := make([]*serviceAccessor, 0, len(g.nodes))
	stack := make([]walkFrame, 0, len(g.nodes))

	// enter marks the node in progress, records its missing dependencies and
	// pushes it onto the traversal stack
	enter := func(node *serviceAccessor) {
		states[node] = nodeInProgress
		pathIdx[node] = len(path)
		path = append(path, node)
		stack = append(stack, walkFrame{node: node})

		for _, edge := range g.edges[node] {
			if len(edge.Targets) == 0 {
				faults = append(faults, &MissingDependencyError{
					RequestingType: nodeType(node),
					DependencyType: edge.DependencyType,
				})
			}
		}
	}

	for _, root := range g.nodes {
		if states[root] != nodeUnvisited {
			continue
		}

		enter(root)

		for len(stack) > 0 {
			frame := &stack[len(stack)-1]
			edges := g.edges[frame.node]

			// Every edge of the node is expanded, so the node is done
			if frame.edgeIdx >= len(edges) {
				states[frame.node] = nodeDone
				delete(pathIdx, frame.node)
				path = path[:len(path)-1]
				stack = stack[:len(stack)-1]
				continue
			}

			edge := edges[frame.edgeIdx]
			if frame.targetIdx >= len(edge.Targets) {
				frame.edgeIdx++
				frame.targetIdx = 0
				continue
			}

			target := edge.Targets[frame.targetIdx]
			frame.targetIdx++

			switch states[target] {
			case nodeInProgress:
				if fault := newCycleFault(path[pathIdx[target]:], seenCycles); fault != nil {
					faults = append(faults, fault)
				}
			case nodeUnvisited:
				// enter may grow the stack, so frame must not be used after it
				enter(target)
			}
		}
	}

	return faults
}

// verify checks the Container registration graph and returns a
// *VerificationError aggregating every fault found, or nil if the graph is
// sound.
//
// Verification is a static analysis of the registrations: it never calls a
// factory and never creates a service instance.
func (c *Container) verify() error {
	faults := newDependencyGraph(c).walk()
	if len(faults) == 0 {
		return nil
	}
	return &VerificationError{
		Faults: faults,
	}
}
