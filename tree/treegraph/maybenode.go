package treegraph

import (
	"tree/adjacencylist"
	"tree/graph"
	"tree/set"
)

// MaybeNode represents somethig that could either be a Node or nothing.
//
// Rationale for creating this is to avoid null dereferencing when it is
// possible for something to be null
//
// Note: avoid using `MaybeNode{} == MaybeNode{}`!
// Use `MaybeNode{}.Equals(MaybeNode{})` instead!
type MaybeNode struct {
	node   *Node
	exists bool
}

// Something initializes a new MaybeNode that will definitely contain a
// node.Node
func Something(node *Node) MaybeNode {
	return MaybeNode{node, true}
}

// Nothing initializes a new MaybeNode that represents an empty node
func Nothing() MaybeNode {
	return MaybeNode{
		node:   nil,
		exists: false,
	}
}

// GetNode attempts to get the value in the MaybeNode.
//
// GetNode two values: the Node (or null), and a boolean, where true is if a value
// existed; false otherwise
func (m MaybeNode) GetNode() (*Node, bool) {
	return m.node, m.exists
}

// Equals determines if this MaybeNode equals to the one provided.
//
// Note: avoid using `MaybeNode{} == MaybeNode{}`!
func (m MaybeNode) Equals(b MaybeNode) bool {
	return (m.exists && b.exists) || (m.node == b.node)
}

func (m MaybeNode) Find(key interface{}) (interface{}, bool) {
	if !m.exists {
		return nil, false
	}

	return (*graph.Node)(m.node).Find(key)
}

func (m MaybeNode) Has(key interface{}) bool {
	if !m.exists {
		return false
	}

	return (*graph.Node)(m.node).Has(key)
}

func (m MaybeNode) AdjacencyList() adjacencylist.AdjacencyList {
	if !m.exists {
		return adjacencylist.AdjacencyList{}
	}
	return (*graph.Node)(m.node).AdjacencyList(set.Set{})
}
