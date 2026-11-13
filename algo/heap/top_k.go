package heap

import (
	"cmp"
	"errors"
	"fmt"
)

// ErrNegativeK: k is less than 0.
var ErrNegativeK = errors.New("heap: k must not be negative")

// TopK returns the k largest elements of values in descending order, or all of them when
// k >= len(values). Duplicates count separately. k = 0 gives an empty, non-nil slice.
// values is read, never modified.
// Errors: ErrNegativeK when k < 0.
//
// ID:      ALG-HEAP-001
// Time:    O(n log k)   n = len(values), k = min(k, n); O(n + k log k) when most values lose to the heap's root
// Space:   O(1) extra; O(k) for the result (one allocation)
// Source:  Sedgewick, Wayne, Algorithms, 4th ed., Addison-Wesley, 2011, §2.4 (client TopM: a min-heap of the M largest)
// Source:  Floyd, "Algorithm 245: Treesort", Communications of the ACM, 1964, 7(12):701, doi:10.1145/355588.365103
// Source:  Williams, "Algorithm 232: Heapsort", Communications of the ACM, 1964, 7(6):347–348, doi:10.1145/512274.3734138
// Tags:    heap, priority-queue, top-k, largest, partial-sort, streaming
func TopK[T cmp.Ordered](values []T, k int) ([]T, error) {
	if k < 0 {
		return nil, fmt.Errorf("%w: k = %d", ErrNegativeK, k)
	}
	k = min(k, len(values))
	if k == 0 {
		return []T{}, nil
	}
	// The result buffer is the heap: a min-heap whose root is the smallest of the k largest so far.
	top := make([]T, k)
	copy(top, values[:k])
	for i := k/2 - 1; i >= 0; i-- {
		siftDown(top, i)
	}
	for _, v := range values[k:] {
		if cmp.Less(top[0], v) {
			top[0] = v
			siftDown(top, 0)
		}
	}
	// Heapsort in place: moving each root (the current minimum) to the end leaves descending order.
	for end := k - 1; end > 0; end-- {
		top[0], top[end] = top[end], top[0]
		siftDown(top[:end], 0)
	}
	return top, nil
}

// siftDown restores the min-heap property of h below index i, assuming both subtrees of i are heaps.
// O(log len(h)), iterative.
func siftDown[T cmp.Ordered](h []T, i int) {
	n := len(h)
	for {
		least := i
		left, right := 2*i+1, 2*i+2
		if left < n && cmp.Less(h[left], h[least]) {
			least = left
		}
		if right < n && cmp.Less(h[right], h[least]) {
			least = right
		}
		if least == i {
			return
		}
		h[i], h[least] = h[least], h[i]
		i = least
	}
}
