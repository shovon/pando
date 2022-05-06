package binarytree

import (
	"encoding/json"
	"sync"
	"tree/binarytree/node"
	"tree/listeners"
)

// BinaryTree represents a spanning tree data structure.

// All public operations on this tree (including reads) are thread-safe
type BinaryTree struct {
	mut  sync.RWMutex
	root *node.Node

	listeners             listeners.KeyedListeners
	ValueMarshalerCreator genericMarshalerCreator
	KeyMarshalerCreator   genericMarshalerCreator
}

var _ json.Marshaler = &BinaryTree{}

// MarshalJSON marshals the tree to a JSON representation
func (t *BinaryTree) MarshalJSON() ([]byte, error) {
	t.mut.RLock()
	defer t.mut.RUnlock()

	if t.root == nil {
		return json.Marshal(nil)
	}

	valueMarshalerCreator := t.ValueMarshalerCreator
	if valueMarshalerCreator == nil {
		valueMarshalerCreator = func(value interface{}) json.Marshaler {
			return &defaultMarshaler{value}
		}
	}

	keyMarshalerCreator := t.KeyMarshalerCreator
	if keyMarshalerCreator == nil {
		keyMarshalerCreator = func(value interface{}) json.Marshaler {
			return &defaultMarshaler{value}
		}
	}

	return marshalNode(t.root, keyMarshalerCreator, valueMarshalerCreator)
}

type marshalableNode struct {
	Key   json.RawMessage `json:"key"`
	Value json.RawMessage `json:"value"`
	Left  json.RawMessage `json:"left"`
	Right json.RawMessage `json:"right"`
}

func marshalNode(
	node *node.Node,
	keyMarshalerCreator genericMarshalerCreator,
	valueMarshalerCreator genericMarshalerCreator,
) (json.RawMessage, error) {
	if node == nil {
		return json.Marshal(nil)
	}

	keyMarshalable := keyMarshalerCreator(node.Key())
	valueMarshalable := valueMarshalerCreator(node.Value())

	keyMarshalled, e := keyMarshalable.MarshalJSON()
	if e != nil {
		return json.RawMessage{}, e
	}

	valueMarshalled, e := valueMarshalable.MarshalJSON()
	if e != nil {
		return json.RawMessage{}, e
	}

	leftMarshalled, e := marshalNode(
		node.Left(),
		keyMarshalerCreator,
		valueMarshalerCreator,
	)
	if e != nil {
		return json.RawMessage{}, e
	}

	rightMarshalled, e := marshalNode(
		node.Right(),
		keyMarshalerCreator,
		valueMarshalerCreator,
	)
	if e != nil {
		return json.RawMessage{}, e
	}

	return json.Marshal(marshalableNode{
		keyMarshalled,
		valueMarshalled,
		leftMarshalled,
		rightMarshalled,
	})
}

// RegisterChangeListener creates a channel that serves as the event listener
func (t *BinaryTree) RegisterChangeListener(key interface{}) <-chan interface{} {
	return t.listeners.RegisterListener(key)
}

// Insert inserts a key/value par into the tree
func (t *BinaryTree) Insert(key, value interface{}) {
	t.mut.Lock()
	defer t.mut.Unlock()

	t.unsafeInsert(key, value)
}

func (t *BinaryTree) unsafeInsert(key, value interface{}) {
	node := node.NewNode(key, value)
	defer t.emitChangeEvent()

	if t.root == nil {
		t.root = node
		return
	}

	t.root.Insert(node)
}

func (t *BinaryTree) UpdateValue(key, value interface{}) bool {
	t.mut.Lock()
	defer t.mut.Unlock()

	if t.root == nil {
		return false
	}

	updated := t.root.UpdateValue(key, value)

	if updated {
		t.emitChangeEvent()
	}

	return updated
}

func (t *BinaryTree) Upsert(key, value interface{}) {
	t.mut.Lock()
	defer t.mut.Unlock()

	if t.root == nil {
		t.unsafeInsert(key, value)
		t.emitChangeEvent()
		return
	}

	node := t.root.Find(key)
	if node == nil {
		t.unsafeInsert(key, value)
		t.emitChangeEvent()
	} else {
		updated := t.root.UpdateValue(key, value)
		if !updated {
			panic("Not should have been updated, but it was not!")
		}
		t.emitChangeEvent()
	}
}

func (t *BinaryTree) Find(key interface{}) (NodeState, bool) {
	t.mut.RLock()
	defer t.mut.RUnlock()

	node := t.unsafeFind(key)
	if node == nil {
		return NodeState{}, false
	}

	return NewNodeState(node), true
}

func (t *BinaryTree) unsafeFind(key interface{}) *node.Node {
	if t.root == nil {
		return nil
	}

	return t.root.Find(key)
}

// Delete delets a node from the tree, given a key
func (t *BinaryTree) Delete(key interface{}) bool {
	t.mut.Lock()
	defer t.mut.Unlock()

	if t.root == nil {
		return false
	}

	if t.root.Key() == key {
		left := t.root.Left()
		right := t.root.Right()
		if left == nil {
			t.root = right
		} else {
			nodes := node.NewScatterer(right).Scatter()
			for node := range nodes {
				left.Insert(node)
			}
			t.emitChangeEvent()
			t.root = left
		}
		return true
	}

	deleted := t.root.Delete(key)
	if deleted {
		t.emitChangeEvent()
	}
	return deleted
}

func (t *BinaryTree) emitChangeEvent() {
	for node := range t.iterateUnsafe() {
		t.listeners.EmitEvent(node.Key(), NewNodeState(node))
	}
}

func (t *BinaryTree) iterateSafe() <-chan *node.Node {
	c := make(chan *node.Node)

	go func() {
		t.mut.RLock()
		defer t.mut.RUnlock()
		defer close(c)

		if t.root == nil {
			return
		}

		for node := range t.root.Iterate() {
			c <- node
		}
	}()

	return c
}

func (t *BinaryTree) iterateUnsafe() <-chan *node.Node {
	c := make(chan *node.Node)

	go func() {
		defer close(c)

		if t.root == nil {
			return
		}

		for node := range t.root.Iterate() {
			c <- node
		}
	}()

	return c
}

// Iterate iterates all key-value pairs in the tree
func (t *BinaryTree) Iterate() <-chan node.Pair {
	t.mut.RLock()
	defer t.mut.RUnlock()

	c := make(chan node.Pair)

	if t.root == nil {
		close(c)
		return c
	}

	go func() {
		t.mut.RLock()
		defer t.mut.RUnlock()

		for n := range t.root.Iterate() {
			pair := node.Pair{Key: n.Key(), Value: n.Value()}
			c <- pair
		}

		close(c)
	}()

	return c
}

func (t *BinaryTree) GetAdjacencyList() []node.AdjacencyNode {
	t.mut.RLock()
	defer t.mut.RUnlock()

	if t.root == nil {
		return []node.AdjacencyNode{}
	}

	return t.root.GetAdjacencyList()
}

func (t *BinaryTree) Cardinality() int {
	t.mut.RLock()
	defer t.mut.RUnlock()

	if t.root == nil {
		return 0
	}

	return t.root.Cardinality()
}

func (t *BinaryTree) IsEmpty() bool {
	t.mut.RLock()
	defer t.mut.RUnlock()

	return t.root == nil
}
