package heap

import (
	"errors"
	"fmt"
	"unsafe"
)

// ErrNilLess: the less function given to NewHeap or Heapify is nil.
var ErrNilLess = errors.New("heap: less function is nil")

// ErrTooLarge: a requested capacity would need more memory than one allocation may take.
var ErrTooLarge = errors.New("heap: requested size is too large")

// maxAllocBytes is half the runtime's largest allocation (2⁴⁷ bytes on 64-bit, 2³¹ on 32-bit); above it
// make panics (runtime/malloc.go maxAlloc). Rule L41.
const maxAllocBytes = 1 << (31 + 16*(^uint(0)>>63))

// checkAlloc returns ErrTooLarge when n elements of T would exceed maxAllocBytes. name is the parameter.
func checkAlloc[T any](name string, n int) error {
	var zero T
	if size := uint64(unsafe.Sizeof(zero)); size > 0 && n > 0 && uint64(n) > maxAllocBytes/size {
		return fmt.Errorf("%w: %s = %d elements of %d bytes", ErrTooLarge, name, n, size)
	}
	return nil
}

// Heap is a binary heap of T ordered by the less function given at construction: Pop returns an
// element e for which less(x, e) is false for every other element x (a min-heap for a "less than"
// function, a max-heap for "greater than"). Elements that compare equal come out in no promised order.
//
// A Heap comes from NewHeap or Heapify. The zero value is an empty heap with no order: Push on it stores
// nothing and returns false; Pop and Peek report empty. No method panics, on a nil *Heap either.
// A Heap is not safe for concurrent use. less must be a strict weak order; one that is not gives an
// unspecified order, never a panic.
//
// No index-based Fix or Remove is offered: an element's slot moves on every Push and Pop, so an index
// is useless to a caller without a swap hook. Callers that need decrease-key (Dijkstra, Prim) push the
// new key and skip stale entries when popped ("lazy deletion").
type Heap[T any] struct {
	items []T
	less  func(a, b T) bool
}

// NewHeap returns an empty heap ordered by less, with room for capacityHint elements before it grows
// (a negative hint is treated as 0).
// Errors: ErrNilLess when less is nil; ErrTooLarge when capacityHint elements of T exceed 2⁴⁷ bytes
// (2³¹ on 32-bit), the most one allocation may take (rule L41).
//
// ID:      ALG-HEAP-002
// Time:    O(1)
// Space:   O(c) for the heap's buffer, c = max(capacityHint, 0) (one allocation when c > 0)
// Source:  Williams, "Algorithm 232: Heapsort", Communications of the ACM, 1964, 7(6):347–348, doi:10.1145/512274.3734138
// Source:  Cormen, Leiserson, Rivest, Stein, Introduction to Algorithms, 3rd ed., MIT Press, 2009, §6.1 (binary heap), §6.5 (priority queues)
// Tags:    heap, priority-queue, binary-heap, generic
func NewHeap[T any](less func(a, b T) bool, capacityHint int) (*Heap[T], error) {
	if less == nil {
		return nil, ErrNilLess
	}
	if err := checkAlloc[T]("capacityHint", capacityHint); err != nil {
		return nil, err
	}
	return &Heap[T]{items: make([]T, 0, max(capacityHint, 0)), less: less}, nil
}

// Heapify returns a heap ordered by less holding a copy of items, built bottom-up in linear time.
// items is read, never modified.
// Errors: ErrNilLess when less is nil.
//
// ID:      ALG-HEAP-003
// Time:    O(n)   n = len(items); at most 2n calls to less
// Space:   O(n) for the heap's buffer (one allocation)
// Source:  Floyd, "Algorithm 245: Treesort", Communications of the ACM, 1964, 7(12):701, doi:10.1145/355588.365103
// Source:  Cormen, Leiserson, Rivest, Stein, Introduction to Algorithms, 3rd ed., MIT Press, 2009, §6.3 (BUILD-MAX-HEAP, O(n) bound)
// Tags:    heap, heapify, build-heap, priority-queue, generic
func Heapify[T any](less func(a, b T) bool, items []T) (*Heap[T], error) {
	if less == nil {
		return nil, fmt.Errorf("%w (Heapify of %d items)", ErrNilLess, len(items))
	}
	h := &Heap[T]{items: make([]T, len(items), max(len(items), 1)), less: less}
	copy(h.items, items)
	for i := len(h.items)/2 - 1; i >= 0; i-- {
		h.down(i)
	}
	return h, nil
}

