package bplus

func NewTree() (*Tree, error) {
	bpt := new(Tree)
	return bpt, nil
}

// Has returns a boolean indicating weather or not
// the provided key and associated record exists.
func (t *Tree) Has(k keyType) bool {
	return t.findEntry(k) != nil
}

// Add inserts a new record using the provided key. It
// only inserts an record if the key does not already exist.
func (t *Tree) Add(k keyType, v valType) {
	// master insertUnique method only inserts if the key
	// does not currently exist in the tree
	t.insertUnique(k, v)
}

// Put is mainly used when you wish to upsert as it assumes the
// data to already be contained the tree. It will  overwrite
// duplicate keys, as it does not check to see if the key exists
func (t *Tree) Put(k keyType, v valType) bool {
	// master insert method treats insertion much like
	// "setting" in a hashmap (an upsert) by default
	return t.insert(k, v)
}

// Get returns the record for a given key if it exists
func (t *Tree) Get(k keyType) (keyType, valType) {
	e := t.findEntry(k)
	if e == nil {
		return *new(keyType), *new(valType)
	}
	return e.Key, e.Value
}

// Del removes the record for the supplied key and attempts
// to return the previous key and value
func (t *Tree) Del(k keyType) (keyType, valType) {
	e := t.delete(k)
	if e == nil {
		return *new(keyType), *new(valType)
	}
	return e.Key, e.Value
}

// Range provides a simple iteration function for the tree
func (t *Tree) Range(iter func(k keyType, v valType) bool) {
	c := findFirstLeaf(t.root)
	if c == nil {
		return
	}
	var e *record
	for {
		for i := 0; i < c.numKeys; i++ {
			e = (*record)(c.ptrs[i])
			if e != nil && !iter(e.Key, e.Value) {
				continue
			}
		}
		if c.ptrs[order-1] != nil {
			c = (*node)(c.ptrs[order-1])
		} else {
			break
		}
	}
}

// Min returns the minimum (lowest) key and value pair in the tree
func (t *Tree) Min() (keyType, valType) {
	c := findFirstLeaf(t.root)
	if c == nil {
		return *new(keyType), *new(valType)
	}
	e := (*record)(c.ptrs[0])
	return e.Key, e.Value
}

// Max returns the maximum (highest) key and value pair in the tree
func (t *Tree) Max() (keyType, valType) {
	c := findLastLeaf(t.root)
	if c == nil {
		return *new(keyType), *new(valType)
	}
	e := (*record)(c.ptrs[c.numKeys-1])
	return e.Key, e.Value
}

// GetClosest attempts to return the closest match in the tree
// if an explicit match cannot be found
func (t *Tree) GetClosest(k keyType) (keyType, valType) {
	l := findLeaf(t.root, k)
	if l == nil {
		return *new(keyType), *new(valType)
	}
	e, ok := l.closest(k)
	if !ok {
		return *new(keyType), *new(valType)
	}
	return e.Key, e.Value
}

// Len returns the a count of the number of items in the tree
func (t *Tree) Len() int {
	var count int
	for n := findFirstLeaf(t.root); n != nil; n = n.nextLeaf() {
		count += n.numKeys
	}
	return count
}

// Close closes the tree
func (t *Tree) Close() {
	// t.destroyTree()
	t.root = nil
}
