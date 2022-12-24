// This is a useless

package main

import (
	"encoding/json"
	"fmt"
	"tree/graph"
	"tree/reactforcegraph"
	"tree/set"
	insertiontree "tree/treegraph"
)

func main() {
	tree := insertiontree.Node{Neighbors: []*graph.Node{}, Key: "cool", Value: 1}

	tree.Upsert("foo", 2, 3, set.Set{})
	tree.Upsert("bar", 2, 3, set.Set{})
	tree.Upsert("baz", 2, 3, set.Set{})
	tree.Upsert("foobar", 2, 3, set.Set{})
	tree.Upsert("widgets", 2, 3, set.Set{})
	tree.Upsert("gadgets", 2, 3, set.Set{})
	tree.Upsert("hello", 2, 3, set.Set{})
	tree.Upsert("world", 2, 3, set.Set{})
	tree.Upsert("sweet", 2, 3, set.Set{})

	fgraph := reactforcegraph.ReactForceGraphMarshaler(graph.Node(tree).AdjacencyList(set.Set{}))

	b, err := json.Marshal(fgraph)

	if err != nil {
		panic(err)
	}

	fmt.Println(string(b))
}
