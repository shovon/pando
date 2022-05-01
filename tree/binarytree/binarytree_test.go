package binarytree

import (
	"sync"
	"testing"
)

func compareTree(t *testing.T, m map[string]int, tree *BinaryTree) {
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

	cardinality := tree.Cardinality()
	if tree.Cardinality() != len(m) {
		t.Errorf("Expected %d, but got %d", len(m), cardinality)
	}
}

func TestInsert(t *testing.T) {
	tree := &BinaryTree{}

	m := map[string]int{
		"hello":   1,
		"world":   2,
		"foo":     3,
		"bar":     4,
		"baz":     5,
		"widgets": 6,
	}

	tree.Insert("hello", 1)
	tree.Insert("world", 2)
	tree.Insert("foo", 3)
	tree.Insert("bar", 4)
	tree.Insert("baz", 5)
	tree.Insert("widgets", 6)

	compareTree(t, m, tree)
}

func TestIterate(t *testing.T) {
	tree := &BinaryTree{}

	tree.Insert("hello", 1)
	tree.Insert("world", 2)
	tree.Insert("foo", 3)
	tree.Insert("bar", 4)
	tree.Insert("baz", 5)
	tree.Insert("widgets", 6)

	expected := map[string]int{
		"hello":   1,
		"world":   2,
		"foo":     3,
		"bar":     4,
		"baz":     5,
		"widgets": 6,
	}

	m := map[string]int{}

	for node := range tree.Iterate() {
		key, ok := node.Key.(string)
		if !ok {
			t.FailNow()
		}
		value, ok := node.Value.(int)
		if !ok {
			t.FailNow()
		}
		m[key] = value
	}

	for key, value := range expected {
		if expected[key] != value {
			t.Errorf("Expected value at key %s to be %d, but got %d", key, expected[key], value)
		}
	}
}

func TestDelete(t *testing.T) {
	tree := &BinaryTree{}

	m := map[string]int{
		"hello":   1,
		"world":   2,
		"foo":     3,
		"widgets": 6,
	}

	tree.Insert("hello", 1)
	tree.Insert("world", 2)
	tree.Insert("foo", 3)
	tree.Insert("bar", 4)
	tree.Insert("baz", 5)
	tree.Insert("widgets", 6)

	tree.Delete("bar")
	tree.Delete("baz")

	compareTree(t, m, tree)
}

func TestTreeEmpty(t *testing.T) {
	tree := &BinaryTree{}

	tree.Insert("cool", 1)
	tree.Insert("nice", 1)
	tree.Insert("amazing", 1)
	tree.Insert("sweet", 1)

	tree.Delete("cool")
	tree.Delete("nice")
	tree.Delete("amazing")
	tree.Delete("sweet")

	if !tree.IsEmpty() {
		t.Errorf("Expected the tree to be determined to be empty, but it was not!")
	}

	cardinality := tree.Cardinality()
	if tree.Cardinality() != 0 {
		t.Errorf("Expected 0, but got %d", cardinality)
	}
}

func TestInsertDelete(t *testing.T) {
	tree := &BinaryTree{}

	tree.Insert("cool", 1)
	tree.Delete("cool")
	tree.Insert("nice", 1)
	tree.Insert("amazing", 1)
	tree.Delete("nice")
	tree.Insert("sweet", 1)

	tree.Delete("amazing")
	tree.Delete("sweet")

	tree.Insert("foo", 2)

	m := map[string]int{
		"foo": 2,
	}

	compareTree(t, m, tree)
}

func TestFind(t *testing.T) {
	tree := &BinaryTree{}

	tree.Insert("hello", 1)
	tree.Insert("world", 2)
	tree.Insert("foo", 3)
	tree.Insert("bar", 4)
	tree.Insert("baz", 5)
	tree.Insert("widgets", 6)

	m := map[string]int{
		"hello":   1,
		"world":   2,
		"foo":     3,
		"bar":     4,
		"baz":     5,
		"widgets": 6,
	}

	state, ok := tree.Find("widgets")
	if !ok {
		t.Error("Expected to find a node with key \"widgetes\", but got nothing")
	}

	if state.Node.Key != "widgets" {
		t.Errorf("Expected node with key \"widgets\" but got %s", state.Node.Key)
	}

	if state.Node.Value != 6 {
		t.Errorf("Expected node with value 6, but got %d", state.Node.Value)
	}

	compareTree(t, m, tree)
}

func TestEvent(t *testing.T) {
	tree := &BinaryTree{}

	listener := tree.RegisterChangeListener("hello")

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		<-listener
		wg.Done()
	}()

	tree.Insert("hello", 1)
	tree.Insert("world", 2)

	m := map[string]int{
		"hello": 1,
		"world": 2,
	}

	wg.Wait()

	compareTree(t, m, tree)
}

func TestUpdate(t *testing.T) {
	tree := &BinaryTree{}

	// listener := tree.RegisterChangeListener("hello")

	tree.Insert("hello", 1)
	tree.Insert("world", 1)
	tree.Insert("foo", 1)
	tree.Insert("bar", 1)
	tree.Insert("baz", 1)
	tree.Insert("widgets", 1)

	tree.UpdateValue("hello", 42)

	state, ok := tree.Find("hello")
	if !ok {
		t.Error("Expected to fidn a node with key \"hello\", but hot nothing")
	}

	if state.Node.Key != "hello" {
		t.Errorf("Expected node with key \"hello\" but got %s", state.Node.Key)
	}

	if state.Node.Value != 42 {
		t.Errorf("Expected node with value 42, but got %d", state.Node.Value)
	}

	state, ok = tree.Find("world")
	if !ok {
		t.Error("Expected to fidn a node with key \"world\", but hot nothing")
	}

	if state.Node.Key != "world" {
		t.Errorf("Expected node with key \"world\" but got %s", state.Node.Key)
	}

	if state.Node.Value != 1 {
		t.Errorf("Expected node with value 1, but got %d", state.Node.Value)
	}

	m := map[string]int{
		"hello":   42,
		"world":   1,
		"foo":     1,
		"bar":     1,
		"baz":     1,
		"widgets": 1,
	}

	compareTree(t, m, tree)
}

func TestUpsert(t *testing.T) {
	tree := &BinaryTree{}

	tree.Insert("hello", 1)
	tree.Upsert("something", 20)

	compareTree(t, map[string]int{
		"hello":     1,
		"something": 20,
	}, tree)

	tree.Upsert("something", 42)

	compareTree(t, map[string]int{
		"hello":     1,
		"something": 42,
	}, tree)
}
