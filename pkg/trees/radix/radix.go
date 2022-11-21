package radix

import (
	"sort"
	"strings"
)

type leafNode struct {
	key string
	val interface{}
}

type edge struct {
	label byte
	node  *node
}

type node struct {
	// leaf stores a possible leaf
	leaf *leafNode

	// prefix contains a common prefix
	prefix string

	// edges should be stored in sorted order
	// for iteration purposes. We are avoiding
	// a fully instantiated array to save memory,
	// since in most cases we expect the set to
	// be rather sparse.
	edges edges
}

func (n *node) isLeaf() bool {
	return n.leaf != nil
}

func (n *node) addEdge(e edge) {
	num := len(n.edges)
	idx := sort.Search(
		num, func(i int) bool {
			return n.edges[i].label >= e.label
		},
	)

	n.edges = append(n.edges, edge{})
	copy(n.edges[idx+1:], n.edges[idx:])
	n.edges[idx] = e
}

func (n *node) updateEdge(label byte, node *node) {
	num := len(n.edges)
	idx := sort.Search(
		num, func(i int) bool {
			return n.edges[i].label >= label
		},
	)
	if idx < num && n.edges[idx].label == label {
		n.edges[idx].node = node
		return
	}
	panic("replacing missing edge")
}

func (n *node) getEdge(label byte) *node {
	num := len(n.edges)
	idx := sort.Search(
		num, func(i int) bool {
			return n.edges[i].label >= label
		},
	)
	if idx < num && n.edges[idx].label == label {
		return n.edges[idx].node
	}
	return nil
}

func (n *node) delEdge(label byte) {
	num := len(n.edges)
	idx := sort.Search(
		num, func(i int) bool {
			return n.edges[i].label >= label
		},
	)
	if idx < num && n.edges[idx].label == label {
		copy(n.edges[idx:], n.edges[idx+1:])
		n.edges[len(n.edges)-1] = edge{}
		n.edges = n.edges[:len(n.edges)-1]
	}
}

func (n *node) mergeChild() {
	e := n.edges[0]
	child := e.node
	n.prefix = n.prefix + child.prefix
	n.leaf = child.leaf
	n.edges = child.edges
}

type edges []edge

func (e edges) Len() int {
	return len(e)
}

func (e edges) Less(i, j int) bool {
	return e[i].label < e[j].label
}

func (e edges) Swap(i, j int) {
	e[i], e[j] = e[j], e[i]
}

func (e edges) sortEdges() {
	sort.Sort(e)
}

// Tree implements a radix tree--which is a space optimized,
// compressed prefix trie. This cna be treated as a dictionary
// abstract data type. The advantage over a standard hash map is
// sorted prefix-based lookups and ordered iteration, but it will
// not be space or time optimized if the data set does not share
// many common prefixes--in which case a hashmap or RedBlackTree
// would be preferred.
type Tree struct {
	root *node
	size int
}

// NewTree returns a new pointer to an empty Tree (radix tree)
func NewTree() *Tree {
	return &Tree{
		root: new(node),
		size: 0,
	}
}

// longestPrefix finds the (longest) length of a shared
// prefix of the two strings provided
func longestPrefix(k1, k2 string) int {
	max := len(k1)
	if l := len(k2); l < max {
		max = l
	}
	var i int
	for i = 0; i < max; i++ {
		if k1[i] != k2[i] {
			break
		}
	}
	return i
}

