package heap

import (
	"errors"
	"fmt"
	"slices"
)

// ErrNotFinite: a floating-point value is NaN or ±Inf.
var ErrNotFinite = errors.New("heap: value is NaN or infinite")

// ErrOverflow: an integer sum does not fit its type.
var ErrOverflow = errors.New("heap: integer sum overflows")

// number is the set of types whose sums KSmallestPairs can rank.
type number interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr |
		~float32 | ~float64
}

// IndexPair is a pair of positions: I in the first slice, J in the second.
type IndexPair struct {
	I, J int
}

// KSmallestPairs appends to dst the k index pairs (i, j) with the smallest sums a[i] + b[j], in
// ascending order of sum, equal sums ordered by i then j; all len(a)·len(b) pairs when k exceeds that.
// a and b must be ascending (not checked: unsorted input gives an unspecified result, never a panic).
// dst == nil allocates once. a and b are read, never modified.
// Errors (dst returned with its length unchanged): ErrNegativeK when k < 0; ErrNotFinite when a float
// value read is NaN or ±Inf; ErrOverflow when an integer sum it forms overflows.
//
// ID:      ALG-HEAP-008
// Time:    O(m + k log m)   m = min(k, len(a)), k = min(k, len(a)·len(b))
// Space:   O(m) extra (the frontier heap); O(k) for the result
// Source:  Cormen, Leiserson, Rivest, Stein, Introduction to Algorithms, 3rd ed., MIT Press, 2009, §6.5, Exercise 6.5-9 (k-way merge with a min-heap; row i is the ascending list a[i]+b[0], a[i]+b[1], …)
// Source:  Frederickson, Johnson, "Generalized selection and ranking: sorted matrices", SIAM Journal on Computing, 1984, 13(1):14–30, doi:10.1137/0213002 (X+Y as a sorted matrix)
// Tags:    heap, pairs, sum, k-smallest, sorted-matrix, x-plus-y
func KSmallestPairs[T number](dst []IndexPair, a, b []T, k int) ([]IndexPair, error) {
	if k < 0 {
		return dst, fmt.Errorf("%w: k = %d", ErrNegativeK, k)
	}
	if len(a) == 0 || len(b) == 0 || k == 0 {
		return dst, nil
	}
	if k/len(a) >= len(b) { // k >= len(a)*len(b), without forming the product
		k = len(a) * len(b)
	}
	m := min(k, len(a))
	type entry struct {
		sum  T
		i, j int
	}
	sumOf := func(i, j int) (T, error) {
		x, y := a[i], b[j]
		if x != x || x-x != 0 {
			return 0, fmt.Errorf("%w: a[%d] = %v", ErrNotFinite, i, x)
		}
		if y != y || y-y != 0 {
			return 0, fmt.Errorf("%w: b[%d] = %v", ErrNotFinite, j, y)
		}
		s := x + y
		if (y > 0 && s < x) || (y < 0 && s > x) {
			return 0, fmt.Errorf("%w: a[%d] + b[%d] = %v + %v", ErrOverflow, i, j, x, y)
		}
		return s, nil
	}
	// Row i of the implicit matrix is a[i]+b[0] <= a[i]+b[1] <= …; only the first m rows can hold one
	// of the k smallest (row i has at least i entries before it, one per earlier row, column 0).
	frontier := make([]entry, m)
	for i := range m {
		s, err := sumOf(i, 0)
		if err != nil {
			return dst, err
		}
		frontier[i] = entry{s, i, 0}
	}
	h := &Heap[entry]{items: frontier, less: func(p, q entry) bool {
		if p.sum != q.sum {
			return p.sum < q.sum
		}
		if p.i != q.i {
			return p.i < q.i
		}
		return p.j < q.j
	}}
	for i := m/2 - 1; i >= 0; i-- {
		h.down(i)
	}
	start := len(dst)
	dst = slices.Grow(dst, k)
	for range k {
		e, _ := h.Pop()
		dst = append(dst, IndexPair{e.i, e.j})
		if e.j+1 < len(b) {
			s, err := sumOf(e.i, e.j+1)
			if err != nil {
				return dst[:start], err
			}
			h.Push(entry{s, e.i, e.j + 1})
		}
	}
	return dst, nil
}
