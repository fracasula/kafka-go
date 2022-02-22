// Code generated by gotemplate. DO NOT EDIT.

package cache

// template type TreeMap(Key, Value)

// Key is a generic key type of the map

// Value is a generic value type of the map

// TreeMap is the red-black tree based map
type index struct {
	endNode   *nodeIndex
	beginNode *nodeIndex
	count     int
	// Less returns a < b
	Less func(a *entry, b *entry) bool
}

type nodeIndex struct {
	right   *nodeIndex
	left    *nodeIndex
	parent  *nodeIndex
	isBlack bool
	key     *entry
	value   struct{}
}

// New creates and returns new TreeMap.
// Parameter less is a function returning a < b.
func newIndex(less func(a *entry, b *entry) bool) *index {
	endNode := &nodeIndex{isBlack: true}
	return &index{beginNode: endNode, endNode: endNode, Less: less}
}

// Len returns total count of elements in a map.
// Complexity: O(1).
func (t *index) Len() int { return t.count }

// Set sets the value and silently overrides previous value if it exists.
// Complexity: O(log N).
func (t *index) Set(key *entry, value struct{}) {
	parent := t.endNode
	current := parent.left
	less := true
	for current != nil {
		parent = current
		switch {
		case t.Less(key, current.key):
			current = current.left
			less = true
		case t.Less(current.key, key):
			current = current.right
			less = false
		default:
			current.value = value
			return
		}
	}
	x := &nodeIndex{parent: parent, value: value, key: key}
	if less {
		parent.left = x
	} else {
		parent.right = x
	}
	if t.beginNode.left != nil {
		t.beginNode = t.beginNode.left
	}
	t.insertFixup(x)
	t.count++
}

// Del deletes the value.
// Complexity: O(log N).
func (t *index) Del(key *entry) {
	z := t.findNode(key)
	if z == nil {
		return
	}
	if t.beginNode == z {
		if z.right != nil {
			t.beginNode = z.right
		} else {
			t.beginNode = z.parent
		}
	}
	t.count--
	removeNodeIndex(t.endNode.left, z)
}

// Clear clears the map.
// Complexity: O(1).
func (t *index) Clear() {
	t.count = 0
	t.beginNode = t.endNode
	t.endNode.left = nil
}

// Get retrieves a value from a map for specified key and reports if it exists.
// Complexity: O(log N).
func (t *index) Get(id *entry) (struct{}, bool) {
	node := t.findNode(id)
	if node == nil {
		node = t.endNode
	}
	return node.value, node != t.endNode
}

// Contains checks if key exists in a map.
// Complexity: O(log N)
func (t *index) Contains(id *entry) bool { return t.findNode(id) != nil }

// Range returns a pair of iterators that you can use to go through all the keys in the range [from, to].
// More specifically it returns iterators pointing to lower bound and upper bound.
// Complexity: O(log N).
func (t *index) Range(from, to *entry) (forwardIteratorIndex, forwardIteratorIndex) {
	return t.LowerBound(from), t.UpperBound(to)
}

// LowerBound returns an iterator pointing to the first element that is not less than the given key.
// Complexity: O(log N).
func (t *index) LowerBound(key *entry) forwardIteratorIndex {
	result := t.endNode
	node := t.endNode.left
	if node == nil {
		return forwardIteratorIndex{tree: t, node: t.endNode}
	}
	for {
		if t.Less(node.key, key) {
			if node.right != nil {
				node = node.right
			} else {
				return forwardIteratorIndex{tree: t, node: result}
			}
		} else {
			result = node
			if node.left != nil {
				node = node.left
			} else {
				return forwardIteratorIndex{tree: t, node: result}
			}
		}
	}
}

// UpperBound returns an iterator pointing to the first element that is greater than the given key.
// Complexity: O(log N).
func (t *index) UpperBound(key *entry) forwardIteratorIndex {
	result := t.endNode
	node := t.endNode.left
	if node == nil {
		return forwardIteratorIndex{tree: t, node: t.endNode}
	}
	for {
		if !t.Less(key, node.key) {
			if node.right != nil {
				node = node.right
			} else {
				return forwardIteratorIndex{tree: t, node: result}
			}
		} else {
			result = node
			if node.left != nil {
				node = node.left
			} else {
				return forwardIteratorIndex{tree: t, node: result}
			}
		}
	}
}

