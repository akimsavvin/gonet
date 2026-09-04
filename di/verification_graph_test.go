// 🔥 GoNet is the first full-fledged framework made for Golang!
// ⚡️ GoNet is inspired by .NET, NestJS and other languages frameworks
// 🤖 GitHub Repository: https://github.com/akimsavvin/gonet

package di

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/suite"
)

// DependencyGraphSuite is the suite for testing the dependencyGraph construction
type DependencyGraphSuite struct {
	suite.Suite
}

// nodesOfType returns the graph nodes registered for the provided identifier type
func (suite *DependencyGraphSuite) nodesOfType(g *dependencyGraph, typ reflect.Type) []*serviceAccessor {
	var res []*serviceAccessor
	for _, node := range g.nodes {
		if node.id.Type == typ {
			res = append(res, node)
		}
	}

	return res
}

// TestMultipleRegistrationsAreNodes tests both registrations of one type are
// present as distinct nodes
func (suite *DependencyGraphSuite) TestMultipleRegistrationsAreNodes() {
	// Arrange
	c := newContainer(
		WithService[testRepository](func() testRepository { return &testRepositoryImpl{} }),
		WithService[testRepository](func() testRepository { return &testRepositoryImpl{} }),
	)

	// Act
	g := newDependencyGraph(c)
	nodes := suite.nodesOfType(g, reflect.TypeFor[testRepository]())

	// Assert
	suite.Len(nodes, 2)
	suite.NotSame(nodes[0], nodes[1])
}

// TestSliceDependencyFansOutToEveryRegistration tests a slice dependency edge
// targets every accessor registered for the element type
func (suite *DependencyGraphSuite) TestSliceDependencyFansOutToEveryRegistration() {
	// Arrange
	c := newContainer(
		WithService[testRepository](func() testRepository { return &testRepositoryImpl{} }),
		WithService[testRepository](func() testRepository { return &testRepositoryImpl{} }),
		WithFactory(func(repos []testRepository) *testCollector { return &testCollector{repos} }),
	)

	// Act
	g := newDependencyGraph(c)
	collector := suite.nodesOfType(g, reflect.TypeFor[*testCollector]())
	edges := g.edges[collector[0]]

	// Assert
	suite.Len(collector, 1)
	suite.Len(edges, 1)
	suite.Equal(reflect.TypeFor[[]testRepository](), edges[0].DependencyType)
	suite.Len(edges[0].Targets, 2)
}

// TestInstanceNodeHasNoEdges tests an accessor registered with an instance has
// no outgoing edges
func (suite *DependencyGraphSuite) TestInstanceNodeHasNoEdges() {
	// Arrange
	c := newContainer(WithValue[testRepository](&testRepositoryImpl{}))

	// Act
	g := newDependencyGraph(c)
	nodes := suite.nodesOfType(g, reflect.TypeFor[testRepository]())

	// Assert
	suite.Len(nodes, 1)
	suite.Empty(g.edges[nodes[0]])
}

// TestMissingDependencyHasNoTargets tests an unregistered dependency produces an
// edge with no targets
func (suite *DependencyGraphSuite) TestMissingDependencyHasNoTargets() {
	// Arrange
	c := newContainer(
		WithFactory(func(repo testRepository) *testService { return &testService{repo} }),
	)

	// Act
	g := newDependencyGraph(c)
	nodes := suite.nodesOfType(g, reflect.TypeFor[*testService]())
	edges := g.edges[nodes[0]]

	// Assert
	suite.Len(edges, 1)
	suite.Empty(edges[0].Targets)
}

// TestNodeOrderIsDeterministic tests the node order does not depend on the
// accessors map iteration order
func (suite *DependencyGraphSuite) TestNodeOrderIsDeterministic() {
	// Arrange
	build := func() []reflect.Type {
		c := newContainer(
			WithService[testRepository](func() testRepository { return &testRepositoryImpl{} }),
			WithService[testLogger](func() testLogger { return &testLoggerImpl{} }),
			WithFactory(func(repo testRepository) *testService { return &testService{repo} }),
			WithKeyedService[testRepository]("keyed", func() testRepository { return &testRepositoryImpl{} }),
		)

		g := newDependencyGraph(c)
		types := make([]reflect.Type, len(g.nodes))
		for i, node := range g.nodes {
			types[i] = node.id.Type
		}

		return types
	}

	// Act
	first := build()

	// Assert
	for range 20 {
		suite.Equal(first, build())
	}
}

// TestDependencyGraph tests the dependencyGraph construction
func TestDependencyGraph(t *testing.T) {
	suite.Run(t, new(DependencyGraphSuite))
}
