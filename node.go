package main

type Node struct {
	key    interface{}
	value  interface{}
	parent *Node
	left   *Node
	right  *Node
	depth  int
}

func (n *Node) Insert(toInsert *Node) {
	if n.left == nil {
		n.left = toInsert
	} else if n.right == nil {
		n.right = toInsert
	} else if n.left.depth < n.right.depth {
		n.left.Insert(toInsert)
	} else {
		n.right.Insert(toInsert)
	}
}

func max(x, y int) int {
	if x < y {
		return y
	}
	return x
}

func (n *Node) Depth() int {
	leftDepth := 0
	if n.left != nil {
		leftDepth = n.left.Depth() + 1
	}
	rightDepth := 0
	if n.right != nil {
		n.right.depth = n.right.Depth() + 1
	}
	n.depth = max(leftDepth, rightDepth)
	return n.depth
}

func (n *Node) Delete(key interface{}) {
	if n.left.key == key {

	}
}
