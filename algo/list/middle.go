package list

import "fmt"

// Middle returns the node at 0-based index ⌊n/2⌋ of the list starting at head, n its length: the middle
// node for odd n, the second of the two middle nodes for even n. nil returns nil. The list is read, never
// modified.
// Errors: ErrCycle when the list has a cycle (the result is then nil).
//
// ID:      ALG-LIST-010
// Time:    O(n)   n = nodes in the list; one pass, about 1.5n link reads
// Space:   O(1)
// Source:  Knuth, The Art of Computer Programming, Vol. 2, 3rd ed., Addison-Wesley, 1997, §3.1, Exercise 6 (Floyd's one-and-two-step walk; when the fast pointer ends after k steps the slow one is at index k = ⌊n/2⌋)
// Tags:    linked-list, middle, two-pointer, tortoise-hare, generic
func Middle[T any](head *Node[T]) (*Node[T], error) {
	slow, fast := head, head
	for fast != nil && fast.Next != nil {
		slow = slow.Next
		fast = fast.Next.Next
		if slow == fast { // on an acyclic list fast stays strictly ahead
			return nil, fmt.Errorf("%w (head)", ErrCycle)
		}
	}
	return slow, nil
}
