package node

type Scatterer struct {
	node *Node
}

func NewScatterer(node *Node) Scatterer {
	return Scatterer{node}
}

func (s Scatterer) Scatter() <-chan *Node {
	return s.node.scatter()
}
