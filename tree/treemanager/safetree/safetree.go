package safetree

import (
	"sync"
	"tree/adjacencylist"
	"tree/set"
	"tree/treegraph"
)

type SafeTree struct {
	mut  *sync.RWMutex
	tree treegraph.Tree
}

func New() SafeTree {
	mut := &sync.RWMutex{}
	return SafeTree{mut, treegraph.Tree{}}
}

func (t *SafeTree) Upsert(key, value interface{}) set.Set {
	t.mut.Lock()
	defer t.mut.Unlock()
	return t.tree.Upsert(key, value)
}

func (t *SafeTree) DeleteByKey(key interface{}) set.Set {
	t.mut.Lock()
	defer t.mut.Unlock()
	return t.tree.DeleteByKey(key)
}

func (t SafeTree) Find(key interface{}) (interface{}, bool) {
	t.mut.RLock()
	defer t.mut.RUnlock()
	return t.tree.Find(key)
}

func (t SafeTree) Has(key interface{}) bool {
	t.mut.RLock()
	defer t.mut.RUnlock()
	return t.tree.Has(key)
}

func (t SafeTree) AdjacencyList() adjacencylist.AdjacencyList {
	t.mut.RLock()
	defer t.mut.RUnlock()
	return t.tree.AdjacencyList()
}

func (t SafeTree) IsEmpty() bool {
	t.mut.RLock()
	defer t.mut.RUnlock()
	return t.tree.IsEmpty()
}
