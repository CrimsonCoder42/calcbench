package tree

// IsBalanced reports whether the tree at root is height-balanced: at every node the heights of the two
// subtrees differ by at most 1. The empty tree is balanced. The tree is read, never modified; the walk
// stops at the first unbalanced node.
//
// ID:      ALG-TREE-011
// Time:    O(n)   n = nodes
// Space:   O(h) extra (iterative postorder; no recursion)   h = height
// Source:  Knuth, The Art of Computer Programming, Vol. 3: Sorting and Searching, 2nd ed., Addison-Wesley, 1998, §6.2.3 (balanced trees: Adel'son-Vel'skii and Landis's condition, subtree heights differ by at most one)
// Tags:    tree, binary-tree, balanced, avl, height
func IsBalanced[T any](root *Node[T]) bool {
	_, ok := foldPostorder(root, -1, func(_ *Node[T], left, right int) (int, bool) {
		if left-right > 1 || right-left > 1 {
			return 0, false
		}
		return 1 + max(left, right), true
	})
	return ok
}
