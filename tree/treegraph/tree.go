package treegraph

import (
	"tree/adjacencylist"
	"tree/graph"
	"tree/set"
)

type Tree struct {
	maybeRoot MaybeNode
}

func (t *Tree) Upsert(key, value interface{}) set.Set {
	n, ok := t.maybeRoot.GetNode()
	if !ok {
		t.maybeRoot = Something(&Node{[]*graph.Node{}, key, value})
		return set.New(key)
	}

	return n.Upsert(key, value, 3, set.Set{})
}

func (t *Tree) DeleteByKey(key interface{}) set.Set {
	n, ok := t.maybeRoot.GetNode()
	if !ok {
		return set.Set{}
	}

	maybeRoot, modifiedNodes := n.DeleteByKey(key, set.Set{})
	t.maybeRoot = maybeRoot
	return modifiedNodes
}

func (t Tree) Find(key interface{}) (interface{}, bool) {
	return t.maybeRoot.Find(key)
}

func (t Tree) Has(key interface{}) bool {
	return t.maybeRoot.Has(key)
}

func (t Tree) AdjacencyList() adjacencylist.AdjacencyList {
	return t.maybeRoot.AdjacencyList()
}

func (t Tree) IsEmpty() bool {
	_, ok := t.maybeRoot.GetNode()
	return !ok
}
