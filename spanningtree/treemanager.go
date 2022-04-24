package main

import (
	"spanningtree/spanningtree"
	"sync"
)

type treeManager struct {
	mut   sync.RWMutex
	trees map[string]*spanningtree.Tree
}

func newTreeManager() treeManager {
	return treeManager{trees: make(map[string]*spanningtree.Tree)}
}

func (t *treeManager) getTree(id string) *spanningtree.Tree {

	tree, ok := t.trees[id]
	if !ok {
		tree = &spanningtree.Tree{}
		trees[id] = tree
	}

	return tree
}

func (t *treeManager) insertNode(treeId string, nodeId string) {

}

func (t *treeManager) deleteNode(treeId string, nodeId string) {

}
