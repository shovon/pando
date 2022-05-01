package binarytree

import (
	"encoding/json"
	"tree/binarytree/node"
)

type NodeState struct {
	Node                  node.Pair
	Neighbors             []node.Pair
	ValueMarshalerCreator genericMarshalerCreator
	KeyMarshalerCreator   genericMarshalerCreator
}

var _ json.Marshaler = &NodeState{}

func (n *NodeState) MarshalJSON() ([]byte, error) {
	keyMarshalerCreator := n.KeyMarshalerCreator
	if keyMarshalerCreator == nil {
		keyMarshalerCreator = func(value interface{}) json.Marshaler {
			return &defaultMarshaler{value}
		}
	}

	valueMarshalerCreator := n.ValueMarshalerCreator
	if valueMarshalerCreator == nil {
		valueMarshalerCreator = func(value interface{}) json.Marshaler {
			return &defaultMarshaler{value}
		}
	}

	type marshalableNode struct {
		Node      json.RawMessage   `json:"node"`
		Neighbors []json.RawMessage `json:"neighbors"`
	}

	node, err := marshalPair(n.Node, keyMarshalerCreator, valueMarshalerCreator)
	if err != nil {
		return []byte{}, err
	}

	neighbors := make([]json.RawMessage, 0)
	for _, node := range n.Neighbors {
		result, err := marshalPair(node, keyMarshalerCreator, valueMarshalerCreator)
		if err != nil {
			return []byte{}, err
		}
		neighbors = append(neighbors, result)
	}

	return json.Marshal(marshalableNode{node, neighbors})
}

func marshalPair(
	nodeState node.Pair,
	valueMarshalerCreator,
	keyMarshalerCreator genericMarshalerCreator,
) (json.RawMessage, error) {
	type pair struct {
		Key   interface{} `json:"key"`
		Value interface{} `json:"value"`
	}

	return json.Marshal(pair(nodeState))
}

func pairFromNode(n *node.Node) node.Pair {
	return node.Pair{n.Key(), n.Value()}
}

// NewNodeState creates a new NodeState object.
//
// Note: this is not thread safe!
func NewNodeState(n *node.Node) NodeState {
	nodes := []node.Pair{}
	parent := n.Parent()
	left := n.Left()
	right := n.Right()
	if parent != nil {
		nodes = append(nodes, pairFromNode(parent))
	}
	if left != nil {
		nodes = append(nodes, pairFromNode(left))
	}
	if right != nil {
		nodes = append(nodes, pairFromNode(right))
	}
	return NodeState{
		Node:      pairFromNode(n),
		Neighbors: nodes,
	}
}
