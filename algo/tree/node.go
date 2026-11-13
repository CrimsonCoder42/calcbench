package tree

// Node is a node of a binary tree: a value and two possibly-nil subtrees. A nil *Node is the empty tree.
//
// Every function in this package that takes a *Node expects a tree: no node is reachable twice from the
// root. A structure with a shared node or a cycle breaks that precondition; it is not detected, and a
// traversal of a cyclic structure does not end. The functions here never modify a tree they are given
// unless their name or doc says so.
type Node[T any] struct {
	Value       T
	Left, Right *Node[T]
}
