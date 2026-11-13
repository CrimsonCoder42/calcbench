package list

import (
	"errors"
	"fmt"
)

// ErrGroupSize: the group size k is less than 1.
var ErrGroupSize = errors.New("list: group size must be at least 1")

// ReverseKGroupInPlace reverses each consecutive block of k nodes of the list starting at head and returns
// the new head; a last block of fewer than k nodes keeps its order. The caller's nodes are relinked; none
// is allocated. k = 1, or k larger than the list, leaves it unchanged. nil returns nil.
// Errors: ErrGroupSize when k < 1; ErrCycle when the list has a cycle. On error the list is unchanged.
//
// ID:      ALG-LIST-007
// Time:    O(n)   n = nodes in the list
// Space:   O(1) extra
// Source:  Knuth, The Art of Computer Programming, Vol. 1, 3rd ed., Addison-Wesley, 1997, §2.2.3, Exercise 7 (inversion of a linked list, applied to each block)
// Tags:    linked-list, reverse, in-place, blocks, generic
func ReverseKGroupInPlace[T any](head *Node[T], k int) (*Node[T], error) {
	if k < 1 {
		return head, fmt.Errorf("%w: k = %d", ErrGroupSize, k)
	}
	n, _, err := walk(head)
	if err != nil {
		return head, fmt.Errorf("%w (head)", err)
	}
	if k == 1 || n < k {
		return head, nil
	}
	var dummy Node[T]
	dummy.Next = head
	tail := &dummy // last node of the part already done
	for blocks := n / k; blocks > 0; blocks-- {
		first := tail.Next
		stop := first
		for range k {
			stop = stop.Next
		}
		tail.Next = reverseRun(first, stop) // first is now the block's last node, linked to stop
		tail = first
	}
	return dummy.Next, nil
}
