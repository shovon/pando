package treegraph

import (
	"encoding/json"
	"testing"
	"tree/adjacencylist"
	"tree/graph"
	"tree/maybeany"
	"tree/reactforcegraph"
	"tree/set"
)

func compareTree(t *testing.T, m map[string]int, node *Node) {
	treeChan := (*graph.Node)(node).Traverse(set.Set{})
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

	cardinality := (*graph.Node)(node).Cardinality()
	if (*graph.Node)(node).Cardinality() != len(m) {
		t.Errorf("Expected the graph to have size %d, but got %d", len(m), cardinality)
	}
}

func TestShortestPathEmpty(t *testing.T) {
	node := &Node{
		[]*graph.Node{},
		"foo",
		10,
	}

	path := node.ShortestPath(set.Set{})

	if len(path) != 1 {
		t.Logf("Expected the path to have 1 node, but actually got %d nodes", len(path))
		t.Fail()
	}
}

func TestShortestPathOneNeighbor(t *testing.T) {
	node := &Node{
		[]*graph.Node{
			{
				Neighbors: []*graph.Node{},
				Key:       "bar",
				Value:     1,
			},
		},
		"foo",
		10,
	}

	path := node.ShortestPath(set.Set{})

	if len(path) != 2 {
		t.Logf("Expected the path to have 2 nodes, but actually got %d nodes", len(path))
		t.Fail()
	}
}

func TestTraversal(t *testing.T) {
	tree := &Node{[]*graph.Node{}, "hello", 1}

	insert := func(key string, value int) {
		tree.Upsert(key, value, 3, set.Set{})
	}

	insert("hello", 1)
	insert("world", 2)
	insert("foo", 3)
	insert("bar", 4)
	insert("baz", 5)
	insert("widgets", 6)

	treeToSlice := (*graph.Node)(tree).ToSlice()
	treeToMap := (*graph.Node)(tree).GetMap()

	if len(treeToSlice) != len(treeToMap) {
		t.Logf("Expected both the map and slice to have the same length, but slice was %d, and map was %d", len(treeToSlice), len(treeToMap))
		t.Fail()
	}
}

func TestInsert(t *testing.T) {
	m := map[string]int{
		"hello":   1,
		"world":   2,
		"foo":     3,
		"bar":     4,
		"baz":     5,
		"widgets": 6,
	}

	tree := Node{[]*graph.Node{}, "hello", 1}

	insert := func(key string, value int) {
		tree.Upsert(key, value, 3, set.Set{})
	}

	insert("hello", 1)
	insert("world", 2)
	insert("foo", 3)
	insert("bar", 4)
	insert("baz", 5)
	insert("widgets", 6)

	compareTree(t, m, &tree)
}