// Iterator returns an iterator for tree map.
// It starts at the first element and goes to the one-past-the-end position.
// You can iterate a map at O(N) complexity.
// Method complexity: O(1)
func (t *index) Iterator() forwardIteratorIndex {
	return forwardIteratorIndex{tree: t, node: t.beginNode}
}

// Reverse returns a reverse iterator for tree map.
// It starts at the last element and goes to the one-before-the-start position.
// You can iterate a map at O(N) complexity.
// Method complexity: O(log N)
func (t *index) Reverse() reverseIteratorIndex {
	node := t.endNode.left
	if node != nil {
		node = mostRightIndex(node)
	}
	return reverseIteratorIndex{tree: t, node: node}
}

func (t *index) findNode(id *entry) *nodeIndex {
	current := t.endNode.left
	for current != nil {
		switch {
		case t.Less(id, current.key):
			current = current.left
		case t.Less(current.key, id):
			current = current.right
		default:
			return current
		}
	}
	return nil
}

func mostLeftIndex(x *nodeIndex) *nodeIndex {
	for x.left != nil {
		x = x.left
	}
	return x
}

func mostRightIndex(x *nodeIndex) *nodeIndex {
	for x.right != nil {
		x = x.right
	}
	return x
}

func successorIndex(x *nodeIndex) *nodeIndex {
	if x.right != nil {
		return mostLeftIndex(x.right)
	}
	for x != x.parent.left {
		x = x.parent
	}
	return x.parent
}

func predecessorIndex(x *nodeIndex) *nodeIndex {
	if x.left != nil {
		return mostRightIndex(x.left)
	}
	for x.parent != nil && x != x.parent.right {
		x = x.parent
	}
	return x.parent
}

func rotateLeftIndex(x *nodeIndex) {
	y := x.right
	x.right = y.left
	if x.right != nil {
		x.right.parent = x
	}
	y.parent = x.parent
	if x == x.parent.left {
		x.parent.left = y
	} else {
		x.parent.right = y
	}
	y.left = x
	x.parent = y
}

func rotateRightIndex(x *nodeIndex) {
	y := x.left
	x.left = y.right
	if x.left != nil {
		x.left.parent = x
	}
	y.parent = x.parent
	if x == x.parent.left {
		x.parent.left = y
	} else {
		x.parent.right = y
	}
	y.right = x
	x.parent = y
}

func (t *index) insertFixup(x *nodeIndex) {
	root := t.endNode.left
	x.isBlack = x == root
	for x != root && !x.parent.isBlack {
		if x.parent == x.parent.parent.left {
			y := x.parent.parent.right
			if y != nil && !y.isBlack {
				x = x.parent
				x.isBlack = true
				x = x.parent
				x.isBlack = x == root
				y.isBlack = true
			} else {
				if x != x.parent.left {
					x = x.parent
					rotateLeftIndex(x)
				}
				x = x.parent
				x.isBlack = true
				x = x.parent
				x.isBlack = false
				rotateRightIndex(x)
				break
			}
		} else {
			y := x.parent.parent.left
			if y != nil && !y.isBlack {
				x = x.parent
				x.isBlack = true
				x = x.parent
				x.isBlack = x == root
				y.isBlack = true
			} else {
				if x == x.parent.left {
					x = x.parent
					rotateRightIndex(x)
				}
				x = x.parent
				x.isBlack = true
				x = x.parent
				x.isBlack = false
				rotateLeftIndex(x)
				break
			}
		}
	}
}

