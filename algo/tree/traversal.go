package tree

import "iter"

// Inorder returns an iterator over the values of the tree at root in inorder (left subtree, node, right
// subtree). For a binary search tree this is ascending order. Breaking out of the range loop stops the
// walk at once. The tree is read, never modified; each range over the iterator walks it afresh.
//
// ID:      ALG-TREE-006
// Time:    O(n) for a full walk, O(h + k) to yield the first k values   n = nodes, h = height
// Space:   O(h) extra (an explicit stack; no recursion, so a degenerate 1e5-node chain is safe)
// Source:  Knuth, The Art of Computer Programming, Vol. 1: Fundamental Algorithms, 3rd ed., Addison-Wesley, 1997, §2.3.1, Algorithm T (traverse binary tree in inorder, with a stack)
// Source:  The Go Programming Language Specification, go1.27, "For statements with range clause" (range over a function; yield returns false after a break and is not called again)
// Tags:    tree, binary-tree, traversal, inorder, iterator, dfs
func Inorder[T any](root *Node[T]) iter.Seq[T] {
	return func(yield func(T) bool) {
		var stack []*Node[T]
		n := root
		for n != nil || len(stack) > 0 {
			for n != nil {
				stack = append(stack, n)
				n = n.Left
			}
			n = stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			if !yield(n.Value) {
				return
			}
			n = n.Right
		}
	}
}

// Preorder returns an iterator over the values of the tree at root in preorder (node, left subtree,
// right subtree). Breaking out of the range loop stops the walk at once. The tree is read, never modified.
//
// ID:      ALG-TREE-007
// Time:    O(n) for a full walk   n = nodes
// Space:   O(h) extra (an explicit stack holding at most h + 1 pending right subtrees)   h = height
// Source:  Knuth, The Art of Computer Programming, Vol. 1: Fundamental Algorithms, 3rd ed., Addison-Wesley, 1997, §2.3.1, exercise 12 (preorder analogue of Algorithm T)
// Source:  The Go Programming Language Specification, go1.27, "For statements with range clause"
// Tags:    tree, binary-tree, traversal, preorder, iterator, dfs
func Preorder[T any](root *Node[T]) iter.Seq[T] {
	return func(yield func(T) bool) {
		if root == nil {
			return
		}
		stack := []*Node[T]{root}
		for len(stack) > 0 {
			n := stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			if !yield(n.Value) {
				return
			}
			if n.Right != nil {
				stack = append(stack, n.Right)
			}
			if n.Left != nil {
				stack = append(stack, n.Left)
			}
		}
	}
}

// Postorder returns an iterator over the values of the tree at root in postorder (left subtree, right
// subtree, node). Breaking out of the range loop stops the walk at once. The tree is read, never modified.
//
// ID:      ALG-TREE-008
// Time:    O(n) for a full walk   n = nodes
// Space:   O(h) extra (an explicit stack)   h = height
// Source:  Knuth, The Art of Computer Programming, Vol. 1: Fundamental Algorithms, 3rd ed., Addison-Wesley, 1997, §2.3.1, exercise 13 (postorder with a stack: remember whether the right subtree is done)
// Source:  The Go Programming Language Specification, go1.27, "For statements with range clause"
// Tags:    tree, binary-tree, traversal, postorder, iterator, dfs
func Postorder[T any](root *Node[T]) iter.Seq[T] {
	return func(yield func(T) bool) {
		var stack []*Node[T]
		var last *Node[T] // the node yielded most recently: tells whether top's right subtree is done
		n := root
		for n != nil || len(stack) > 0 {
			if n != nil {
				stack = append(stack, n)
				n = n.Left
				continue
			}
			top := stack[len(stack)-1]
			if top.Right != nil && top.Right != last {
				n = top.Right
				continue
			}
			if !yield(top.Value) {
				return
			}
			last = top
			stack = stack[:len(stack)-1]
		}
	}
}

// LevelOrder returns an iterator over the values of the tree at root level by level from the root,
// left to right within a level (breadth-first). Breaking out of the range loop stops the walk at once.
// The tree is read, never modified.
//
// ID:      ALG-TREE-009
// Time:    O(n) for a full walk   n = nodes
// Space:   O(w) extra (the current and next levels)   w = the largest number of nodes on one level
// Source:  Cormen, Leiserson, Rivest, Stein, Introduction to Algorithms, 3rd ed., MIT Press, 2009, §22.2 (breadth-first search: vertices in order of distance from the source)
// Source:  The Go Programming Language Specification, go1.27, "For statements with range clause"
// Tags:    tree, binary-tree, traversal, level-order, bfs, iterator
func LevelOrder[T any](root *Node[T]) iter.Seq[T] {
	return func(yield func(T) bool) {
		if root == nil {
			return
		}
		level := []*Node[T]{root}
		var next []*Node[T]
		for len(level) > 0 {
			next = next[:0]
			for _, n := range level {
				if !yield(n.Value) {
					return
				}
				if n.Left != nil {
					next = append(next, n.Left)
				}
				if n.Right != nil {
					next = append(next, n.Right)
				}
			}
			level, next = next, level
		}
	}
}
