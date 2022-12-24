package adjacencylist

import (
	"fmt"
	"tree/maybeany"
	"tree/set"
)

// AdjacencyListNode represents a pairing from source node to a list of target
// nodes
type AdjacencyListNode struct {
	Value     interface{}
	Neighbors set.Set
}

// UnionLinks creates a new AdjacencyListNode with links added to the pre-
// existing set of links
func (a AdjacencyListNode) UnionLinks(links set.Set) AdjacencyListNode {
	return AdjacencyListNode{Value: a.Value, Neighbors: links.Union(a.Neighbors)}
}

// UnionLinks creates a new AdjacencyListNode with the link added to the pre-
// existing set of links
func (a AdjacencyListNode) AddLink(link interface{}) AdjacencyListNode {
	return AdjacencyListNode{Value: a.Value, Neighbors: a.Neighbors.Union(set.New(link))}
}

// AdjacencyList represents a mapping of all nodes to other links in the network
type AdjacencyList map[interface{}]AdjacencyListNode

// Sets an adjacency link to the AdjacencyList
func (a *AdjacencyList) SetLink(key, value interface{}) {
	link, ok := (*a)[key]
	if !ok {
		(*a)[key] = AdjacencyListNode{value, set.Set{}}
		return
	}
	(*a)[key] = AdjacencyListNode{value, link.Neighbors}
}

func (a *AdjacencyList) AddLinks(key interface{}, link set.Set, defaultValue interface{}) {
	l, ok := (*a)[key]
	if !ok {
		(*a)[key] = AdjacencyListNode{defaultValue, link}
	} else {
		(*a)[key] = AdjacencyListNode{l.Value, l.Neighbors.Union(link)}
	}
}

// GetReversed gets the graph represented by the adjacencylist, but with the edges
// reversed
func (a AdjacencyList) GetReversed() AdjacencyList {
	newList := AdjacencyList{}

	// Iterate through the original list of ndoes
	for key, node := range a {

		// Then, for every link, point it to key
		for link := range node.Neighbors {
			newList.AddLinks(link, set.New(key), node.Value)
		}

	}

	return newList
}

// Union combines two adjacency lists into a single graph. Especially useful if
// combined with `Reversed` to drive an bidirectional (undirected) graph.
func (a AdjacencyList) Union(b AdjacencyList) AdjacencyList {
	newList := AdjacencyList{}

	for key, node := range a {
		newList.AddLinks(key, node.Neighbors, node.Value)
	}

	for key, node := range b {
		newList.AddLinks(key, node.Neighbors, node.Value)
	}

	return newList
}

// Equal tests the equality of two adjacency lists. The graph not only must be
// homomorphic from each other, but they must key-by-key (e.g. it is not
// sufficient for two graphs to have the same shape, but each node must have
// the same keys, and must direct to the same keys)
func (a AdjacencyList) Equal(a1 AdjacencyList) bool {
	if len(a) != len(a1) {
		return false
	}

	for key, n := range a {
		node, ok := a1[key]
		if !ok {
			return false
		}

		if !n.Neighbors.Equals(node.Neighbors) {
			return false
		}
	}

	return true
}

// GetKeys gets the list of keys in the graph
func (a AdjacencyList) GetKeys() set.Set {
	result := set.Set{}
	for k := range a {
		result.Add(k)
	}
	return result
}

func (a AdjacencyList) GetAnyKey() maybeany.MaybeAny {
	for k := range a {
		return maybeany.Something(k)
	}

	return maybeany.Nothing()
}

type Pair struct {
	Key   interface{}
	Value interface{}
}

func (a AdjacencyList) Traverse(currentNode interface{}, visited set.Set) <-chan Pair {
	c := make(chan Pair)

	go func() {
		defer close(c)

		node, ok := a[currentNode]
		if !ok {
			return
		}

		visited.Add(currentNode)
		c <- Pair{currentNode, node.Value}

		// Iterate through each of the keys of the neighboring nodes
		for neighborKey := range node.Neighbors {
			if !visited.Has(neighborKey) {
				next := a.Traverse(neighborKey, visited)

				for pair := range next {
					c <- pair
				}
			}
		}
	}()

	return c
}

func GetKeysFromTraversal(traversal <-chan Pair) set.Set {
	s := set.Set{}

	for pair := range traversal {
		s.Add(pair.Key)
	}

	fmt.Println("The visited set", s)

	return s
}
