// 🔥 GoNet is the first full-fledged framework made for Golang!
// ⚡️ GoNet is inspired by .NET, NestJS and other languages frameworks
// 🤖 GitHub Repository: https://github.com/akimsavvin/gonet

package di

import "reflect"

type (
	// serviceDependencyTreeNode describes a service dependency tree node
	serviceDependencyTreeNode struct {
		// Type is the type of the dependency
		Type reflect.Type

		// Parent is the node parent
		Parent *serviceDependencyTreeNode

		// Deps is the slice of dependency nodes
		Deps []*serviceDependencyTreeNode
	}

	// serviceDependencyTree describes a service dependency tree
	serviceDependencyTree struct {
		// Root is the root node of the tree
		Root *serviceDependencyTree
	}
)
