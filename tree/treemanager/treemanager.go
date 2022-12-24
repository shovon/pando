package treemanager

import (
	"sync"
	"tree/treemanager/listeners"
	"tree/treemanager/safetree"
)

type treeManager struct {
	mut       *sync.RWMutex
	trees     map[string]*safetree.SafeTree
	listeners listeners.KeyedListeners
}

func NewTreeManager() treeManager {
	managerMut := &sync.RWMutex{}
	return treeManager{
		mut:   managerMut,
		trees: make(map[string]*safetree.SafeTree),
	}
}

func (t *treeManager) GetTree(id string) *safetree.SafeTree {
	t.mut.Lock()
	defer t.mut.Unlock()

	tree, ok := t.trees[id]
	if !ok {
		newTree := safetree.New()
		t.trees[id] = &newTree
		tree = &newTree
	}

	return tree
}

func (t *treeManager) Upsert(treeId, nodeId string, p interface{}) {
	t.mut.Lock()
	defer t.mut.Unlock()
	tree := t.GetTree(treeId)

	changedNodes := tree.Upsert(nodeId, p)
	t.listeners.EmitEvent(treeId, changedNodes)
}

func (t *treeManager) DeleteNode(treeId, nodeId string) {
	t.mut.Lock()
	defer t.mut.Unlock()
	tree := t.GetTree(treeId)

	changedNodes := tree.DeleteByKey(nodeId)
	if tree.IsEmpty() {
		delete(t.trees, treeId)
	}

	t.listeners.EmitEvent(treeId, changedNodes)
}

func (t *treeManager) RegisterChangeListener(
	treeId interface{},
) <-chan interface{} {
	return t.listeners.RegisterListener(treeId)
}

func (t *treeManager) UnregisterChangeListener(
	treeId interface{},
	listener <-chan interface{},
) {
	t.listeners.UnregisterListener(treeId, listener)
}
