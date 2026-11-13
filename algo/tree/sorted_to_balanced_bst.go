package tree

import (
	"cmp"
	"errors"
	"fmt"
)

// ErrNotSorted: the keys are not strictly ascending in cmp.Compare order.
var ErrNotSorted = errors.New("tree: keys are not strictly ascending")

// SortedToBalancedBST returns a BST holding the keys of sorted, perfectly balanced: at every node the
// two subtrees' sizes differ by at most one (the left one gets ⌊m/2⌋ of a subtree's m nodes), so the
// height is ⌊log2 n⌋. sorted must be strictly ascending in cmp.Compare order; it is read, never modified.
// Empty input gives an empty BST.
// Errors: ErrNotSorted when some sorted[i] >= sorted[i+1].
//
// ID:      ALG-TREE-016
// Time:    O(n)   n = len(sorted)
// Space:   O(log n) extra (an explicit stack of ranges); O(n) for the result (one allocation of n nodes)
// Source:  Wirth, Algorithms and Data Structures, Oberon version, 2004 (author's edition AD2012), §4.4.1, pp. 151–152 (perfectly balanced tree: root, n DIV 2 nodes left, n − n DIV 2 − 1 right)
// Source:  Cormen, Leiserson, Rivest, Stein, Introduction to Algorithms, 3rd ed., MIT Press, 2009, §12.1 (an inorder walk of a BST lists its keys in sorted order)
// Tags:    tree, bst, balanced, sorted, construct, build
func SortedToBalancedBST[T cmp.Ordered](sorted []T) (*BST[T], error) {
	for i := 1; i < len(sorted); i++ {
		if cmp.Compare(sorted[i-1], sorted[i]) >= 0 {
			return nil, fmt.Errorf("%w: sorted[%d] = %v, sorted[%d] = %v", ErrNotSorted, i-1, sorted[i-1], i, sorted[i])
		}
	}
	n := len(sorted)
	if n == 0 {
		return &BST[T]{}, nil
	}
	nodes := make([]Node[T], n)
	for i, v := range sorted {
		nodes[i].Value = v
	}
	// span [lo, hi) of sorted becomes the subtree hung on *link; its root is index lo + (hi−lo)/2, which
	// leaves (hi−lo)/2 nodes on the left as in Wirth's rule.
	type span struct {
		lo, hi int
		link   **Node[T]
	}
	var root *Node[T]
	stack := []span{{0, n, &root}}
	for len(stack) > 0 {
		s := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		mid := s.lo + (s.hi-s.lo)/2
		nd := &nodes[mid]
		*s.link = nd
		if mid+1 < s.hi {
			stack = append(stack, span{mid + 1, s.hi, &nd.Right})
		}
		if s.lo < mid {
			stack = append(stack, span{s.lo, mid, &nd.Left})
		}
	}
	return &BST[T]{root: root, size: n}, nil
}