func TestIterate(t *testing.T) {
	tree := Node{[]*graph.Node{}, "hello", 1}

	insert := func(key string, value int) {
		tree.Upsert(key, value, 3, set.Set{})
	}

	insert("world", 2)
	insert("foo", 3)
	insert("bar", 4)
	insert("baz", 5)
	insert("widgets", 6)

	expected := map[string]int{
		"hello":   1,
		"world":   2,
		"foo":     3,
		"bar":     4,
		"baz":     5,
		"widgets": 6,
	}

	m := map[string]int{}

	for node := range (graph.Node)(tree).Traverse(set.Set{}) {
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
	maybeTree := Something(&Node{[]*graph.Node{}, "hello", 1})

	expectedPairs := map[string]int{
		"hello":   1,
		"world":   2,
		"foo":     3,
		"widgets": 6,
	}

	insert := func(key string, value int) {
		if tree, ok := maybeTree.GetNode(); ok {
			tree.Upsert(key, value, 3, set.Set{})
		} else {
			maybeTree = Something(&Node{[]*graph.Node{}, key, value})
		}
	}

	deleteByKey := func(key string) set.Set {
		if tree, ok := maybeTree.GetNode(); ok {
			newTree, modifiedNodes := tree.DeleteByKey(key, set.Set{})
			maybeTree = newTree
			return modifiedNodes
		}
		return set.Set{}
	}

	insert("hello", 1)
	insert("world", 2)
	insert("foo", 3)
	insert("bar", 4)
	insert("baz", 5)
	insert("widgets", 6)

	deletedKeys := deleteByKey("bar")
	if !deletedKeys.Has("bar") {
		t.Error("The set of modified nodes should have had 'bar', but it doesn't!")
	}
	deletedKeys = deleteByKey("baz")
	if !deletedKeys.Has("baz") {
		t.Error("The set of modified nodes should have had 'baz', but it doesn't!")
	}

	tree, ok := maybeTree.GetNode()
	if !ok {
		t.Fail()
	}

	compareTree(t, expectedPairs, tree)
}

func TestTreeEmpty(t *testing.T) {
	maybeTree := Something(&Node{[]*graph.Node{}, "cool", 1})

	insert := func(key string, value int) {
		if tree, ok := maybeTree.GetNode(); ok {
			tree.Upsert(key, value, 3, set.Set{})
		} else {
			maybeTree = Something(&Node{[]*graph.Node{}, key, value})
		}
	}

	deleteByKey := func(key string) {
		if tree, ok := maybeTree.GetNode(); ok {
			maybeTree, _ = tree.DeleteByKey(key, set.Set{})
		}
	}

	insert("cool", 1)
	insert("nice", 1)

	insert("amazing", 1)

	insert("sweet", 1)

	deleteByKey("cool")

	deleteByKey("nice")

	deleteByKey("amazing")

	deleteByKey("sweet")

	if _, ok := maybeTree.GetNode(); ok {
		t.Errorf("Expected the tree to be determined to be empty, but it was not!")
	}
}

func TestInsertDelete(t *testing.T) {
	maybeTree := Something(&Node{[]*graph.Node{}, "cool", 1})

	insert := func(key string, value int) {
		if tree, ok := maybeTree.GetNode(); ok {
			tree.Upsert(key, value, 3, set.Set{})
		} else {
			maybeTree = Something(&Node{[]*graph.Node{}, key, value})
		}
	}

	deleteByKey := func(key string) {
		if tree, ok := maybeTree.GetNode(); ok {
			maybeTree, _ = tree.DeleteByKey(key, set.Set{})
		}
	}

	deleteByKey("cool")
	insert("nice", 1)
	insert("amazing", 1)
	deleteByKey("nice")
	insert("sweet", 1)

	deleteByKey("amazing")
	deleteByKey("sweet")

	insert("foo", 2)

	m := map[string]int{
		"foo": 2,
	}

	tree, ok := maybeTree.GetNode()
	if !ok {
		t.Fail()
	}

	compareTree(t, m, tree)
}

func TestFind(t *testing.T) {
	tree := Node{[]*graph.Node{}, "hello", 1}

	insert := func(key string, value int) {
		tree.Upsert(key, value, 3, set.Set{})
	}

	insert("hello", 1)
	insert("world", 2)
	insert("foo", 3)
	insert("bar", 4)
	insert("baz", 5)
	insert("widgets", 6)

	m := map[string]int{
		"hello":   1,
		"world":   2,
		"foo":     3,
		"bar":     4,
		"baz":     5,
		"widgets": 6,
	}

	value, ok := graph.Node(tree).Find("widgets")
	if !ok {
		t.Error("Expected to find a node with key \"widgetes\", but got nothing")
	}

	if value != 6 {
		t.Errorf("Expected node with value 6, but got %d", value)
	}

	compareTree(t, m, &tree)
}

func TestUpsert(t *testing.T) {
	tree := Node{[]*graph.Node{}, "hello", 1}

	insert := func(key string, value int) {
		tree.Upsert(key, value, 3, set.Set{})
	}

	insert("hello", 1)
	insert("world", 2)
	insert("bar", 3)
	insert("baz", 4)
	insert("widgets", 5)
	insert("gadgets", 6)
	insert("something", 20)

	compareTree(t, map[string]int{
		"hello":     1,
		"world":     2,
		"bar":       3,
		"baz":       4,
		"widgets":   5,
		"gadgets":   6,
		"something": 20,
	}, &tree)

	insert("world", 42)

	compareTree(t, map[string]int{
		"hello":     1,
		"world":     42,
		"bar":       3,
		"baz":       4,
		"widgets":   5,
		"gadgets":   6,
		"something": 20,
	}, &tree)
}

func DFS(list adjacencylist.AdjacencyList, startingKey interface{}, visited set.Set) set.Set {
	v := set.Set{}.Union(visited)

	node, ok := list[startingKey]
	if !ok {
		return v
	}

	v.Add(startingKey)

	for _, link := range node.Neighbors {
		if !v.Has(link) {
			visits := DFS(list, link, v)
			v = v.Union(visits)
		}
	}

	return v
}

func getKeys(m map[interface{}]adjacencylist.AdjacencyListNode) set.Set {
	s := set.Set{}
	for k := range m {
		s.Add(k)
	}
	return s
}

func IsGraphConnected(list adjacencylist.AdjacencyList) bool {
	first := list.GetAnyKey()
	key, ok := first.Get()
	if !ok {
		return false
	}

	traversalKeys := adjacencylist.GetKeysFromTraversal(list.Union(list.GetReversed()).Traverse(key, set.Set{}))
	return traversalKeys.Equals(list.GetKeys())
}

func IsGraphUndirected(list adjacencylist.AdjacencyList) bool {
	return list.GetReversed().Equal(list)
}

func PathHasCycle(list adjacencylist.AdjacencyList, startingKey interface{}, visited set.Set, maybePredecessor maybeany.MaybeAny) bool {
	if visited.Has(startingKey) {
		return true
	}

	visits := set.Set{}
	visits.Add(startingKey)

	node, ok := list[startingKey]
	if !ok {
		return false
	}

	for link := range node.Neighbors {
		predecessor, ok := maybePredecessor.Get()
		if !ok || link != predecessor {
			if PathHasCycle(list, link, visits, maybeany.Something(startingKey)) {
				return true
			}
		}
	}

	return false
}

func GraphHasCycle(list adjacencylist.AdjacencyList) bool {
	for k := range list {
		if PathHasCycle(list, k, set.Set{}, maybeany.Nothing()) {
			return true
		}
	}
	return false
}

func IsTree(list adjacencylist.AdjacencyList) bool {
	return !GraphHasCycle(list) && IsGraphConnected(list)
}

func TestAdjacencyListTree(t *testing.T) {
	tree := Node{[]*graph.Node{}, "hello", 1}

	insert := func(key string, value int) {
		tree.Upsert(key, value, 3, set.Set{})
	}

	insert("foo", 1)
	insert("bar", 2)
	insert("baz", 3)
	insert("foobar", 4)

	actualList := graph.Node(tree).AdjacencyList(set.Set{})

	if !IsTree(actualList) {
		t.Error("Expected the graph to be a tree, but ended up not being a tree")
	}
}

func TestAdjacencyList(t *testing.T) {
	maybeTree := Nothing()

	insert := func(key string, value byte) {
		if tree, ok := maybeTree.GetNode(); ok {
			tree.Upsert(key, value, 3, set.Set{})
		} else {
			maybeTree = Something(&Node{[]*graph.Node{}, key, value})
		}
	}

	expectedMapping := map[interface{}]interface{}{}

	var i, j byte
	for i = 97; i <= 122; i++ {
		for j = 97; j <= 122; j++ {
			str := string([]byte{i, j})
			val := (i-97)*(122-97) + (j - 97)
			insert(str, val)
			expectedMapping[str] = val
		}
	}

	tree, ok := maybeTree.GetNode()
	if !ok {
		t.Fail()
	}

	actualMapping := (*graph.Node)(tree).GetMap()

	logForceGraph := func() {
		j, err := json.Marshal(reactforcegraph.ReactForceGraphMarshaler((*graph.Node)(tree).AdjacencyList(set.Set{})))
		if err != nil {
			t.Error(err)
		}
		t.Log(string(j))
	}

	if len(actualMapping) != len(expectedMapping) {
		t.Logf("Expected length of mapping to be %d, but got %d", len(expectedMapping), len(actualMapping))
		logForceGraph()
		t.Fail()
	}

	for k := range expectedMapping {
		if actualMapping[k] != expectedMapping[k] {
			t.Logf("Expected %d for key %s, but got %d", expectedMapping[k], k, actualMapping[k])
			found, ok := (*graph.Node)(tree).Find(k)
			if ok {
				t.Logf("%d", found)
			} else {
				t.Logf("Node with key %s not found", k)
			}
			t.Fail()
		}
	}

	list := (*graph.Node)(tree).AdjacencyList(set.Set{})
	if len(list) != (*graph.Node)(tree).Cardinality() {
		t.Errorf("Expected the number of keys in the adjacency to be %d (the tree's cardinality) but got %d", len(list), (*graph.Node)(tree).Cardinality())
	}

	if !IsGraphUndirected(list) {
		t.Errorf("The graph is not an undirected graph")
	}

	if !IsTree(list) {
		t.Logf("The adjacency list is a graph, but not a tree")
		if GraphHasCycle(list) {
			t.Logf("It seems because the list represents a graph with cycles in it")
		}
		if !IsGraphConnected(list) {
			t.Logf("It seems because the list represents disjointed graph")
		}
		logForceGraph()
		t.Fail()
	}
}