// nolint: gocyclo
func removeNodeIndex(root *nodeIndex, z *nodeIndex) {
	var y *nodeIndex
	if z.left == nil || z.right == nil {
		y = z
	} else {
		y = successorIndex(z)
	}
	var x *nodeIndex
	if y.left != nil {
		x = y.left
	} else {
		x = y.right
	}
	var w *nodeIndex
	if x != nil {
		x.parent = y.parent
	}
	if y == y.parent.left {
		y.parent.left = x
		if y != root {
			w = y.parent.right
		} else {
			root = x // w == nil
		}
	} else {
		y.parent.right = x
		w = y.parent.left
	}
	removedBlack := y.isBlack
	if y != z {
		y.parent = z.parent
		if z == z.parent.left {
			y.parent.left = y
		} else {
			y.parent.right = y
		}
		y.left = z.left
		y.left.parent = y
		y.right = z.right
		if y.right != nil {
			y.right.parent = y
		}
		y.isBlack = z.isBlack
		if root == z {
			root = y
		}
	}
	if removedBlack && root != nil {
		if x != nil {
			x.isBlack = true
		} else {
			for {
				if w != w.parent.left {
					if !w.isBlack {
						w.isBlack = true
						w.parent.isBlack = false
						rotateLeftIndex(w.parent)
						if root == w.left {
							root = w
						}
						w = w.left.right
					}
					if (w.left == nil || w.left.isBlack) && (w.right == nil || w.right.isBlack) {
						w.isBlack = false
						x = w.parent
						if x == root || !x.isBlack {
							x.isBlack = true
							break
						}
						if x == x.parent.left {
							w = x.parent.right
						} else {
							w = x.parent.left
						}
					} else {
						if w.right == nil || w.right.isBlack {
							w.left.isBlack = true
							w.isBlack = false
							rotateRightIndex(w)
							w = w.parent
						}
						w.isBlack = w.parent.isBlack
						w.parent.isBlack = true
						w.right.isBlack = true
						rotateLeftIndex(w.parent)
						break
					}
				} else {
					if !w.isBlack {
						w.isBlack = true
						w.parent.isBlack = false
						rotateRightIndex(w.parent)
						if root == w.right {
							root = w
						}
						w = w.right.left
					}
					if (w.left == nil || w.left.isBlack) && (w.right == nil || w.right.isBlack) {
						w.isBlack = false
						x = w.parent
						if !x.isBlack || x == root {
							x.isBlack = true
							break
						}
						if x == x.parent.left {
							w = x.parent.right
						} else {
							w = x.parent.left
						}
					} else {
						if w.left == nil || w.left.isBlack {
							w.right.isBlack = true
							w.isBlack = false
							rotateLeftIndex(w)
							w = w.parent
						}
						w.isBlack = w.parent.isBlack
						w.parent.isBlack = true
						w.left.isBlack = true
						rotateRightIndex(w.parent)
						break
					}
				}
			}
		}
	}
}

// ForwardIterator represents a position in a tree map.
// It is designed to iterate a map in a forward order.
// It can point to any position from the first element to the one-past-the-end element.
type forwardIteratorIndex struct {
	tree *index
	node *nodeIndex
}

// Valid reports if an iterator's position is valid.
// In other words it returns true if an iterator is not at the one-past-the-end position.
func (i forwardIteratorIndex) Valid() bool { return i.node != i.tree.endNode }

// Next moves an iterator to the next element.
// It panics if goes out of bounds.
func (i *forwardIteratorIndex) Next() {
	if i.node == i.tree.endNode {
		panic("out of bound iteration")
	}
	i.node = successorIndex(i.node)
}

// Prev moves an iterator to the previous element.
// It panics if goes out of bounds.
func (i *forwardIteratorIndex) Prev() {
	i.node = predecessorIndex(i.node)
	if i.node == nil {
		panic("out of bound iteration")
	}
}

// Key returns a key at an iterator's position
func (i forwardIteratorIndex) Key() *entry { return i.node.key }

// Value returns a value at an iterator's position
func (i forwardIteratorIndex) Value() struct{} { return i.node.value }

// ReverseIterator represents a position in a tree map.
// It is designed to iterate a map in a reverse order.
// It can point to any position from the one-before-the-start element to the last element.
type reverseIteratorIndex struct {
	tree *index
	node *nodeIndex
}

// Valid reports if an iterator's position is valid.
// In other words it returns true if an iterator is not at the one-before-the-start position.
func (i reverseIteratorIndex) Valid() bool { return i.node != nil }

// Next moves an iterator to the next element in reverse order.
// It panics if goes out of bounds.
func (i *reverseIteratorIndex) Next() {
	if i.node == nil {
		panic("out of bound iteration")
	}
	i.node = predecessorIndex(i.node)
}

// Prev moves an iterator to the previous element in reverse order.
// It panics if goes out of bounds.
func (i *reverseIteratorIndex) Prev() {
	if i.node != nil {
		i.node = successorIndex(i.node)
	} else {
		i.node = i.tree.beginNode
	}
	if i.node == i.tree.endNode {
		panic("out of bound iteration")
	}
}

// Key returns a key at an iterator's position
func (i reverseIteratorIndex) Key() *entry { return i.node.key }

// Value returns a value at an iterator's position
func (i reverseIteratorIndex) Value() struct{} { return i.node.value }
