package list

// DetectCycle reports whether the list starting at head loops back on itself. nil is acyclic.
// The list is read, never modified.
//
// ID:      ALG-LIST-008
// Time:    O(μ + λ)   μ = nodes before the cycle (n for an acyclic list), λ = cycle length
// Space:   O(1)
// Source:  Knuth, The Art of Computer Programming, Vol. 2, 3rd ed., Addison-Wesley, 1997, §3.1, Exercise 6 and its answer (X_n = X_2n; the method Knuth ascribes to R. W. Floyd)
// Tags:    linked-list, cycle-detection, floyd, tortoise-hare, two-pointer, generic
func DetectCycle[T any](head *Node[T]) bool {
	return meetingPoint(head) != nil
}

// CycleStart returns the first node of the list's cycle (the node two links point to) and the cycle's
// length; (nil, 0) when the list starting at head is acyclic or nil. The list is read, never modified.
//
// ID:      ALG-LIST-009
// Time:    O(μ + λ)   μ = nodes before the cycle, λ = cycle length
// Space:   O(1)
// Source:  Knuth, The Art of Computer Programming, Vol. 2, 3rd ed., Addison-Wesley, 1997, §3.1, Exercise 6 and its answer (μ by restarting one pointer at X_0, λ by one lap; Floyd's method)
// Tags:    linked-list, cycle-detection, floyd, tortoise-hare, two-pointer, generic
func CycleStart[T any](head *Node[T]) (start *Node[T], cycleLen int) {
	meet := meetingPoint(head)
	if meet == nil {
		return nil, 0
	}
	// The meeting point is a multiple of λ steps from head, so stepping one pointer from head and one
	// from the meeting point together brings both to the cycle's first node after μ steps.
	slow := head
	for slow != meet {
		slow = slow.Next
		meet = meet.Next
	}
	cycleLen = 1
	for p := slow.Next; p != slow; p = p.Next {
		cycleLen++
	}
	return slow, cycleLen
}

// meetingPoint returns a node where Floyd's slow (one link) and fast (two links) pointers meet, or nil
// when the fast pointer reaches the end.
func meetingPoint[T any](head *Node[T]) *Node[T] {
	slow, fast := head, head
	for fast != nil && fast.Next != nil {
		slow = slow.Next
		fast = fast.Next.Next
		if slow == fast {
			return slow
		}
	}
	return nil
}
