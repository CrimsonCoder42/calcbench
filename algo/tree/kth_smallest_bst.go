package tree

import "fmt"

// KthSmallestBST returns the k-th smallest key (k = 1 is the minimum) of the binary search tree at root,
// read off its inorder walk, which stops at the k-th node. For a BST from this package pass b.Root().
// The tree is read, never modified.
// Errors: ErrRange when k < 1 or k exceeds the number of nodes.
//
// ID:      ALG-TREE-017
// Time:    O(h + k)   h = height, k = the rank asked for
// Space:   O(h) extra (the inorder stack)
// Source:  Cormen, Leiserson, Rivest, Stein, Introduction to Algorithms, 3rd ed., MIT Press, 2009, §12.1, Theorem 12.1 (inorder walk of a BST yields its keys in sorted order)
// Source:  Knuth, The Art of Computer Programming, Vol. 1: Fundamental Algorithms, 3rd ed., Addison-Wesley, 1997, §2.3.1, Algorithm T
// Tags:    tree, bst, kth-smallest, order-statistic, rank, inorder
func KthSmallestBST[T any](root *Node[T], k int) (T, error) {
	var zero T
	if k < 1 {
		return zero, fmt.Errorf("%w: k = %d, want k >= 1", ErrRange, k)
	}
	seen := 0
	for v := range Inorder(root) {
		seen++
		if seen == k {
			return v, nil
		}
	}
	return zero, fmt.Errorf("%w: k = %d, tree has %d nodes", ErrRange, k, seen)
}
