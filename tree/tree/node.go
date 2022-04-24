package tree

import (
	"encoding/json"
	"tree/tree/node"
)

type NodeState struct {
	Node                  Pair
	Neighbors             []Pair
	ValueMarshalerCreator genericMarshalerCreator
	KeyMarshalerCreator   genericMarshalerCreator
}

var _ json.Marshaler = &NodeState{}

func (n *NodeState) MarshalJSON() ([]byte, error) {
	valueMarshalerCreator := n.ValueMarshalerCreator
	if valueMarshalerCreator == nil {
		valueMarshalerCreator = func(value interface{}) json.Marshaler {
			return &defaultMarshaler{value}
		}
	}

	keyMarshalerCreator := n.KeyMarshalerCreator
	if keyMarshalerCreator == nil {
		keyMarshalerCreator = func(value interface{}) json.Marshaler {
			return &defaultMarshaler{value}
		}
	}

}

func pairFromNode(n *node.Node) Pair {
	return Pair{n.Key(), n.Value()}
}

// NewNodeState creates a new NodeState object.
//
// Note: this is not thread safe!
func NewNodeState(n *node.Node) NodeState {
	nodes := []Pair{}
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
		pairFromNode(n),
		nodes,
	}
}
