package list

import (
	"errors"
	"fmt"
	"slices"
)

// ErrCycle: a singly linked list given to a function that walks to its end loops back on itself.
// Wrapped with the parameter's name; test with errors.Is. The list is left unchanged.
var ErrCycle = errors.New("list: linked list has a cycle")

// Node is one node of a singly linked list: a value and the link to the next node, nil at the end.
// A list is named by its first node; nil is the empty list. Fields are exported so callers build and
// relink lists directly; every function of this package that walks to the end of a list checks for a
// cycle first (Floyd) and returns ErrCycle instead of looping, so a malformed list never hangs a call.
//
// ID:      ALG-LIST-001
// Time:    O(1) per field access
// Space:   one node per element
// Source:  Knuth, The Art of Computer Programming, Vol. 1, 3rd ed., Addison-Wesley, 1997, §2.2.3 (linked allocation: INFO and LINK fields, Λ the null link)
// Source:  Cormen, Leiserson, Rivest, Stein, Introduction to Algorithms, 3rd ed., MIT Press, 2009, §10.2 (singly linked lists)
// Tags:    linked-list, singly-linked, node, data-structure, generic
type Node[T any] struct {
	Value T
	Next  *Node[T]
}

// FromSlice returns a singly linked list holding values in order (nil for an empty slice).
// The nodes are allocated as one block, so the block stays reachable while any of its nodes is.
// values is read, never modified.
//
// ID:      ALG-LIST-002
// Time:    O(n)   n = len(values)
// Space:   O(n) for the result, one allocation
// Source:  Knuth, The Art of Computer Programming, Vol. 1, 3rd ed., Addison-Wesley, 1997, §2.2.3 (linked allocation)
// Tags:    linked-list, singly-linked, build, generic
func FromSlice[T any](values []T) *Node[T] {
	if len(values) == 0 {
		return nil
	}
	nodes := make([]Node[T], len(values))
	for i, v := range values {
		nodes[i].Value = v
		if i+1 < len(nodes) {
			nodes[i].Next = &nodes[i+1]
		}
	}
	return &nodes[0]
}

// AppendValues appends the values of the list starting at head to dst, in list order, and returns the
// extended slice. The list is read, never modified; nil appends nothing.
// Errors: ErrCycle when the list has a cycle; dst is then returned with its length unchanged.
//
// ID:      ALG-LIST-003
// Time:    O(n)   n = nodes in the list (the cycle check reads about 1.5n links more)
// Space:   O(1) extra; O(n) for the result, one allocation only if dst lacks capacity
// Source:  Knuth, The Art of Computer Programming, Vol. 1, 3rd ed., Addison-Wesley, 1997, §2.2.3 (traversal by LINK)
// Source:  Knuth, The Art of Computer Programming, Vol. 2, 3rd ed., Addison-Wesley, 1997, §3.1, Exercise 6 (Floyd's cycle test, the guard)
// Tags:    linked-list, singly-linked, traversal, generic
func AppendValues[T any](dst []T, head *Node[T]) ([]T, error) {
	n, _, err := walk(head)
	if err != nil {
		return dst, fmt.Errorf("%w (head)", err)
	}
	if n == 0 {
		return dst, nil
	}
	out := slices.Grow(dst, n)
	for p := head; p != nil; p = p.Next {
		out = append(out, p.Value)
	}
	return out, nil
}

// Length returns the number of nodes in the list starting at head (0 for nil). The list is read, never
// modified.
// Errors: ErrCycle when the list has a cycle (the length is then 0).
//
// ID:      ALG-LIST-004
// Time:    O(n)   n = nodes in the list; about 1.5n link reads
// Space:   O(1)
// Source:  Knuth, The Art of Computer Programming, Vol. 2, 3rd ed., Addison-Wesley, 1997, §3.1, Exercise 6 (Floyd's two-speed walk, R. W. Floyd's method as Knuth ascribes it)
// Tags:    linked-list, singly-linked, length, cycle-detection, generic
func Length[T any](head *Node[T]) (int, error) {
	n, _, err := walk(head)
	if err != nil {
		return 0, fmt.Errorf("%w (head)", err)
	}
	return n, nil
}

// walk returns the number of nodes of the list at head and its last node, or ErrCycle (unwrapped) when
// the list loops. A fast pointer takes two links per step and a slow one takes one; the fast one reaches
// the end of an acyclic list, and on a cyclic list the two meet (Knuth Vol. 2 §3.1 Ex. 6). The fast
// pointer's position counts the nodes, so no second pass is needed.
func walk[T any](head *Node[T]) (n int, last *Node[T], err error) {
	slow, fast := head, head
	for fast != nil {
		next := fast.Next
		if next == nil {
			return n + 1, fast, nil
		}
		if next.Next == nil {
			return n + 2, next, nil
		}
		fast = next.Next
		slow = slow.Next
		n += 2
		if fast == slow {
			return 0, nil, ErrCycle
		}
	}
	return 0, nil, nil
}

// nodeAt returns the node at 0-based index i of an acyclic list known to have more than i nodes.
func nodeAt[T any](head *Node[T], i int) *Node[T] {
	for ; i > 0; i-- {
		head = head.Next
	}
	return head
}
