package window

import (
	"cmp"
	"fmt"
	"slices"
)

// stackRingMax is the largest windowLen whose index deque is a stack array (no allocation).
const stackRingMax = 64

// appendWindowExtremes appends the least (wantMax false) or greatest (wantMax true) element of each
// window to dst. It keeps a ring-buffer deque of indices whose values are strictly increasing (min) or
// strictly decreasing (max) from front to back, so the front is always the window's extreme
// (Lemire 2006). Each index is pushed and popped at most once: O(n) time in total, at most 2n
// comparisons, and the deque never holds more than windowLen indices.
func appendWindowExtremes[T cmp.Ordered](dst, values []T, windowLen int, wantMax bool) ([]T, error) {
	if windowLen < 1 {
		return dst, fmt.Errorf("%w: windowLen = %d", ErrWindowLen, windowLen)
	}
	n := len(values)
	if windowLen > n {
		return dst, nil
	}
	out := slices.Grow(dst, n-windowLen+1)
	var stackRing [stackRingMax]int
	var ring []int
	if windowLen <= stackRingMax {
		ring = stackRing[:windowLen]
	} else {
		ring = make([]int, windowLen)
	}
	head, size := 0, 0
	for i, v := range values {
		// The front is the oldest index; at most one index leaves per step.
		if size > 0 && ring[head] == i-windowLen {
			head++
			if head == windowLen {
				head = 0
			}
			size--
		}
		// Drop back entries that v makes useless: no better than v and leaving sooner.
		for size > 0 {
			back := head + size - 1
			if back >= windowLen {
				back -= windowLen
			}
			c := cmp.Compare(values[ring[back]], v)
			if wantMax {
				c = -c
			}
			if c < 0 {
				break
			}
			size--
		}
		tail := head + size
		if tail >= windowLen {
			tail -= windowLen
		}
		ring[tail] = i
		size++
		if i >= windowLen-1 {
			out = append(out, values[ring[head]])
		}
	}
	return out, nil
}

// MonotonicDeque is a double-ended queue whose elements are kept in monotone order from front to back:
// non-decreasing for a min-deque, non-increasing for a max-deque. PushBack first evicts from the back
// every element that the new one beats (strictly greater, for a min-deque), so the front is always the
// least (or greatest) of the elements pushed since they were last popped. It is the monotonic queue of
// the sliding-window minimum and maximum: push each value as it enters the window and, as a value
// leaves, call PopFrontIfEqual with it. Equal elements are all kept, which is what makes that
// value-only eviction correct. Order is cmp.Compare's, so NaN is the least value.
// The zero value is an empty min-deque, ready to use; NewMonotonicDeque makes a max-deque. Front,
// PopFront, PopFrontIfEqual and Len accept a nil *MonotonicDeque (it is empty); PushBack needs a
// non-nil one. A MonotonicDeque is not safe for concurrent use. This is the one home of the monotonic queue in
// calcbench.
type MonotonicDeque[T cmp.Ordered] struct {
	ring       []T // ring buffer; len(ring) is the capacity, a power of two or 0
	head, size int
	wantMax    bool
}

// NewMonotonicDeque returns an empty MonotonicDeque: a max-deque (front = greatest) when wantMax is
// true, otherwise a min-deque (front = least).
//
// ID:      ALG-WINDOW-005
// Time:    O(1)
// Space:   O(1); the deque grows on demand to O(m), m = the most elements it has held at once
// Source:  Lemire, "Streaming maximum-minimum filter using no more than three comparisons per element", Nordic Journal of Computing, 2006, 13(4):328–339, arXiv:cs/0610046
// Source:  Laaksonen, Competitive Programmer's Handbook, draft of 3 July 2018, §8.3 "Sliding window minimum", pp. 81–82, https://cses.fi/book/book.pdf (read 2026-11-12)
// Tags:    monotonic-deque, monotonic-queue, deque, sliding-window, minimum, maximum, streaming
func NewMonotonicDeque[T cmp.Ordered](wantMax bool) *MonotonicDeque[T] {
	return &MonotonicDeque[T]{wantMax: wantMax}
}

