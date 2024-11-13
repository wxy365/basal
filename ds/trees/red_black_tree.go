package trees

import "github.com/wxy365/basal/fn"

type Color bool

const (
	Red   Color = true
	Black Color = false
)

func NewRBTree[K, V any](comparator fn.Comparator[K]) *RBTree[K, V] {
	return &RBTree[K, V]{
		comparator: comparator,
	}
}

type Node[K any, V any] struct {
	Key    K
	Value  V
	Parent *Node[K, V]
	Left   *Node[K, V]
	Right  *Node[K, V]
	Color  Color
}

type RBTree[K any, V any] struct {
	root       *Node[K, V]
	comparator fn.Comparator[K]
}

func (rbt *RBTree[K, V]) Insert(key K, value V) {
	newNode := &Node[K, V]{
		Key:   key,
		Value: value,
		Color: Red,
	}
	if rbt.root == nil {
		rbt.root = newNode
	} else {
		currentNode := rbt.root
		var parentNode *Node[K, V]
		for currentNode != nil {
			parentNode = currentNode
			if rbt.comparator(newNode.Key, currentNode.Key) < 0 {
				currentNode = currentNode.Left
			} else {
				currentNode = currentNode.Right
			}
		}
		newNode.Parent = parentNode
		if rbt.comparator(newNode.Key, parentNode.Key) < 0 {
			parentNode.Left = newNode
		} else {
			parentNode.Right = newNode
		}
	}
	rbt.fixInsert(newNode)
}

func (rbt *RBTree[K, V]) fixInsert(n *Node[K, V]) {
	parentNode := n.Parent

	// While the parent node is red, restructure.
	for parentNode != nil && parentNode.Color {
		grandparentNode := parentNode.Parent

		// Parent is the left child of the grandparent.
		if parentNode == grandparentNode.Left {
			uncle := grandparentNode.Right

			// Case 1: Uncle is red. Just recolor.
			if uncle != nil && uncle.Color {
				grandparentNode.Color = Red
				parentNode.Color = Black
				uncle.Color = Black
				n = grandparentNode
			} else {
				// Case 2: Node is the right child. Left rotation on parent.
				if n == parentNode.Right {
					n = parentNode
					rbt.rotateLeft(n)
				}

				// Case 3: Node is the left child. Right rotation on grandparent.
				parentNode.Color = Black
				grandparentNode.Color = Red
				rbt.rotateRight(grandparentNode)
			}
		} else { // Parent is the right child of the grandparent.
			uncle := grandparentNode.Left

			// Case 1: Uncle is red. Just recolor.
			if uncle != nil && uncle.Color {
				grandparentNode.Color = Red
				parentNode.Color = Black
				uncle.Color = Black
				n = grandparentNode
			} else {
				// Case 2: Node is the left child. Right rotation on parent.
				if n == parentNode.Left {
					n = parentNode
					rbt.rotateRight(n)
				}

				// Case 3: Node is the right child. Left rotation on grandparent.
				parentNode.Color = Black
				grandparentNode.Color = Red
				rbt.rotateLeft(grandparentNode)
			}
		}
		parentNode = n.Parent
	}

	// the root is always black.
	rbt.root.Color = Black
}

func (rbt *RBTree[K, V]) rotateLeft(n *Node[K, V]) {
	rChild := n.Right
	n.Right = rChild.Left
	if rChild.Left != nil {
		rChild.Left.Parent = n
	}
	rChild.Parent = n.Parent
	if n.Parent == nil {
		// n is root
	} else if n == n.Parent.Left {
		n.Parent.Left = rChild
	} else {
		n.Parent.Right = rChild
	}
	rChild.Left = n
	n.Parent = rChild
	rbt.root = rChild
}

func (rbt *RBTree[K, V]) rotateRight(n *Node[K, V]) {
	lChild := n.Left
	n.Left = lChild.Right
	if lChild.Right != nil {
		lChild.Right.Parent = n
	}
	lChild.Parent = n.Parent
	if n.Parent == nil {
		// n is root
	} else if n == n.Parent.Right {
		n.Parent.Right = lChild
	} else {
		n.Parent.Left = lChild
	}
	lChild.Right = n
	n.Parent = lChild
	rbt.root = lChild
}
