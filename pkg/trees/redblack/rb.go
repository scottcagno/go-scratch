package redblack

import (
	"fmt"
	"runtime"
	"sync"
)

// Item is the myItem that a red-black tree stores
type Item interface {
	Compare(other Item) int
	Size() int
}

var emptyItem = *new(Item)

type nodeColor uint8

const (
	_ nodeColor = iota
	red
	black
)

// node is a node of the red-black tree
type node struct {
	parent *node
	left   *node
	right  *node
	color  nodeColor
	item   Item
}

// Tree is a red-black tree implementation
type Tree struct {
	lock  sync.Mutex
	root  *node
	NIL   *node
	count int
	size  int
}

// NewTree creates and returns a new red-black tree instance.
func NewTree() *Tree {
	n := &node{
		parent: nil,
		left:   nil,
		right:  nil,
		color:  black,
		item:   emptyItem,
	}
	return &Tree{
		root:  n,
		NIL:   n,
		count: 0,
		size:  0,
	}
}

// Has checks and returns a boolean indicating true if the
// provided item exists in the tree, and false if it cannot
// be located.
func (t *Tree) Has(item Item) bool {
	res := t.search(
		&node{
			parent: t.NIL,
			left:   t.NIL,
			right:  t.NIL,
			color:  red,
			item:   item,
		},
	)
	return res.item != nil
}

// Add inserts the provided item if and only if the item
// does not already exist in the tree. It returns a boolean
// indicating false if the item was already present in the
// tree and could therefore not be added, otherwise returning
// true for a successful insertion.
func (t *Tree) Add(item Item) bool {
	n := &node{
		parent: t.NIL,
		left:   t.NIL,
		right:  t.NIL,
		color:  red,
		item:   item,
	}
	res := t.search(n)
	if res.item != nil {
		// item already exists in the tree, so we will not add
		return false
	}
	// otherwise, the item does not exist in the tree, so we can
	// insert it.
	_, updated := t.insert(n)
	return !updated
}

// Put inserts the provided item. If the item already exists
// in the tree, put will overwrite the existing item with the
// newly provided item. If the item does not yet exist in the
// tree, it's inserted as a new item. It returns a boolean
// indicating true if the item already existed and was therefore
// updated, and false if there was no prior matching item in
// the tree.
func (t *Tree) Put(item Item) bool {
	if item == nil {
		return false
	}
	n := &node{
		parent: t.NIL,
		left:   t.NIL,
		right:  t.NIL,
		color:  red,
		item:   item,
	}
	_, updated := t.insert(n)
	return updated
}

// Del locates and removes the item matching the provided item
// key and returns the removed item and a boolean indicating true
// if the item was successfully located and removed and false
// if the item could not be found or removed.
func (t *Tree) Del(item Item) (Item, bool) {
	if item == nil {
		return nil, false
	}
	cnt := t.count
	res := t.remove(
		&node{
			parent: t.NIL,
			left:   t.NIL,
			right:  t.NIL,
			color:  red,
			item:   item,
		},
	)
	return res.item, cnt == t.count+1
}

// Get performs a search and attempts to return the item that
// contains a matching key. It returns the item along with a
// boolean indicating true if the item was successfully found,
// and false if the item could not be located.
func (t *Tree) Get(item Item) (Item, bool) {
	if item == nil {
		return nil, false
	}
	res := t.search(
		&node{
			parent: t.NIL,
			left:   t.NIL,
			right:  t.NIL,
			color:  red,
			item:   item,
		},
	)
	return res.item, res.item != nil
}

// GetNearMin performs an approximate search for the specified item
// key and returns the closest item that is less than (the predecessor)
// of the searched item key. It returns a boolean indicating true if
// an exact match was found for the key, and false if it is unknown
// or an exact match was not found.
func (t *Tree) GetNearMin(item Item) (Item, bool) {
	if item == nil {
		return nil, false
	}
	n := &node{
		parent: t.NIL,
		left:   t.NIL,
		right:  t.NIL,
		color:  red,
		item:   item,
	}
	res := t.searchApprox(n)
	prev := t.predecessor(res)
	if prev == nil {
		prev = t.min(prev)
	}
	return prev.item, prev.item.Compare(item) == 0
}