// PushBack evicts from the back every element that v beats (strictly greater than v for a min-deque,
// strictly less for a max-deque), then appends v. Returns how many elements were evicted.
//
// ID:      ALG-WINDOW-006
// Time:    O(1) amortized   each element is evicted at most once; O(m) for the push that grows the ring, m = Len()
// Space:   O(1) amortized; the ring doubles when full
// Source:  Lemire, "Streaming maximum-minimum filter using no more than three comparisons per element", Nordic Journal of Computing, 2006, 13(4):328–339, arXiv:cs/0610046
// Source:  Laaksonen, Competitive Programmer's Handbook, draft of 3 July 2018, §8.3 "Sliding window minimum", pp. 81–82 ("remove elements from the end of the queue until the last queue element is smaller than the new window element"), https://cses.fi/book/book.pdf (read 2026-11-12)
// Tags:    monotonic-deque, monotonic-queue, push, eviction, amortized
func (d *MonotonicDeque[T]) PushBack(v T) int {
	evicted := 0
	for d.size > 0 {
		back := (d.head + d.size - 1) & (len(d.ring) - 1)
		c := cmp.Compare(d.ring[back], v)
		if d.wantMax {
			c = -c
		}
		if c <= 0 {
			break
		}
		var zero T
		d.ring[back] = zero
		d.size--
		evicted++
	}
	if d.size == len(d.ring) {
		d.grow()
	}
	d.ring[(d.head+d.size)&(len(d.ring)-1)] = v
	d.size++
	return evicted
}

// grow doubles the ring (minimum 8), keeping the elements in order from index 0.
func (d *MonotonicDeque[T]) grow() {
	ring := make([]T, max(8, 2*len(d.ring)))
	for i := range d.size {
		ring[i] = d.ring[(d.head+i)&(len(d.ring)-1)]
	}
	d.ring, d.head = ring, 0
}

// Front returns the front element, the least (min-deque) or greatest (max-deque) of those held, and
// true; or the zero value and false when the deque is empty.
//
// ID:      ALG-WINDOW-007
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Lemire, "Streaming maximum-minimum filter using no more than three comparisons per element", Nordic Journal of Computing, 2006, 13(4):328–339, arXiv:cs/0610046
// Tags:    monotonic-deque, monotonic-queue, front, minimum, maximum
func (d *MonotonicDeque[T]) Front() (T, bool) {
	if d == nil || d.size == 0 {
		var zero T
		return zero, false
	}
	return d.ring[d.head], true
}

// PopFront removes and returns the front element and true, or the zero value and false when the deque
// is empty.
//
// ID:      ALG-WINDOW-008
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Lemire, "Streaming maximum-minimum filter using no more than three comparisons per element", Nordic Journal of Computing, 2006, 13(4):328–339, arXiv:cs/0610046
// Tags:    monotonic-deque, monotonic-queue, pop
func (d *MonotonicDeque[T]) PopFront() (T, bool) {
	var zero T
	if d == nil || d.size == 0 {
		return zero, false
	}
	v := d.ring[d.head]
	d.ring[d.head] = zero
	d.head = (d.head + 1) & (len(d.ring) - 1)
	d.size--
	return v, true
}

// PopFrontIfEqual removes the front element when it equals leaving (by cmp.Compare) and reports
// whether it did. Call it with each value as it leaves a sliding window: if that value is still held
// it is at the front, and if it was evicted the front differs from it (it was beaten by a later
// element), so the deque always holds exactly the window's candidates.
//
// ID:      ALG-WINDOW-009
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Laaksonen, Competitive Programmer's Handbook, draft of 3 July 2018, §8.3 "Sliding window minimum", pp. 81–82 ("remove the first queue element if it is not inside the window anymore"), https://cses.fi/book/book.pdf (read 2026-11-12)
// Tags:    monotonic-deque, monotonic-queue, sliding-window, eviction
func (d *MonotonicDeque[T]) PopFrontIfEqual(leaving T) bool {
	if d == nil || d.size == 0 || cmp.Compare(d.ring[d.head], leaving) != 0 {
		return false
	}
	d.PopFront()
	return true
}

// Len returns the number of elements held.
//
// ID:      ALG-WINDOW-010
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Lemire, "Streaming maximum-minimum filter using no more than three comparisons per element", Nordic Journal of Computing, 2006, 13(4):328–339, arXiv:cs/0610046
// Tags:    monotonic-deque, monotonic-queue, length
func (d *MonotonicDeque[T]) Len() int {
	if d == nil {
		return 0
	}
	return d.size
}
