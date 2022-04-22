package node

// Note: NOT THREAD SAFE!
type Node struct {
	key    interface{}
	value  interface{}
	parent *Node
	left   *Node
	right  *Node
	height int
}

func NewNode(key interface{}, value interface{}) Node {
	return Node{key: key, value: value}
}

func max(x, y int) int {
	if x < y {
		return y
	}
	return x
}

// Insert inserts the given node as a child node
func (n *Node) Insert(toInsert *Node) {
	if n.leftHeight() < n.rightHeight() {
		n.insertLeft(toInsert)
	} else {
		n.insertRight(toInsert)
	}

	n.height = max(n.leftHeight(), n.rightHeight()) + 1
}

func (n *Node) insertLeft(toInsert *Node) {
	if n.left == nil {
		n.left = toInsert
		n.left.parent = n
	}

	n.left.Insert(toInsert)
}

func (n *Node) insertRight(toInsert *Node) {
	if n.right == nil {
		n.right = toInsert
		n.right.parent = n
	}

	n.right.Insert(toInsert)
}

func (n *Node) leftHeight() int {
	if n.left == nil {
		return 0
	}

	return n.left.height
}

func (n *Node) rightHeight() int {
	if n.left == nil {
		return 0
	}

	return n.right.height
}

// Height retreives the height associated with the node
func (n *Node) Height() int {
	return n.height
}

// Deletes any node that has a key specified by the key parameter, except for
//   the node itself
func (n *Node) Delete(key interface{}) bool {
	if n.left != nil && n.left.key == key {
		n.left.deleteLeft()
		return true
	} else if n.right != nil && n.right.key == key {
		n.right.deleteRight()
		return true
	}

	return (n.left != nil && n.left.Delete(key)) ||
		(n.right != nil && n.right.Delete(key))
}

func (n *Node) deleteLeft() {
	if n.left != nil {
		left := n.left.left
		right := n.left.right

		n.left = nil

		n.scatterAndInsert(left)
		n.scatterAndInsert(right)
	}
}

func (n *Node) deleteRight() {
	if n.right != nil {
		left := n.right.left
		right := n.right.right

		n.right = nil

		n.scatterAndInsert(left)
		n.scatterAndInsert(right)
	}
}

func (n *Node) scatterAndInsert(node *Node) {
	if node != nil {
		nodes := node.scatter()
		for child := range nodes {
			n.Insert(child)
		}
	}
}

func (n *Node) scatter() <-chan *Node {
	c := make(chan *Node)

	left := n.left
	n.left = nil
	right := n.right
	n.right = nil

	if left != nil {
		left.parent = nil
		cl := left.scatter()
		go func() {
			for node := range cl {
				c <- node
			}

			if right != nil {
				right.parent = nil
				cr := right.scatter()
				go func() {
					for node := range cr {
						c <- node
					}
				}()
			}
		}()
	}

	return c
}

// Find performs a depth-first-search of the node that is represented by the key
func (n *Node) Find(key interface{}) *Node {
	if n.key == key {
		return n
	}
	if n.left != nil {
		node := n.left.Find(key)
		if node != nil {
			return node
		}
	}

	if n.right != nil {
		return n.right.Find(key)
	}

	return nil
}