// GetNearMax performs an approximate search for the specified item
// key and returns the closest item that is greater than (the successor)
// of the searched item key. It returns a boolean indicating true if
// an exact match was found for the key, and false if it is unknown
// or an exact match was not found.
func (t *Tree) GetNearMax(item Item) (Item, bool) {
	if item == nil {
		return nil, false
	}
	n := &node{
		parent: t.NIL,
		left:   t.NIL,
		right:  t.NIL,
		color:  red,
		item:   item,
	}
	res := t.searchApprox(n)
	next := t.successor(res)
	if next == nil {
		next = t.min(next)
	}
	return next.item, next.item.Compare(item) == 0
}

// Scan iterates the items in the tree from min to max
func (t *Tree) Scan(iter func(item Item) bool) {
	t.ascend(t.root, t.min(t.root).item, iter)
}

func (t *Tree) Close() {
	t.NIL = nil
	t.root = nil
	t.count = 0
	t.size = 0
	return
}

func (t *Tree) Reset() {
	t.NIL = nil
	t.root = nil
	t.count = 0
	t.size = 0
	runtime.GC()
	n := &node{
		left:   nil,
		right:  nil,
		parent: nil,
		color:  black,
		item:   emptyItem,
	}
	t.NIL = n
	t.root = n
	t.count = 0
	t.size = 0
}

func (t *Tree) Print() {
	// todo: need to fix this
	panic("need to fix this")
	p := t.min(t.root)
	var prevItem Item
	// minItem := t.min(t.root).item
	maxItem := t.max(t.root).item
	for p != t.NIL {
		if p.left == t.NIL && p.right == t.NIL {
			fmt.Printf(
				"%s (prev=%s)\n",
				p.item, prevItem,
			)
		}
		prevItem = p.item
		if p.item == maxItem {
			p = t.predecessor(p)
		} else {
			p = t.successor(p)
		}
	}
}

// insert will insert the provided node into the tree, updating
// the existing node if it is already present in the tree. It
// returns the newly inserted (or updated) node and a boolean
// indicating true if the node was updated, and false if the
// node was newly inserted.
func (t *Tree) insert(n *node) (*node, bool) {
	x := t.root
	y := t.NIL
	for x != t.NIL {
		y = x
		if n.item.Compare(x.item) == -1 {
			x = x.left
		} else if x.item.Compare(n.item) == -1 {
			x = x.right
		} else {
			t.size -= x.item.Size()
			t.size += n.item.Size()
			// During the first implementation we were simply
			// retuning x without updating the item, but then
			// I decided that I wanted to support updating much
			// like with a hashmap, so we need to update any
			// entries that already exist in the tree before
			// returning.
			x.item = n.item
			return x, true // true = updated existing item
			//
			// It should be noted that we don't need to re-balance
			// the tree because the keys for the item have not
			// been changed and the tree is balance is maintained
			// by the item keys, and not by their values.
		}
	}
	n.parent = y
	if y == t.NIL {
		t.root = n
	} else if n.item.Compare(y.item) == -1 {
		y.left = n
	} else {
		y.right = n
	}
	// Increase the count and size because we have just inserted new
	t.count++
	t.size += n.item.Size()
	// And now, we run any fix-ups for re-balancing
	t.insertFixup(n)
	return n, false
}

func (t *Tree) insertFixup(n *node) {
	for n.parent.color == red {
		if n.parent == n.parent.parent.left {
			y := n.parent.parent.right
			if y.color == red {
				n.parent.color = black
				y.color = black
				n.parent.parent.color = red
				n = n.parent.parent
			} else {
				if n == n.parent.right {
					n = n.parent
					t.rotateLeft(n)
				}
				n.parent.color = black
				n.parent.parent.color = red
				t.rotateRight(n.parent.parent)
			}
		} else {
			y := n.parent.parent.left
			if y.color == red {
				n.parent.color = black
				y.color = black
				n.parent.parent.color = red
				n = n.parent.parent
			} else {
				if n == n.parent.left {
					n = n.parent
					t.rotateRight(n)
				}
				n.parent.color = black
				n.parent.parent.color = red
				t.rotateLeft(n.parent.parent)
			}
		}
	}
	t.root.color = black
}

