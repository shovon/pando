package spanningtree

import (
	"encoding/json"
	"spanningtree/listeners"
	"spanningtree/spanningtree/node"
	"sync"
)

type Tree struct {
	mut  sync.RWMutex
	root *node.Node

	listeners             listeners.KeyedListeners
	ValueMarshalerCreator genericMarshalerCreator
	KeyMarshalerCreator   genericMarshalerCreator
}

var _ json.Marshaler = &Tree{}

func (t *Tree) MarshalJSON() ([]byte, error) {
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

func (t *Tree) RegisterChangeListener(key interface{}) <-chan interface{} {
	return t.listeners.RegisterListener(key)
}

func (t *Tree) Insert(pair Pair) {
	t.mut.Lock()
	defer t.mut.Unlock()

	node := node.NewNode(pair.Key, pair.Value)
	defer t.listeners.EmitEvent(node.Key(), NewNodeState(node))

	if t.root == nil {
		t.root = node
		return
	}

	t.root.Insert(node)
}

func (t *Tree) Delete(key interface{}) bool {
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
			t.listeners.EmitEventToAll(Deleted{})
			t.root = left
		}
		return true
	}

	deleted := t.root.Delete(key)
	if deleted {
		t.listeners.EmitEventToAll(Deleted{})
	}
	return deleted
}

func (t *Tree) emitEvent() {
	for node := range t.iterate() {
		t.listeners.EmitEvent(node.Key(), NewNodeState(node))
	}
}

func (t *Tree) iterate() <-chan *node.Node {
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

func (t *Tree) Iterate() <-chan Pair {
	t.mut.RLock()
	defer t.mut.RUnlock()

	c := make(chan Pair)

	if t.root == nil {
		close(c)
		return c
	}

	go func() {
		t.mut.RLock()
		defer t.mut.RUnlock()

		for node := range t.root.Iterate() {
			pair := Pair{node.Key(), node.Value()}
			c <- pair
		}

		close(c)
	}()

	return c
}

func (t *Tree) Cardinality() int {
	if t.root == nil {
		return 0
	}

	return t.root.Cardinality()
}

func (t *Tree) IsEmpty() bool {
	t.mut.RLock()
	defer t.mut.RUnlock()
	return t.root == nil
}