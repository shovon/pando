package main

import (
	"sync"
	"tree/binarytree"
)

type treeManager struct {
	treeMut    sync.RWMutex
	managerMut sync.RWMutex
	trees      map[string]*binarytree.BinaryTree
}

func newTreeManager() treeManager {
	return treeManager{trees: make(map[string]*binarytree.BinaryTree)}
}

func (t *treeManager) getTree(id string) *binarytree.BinaryTree {
	t.managerMut.Lock()
	defer t.managerMut.Unlock()

	tree, ok := t.trees[id]
	if !ok {
		tree = &binarytree.BinaryTree{}
		t.trees[id] = tree
	}

	return tree
}

func (t *treeManager) upsert(treeId, nodeId string, p participant) {
	t.managerMut.Lock()
	defer t.managerMut.Unlock()

	tree := t.getTree(treeId)
	tree.Upsert(nodeId, p)
}

func (t *treeManager) update(treeId, nodeId string, p participant) {
	t.managerMut.Lock()
	defer t.managerMut.Unlock()

	tree := t.getTree(treeId)
	tree.UpdateValue(nodeId, p)
}

func (t *treeManager) insertNode(treeId, nodeId string, p participant) {
	t.treeMut.Lock()
	defer t.treeMut.Unlock()
	tree := t.getTree(treeId)
	tree.Insert(nodeId, p)
}

func (t *treeManager) deleteNode(treeId, nodeId string) {
	t.treeMut.Lock()
	defer t.treeMut.Unlock()
	tree := t.getTree(treeId)
	tree.Delete(nodeId)
	if tree.IsEmpty() {
		delete(t.trees, treeId)
	}
}
