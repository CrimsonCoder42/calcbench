package tree

// Height returns the number of edges on the longest path from root down to a leaf: −1 for the empty tree,
// 0 for a single node. The maximum depth counted in nodes is Height + 1. The tree is read, never modified.
//
// ID:      ALG-TREE-010
// Time:    O(n)   n = nodes
// Space:   O(h) extra (iterative postorder; no recursion)   h = the result
// Source:  Cormen, Leiserson, Rivest, Stein, Introduction to Algorithms, 3rd ed., MIT Press, 2009, Appendix B.5.3 (height of a rooted tree: edges on the longest simple downward path to a leaf)
// Tags:    tree, binary-tree, height, depth, max-depth
func Height[T any](root *Node[T]) int {
	h, _ := foldPostorder(root, -1, func(_ *Node[T], left, right int) (int, bool) {
		return 1 + max(left, right), true
	})
	return h
}
