package tree

// RightSideView appends to dst the last value of each level of the tree at root, from the root level
// down (what is seen looking at the tree from its right side), and returns the extended slice.
// dst == nil allocates. The tree is read, never modified.
//
// ID:      ALG-TREE-022
// Time:    O(n)   n = nodes
// Space:   O(w) extra (the current and next levels)   w = the largest number of nodes on one level; O(h + 1) appended to dst, h = height
// Source:  Cormen, Leiserson, Rivest, Stein, Introduction to Algorithms, 3rd ed., MIT Press, 2009, §22.2 (breadth-first search by distance from the root)
// Tags:    tree, binary-tree, right-side-view, level-order, bfs
func RightSideView[T any](dst []T, root *Node[T]) []T {
	if root == nil {
		return dst
	}
	level := []*Node[T]{root}
	var next []*Node[T]
	for len(level) > 0 {
		dst = append(dst, level[len(level)-1].Value)
		next = next[:0]
		for _, n := range level {
			if n.Left != nil {
				next = append(next, n.Left)
			}
			if n.Right != nil {
				next = append(next, n.Right)
			}
		}
		level, next = next, level
	}
	return dst
}
