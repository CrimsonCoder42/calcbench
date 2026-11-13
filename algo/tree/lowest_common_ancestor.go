package tree

import "cmp"

// LowestCommonAncestorBST returns the deepest node of the binary search tree at root whose subtree holds
// both keys p and q (a node is its own ancestor), and true; or nil and false when either key is absent.
// The order of p and q does not matter. The tree must be a BST in cmp.Compare order (as BST.Root gives);
// it is read, never modified.
//
// ID:      ALG-TREE-012
// Time:    O(h)   h = height of the tree
// Space:   O(1) extra, no allocation
// Source:  Aho, Hopcroft, Ullman, "On finding lowest common ancestors in trees", SIAM Journal on Computing, 1976, 5(1):115–132, doi:10.1137/0205011 (the problem)
// Source:  Cormen, Leiserson, Rivest, Stein, Introduction to Algorithms, 3rd ed., MIT Press, 2009, §12.2 (search path; the LCA is the first node on it with p <= key <= q, proved in doc/algorithms_core_2026-11-12/catalog/tree.md)
// Tags:    tree, bst, lowest-common-ancestor, lca
func LowestCommonAncestorBST[T cmp.Ordered](root *Node[T], p, q T) (*Node[T], bool) {
	if cmp.Less(q, p) {
		p, q = q, p
	}
	n := root
	for n != nil {
		if cmp.Less(q, n.Value) {
			n = n.Left
		} else if cmp.Less(n.Value, p) {
			n = n.Right
		} else {
			break
		}
	}
	if n == nil || searchNode(n, p) == nil || searchNode(n, q) == nil {
		return nil, false
	}
	return n, true
}

// lcaState is the per-subtree value of LowestCommonAncestor's fold: how many of p and q the subtree holds,
// and the ancestor once both are found.
type lcaState[T any] struct {
	found int
	lca   *Node[T]
}

// LowestCommonAncestor returns the deepest node of the binary tree at root whose subtree contains both
// nodes p and q (identified by pointer; a node is its own ancestor), or nil when p or q is nil or not in
// the tree. p == q returns p. The tree is read, never modified.
//
// ID:      ALG-TREE-013
// Time:    O(n)   n = nodes; the walk stops once both are found below one node
// Space:   O(h) extra (iterative postorder; no recursion)   h = height
// Source:  Aho, Hopcroft, Ullman, "On finding lowest common ancestors in trees", SIAM Journal on Computing, 1976, 5(1):115–132, doi:10.1137/0205011 (the problem)
// Source:  Laaksonen, Competitive Programmer's Handbook, draft of 3 July 2018, §18.3 "Lowest common ancestor", p. 167, https://cses.fi/book/book.pdf (read 2026-11-12)
// Tags:    tree, binary-tree, lowest-common-ancestor, lca
func LowestCommonAncestor[T any](root, p, q *Node[T]) *Node[T] {
	if p == nil || q == nil {
		return nil
	}
	st, _ := foldPostorder(root, lcaState[T]{}, func(n *Node[T], left, right lcaState[T]) (lcaState[T], bool) {
		found := left.found + right.found
		if n == p {
			found++
		}
		if n == q {
			found++
		}
		if found >= 2 {
			return lcaState[T]{found: 2, lca: n}, false // stop: n is the first node holding both
		}
		return lcaState[T]{found: found}, true
	})
	return st.lca
}
