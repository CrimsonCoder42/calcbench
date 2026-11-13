package hashing

import (
	"cmp"
	"errors"
	"fmt"
)

// ErrNegativeK: k is less than 0.
var ErrNegativeK = errors.New("hashing: k must not be negative")

// Count is one value and the number of times it occurs.
type Count[T any] struct {
	Value T
	Count int
}

// TopKByCount returns the k entries of counts with the largest counts, ranked by count descending and,
// for equal counts, by value ascending (cmp.Compare), or all entries when k >= len(counts). Map keys are
// distinct, so the ranking is a total order and map iteration order never reaches the result.
// k = 0 gives an empty, non-nil slice. counts is read, never modified.
// Errors: ErrNegativeK when k < 0.
//
// ID:      ALG-HASH-003
// Time:    O(m log k + m) expected   m = len(counts), k = min(k, m)
// Space:   O(1) extra; O(k) for the result (one allocation)
// Source:  Sedgewick, Wayne, Algorithms, 4th ed., Addison-Wesley, 2011, §2.4 (client TopM: a min-heap of the M largest)
// Source:  Williams, "Algorithm 232: Heapsort", Communications of the ACM, 1964, 7(6):347–348, doi:10.1145/512274.3734138
// Tags:    hash, heap, top-k, frequency, ranking, deterministic
func TopKByCount[T cmp.Ordered](counts map[T]int, k int) ([]Count[T], error) {
	if k < 0 {
		return nil, fmt.Errorf("%w: k = %d", ErrNegativeK, k)
	}
	k = min(k, len(counts))
	if k == 0 {
		return []Count[T]{}, nil
	}
	// The result buffer is the heap: a min-heap by rank whose root is the weakest of the k best so far.
	top := make([]Count[T], 0, k)
	for v, c := range counts {
		e := Count[T]{Value: v, Count: c}
		if len(top) < k {
			top = append(top, e)
			if len(top) == k {
				for i := k/2 - 1; i >= 0; i-- {
					siftDownCount(top, i)
				}
			}
			continue
		}
		if ranksAbove(e, top[0]) {
			top[0] = e
			siftDownCount(top, 0)
		}
	}
	// Heapsort in place: moving each root (the weakest) to the end leaves the best first.
	for end := k - 1; end > 0; end-- {
		top[0], top[end] = top[end], top[0]
		siftDownCount(top[:end], 0)
	}
	return top, nil
}

// ranksAbove reports whether a ranks strictly before b: a larger count, or an equal count and a smaller value.
func ranksAbove[T cmp.Ordered](a, b Count[T]) bool {
	if a.Count != b.Count {
		return a.Count > b.Count
	}
	return cmp.Less(a.Value, b.Value)
}

// siftDownCount restores the min-heap (weakest at the root) below index i. O(log len(h)), iterative.
func siftDownCount[T cmp.Ordered](h []Count[T], i int) {
	n := len(h)
	for {
		weakest := i
		left, right := 2*i+1, 2*i+2
		if left < n && ranksAbove(h[weakest], h[left]) {
			weakest = left
		}
		if right < n && ranksAbove(h[weakest], h[right]) {
			weakest = right
		}
		if weakest == i {
			return
		}
		h[i], h[weakest] = h[weakest], h[i]
		i = weakest
	}
}
