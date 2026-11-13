package list

import (
	"cmp"
	"fmt"
)

// PartitionInPlace relinks the list starting at head so that every node whose value is less than pivot
// comes before every node whose value is not, and returns the new head. Stable: within each part the
// nodes keep their input order. Values are compared with cmp.Less (a NaN is less than every other value;
// nothing is less than a NaN pivot, so a NaN pivot leaves the order unchanged). The caller's nodes are
// relinked; none is allocated. nil returns nil.
// Errors: ErrCycle when the list has a cycle; the list is then unchanged.
//
// ID:      ALG-LIST-014
// Time:    O(n)   n = nodes in the list; n comparisons
// Space:   O(1) extra
// Source:  Knuth, The Art of Computer Programming, Vol. 1, 3rd ed., Addison-Wesley, 1997, §2.2.3 (linked allocation: building two lists by relinking, then joining them)
// Source:  Cormen, Leiserson, Rivest, Stein, Introduction to Algorithms, 3rd ed., MIT Press, 2009, §7.1 (PARTITION around a pivot; the linked form keeps order in O(1) space)
// Tags:    linked-list, partition, stable, pivot, in-place, generic
func PartitionInPlace[T cmp.Ordered](head *Node[T], pivot T) (*Node[T], error) {
	if _, _, err := walk(head); err != nil {
		return head, fmt.Errorf("%w (head)", err)
	}
	var lessHead, restHead Node[T]
	lessTail, restTail := &lessHead, &restHead
	for p := head; p != nil; p = p.Next {
		if cmp.Less(p.Value, pivot) {
			lessTail.Next = p
			lessTail = p
		} else {
			restTail.Next = p
			restTail = p
		}
	}
	restTail.Next = nil
	lessTail.Next = restHead.Next
	return lessHead.Next, nil
}
