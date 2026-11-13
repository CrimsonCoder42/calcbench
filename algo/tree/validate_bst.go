package tree

import (
	"cmp"
	"errors"
	"fmt"
)

// ErrNotBST: the tree's inorder walk is not strictly ascending, so it is not a binary search tree.
var ErrNotBST = errors.New("tree: not a binary search tree")

// ValidateBST returns nil when the tree at root is a binary search tree with distinct keys in cmp.Compare
// order (every left descendant less than its ancestor, every right descendant greater), and ErrNotBST,
// naming the first two inorder neighbours out of order, when it is not. The empty tree is valid.
// The tree is read, never modified; the walk stops at the first violation.
//
// ID:      ALG-TREE-018
// Time:    O(n)   n = nodes
// Space:   O(h) extra (the inorder stack)   h = height
// Source:  Cormen, Leiserson, Rivest, Stein, Introduction to Algorithms, 3rd ed., MIT Press, 2009, §12.1 (binary-search-tree property; Theorem 12.1, the inorder walk of a BST is sorted)
// Tags:    tree, bst, validate, binary-search-tree, inorder
func ValidateBST[T cmp.Ordered](root *Node[T]) error {
	first := true
	var prev T
	i := 0
	for v := range Inorder(root) {
		if !first && cmp.Compare(prev, v) >= 0 {
			return fmt.Errorf("%w: inorder position %d holds %v after %v", ErrNotBST, i, v, prev)
		}
		first = false
		prev = v
		i++
	}
	return nil
}
