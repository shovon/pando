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

func NewNode(key interface{}, value interface{}) *Node {
	return &Node{key: key, value: value}
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
		return
	}

	n.left.Insert(toInsert)
}

func (n *Node) insertRight(toInsert *Node) {
	if n.right == nil {
		n.right = toInsert
		n.right.parent = n
		return
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

// Delete deletes any node that has a key specified by the key parameter, except
// for the node itself.
//
// Note: this is a BFS algorithm
func (n *Node) Delete(key interface{}) bool {
	if n.left != nil && n.left.key == key {
		n.deleteLeft()
		return true
	} else if n.right != nil && n.right.key == key {
		n.deleteRight()
		return true
	}

	return (n.left != nil && n.left.Delete(key)) ||
		(n.right != nil && n.right.Delete(key))
}

func (n *Node) deleteLeft() {
	left := n.left
	n.left = nil
	if left != nil {
		n.deleteNode(left)
	}
}

func (n *Node) deleteRight() {
	right := n.right
	n.right = nil
	if right != nil {
		n.deleteNode(right)
	}
}

func (n *Node) deleteNode(node *Node) {
	left := node.left
	right := node.right

	n.scatterAndInsert(left)
	n.scatterAndInsert(right)
}

func (n *Node) scatterAndInsert(node *Node) {
	if node != nil {
		node.parent = nil
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

	iterate := func(node *Node) {
		node.parent = nil
		cr := node.scatter()
		for node := range cr {
			c <- node
		}
		c <- node
	}

	go func() {
		if left != nil {
			iterate(left)
		}
		if right != nil {
			iterate(right)
		}
		close(c)
	}()

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

// Key gets the key that the node holds
func (n Node) Key() interface{} {
	return n.key
}

// Value gets the value that the node holds
func (n Node) Value() interface{} {
	return n.value
}

func (n Node) Parent() *Node {
	return n.parent
}

// Left gets the left subtree that the node points to
func (n Node) Left() *Node {
	return n.left
}

// Right gets the right subtree that the node points to
func (n Node) Right() *Node {
	return n.right
}

// Iterates through all nodes in all subtrees
func (n Node) Iterate() <-chan *Node {
	c := make(chan *Node)

	go func() {
		iterate := func(node *Node) {
			for node := range node.Iterate() {
				c <- node
			}
			c <- node
		}

		if n.left != nil {
			iterate(n.left)
		}

		if n.right != nil {
			iterate(n.right)
		}

		c <- &n

		close(c)
	}()

	return c
}

func (n Node) Cardinality() int {
	return 1 + nodeCardinality(n.left) + nodeCardinality(n.right)
}

func nodeCardinality(node *Node) int {
	if node == nil {
		return 0
	}

	return node.Cardinality()
}
