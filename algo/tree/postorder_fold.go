package tree

// foldFrame is one pending node of foldPostorder: stage 0 visits the left child, 1 the right, 2 the node.
type foldFrame[T any] struct {
	n     *Node[T]
	stage uint8
}

// foldPostorder computes a value for every subtree bottom-up without recursion: an empty subtree has the
// value empty, and a node's value is f(node, value of left, value of right). It returns the root's value
// and true, or, as soon as f returns false, that call's value and false (the walk stops there).
// O(n) time, O(h) extra (two explicit stacks), n = nodes, h = height.
func foldPostorder[T, V any](root *Node[T], empty V, f func(n *Node[T], left, right V) (V, bool)) (V, bool) {
	if root == nil {
		return empty, true
	}
	stack := []foldFrame[T]{{n: root}}
	var vals []V
	for len(stack) > 0 {
		top := &stack[len(stack)-1]
		switch top.stage {
		case 0:
			top.stage = 1
			if l := top.n.Left; l != nil {
				stack = append(stack, foldFrame[T]{n: l})
			} else {
				vals = append(vals, empty)
			}
		case 1:
			top.stage = 2
			if r := top.n.Right; r != nil {
				stack = append(stack, foldFrame[T]{n: r})
			} else {
				vals = append(vals, empty)
			}
		default:
			right, left := vals[len(vals)-1], vals[len(vals)-2]
			vals = vals[:len(vals)-2]
			v, ok := f(top.n, left, right)
			if !ok {
				return v, false
			}
			stack = stack[:len(stack)-1]
			vals = append(vals, v)
		}
	}
	return vals[0], true
}
