package spanningtree

import (
	"testing"
)

func TestInsert(t *testing.T) {
	tree := &Tree{}

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

	cardinality := tree.Cardinality()
	if tree.Cardinality() != 6 {
		t.Errorf("Expected 6, but got %d", cardinality)
	}
}

func TestDelete(t *testing.T) {
	tree := &Tree{}

	m := map[string]int{
		"hello":   1,
		"world":   2,
		"foo":     3,
		"widgets": 6,
	}

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

	cardinality := tree.Cardinality()
	if tree.Cardinality() != 4 {
		t.Errorf("Expected 4, but got %d", cardinality)
	}
}

func TestTreeEmpty(t *testing.T) {
	tree := &Tree{}

	tree.Insert(Pair{"cool", 1})
	tree.Insert(Pair{"nice", 1})
	tree.Insert(Pair{"amazing", 1})
	tree.Insert(Pair{"sweet", 1})

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
	tree := &Tree{}

	tree.Insert(Pair{"cool", 1})
	tree.Delete("cool")
	tree.Insert(Pair{"nice", 1})
	tree.Insert(Pair{"amazing", 1})
	tree.Delete("nice")
	tree.Insert(Pair{"sweet", 1})

	tree.Delete("amazing")
	tree.Delete("sweet")

	tree.Insert(Pair{"foo", 2})

	cardinality := tree.Cardinality()
	if tree.Cardinality() != 1 {
		t.Errorf("Expected 1, but got %d", cardinality)
	}
}
