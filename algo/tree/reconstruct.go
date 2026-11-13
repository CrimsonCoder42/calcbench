package tree

import (
	"errors"
	"fmt"
	"iter"
)

var (
	// ErrLength: two traversals given together have different lengths.
	ErrLength = errors.New("tree: traversals have different lengths")
	// ErrDuplicate: a value occurs twice in a traversal; reconstruction needs distinct values.
	ErrDuplicate = errors.New("tree: duplicate value in traversal")
	// ErrInconsistent: the traversals are not those of one binary tree.
	ErrInconsistent = errors.New("tree: traversals do not describe one tree")
)

// FromPreIn returns the unique binary tree whose preorder is preorder and whose inorder is inorder.
// Values must be distinct; a float NaN never equals itself, so it can never be matched (ErrInconsistent).
// Empty inputs give the empty tree (nil, nil). The inputs are read, never modified.
// Errors: ErrLength when the lengths differ; ErrDuplicate when inorder repeats a value; ErrInconsistent
// when the two sequences are not the preorder and inorder of one tree (checked exactly).
//
// ID:      ALG-TREE-014
// Time:    O(n) expected (one hash set of the values; one stack pass; one check by both traversals)   n = len(inorder)
// Space:   O(n) extra (the hash set); O(n) for the result (one allocation of n nodes)
// Source:  Knuth, The Art of Computer Programming, Vol. 1: Fundamental Algorithms, 3rd ed., Addison-Wesley, 1997, §2.3.1, exercise 7 (preorder and inorder determine the binary tree)
// Source:  Laaksonen, Competitive Programmer's Handbook, draft of 3 July 2018, §14.4 "Binary trees", p. 139, https://cses.fi/book/book.pdf (read 2026-11-12)
// Tags:    tree, binary-tree, reconstruction, preorder, inorder, construct
func FromPreIn[T comparable](preorder, inorder []T) (*Node[T], error) {
	n := len(inorder)
	if len(preorder) != n {
		return nil, fmt.Errorf("%w: len(preorder) = %d, len(inorder) = %d", ErrLength, len(preorder), n)
	}
	if err := checkSameDistinct(preorder, inorder, "preorder"); err != nil {
		return nil, err
	}
	root := buildRootFirst(n, func(i int) T { return preorder[i] }, func(j int) T { return inorder[j] }, false)
	if !matches(Preorder(root), preorder) || !matches(Inorder(root), inorder) {
		return nil, fmt.Errorf("%w (preorder and inorder of %d values)", ErrInconsistent, n)
	}
	return root, nil
}

// FromInPost returns the unique binary tree whose inorder is inorder and whose postorder is postorder.
// Values must be distinct; a float NaN never equals itself, so it can never be matched (ErrInconsistent).
// Empty inputs give the empty tree (nil, nil). The inputs are read, never modified.
// Errors: ErrLength when the lengths differ; ErrDuplicate when inorder repeats a value; ErrInconsistent
// when the two sequences are not the inorder and postorder of one tree (checked exactly).
//
// ID:      ALG-TREE-015
// Time:    O(n) expected (one hash set of the values; one stack pass; one check by both traversals)   n = len(inorder)
// Space:   O(n) extra (the hash set); O(n) for the result (one allocation of n nodes)
// Source:  Laaksonen, Competitive Programmer's Handbook, draft of 3 July 2018, §14.4 "Binary trees", p. 139 (postorder and inorder also determine the tree), https://cses.fi/book/book.pdf (read 2026-11-12)
// Source:  Knuth, The Art of Computer Programming, Vol. 1: Fundamental Algorithms, 3rd ed., Addison-Wesley, 1997, §2.3.1, exercise 7 (the preorder case; postorder is its mirror read backwards)
// Tags:    tree, binary-tree, reconstruction, postorder, inorder, construct
func FromInPost[T comparable](inorder, postorder []T) (*Node[T], error) {
	n := len(inorder)
	if len(postorder) != n {
		return nil, fmt.Errorf("%w: len(inorder) = %d, len(postorder) = %d", ErrLength, n, len(postorder))
	}
	if err := checkSameDistinct(postorder, inorder, "postorder"); err != nil {
		return nil, err
	}
	// Postorder read backwards is root, right, left: the preorder of the mirror image, whose inorder is
	// inorder read backwards. Build the mirror with the children swapped back.
	root := buildRootFirst(n, func(i int) T { return postorder[n-1-i] }, func(j int) T { return inorder[n-1-j] }, true)
	if !matches(Inorder(root), inorder) || !matches(Postorder(root), postorder) {
		return nil, fmt.Errorf("%w (inorder and postorder of %d values)", ErrInconsistent, n)
	}
	return root, nil
}

// checkSameDistinct checks that inorder has distinct values and that other holds exactly the same ones.
func checkSameDistinct[T comparable](other, inorder []T, otherName string) error {
	seen := make(map[T]bool, len(inorder)) // value -> already matched in other
	for j, v := range inorder {
		if _, dup := seen[v]; dup {
			return fmt.Errorf("%w: inorder[%d] = %v", ErrDuplicate, j, v)
		}
		seen[v] = false
	}
	for i, v := range other {
		used, ok := seen[v]
		if !ok {
			return fmt.Errorf("%w: %s[%d] = %v is not in inorder", ErrInconsistent, otherName, i, v)
		}
		if used {
			return fmt.Errorf("%w: %s[%d] = %v", ErrDuplicate, otherName, i, v)
		}
		seen[v] = true
	}
	return nil
}

// buildRootFirst builds the tree from a root-first order rf(0..n-1) (preorder) and the matching
// left-to-right order in(0..n-1) (inorder) in one stack pass: each value becomes the left child of the
// stack top unless the top is the next inorder value, in which case finished nodes are popped and the
// value becomes the right child of the last one popped. mirror swaps left and right. It never panics on
// inconsistent input; the caller checks the result.
func buildRootFirst[T comparable](n int, rf, in func(int) T, mirror bool) *Node[T] {
	if n == 0 {
		return nil
	}
	nodes := make([]Node[T], n)
	nodes[0].Value = rf(0)
	stack := []*Node[T]{&nodes[0]}
	j := 0
	for i := 1; i < n; i++ {
		nd := &nodes[i]
		nd.Value = rf(i)
		top := stack[len(stack)-1]
		if top.Value != in(j) {
			if mirror {
				top.Right = nd
			} else {
				top.Left = nd
			}
		} else {
			var parent *Node[T]
			for len(stack) > 0 && j < n && stack[len(stack)-1].Value == in(j) {
				parent = stack[len(stack)-1]
				stack = stack[:len(stack)-1]
				j++
			}
			if mirror {
				parent.Left = nd
			} else {
				parent.Right = nd
			}
		}
		stack = append(stack, nd)
	}
	return &nodes[0]
}

// matches reports whether seq yields exactly the values of want, in order.
func matches[T comparable](seq iter.Seq[T], want []T) bool {
	i := 0
	for v := range seq {
		if i >= len(want) || v != want[i] {
			return false
		}
		i++
	}
	return i == len(want)
}
