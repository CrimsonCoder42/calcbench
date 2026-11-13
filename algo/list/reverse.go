package list

import "fmt"

// ReverseInPlace reverses the list starting at head by relinking its nodes and returns the new head (the
// old last node). The caller's nodes are modified; no node is allocated. nil returns nil.
// Errors: ErrCycle when the list has a cycle; the list is then left unchanged.
//
// ID:      ALG-LIST-006
// Time:    O(n)   n = nodes in the list (one guard pass, one relinking pass)
// Space:   O(1) extra
// Source:  Knuth, The Art of Computer Programming, Vol. 1, 3rd ed., Addison-Wesley, 1997, §2.2.3, Exercise 7 (invert a linked linear list by changing its links)
// Source:  Cormen, Leiserson, Rivest, Stein, Introduction to Algorithms, 3rd ed., MIT Press, 2009, Exercise 10.2-7 (Θ(n) nonrecursive reversal in constant storage)
// Tags:    linked-list, reverse, in-place, generic
func ReverseInPlace[T any](head *Node[T]) (*Node[T], error) {
	if _, _, err := walk(head); err != nil {
		return head, fmt.Errorf("%w (head)", err)
	}
	return reverseRun(head, nil), nil
}

// reverseRun reverses the acyclic run from head up to (not including) stop and returns the run's new
// first node; the old head, now last, links to stop.
func reverseRun[T any](head, stop *Node[T]) *Node[T] {
	prev := stop
	for p := head; p != stop; {
		next := p.Next
		p.Next = prev
		prev, p = p, next
	}
	return prev
}