// Push adds v to the heap and returns true; on the zero value or a nil *Heap (no order) it stores
// nothing and returns false.
//
// ID:      ALG-HEAP-004
// Time:    O(log n)   n = h.Len(); amortized when the buffer grows
// Space:   O(1) extra amortized (the buffer doubles when full)
// Source:  Williams, "Algorithm 232: Heapsort", Communications of the ACM, 1964, 7(6):347–348, doi:10.1145/512274.3734138 (procedure INHEAP)
// Source:  Cormen, Leiserson, Rivest, Stein, Introduction to Algorithms, 3rd ed., MIT Press, 2009, §6.5 (MAX-HEAP-INSERT)
// Tags:    heap, push, insert, priority-queue
func (h *Heap[T]) Push(v T) bool {
	if h == nil || h.less == nil {
		return false
	}
	h.items = append(h.items, v)
	h.up(len(h.items) - 1)
	return true
}

// Pop removes and returns the heap's first element (see Heap), with ok true; on an empty heap it
// returns the zero T and false.
//
// ID:      ALG-HEAP-005
// Time:    O(log n)   n = h.Len()
// Space:   O(1) extra; the buffer never shrinks
// Source:  Williams, "Algorithm 232: Heapsort", Communications of the ACM, 1964, 7(6):347–348, doi:10.1145/512274.3734138 (procedure OUTHEAP)
// Source:  Cormen, Leiserson, Rivest, Stein, Introduction to Algorithms, 3rd ed., MIT Press, 2009, §6.5 (HEAP-EXTRACT-MAX)
// Tags:    heap, pop, extract-min, extract-max, priority-queue
func (h *Heap[T]) Pop() (v T, ok bool) {
	if h == nil || len(h.items) == 0 {
		return v, false
	}
	last := len(h.items) - 1
	v = h.items[0]
	h.items[0] = h.items[last]
	var zero T
	h.items[last] = zero // drop the reference so the garbage collector can free it
	h.items = h.items[:last]
	if last > 0 {
		h.down(0)
	}
	return v, true
}

// Peek returns the heap's first element (the one Pop would return) without removing it, with ok true;
// on an empty heap it returns the zero T and false.
//
// ID:      ALG-HEAP-006
// Time:    O(1)
// Space:   O(1)
// Source:  Cormen, Leiserson, Rivest, Stein, Introduction to Algorithms, 3rd ed., MIT Press, 2009, §6.5 (HEAP-MAXIMUM)
// Tags:    heap, peek, find-min, find-max, priority-queue
func (h *Heap[T]) Peek() (v T, ok bool) {
	if h == nil || len(h.items) == 0 {
		return v, false
	}
	return h.items[0], true
}

// Len returns the number of elements in the heap (0 for a nil *Heap).
//
// ID:      ALG-HEAP-007
// Time:    O(1)
// Space:   O(1)
// Source:  Cormen, Leiserson, Rivest, Stein, Introduction to Algorithms, 3rd ed., MIT Press, 2009, §6.1 (heap-size)
// Tags:    heap, size, priority-queue
func (h *Heap[T]) Len() int {
	if h == nil {
		return 0
	}
	return len(h.items)
}

// up moves the element at i toward the root until its parent does not order after it. O(log n), iterative.
func (h *Heap[T]) up(i int) {
	v := h.items[i]
	for i > 0 {
		parent := (i - 1) / 2
		if !h.less(v, h.items[parent]) {
			break
		}
		h.items[i] = h.items[parent]
		i = parent
	}
	h.items[i] = v
}

// down moves the element at i toward the leaves until neither child orders before it. O(log n), iterative.
func (h *Heap[T]) down(i int) {
	n := len(h.items)
	v := h.items[i]
	for {
		child := 2*i + 1
		if child >= n || child < 0 { // child < 0: index overflow, unreachable for real slices
			break
		}
		if right := child + 1; right < n && h.less(h.items[right], h.items[child]) {
			child = right
		}
		if !h.less(h.items[child], v) {
			break
		}
		h.items[i] = h.items[child]
		i = child
	}
	h.items[i] = v
}
