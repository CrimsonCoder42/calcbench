package heap

import (
	"cmp"
	"errors"
)

// ErrEmpty: there is no element to take a median of.
var ErrEmpty = errors.New("heap: no elements")

// MedianStream keeps the running median of values added one at a time. The zero value is an empty,
// ready-to-use stream. Ordering is cmp.Compare's total order (NaN least, doc/DECISIONS.md D-002).
// Values are kept, so memory grows with the count. Not safe for concurrent use.
//
// The lower half sits in a max-heap and the upper half in a min-heap; the lower half holds the extra
// value when the count is odd (Sedgewick & Wayne, §2.4, Exercise 2.4.30 "Dynamic median-finding").
type MedianStream[T cmp.Ordered] struct {
	lower Heap[T] // max-heap: every value <= every value in upper
	upper Heap[T] // min-heap
}

// NewMedianStream returns an empty stream with room for capacityHint values before it grows
// (a negative hint is treated as 0). The zero MedianStream is equally valid; this only presizes it.
// Errors: ErrTooLarge when the two heaps' buffers (capacityHint/2 + 1 values of T each) would exceed
// 2⁴⁷ bytes (2³¹ on 32-bit), the most one allocation may take (rule L41).
//
// ID:      ALG-HEAP-010
// Time:    O(1)
// Space:   O(c) for the two heaps' buffers, c = max(capacityHint, 0)
// Source:  Sedgewick, Wayne, Algorithms, 4th ed., Addison-Wesley, 2011, §2.4, Exercise 2.4.30 "Dynamic median-finding" (a max-heap below the median, a min-heap above; algs4.cs.princeton.edu/24pq, read 2026-11-12)
// Tags:    heap, median, running-median, streaming, two-heaps, order-statistic
func NewMedianStream[T cmp.Ordered](capacityHint int) (*MedianStream[T], error) {
	c := max(capacityHint, 0)
	if err := checkAlloc[T]("capacityHint/2+1", c/2+1); err != nil {
		return nil, err
	}
	m := &MedianStream[T]{}
	m.lower.items = make([]T, 0, c/2+1)
	m.upper.items = make([]T, 0, c/2+1)
	m.init()
	return m, nil
}

func (m *MedianStream[T]) init() {
	if m.lower.less == nil {
		m.lower.less = func(a, b T) bool { return cmp.Less(b, a) }
		m.upper.less = cmp.Less[T]
	}
}

// Add puts v into the stream and returns true (false, storing nothing, on a nil *MedianStream).
//
// ID:      ALG-HEAP-011
// Time:    O(log n)   n = values added so far; amortized when a buffer grows
// Space:   O(1) extra amortized; the stream keeps every value, O(n) in all
// Source:  Sedgewick, Wayne, Algorithms, 4th ed., Addison-Wesley, 2011, §2.4, Exercise 2.4.30 "Dynamic median-finding" (insert in logarithmic time)
// Tags:    heap, median, running-median, streaming, insert
func (m *MedianStream[T]) Add(v T) bool {
	if m == nil {
		return false
	}
	m.init()
	if top, ok := m.lower.Peek(); !ok || cmp.Compare(v, top) <= 0 {
		m.lower.Push(v)
	} else {
		m.upper.Push(v)
	}
	// Restore len(lower) - len(upper) ∈ {0, 1}; each Push moved the difference by one.
	if m.lower.Len() > m.upper.Len()+1 {
		x, _ := m.lower.Pop()
		m.upper.Push(x)
	} else if m.upper.Len() > m.lower.Len() {
		x, _ := m.upper.Pop()
		m.lower.Push(x)
	}
	return true
}

// Median returns the lower and upper medians of the values added so far: the ⌊(n+1)/2⌋-th and
// ⌈(n+1)/2⌉-th smallest, equal when n is odd (the same shape as ALG-SELECT-002 Median).
// Errors: ErrEmpty when nothing has been added (or m is nil).
//
// ID:      ALG-HEAP-012
// Time:    O(1)
// Space:   O(1)
// Source:  Sedgewick, Wayne, Algorithms, 4th ed., Addison-Wesley, 2011, §2.4, Exercise 2.4.30 "Dynamic median-finding" (find the median in constant time)
// Source:  Cormen, Leiserson, Rivest, Stein, Introduction to Algorithms, 3rd ed., MIT Press, 2009, ch. 9 introduction (lower and upper median)
// Tags:    heap, median, running-median, streaming, order-statistic
func (m *MedianStream[T]) Median() (lower, upper T, err error) {
	if m == nil || m.lower.Len() == 0 {
		return lower, upper, ErrEmpty
	}
	lower, _ = m.lower.Peek()
	if m.lower.Len() > m.upper.Len() {
		return lower, lower, nil
	}
	upper, _ = m.upper.Peek()
	return lower, upper, nil
}

// Len returns the number of values added so far (0 for a nil *MedianStream).
//
// ID:      ALG-HEAP-013
// Time:    O(1)
// Space:   O(1)
// Source:  Sedgewick, Wayne, Algorithms, 4th ed., Addison-Wesley, 2011, §2.4, Exercise 2.4.30 "Dynamic median-finding"
// Tags:    heap, median, streaming, size
func (m *MedianStream[T]) Len() int {
	if m == nil {
		return 0
	}
	return m.lower.Len() + m.upper.Len()
}
