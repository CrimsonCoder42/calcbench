package list

import (
	"errors"
	"fmt"
)

// ErrPosition: a position is outside the list (fromEnd < 1 or larger than the length).
var ErrPosition = errors.New("list: position is outside the list")

// RemoveNthFromEndInPlace unlinks the fromEnd-th node counted from the end (fromEnd = 1 is the last node)
// of the list starting at head and returns the list's head, nil when the only node was removed. The
// removed node's Next is set to nil. The caller's nodes are relinked; none is allocated.
// Errors: ErrPosition when fromEnd < 1 or fromEnd > the list's length; ErrCycle when the list has a cycle.
// On error the list is unchanged and head is returned.
//
// ID:      ALG-LIST-012
// Time:    O(n)   n = nodes in the list (one guard pass that also counts, one walk to the predecessor)
// Space:   O(1) extra
// Source:  Knuth, The Art of Computer Programming, Vol. 1, 3rd ed., Addison-Wesley, 1997, §2.2.3 (deletion from a linked list by relinking the predecessor)
// Tags:    linked-list, remove, two-pointer, in-place, generic
func RemoveNthFromEndInPlace[T any](head *Node[T], fromEnd int) (*Node[T], error) {
	n, _, err := walk(head)
	if err != nil {
		return head, fmt.Errorf("%w (head)", err)
	}
	if fromEnd < 1 || fromEnd > n {
		return head, fmt.Errorf("%w: fromEnd = %d, length = %d", ErrPosition, fromEnd, n)
	}
	idx := n - fromEnd // 0-based index of the node to remove
	if idx == 0 {
		next := head.Next
		head.Next = nil
		return next, nil
	}
	prev := nodeAt(head, idx-1)
	removed := prev.Next
	prev.Next = removed.Next
	removed.Next = nil
	return head, nil
}