// Insert is used to add a new entry or update an existing entry.
// Returns a boolean indicating true if an old value was updated.
func (t *Tree) Insert(k string, v any) (any, bool) {
	var parent *node
	n := t.root
	search := k
	for {
		// Handle key exhaustion
		if len(search) == 0 {
			if n.isLeaf() {
				// update old value
				old := n.leaf.val
				n.leaf.val = v
				return old, true
			}
			// otherwise, create a
			// new leaf, and insert
			n.leaf = &leafNode{
				key: k,
				val: v,
			}
			t.size++
			return nil, false
		}

		// Look for the edge
		parent = n
		n = n.getEdge(search[0])

		// No edge found, create a new one
		if n == nil {
			e := edge{
				label: search[0],
				node: &node{
					leaf: &leafNode{
						key: k,
						val: v,
					},
					prefix: search,
				},
			}
			parent.addEdge(e)
			t.size++
			return nil, false
		}

		// Determine the longest prefix match for the search key
		common := longestPrefix(search, n.prefix)
		if common == len(n.prefix) {
			search = search[common:]
			continue
		}

		// Split the node
		t.size++
		child := &node{
			prefix: search[:common],
		}
		parent.updateEdge(search[0], child)

		// Restore the existing node
		child.addEdge(
			edge{
				label: n.prefix[common],
				node:  n,
			},
		)
		n.prefix = n.prefix[common:]

		// Create a new leaf node
		leaf := &leafNode{
			key: k,
			val: v,
		}

		// If the new key is a subset, add it to this node
		search = search[common:]
		if len(search) == 0 {
			child.leaf = leaf
			return nil, false
		}

		// Create a new edge for the node
		child.addEdge(
			edge{
				label: search[0],
				node: &node{
					leaf:   leaf,
					prefix: search,
				},
			},
		)
		return nil, false
	}
}

// Delete is used to delete a key. It will return the previous
// value and a boolean indicating true if it was deleted.
func (t *Tree) Delete(k string) (any, bool) {
	var parent *node
	var label byte
	n := t.root
	search := k
	for {
		// Check for key exhaustion
		if len(search) == 0 {
			if !n.isLeaf() {
				break
			}
			goto delete
		}

		// Look for an edge
		parent = n
		label = search[0]
		n = n.getEdge(label)
		if n == nil {
			break
		}

		// Consume the search prefix
		if !(len(search) >= len(n.prefix) && search[0:len(n.prefix)] == n.prefix) {
			// inlined version of !strings.HasPrefix(search, n.prefix)
			break
		}
		search = search[len(n.prefix):]
	}
	return nil, false

delete:
	// Delete the leaf
	leaf := n.leaf
	n.leaf = nil
	t.size--

	// Check if we need to delete this node (from the parent)
	if parent != nil && len(n.edges) == 0 {
		parent.delEdge(label)
	}

	// Check if we need to merge this node
	if n != t.root && len(n.edges) == 1 {
		n.mergeChild()
	}

	// Check if we need to merge the sibling
	if parent != nil && parent != t.root && len(parent.edges) == 1 && !parent.isLeaf() {
		parent.mergeChild()
	}
	return leaf.val, true
}

// DeletePrefix is used to remove the subtree under a given prefix. It
// returns the number of nodes were deleted. This method can be used to
// remove a large subtree efficiently.
func (t *Tree) DeletePrefix(k string) int {
	return t.deletePrefixRecursive(nil, t.root, k)
}

// deletePrefixRecursive does a recursive subtree removal
func (t *Tree) deletePrefixRecursive(parent, n *node, prefix string) int {
	// Check for key exhaustion
	if len(prefix) == 0 {
		// Remove leaf node
		subTreeSize := 0
		// Recursively walk from all edges of the node (to be deleted)
		recursiveWalk(
			n, func(k string, v any) bool {
				subTreeSize++
				return false
			},
		)
		if n.isLeaf() {
			n.leaf = nil
		}
		n.edges = nil

		// Check if we need to marge the sibling
		if parent != nil && parent != t.root && len(parent.edges) == 1 && !parent.isLeaf() {
			parent.mergeChild()
		}
		t.size -= subTreeSize
		return subTreeSize
	}

	// Look for an edge
	label := prefix[0]
	child := n.getEdge(label)
	if child == nil || (!strings.HasPrefix(child.prefix, prefix) && !strings.HasPrefix(prefix, child.prefix)) {
		return 0
	}

	// Consume the search prefix
	if len(child.prefix) > len(prefix) {
		prefix = prefix[len(prefix):]
	} else {
		prefix = prefix[len(child.prefix):]
	}

	// Recursively call self
	return t.deletePrefixRecursive(n, child, prefix)
}

// recursiveWalk walks the tree recursively from node n, using the WalkFn fn.
func recursiveWalk(n *node, fn WalkFn) bool {
	// Visit the leaf values, if there are any
	if n.leaf != nil && fn(n.leaf.key, n.leaf.val) {
		return true
	}

	// Recurse on the children...
	for _, e := range n.edges {
		if recursiveWalk(e.node, fn) {
			return true
		}
	}
	return false
}

