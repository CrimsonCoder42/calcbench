package tree

import "fmt"

// FenwickRangeUpdate holds an array of n numbers so that adding a constant to a whole range and reading
// one element each cost O(log n): a Fenwick tree over the difference array d[0] = a[0],
// d[i] = a[i] − a[i−1], whose prefix sums are the elements.
//
// T is a signed integer or float type. Integer arithmetic is exact or ErrOverflow: a difference, node
// sum or partial sum that does not fit T is reported, never wrapped (it can occur while every element
// fits, when neighbouring values are far apart near the type's limits; use a wider T). Floats must be
// finite. The zero value is an empty array (Len 0). Methods on a nil *FenwickRangeUpdate never panic.
// RangeAdd is atomic: on error the array is unchanged. Not safe for concurrent use.
type FenwickRangeUpdate[T signed] struct {
	diff FenwickTree[T]
}

// NewFenwickRangeUpdate returns a range-update structure holding a copy of values (pass make([]T, n)
// for n zeros). values is read, never modified.
// Errors: ErrNotFinite when a float value is NaN or ±Inf; ErrOverflow when a difference of neighbours
// or a node's sum does not fit T.
//
// ID:      ALG-TREE-028
// Time:    O(n)   n = len(values)
// Space:   O(n) for the result (one allocation for the differences, reused by the tree)
// Source:  Laaksonen, Competitive Programmer's Handbook, draft of 3 July 2018, §9.4 "Additional techniques", "Range updates", pp. 93–94 (a binary indexed tree over the difference array), https://cses.fi/book/book.pdf (read 2026-11-12)
// Source:  Fenwick, "A new data structure for cumulative frequency tables", Software: Practice and Experience, 1994, 24(3):327–336, doi:10.1002/spe.4380240306
// Tags:    fenwick-tree, binary-indexed-tree, range-update, point-query, difference-array
func NewFenwickRangeUpdate[T signed](values []T) (*FenwickRangeUpdate[T], error) {
	d := make([]T, len(values))
	var prev T
	for i, v := range values {
		if !finite(v) {
			return nil, fmt.Errorf("%w: values[%d] = %v", ErrNotFinite, i, v)
		}
		diff, ok := subChecked(v, prev)
		if !ok {
			return nil, fmt.Errorf("%w: values[%d] − values[%d]", ErrOverflow, i, i-1)
		}
		d[i] = diff
		prev = v
	}
	// Build in place over d (NewFenwickTree would copy it a second time).
	n := len(d)
	for k := 1; k <= n; k++ {
		if parent := k + k&-k; parent <= n {
			s, ok := addChecked(d[parent-1], d[k-1])
			if !ok {
				return nil, fmt.Errorf("%w: a node sum of the differences up to index %d", ErrOverflow, parent-1)
			}
			d[parent-1] = s
		}
	}
	return &FenwickRangeUpdate[T]{diff: FenwickTree[T]{t: d}}, nil
}

// Len returns the number of elements n (0 for a nil or zero value).
//
// ID:      ALG-TREE-031
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Laaksonen, Competitive Programmer's Handbook, draft of 3 July 2018, §9.4 "Range updates", p. 93, https://cses.fi/book/book.pdf (read 2026-11-12)
// Tags:    fenwick-tree, range-update, length
func (f *FenwickRangeUpdate[T]) Len() int {
	if f == nil {
		return 0
	}
	return len(f.diff.t)
}

// RangeAdd adds delta to every element of the half-open range [lo, hi); lo == hi is a no-op. It adds
// delta to d[lo] and subtracts it from d[hi] (when hi < n); nodes covering both are left as they are.
// Errors: ErrRange unless 0 <= lo <= hi <= Len(); ErrNotFinite when a float delta is NaN or ±Inf;
// ErrOverflow when a node's new sum would not fit T. On error nothing changes.
//
// ID:      ALG-TREE-029
// Time:    O(log n)   n = Len(); two passes over at most 2(⌊log2 n⌋ + 1) nodes (check, then write)
// Space:   O(1) extra, no allocation
// Source:  Laaksonen, Competitive Programmer's Handbook, draft of 3 July 2018, §9.4 "Range updates", p. 94 (increase position a by x and decrease position b + 1 by x), https://cses.fi/book/book.pdf (read 2026-11-12)
// Tags:    fenwick-tree, binary-indexed-tree, range-update, range-add
func (f *FenwickRangeUpdate[T]) RangeAdd(lo, hi int, delta T) error {
	n := f.Len()
	if lo < 0 || hi < lo || hi > n {
		return fmt.Errorf("%w: [%d, %d), Len = %d", ErrRange, lo, hi, n)
	}
	if !finite(delta) {
		return fmt.Errorf("%w: delta = %v", ErrNotFinite, delta)
	}
	if lo == hi || delta == 0 {
		return nil
	}
	t := f.diff.t
	// Walk the update chains of lo (+delta) and hi (−delta) in increasing node order; once they meet,
	// every later node is on both and its net change is 0.
	for write := range 2 {
		a, b := lo+1, hi+1 // b > n means no subtraction (hi == n)
		for (a <= n || b <= n) && a != b {
			var s T
			var ok bool
			k := a
			if b <= n && b < a {
				k = b
				s, ok = subChecked(t[k-1], delta)
				b += b & -b
			} else {
				s, ok = addChecked(t[k-1], delta)
				a += a & -a
			}
			if !ok {
				return fmt.Errorf("%w: adding %v to [%d, %d)", ErrOverflow, delta, lo, hi)
			}
			if write == 1 {
				t[k-1] = s
			}
		}
	}
	return nil
}

// Get returns element i (0-based): the prefix sum d[0] + … + d[i] of the differences.
// Errors: ErrRange unless 0 <= i < Len(); ErrOverflow when a partial sum does not fit T.
//
// ID:      ALG-TREE-030
// Time:    O(log n)   n = Len()
// Space:   O(1) extra, no allocation
// Source:  Laaksonen, Competitive Programmer's Handbook, draft of 3 July 2018, §9.4 "Range updates", p. 93 (the original array is the prefix sum array of the difference array), https://cses.fi/book/book.pdf (read 2026-11-12)
// Tags:    fenwick-tree, binary-indexed-tree, point-query, range-update
func (f *FenwickRangeUpdate[T]) Get(i int) (T, error) {
	if i < 0 || i >= f.Len() {
		var zero T
		return zero, fmt.Errorf("%w: i = %d, Len = %d", ErrRange, i, f.Len())
	}
	return f.diff.RangeSum(0, i+1)
}
