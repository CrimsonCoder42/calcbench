package list

import (
	"cmp"
	"errors"
	"fmt"
)

// ErrSharedNodes: the two lists given to a merge share nodes (one is a tail of the other, or both end in
// the same nodes); relinking them would build a cycle.
var ErrSharedNodes = errors.New("list: the two lists share nodes")

// MergeTwoSortedInPlace merges two ascending lists into one ascending list by relinking their nodes and
// returns its head. Stable: equal values keep their order, those of a before those of b. Values are
// compared with cmp.Less (a NaN sorts before every other value, D-002). Lists that are not ascending
// are merged all the same (no error); the result then holds every node once, in no promised order.
// Both lists' nodes are modified; none is allocated.
// Errors: ErrCycle when either list has a cycle; ErrSharedNodes when the lists share nodes. On error
// neither list is changed and the result is nil.
//
// ID:      ALG-LIST-011
// Time:    O(m + n)   m, n = nodes in a and b; at most m + n − 1 comparisons
// Space:   O(1) extra
// Source:  Knuth, The Art of Computer Programming, Vol. 3, 2nd ed., Addison-Wesley, 1998, §5.2.4, Algorithm L (list merge sort; its merge relinks two sorted linked lists)
// Source:  Cormen, Leiserson, Rivest, Stein, Introduction to Algorithms, 3rd ed., MIT Press, 2009, §2.3.1 (MERGE)
// Tags:    linked-list, merge, sorted, stable, in-place, generic
func MergeTwoSortedInPlace[T cmp.Ordered](a, b *Node[T]) (*Node[T], error) {
	_, lastA, err := walk(a)
	if err != nil {
		return nil, fmt.Errorf("%w (a)", err)
	}
	_, lastB, err := walk(b)
	if err != nil {
		return nil, fmt.Errorf("%w (b)", err)
	}
	if lastA != nil && lastA == lastB { // two acyclic lists share a node iff they share the last one
		return nil, ErrSharedNodes
	}
	var dummy Node[T]
	tail := &dummy
	for a != nil && b != nil {
		if cmp.Less(b.Value, a.Value) {
			tail.Next, b = b, b.Next
		} else {
			tail.Next, a = a, a.Next
		}
		tail = tail.Next
	}
	if a != nil {
		tail.Next = a
	} else {
		tail.Next = b
	}
	return dummy.Next, nil
}