// Find is used to look up a specific key, returning the
// associated value a boolean indicating true if it was found.
func (t *Tree) Find(k string) (any, bool) {
	n := t.root
	search := k
	for {
		// Check for key exhaustion
		if len(search) == 0 {
			if n.isLeaf() {
				// Found value, return
				return n.leaf.val, true
			}
			// We are at the end, and did not
			// find anything, time to break
			break
		}

		// Look for an edge
		n = n.getEdge(search[0])
		if n == nil {
			// No more edges, time to break
			break
		}

		// Consume the search prefix
		if !(len(search) >= len(n.prefix) && search[0:len(n.prefix)] == n.prefix) {
			// inlined version of !strings.HasPrefix(search, n.prefix)
			break
		}
		search = search[len(n.prefix):]
	}
	return nil, false
}

// FindLongestPrefix is very much like Find, but instead looking
// for an exact match, it attempts to locate the longest prefix
// match. Upon success, it will return the last matched key, value
// and a boolean indicating true, otherwise "", nil and false.
func (t *Tree) FindLongestPrefix(k string) (string, any, bool) {
	var last *leafNode
	n := t.root
	search := k
	for {

		// Look for a leaf node
		if n.isLeaf() {
			last = n.leaf
		}

		// Check for key exhaustion
		if len(search) == 0 {
			break
		}

		// Look for an edge
		n = n.getEdge(search[0])
		if n == nil {
			break
		}

		// Consume the search prefix
		if !(len(search) >= len(n.prefix) && search[0:len(n.prefix)] == n.prefix) {
			// inlined version of !strings.HasPrefix(search, n.prefix)
			break
		}
		search = search[len(n.prefix):]
	}
	if last != nil {
		return last.key, last.val, true
	}
	return "", nil, false
}

// Len returns the number of elements in the tree.
func (t *Tree) Len() int {
	return t.size
}

// Min returns the minimum key, and value in the tree.
func (t *Tree) Min() (string, any, bool) {
	n := t.root
	for {
		if n.isLeaf() {
			return n.leaf.key, n.leaf.val, true
		}
		if len(n.edges) < 1 {
			break
		}
		n = n.edges[0].node
	}
	return "", nil, false
}

// Max returns the maximum key, and value in the tree.
func (t *Tree) Max() (string, any, bool) {
	n := t.root
	for {
		if num := len(n.edges); num > 0 {
			n = n.edges[num-1].node
			continue
		}
		if n.isLeaf() {
			return n.leaf.key, n.leaf.val, true
		}
		break
	}
	return "", nil, false
}

type WalkFn func(s string, v any) bool

// Walk recursively walks the tree using the WalkFn fn.
func (t *Tree) Walk(fn WalkFn) {
	recursiveWalk(t.root, fn)
}

// WalkPrefix recursively walks the tree using the supplied WalkFn fn, under
// a specific prefix supplied by the prefix string.
func (t *Tree) WalkPrefix(prefix string, fn WalkFn) {
	n := t.root
	search := prefix
	for {

		// Check for key exhaustion
		if len(search) == 0 {
			recursiveWalk(n, fn)
			return
		}

		// Look for an edge
		n = n.getEdge(search[0])
		if n == nil {
			return
		}

		// Consume the search prefix
		if strings.HasPrefix(n.prefix, search) {
			recursiveWalk(n, fn) // The child may be under our search prefix
			return
		}
		if !strings.HasPrefix(search, n.prefix) {
			break
		}
		search = search[len(n.prefix):]
	}
}

// WalkPath recursively walks the tree using the supplied WalkFn fn, under
// a specific path supplied by the path string. It is like WalkPrefix, but
// instead of visiting all the entries under a given prefix, this walks the
// entries above the path.
func (t *Tree) WalkPath(path string, fn WalkFn) {
	n := t.root
	search := path
	for {

		// Visit the leaf values, if there are any.
		if n.leaf != nil && fn(n.leaf.key, n.leaf.val) {
			return
		}

		// Check for key exhaustion
		if len(search) == 0 {
			return
		}

		// Look for an edge
		n = n.getEdge(search[0])
		if n == nil {
			return
		}

		// Consume the search prefix
		if !(len(search) >= len(n.prefix) && search[0:len(n.prefix)] == n.prefix) {
			// inlined version of !strings.HasPrefix(search, n.prefix)
			break
		}
		search = search[len(n.prefix):]
	}
}
