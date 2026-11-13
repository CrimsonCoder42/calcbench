package list

import "fmt"

// ReorderInPlace relinks the list L0, L1, …, Ln−1 starting at head into the order L0, Ln−1, L1, Ln−2, …
// (first, last, second, second-to-last, …) and returns its head, which is still L0. The caller's nodes are
// relinked; none is allocated. nil and lists of one or two nodes are unchanged.
// Errors: ErrCycle when the list has a cycle; the list is then unchanged.
//
// ID:      ALG-LIST-013
// Time:    O(n)   n = nodes in the list
// Space:   O(1) extra
// Source:  Knuth, The Art of Computer Programming, Vol. 1, 3rd ed., Addison-Wesley, 1997, §2.2.3, Exercise 7 (inversion of the second half) and §2.2.3 (relinking); the split is ALG-LIST-010's walk. Proof of the order: catalog/list.md ALG-LIST-013
// Tags:    linked-list, reorder, interleave, two-pointer, in-place, generic
func ReorderInPlace[T any](head *Node[T]) (*Node[T], error) {
	n, _, err := walk(head)
	if err != nil {
		return head, fmt.Errorf("%w (head)", err)
	}
	if n < 3 {
		return head, nil
	}
	// First half: L0 … L⌈n/2⌉−1; second half reversed: Ln−1 … L⌈n/2⌉.
	firstEnd := nodeAt(head, (n-1)/2)
	second := reverseRun(firstEnd.Next, nil)
	firstEnd.Next = nil
	for p := head; second != nil; {
		pNext, sNext := p.Next, second.Next
		p.Next = second
		second.Next = pNext
		p, second = pNext, sNext
	}
	return head, nil
}