func (t *Tree) rotateLeft(n *node) {
	if n.right == t.NIL {
		return
	}
	y := n.right
	n.right = y.left
	if y.left != t.NIL {
		y.left.parent = n
	}
	y.parent = n.parent
	if n.parent == t.NIL {
		t.root = y
	} else if n == n.parent.left {
		n.parent.left = y
	} else {
		n.parent.right = y
	}
	y.left = n
	n.parent = y
}

func (t *Tree) rotateRight(n *node) {
	if n.left == t.NIL {
		return
	}
	y := n.left
	n.left = y.right
	if y.right != t.NIL {
		y.right.parent = n
	}
	y.parent = n.parent
	if n.parent == t.NIL {
		t.root = y
	} else if n == n.parent.left {
		n.parent.left = y
	} else {
		n.parent.right = y
	}
	y.right = n
	n.parent = y
}

// search attempts to locate the node where the provided item
// key resides and returns the matching node, or nil if a node
// with a matching item could not be located.
func (t *Tree) search(n *node) *node {
	p := t.root
	for p != t.NIL {
		if p.item.Compare(n.item) == -1 {
			p = p.right
		} else if n.item.Compare(p.item) == -1 {
			p = p.left
		} else {
			break
		}
	}
	return p
}

// searchApprox attempts to locate the node where the provided
// item key resides, but if an exact match cannot be found it
// will return the closest match found.
func (t *Tree) searchApprox(n *node) *node {
	p := t.root
	for p != t.NIL {
		if p.item.Compare(n.item) == -1 {
			if p.right == t.NIL {
				break
			}
			p = p.right
		} else if n.item.Compare(p.item) == -1 {
			if p.left == t.NIL {
				break
			}
			p = p.left
		} else {
			break
		}
	}
	return p
}

// min traverses from root to left recursively until left is NIL
func (t *Tree) min(n *node) *node {
	if n == t.NIL {
		return t.NIL
	}
	for n.left != t.NIL {
		n = n.left
	}
	return n
}

// max traverses from root to right recursively until right is NIL
func (t *Tree) max(n *node) *node {
	if n == t.NIL {
		return t.NIL
	}
	for n.right != t.NIL {
		n = n.right
	}
	return n
}

// predecessor locates the node that precedes the provided one
func (t *Tree) predecessor(n *node) *node {
	if n == t.NIL {
		return t.NIL
	}
	if n.left != t.NIL {
		return t.max(n.left)
	}
	y := n.parent
	for y != t.NIL && n == y.left {
		n = y
		y = y.parent
	}
	return y
}

// successor locates the node that succeeds the provided one
func (t *Tree) successor(n *node) *node {
	if n == t.NIL {
		return t.NIL
	}
	if n.right != t.NIL {
		return t.min(n.right)
	}
	y := n.parent
	for y != t.NIL && n == y.right {
		n = y
		y = y.parent
	}
	return y
}

// remove is the internal method that deletes the provided
// node from the tree. It returns the removed node or nil if
// the node could not be found or removed for some reason.
func (t *Tree) remove(n *node) *node {
	// locate
	p := t.root
	for p != t.NIL {
		if p.item.Compare(n.item) == -1 {
			p = p.right
		} else if n.item.Compare(p.item) == -1 {
			p = p.left
		} else {
			break
		}
	}
	// remove
	z := p
	if z == t.NIL {
		return t.NIL
	}
	res := &node{
		t.NIL,
		t.NIL,
		t.NIL,
		z.color,
		z.item,
	}
	var y, x *node
	if z.left == t.NIL || z.right == t.NIL {
		y = z
	} else {
		y = t.successor(z)
	}
	if y.left != t.NIL {
		x = y.left
	} else {
		x = y.right
	}
	x.parent = y.parent

	if y.parent == t.NIL {
		t.root = x
	} else if y == y.parent.left {
		y.parent.left = x
	} else {
		y.parent.right = x
	}
	if y != z {
		z.item = y.item
	}
	if y.color == BLACK {
		t.removeFixup(x)
	}
	t.size -= res.item.Size()
	t.count--
	return res
}

