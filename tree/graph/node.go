package graph

import (
	"sync"
	"tree/adjacencylist"
	"tree/set"
)

// Node represents a node in a graph.
//
// This structure does not care the least about what the graph represents.
// It could be a directed tree. Directed acyclic lattice. Undirected tree.
// Full mesh, etc. It should not matter. This structure does not care.
//
// Other entities can, however, use this structure to enforce its use in a
// particular graph configuration, but this structure itself is not responsible
// for anything, other than to provide basic primitive that graph theorists may
// depend on
type Node struct {
	Neighbors []*Node
	Key       interface{}
	Value     interface{}
}

// Pair representing a tuple containing the key and value that would typically
// be held by a node
type Pair struct {
	Key   interface{}
	Value interface{}
}

// Cleave cleaves the node, and returns its neighbors, and a set of keys
// associated with all modified nodes.
func (n *Node) Cleave() ([]*Node, set.Set) {
	// Trivially, all neighbors will be modified, so grab their key set
	modified := n.GetNeighborKeys()

	// Let's not forget the newly cleaved node as well
	modified.Add(n.Key)

	// Grab the neigbors
	neighbors := n.Neighbors

	for _, neighbor := range neighbors {
		// Remove references of the cleaved node from all neighboring nodes
		neighbor.Neighbors =
			ExcludeNodesByKeys(neighbor.Neighbors, set.New(n.Key))
	}

	// The cleaved node is now a lone node
	n.Neighbors = []*Node{}

	// Return the neighbors and the set of modified nodes
	return neighbors, modified
}

// Interject will take a slice of distinct graphs (each pointed to by an
// orginating node), and attach them
func (n *Node) Interject(newNeighbors []*Node) set.Set {
	s := set.New(n.Key)
	for _, neighbor := range newNeighbors {
		s.Add(neighbor.Key)
		neighbor.Neighbors = append(n.Neighbors, n)
		n.Neighbors = append(n.Neighbors, neighbor)
	}
	return s
}

// Gets the keys of the neighboring nodes
func (n Node) GetNeighborKeys() set.Set {
	keys := set.Set{}

	for _, neighbor := range n.Neighbors {
		keys.Add(neighbor.Key)
	}

	return keys
}

// AdjacencyList performs a breadth-first-search of the node and creates an
// adjacency list of all neighboring nodes of all nodes that the BFS encountered
func (n Node) AdjacencyList(visited set.Set) adjacencylist.AdjacencyList {
	list := adjacencylist.AdjacencyList{}
	list[n.Key] = adjacencylist.AdjacencyListNode{
		Value:     n.Value,
		Neighbors: n.GetNeighborKeys(),
	}

	for _, neighbor := range n.Neighbors {
		if !visited.Has(neighbor.Key) {
			list = list.Union(neighbor.AdjacencyList(visited.Union(set.New(n.Key))))
		}
	}

	return list
}

// Traverse iterates through all nodes in the graph, on a depth-first-search
// basis, ensuring to avoid traversing the same node more than once
func (n Node) Traverse(visited set.Set) <-chan Pair {
	visited.Add(n.Key)

	c := make(chan Pair, 3)

	go func() {
		// Emit the current node to the listener
		c <- Pair{n.Key, n.Value}

		// Iterate through all the neighbors
		for _, neighbor := range n.Neighbors {

			// Check to see if we have not visited the neighbor yet
			if !visited.Has(neighbor.Key) {

				// If not, add it to our visited list
				visited.Add(neighbor.Key)

				// Do a traversal on the neighbouring subtree
				for child := range neighbor.Traverse(visited) {
					c <- child
				}
			}
		}

		close(c)
	}()

	return c
}

// ToSlice gets all key/value pairs in the graph, by traversing each node via
// neighboring nodes on a DFS basis
func (n Node) ToSlice() []Pair {
	result := []Pair{}

	for c := range n.Traverse(set.Set{}) {
		result = append(result, c)
	}

	return result
}

// Cardinality gets the count of all reachable nodes, through a DFS traversal
func (n Node) Cardinality() int {
	return len(n.ToSlice())
}

// Find gets the value associated with the supplied key
func (n Node) Find(key interface{}) (interface{}, bool) {
	for pair := range n.Traverse(set.Set{}) {
		if pair.Key == key {
			return pair.Value, true
		}
	}

	return nil, false
}

// GetNodesFromSet gets all nodes, that have keys that is in the keys set
func (n *Node) GetNodesFromSet(keys set.Set) []*Node {
	c := make(chan []*Node, 3)
	var wg sync.WaitGroup

	for _, node := range n.Neighbors {
		wg.Add(1)
		go func(node *Node) {
			defer wg.Done()
			c <- node.GetNodesFromSet(keys)
		}(node)
	}

	wg.Wait()

	close(c)

	result := []*Node{}

	for slice := range c {
		result = append(result, slice...)
	}

	if keys.Has(n.Key) {
		result = append(result, n)
	}

	return result
}

// Has determines whether a node with the supplied key exists in the graph
func (n Node) Has(key interface{}) bool {
	_, ok := n.Find(key)
	return ok
}

func (n Node) GetMap() map[interface{}]interface{} {
	m := map[interface{}]interface{}{}

	for c := range n.Traverse(set.Set{}) {
		m[c.Key] = c.Value
	}

	return m
}
