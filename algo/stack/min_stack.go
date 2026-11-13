package stack

import "cmp"

// MinStack is a last-in, first-out stack that also reports its least element in O(1). Each slot keeps
// its value and the least value at or below it, so Min never searches. The least value is by
// cmp.Compare (a NaN is least). The zero value is an empty stack, ready to use. A MinStack is not safe
// for concurrent use.
//
// Methods: Push(v) adds v on top. Pop removes and returns the top. Top returns the top. Min returns the
// least element. Pop, Top and Min return (zero, false) on an empty stack. Len is the element count.
//
// ID:      ALG-STACK-008
// Time:    O(1) per Pop, Top, Min and Len; O(1) amortized per Push (slice growth)
// Space:   O(n) for n elements (two values per element)
// Source:  Gajewska, Tarjan, "Deques with heap order", Information Processing Letters, 1986, 22(4):197–200, doi:10.1016/0020-0190(86)90028-1
// Source:  Knuth, The Art of Computer Programming, Vol. 1, 3rd ed., Addison-Wesley, 1997, §2.2.1 (stacks)
// Tags:    stack, min-stack, find-min, data-structure, lifo
type MinStack[T cmp.Ordered] struct {
	slots []minSlot[T]
}

// minSlot is one element and the least element at or below it.
type minSlot[T cmp.Ordered] struct {
	value, min T
}

// Push adds v on top of the stack. O(1) amortized.
func (s *MinStack[T]) Push(v T) {
	least := v
	if n := len(s.slots); n > 0 && cmp.Less(s.slots[n-1].min, v) {
		least = s.slots[n-1].min
	}
	s.slots = append(s.slots, minSlot[T]{value: v, min: least})
}

// Pop removes and returns the top element, or (zero, false) when the stack is empty. O(1).
func (s *MinStack[T]) Pop() (T, bool) {
	n := len(s.slots)
	if n == 0 {
		var zero T
		return zero, false
	}
	v := s.slots[n-1].value
	s.slots[n-1] = minSlot[T]{} // release a string's memory
	s.slots = s.slots[:n-1]
	return v, true
}

// Top returns the top element without removing it, or (zero, false) when the stack is empty. O(1).
func (s *MinStack[T]) Top() (T, bool) {
	if n := len(s.slots); n > 0 {
		return s.slots[n-1].value, true
	}
	var zero T
	return zero, false
}

// Min returns the least element, or (zero, false) when the stack is empty. O(1).
func (s *MinStack[T]) Min() (T, bool) {
	if n := len(s.slots); n > 0 {
		return s.slots[n-1].min, true
	}
	var zero T
	return zero, false
}

// Len returns the number of elements. O(1).
func (s *MinStack[T]) Len() int {
	return len(s.slots)
}