func (t *Tree) removeFixup(n *node) {
	for n != t.root && n.color == BLACK {
		if n == n.parent.left {
			w := n.parent.right
			if w.color == RED {
				w.color = BLACK
				n.parent.color = RED
				t.rotateLeft(n.parent)
				w = n.parent.right
			}
			if w.left.color == BLACK && w.right.color == BLACK {
				w.color = RED
				n = n.parent
			} else {
				if w.right.color == BLACK {
					w.left.color = BLACK
					w.color = RED
					t.rotateRight(w)
					w = n.parent.right
				}
				w.color = n.parent.color
				n.parent.color = BLACK
				w.right.color = BLACK
				t.rotateLeft(n.parent)
				// this is to exit while loop
				n = t.root
			}
		} else {
			w := n.parent.left
			if w.color == RED {
				w.color = BLACK
				n.parent.color = RED
				t.rotateRight(n.parent)
				w = n.parent.left
			}
			if w.left.color == BLACK && w.right.color == BLACK {
				w.color = RED
				n = n.parent
			} else {
				if w.left.color == BLACK {
					w.right.color = BLACK
					w.color = RED
					t.rotateLeft(w)
					w = n.parent.left
				}
				w.color = n.parent.color
				n.parent.color = BLACK
				w.left.color = BLACK
				t.rotateRight(n.parent)
				n = t.root
			}
		}
	}
	n.color = BLACK
}

func (t *Tree) ascend(x *node, item Item, iter func(item Item) bool) bool {
	if x == t.NIL {
		return true
	}
	if !(x.item.Compare(item) == -1) {
		if !t.ascend(x.left, item, iter) {
			return false
		}
		if !iter(x.item) {
			return false
		}
	}
	return t.ascend(x.right, item, iter)
}

func (t *Tree) descend(x *node, pivot Item, iter func(item Item) bool) bool {
	if x == t.NIL {
		return true
	}
	if !(pivot.Compare(x.item) == -1) {
		if !t.descend(x.right, pivot, iter) {
			return false
		}
		if !iter(x.item) {
			return false
		}
	}
	return t.descend(x.left, pivot, iter)
}

func (t *Tree) ascendRange(x *node, inf, sup Item, iter func(item Item) bool) bool {
	if x == t.NIL {
		return true
	}
	if !(x.item.Compare(sup) == -1) {
		return t.ascendRange(x.left, inf, sup, iter)
	}
	if x.item.Compare(inf) == -1 {
		return t.ascendRange(x.right, inf, sup, iter)
	}
	if !t.ascendRange(x.left, inf, sup, iter) {
		return false
	}
	if !iter(x.item) {
		return false
	}
	return t.ascendRange(x.right, inf, sup, iter)
}

// Iterator is an iteration type for the (red-black) Tree
type Iterator struct {
	tree    *Tree
	current *node
	index   int
}

func (t *Tree) NewIterator() *Iterator {
	n := t.min(t.root)
	if n == t.NIL {
		return nil
	}
	it := &Iterator{
		tree:    t,
		current: n,
		index:   t.size,
	}
	return it
}

func (it *Iterator) First() Item {
	n := it.tree.min(it.tree.root)
	if n == it.tree.NIL {
		return nil
	}
	if it.current != nil && it.current == n {
		return it.current.item
	}
	it.current = n
	it.index = it.tree.size
	return it.current.item
}

func (it *Iterator) Last() Item {
	n := it.tree.max(it.tree.root)
	if n == it.tree.NIL {
		return nil
	}
	if it.current != nil && it.current == n {
		return it.current.item
	}
	it.current = n
	it.index = it.tree.size
	return it.current.item
}

func (it *Iterator) Next() Item {
	next := it.tree.successor(it.current)
	if next == it.tree.NIL {
		return nil
	}
	it.index--
	it.current = next
	return next.item
}

func (it *Iterator) Prev() Item {
	prev := it.tree.predecessor(it.current)
	if prev == it.tree.NIL {
		return nil
	}
	it.index--
	it.current = prev
	return prev.item
}

func (it *Iterator) More() bool {
	return it.index > 1
}
