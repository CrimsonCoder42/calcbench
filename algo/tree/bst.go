package tree

import "cmp"

// BST is an unbalanced binary search tree holding a set of distinct keys ordered by cmp.Compare (a NaN
// sorts before every other float and equals itself, -0 equals +0; doc/DECISIONS.md D-002). Every left
// descendant of a node is less than it and every right descendant greater.
//
// The zero value is an empty set, ready to use. Methods on a nil *BST never panic: Search and Delete
// report false, Len is 0, Root is nil, Insert reports false and stores nothing. Operations cost O(h), h the
// height: O(log n) expected for keys inserted in random order, n in the worst case (sorted insertion
// builds a chain); every method is iterative, so a chain of any length is safe.
// A BST is not safe for concurrent use.
type BST[T cmp.Ordered] struct {
	root *Node[T]
	size int
}

// Insert adds key to the set and returns true, or returns false when key is already present (the tree is
// then unchanged) or b is nil.
//
// ID:      ALG-TREE-001
// Time:    O(h)   h = height of the tree
// Space:   O(1) extra; one new node
// Source:  Cormen, Leiserson, Rivest, Stein, Introduction to Algorithms, 3rd ed., MIT Press, 2009, §12.3 (TREE-INSERT)
// Tags:    tree, bst, binary-search-tree, insert, set
func (b *BST[T]) Insert(key T) bool {
	if b == nil {
		return false
	}
	link := &b.root
	for *link != nil {
		switch c := cmp.Compare(key, (*link).Value); {
		case c < 0:
			link = &(*link).Left
		case c > 0:
			link = &(*link).Right
		default:
			return false
		}
	}
	*link = &Node[T]{Value: key}
	b.size++
	return true
}

// Search reports whether key is in the set.
//
// ID:      ALG-TREE-002
// Time:    O(h)   h = height of the tree
// Space:   O(1) extra, no allocation
// Source:  Cormen, Leiserson, Rivest, Stein, Introduction to Algorithms, 3rd ed., MIT Press, 2009, §12.2 (ITERATIVE-TREE-SEARCH)
// Tags:    tree, bst, binary-search-tree, search, contains, lookup
func (b *BST[T]) Search(key T) bool {
	if b == nil {
		return false
	}
	return searchNode(b.root, key) != nil
}

// searchNode returns the node holding key in the BST at n, or nil.
func searchNode[T cmp.Ordered](n *Node[T], key T) *Node[T] {
	for n != nil {
		switch c := cmp.Compare(key, n.Value); {
		case c < 0:
			n = n.Left
		case c > 0:
			n = n.Right
		default:
			return n
		}
	}
	return nil
}

// Delete removes key from the set and returns true, or returns false when key is absent. A node with two
// children is replaced by its inorder successor node (relinked, not copied), so every other node keeps
// its identity: a *Node obtained earlier from Root still holds the same key.
//
// ID:      ALG-TREE-003
// Time:    O(h)   h = height of the tree
// Space:   O(1) extra, no allocation
// Source:  Hibbard, "Some combinatorial properties of certain trees with applications to searching and sorting", Journal of the ACM, 1962, 9(1):13–28, doi:10.1145/321105.321108 (deletion by the successor)
// Source:  Cormen, Leiserson, Rivest, Stein, Introduction to Algorithms, 3rd ed., MIT Press, 2009, §12.3 (TRANSPLANT, TREE-DELETE)
// Tags:    tree, bst, binary-search-tree, delete, hibbard-deletion, set
func (b *BST[T]) Delete(key T) bool {
	if b == nil {
		return false
	}
	link := &b.root // the pointer that holds z
	for *link != nil {
		c := cmp.Compare(key, (*link).Value)
		if c == 0 {
			break
		}
		if c < 0 {
			link = &(*link).Left
		} else {
			link = &(*link).Right
		}
	}
	z := *link
	if z == nil {
		return false
	}
	switch {
	case z.Left == nil:
		*link = z.Right
	case z.Right == nil:
		*link = z.Left
	default:
		// The successor y is the leftmost node of z's right subtree; it has no left child.
		ylink := &z.Right
		for (*ylink).Left != nil {
			ylink = &(*ylink).Left
		}
		y := *ylink
		if ylink != &z.Right {
			*ylink = y.Right // y's parent adopts y's right subtree
			y.Right = z.Right
		}
		y.Left = z.Left
		*link = y
	}
	z.Left, z.Right = nil, nil
	b.size--
	return true
}

// Len returns the number of keys in the set (0 for a nil *BST).
//
// ID:      ALG-TREE-004
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Cormen, Leiserson, Rivest, Stein, Introduction to Algorithms, 3rd ed., MIT Press, 2009, ch. 12 (binary search trees)
// Tags:    tree, bst, size, length
func (b *BST[T]) Len() int {
	if b == nil {
		return 0
	}
	return b.size
}

// Root returns the root node of the tree (nil when empty), for use with the traversal and query
// functions of this package. The nodes belong to the BST: a caller that changes a Value, Left or Right
// through it breaks the search-tree order, and later results are then unspecified (never a panic).
//
// ID:      ALG-TREE-005
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Cormen, Leiserson, Rivest, Stein, Introduction to Algorithms, 3rd ed., MIT Press, 2009, §12.1 (T.root)
// Tags:    tree, bst, root
func (b *BST[T]) Root() *Node[T] {
	if b == nil {
		return nil
	}
	return b.root
}
