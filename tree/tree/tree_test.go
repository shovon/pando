package tree

import (
	"sync"
	"testing"
)

func TestInsert(t *testing.T) {
	tree := &Tree{}

	listener := tree.RegisterChangeListener()

	var wg sync.WaitGroup
	wg.Add(6)

	go func() {
		for {
			switch (<-listener).(type) {
			case NodeState:
				wg.Done()
			default:
				t.Error("Expected a NodeState, but got something else")
			}
		}
	}()

	m := map[string]int{
		"hello":   1,
		"world":   2,
		"foo":     3,
		"bar":     4,
		"baz":     5,
		"widgets": 6,
	}

	tree.Insert(Pair{"hello", 1})
	tree.Insert(Pair{"world", 2})
	tree.Insert(Pair{"foo", 3})
	tree.Insert(Pair{"bar", 4})
	tree.Insert(Pair{"baz", 5})
	tree.Insert(Pair{"widgets", 6})

	treeChan := tree.Iterate()
	for pair := range treeChan {
		key, value := pair.Key, pair.Value
		keyStr, ok := key.(string)
		if !ok {
			t.Error("Key should have been a string, but was something else")
		}
		if v, ok := m[keyStr]; ok {
			if v != value {
				t.Errorf("Expected %d, but got %d", v, value)
			}
		} else {
			t.Errorf("Item of %s not found", keyStr)
		}
	}

	wg.Wait()
}

func TestDelete(t *testing.T) {
	tree := &Tree{}

	m := map[string]int{
		"hello":   1,
		"world":   2,
		"foo":     3,
		"widgets": 6,
	}

	var insertWg, deleteWg sync.WaitGroup

	listener := tree.RegisterChangeListener()

	go func() {
		for {
			switch (<-listener).(type) {
			case NodeState:
				insertWg.Done()
			case Deleted:
				deleteWg.Done()
			default:
				t.Error("Expected either a NodeState or Deleted, but got something else")
			}
		}
	}()

	insertWg.Add(6)
	deleteWg.Add(2)

	tree.Insert(Pair{"hello", 1})
	tree.Insert(Pair{"world", 2})
	tree.Insert(Pair{"foo", 3})
	tree.Insert(Pair{"bar", 4})
	tree.Insert(Pair{"baz", 5})
	tree.Insert(Pair{"widgets", 6})

	tree.Delete("bar")
	tree.Delete("baz")

	treeChan := tree.Iterate()
	for pair := range treeChan {
		key, value := pair.Key, pair.Value
		keyStr, ok := key.(string)
		if !ok {
			t.Error("Key should have been a string, but was something else")
		}
		if v, ok := m[keyStr]; ok {
			if v != value {
				t.Errorf("Expected %d, but got %d", v, value)
			}
		} else {
			t.Errorf("Item of %s not found", keyStr)
		}
	}

	insertWg.Wait()
	insertWg.Wait()
}
