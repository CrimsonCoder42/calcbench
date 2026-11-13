package tree

import (
	"errors"
	"fmt"
)

// ErrEmptyTree: the tree is empty and the function has no value for it.
var ErrEmptyTree = errors.New("tree: tree is empty")

// MaxPathSum returns the largest sum of the values along any path of the tree at root, where a path is a
// non-empty sequence of nodes joined by edges, each node used at most once (it need not pass through the
// root or end at a leaf). With all values negative it is the largest single value. The tree is read,
// never modified.
// Errors: ErrEmptyTree when root is nil; ErrNotFinite when a float value is NaN or ±Inf; ErrOverflow when
// the maximum does not fit T (only a too-large sum can occur, since a path's sum is formed as a value plus
// non-negative gains).
//
// ID:      ALG-TREE-021
// Time:    O(n)   n = nodes
// Space:   O(h) extra (iterative postorder; no recursion)   h = height
// Source:  Laaksonen, Competitive Programmer's Handbook, draft of 3 July 2018, §14.2 "Diameter", Algorithm 1, pp. 135–136 (every path has a highest node; per node the best downward path and the best path through it), weighted here with a downward gain clipped at 0, https://cses.fi/book/book.pdf (read 2026-11-12)
// Tags:    tree, binary-tree, max-path-sum, dynamic-programming, tree-dp
func MaxPathSum[T number](root *Node[T]) (T, error) {
	var zero T
	if root == nil {
		return zero, ErrEmptyTree
	}
	var best T
	haveBest := false
	var failure error
	_, ok := foldPostorder(root, zero, func(n *Node[T], left, right T) (T, bool) {
		v := n.Value
		if !finite(v) {
			failure = fmt.Errorf("%w: node value %v", ErrNotFinite, v)
			return zero, false
		}
		l, r := max(left, 0), max(right, 0) // a path may stop at n instead of continuing down
		through, ok1 := addChecked(v, l)
		if ok1 {
			through, ok1 = addChecked(through, r)
		}
		gain, ok2 := addChecked(v, max(l, r))
		if !ok1 || !ok2 {
			failure = fmt.Errorf("%w: a path through the node holding %v", ErrOverflow, v)
			return zero, false
		}
		if !haveBest || through > best {
			best, haveBest = through, true
		}
		return gain, true
	})
	if !ok {
		return zero, failure
	}
	return best, nil
}
