package tree

import "fmt"

// FenwickTree (binary indexed tree) holds an array of n numbers so that adding to one element and
// summing a prefix or a range each cost O(log n). Internally t[k−1] = a[k−p(k)] + … + a[k−1] for
// k = 1…n, where p(k) = k & −k is the largest power of two dividing k.
//
// T is a signed integer or float type. Integer sums are exact or ErrOverflow: a sum the tree forms (a
// node, or a partial sum while answering) that does not fit T is reported, never wrapped; this can occur
// for a range whose total fits when values of both signs sit near the type's limits (use a wider T).
// Float values must be finite; float range sums are differences of partial sums and carry their rounding.
// The zero value is an empty tree (Len 0). Methods on a nil *FenwickTree never panic. Every method that
// changes the tree is atomic: on error the tree is unchanged. A FenwickTree is not safe for concurrent use.
type FenwickTree[T signed] struct {
	t []T
}

// NewFenwickTree returns a Fenwick tree holding a copy of values (pass make([]T, n) for n zeros), built
// in linear time: each node adds its sum once into the next node that covers it, k + p(k).
// values is read, never modified.
// Errors: ErrNotFinite when a float value is NaN or ±Inf; ErrOverflow when a node's sum does not fit T.
//
// ID:      ALG-TREE-023
// Time:    O(n)   n = len(values)
// Space:   O(n) for the result (one allocation)
// Source:  Fenwick, "A new data structure for cumulative frequency tables", Software: Practice and Experience, 1994, 24(3):327–336, doi:10.1002/spe.4380240306
// Source:  Laaksonen, Competitive Programmer's Handbook, draft of 3 July 2018, §9.2 "Binary indexed tree", pp. 86–88 (tree[k] = sum of (k − p(k), k], p(k) = k & −k), https://cses.fi/book/book.pdf (read 2026-11-12)
// Tags:    fenwick-tree, binary-indexed-tree, bit, prefix-sum, range-sum, construct
func NewFenwickTree[T signed](values []T) (*FenwickTree[T], error) {
	t := make([]T, len(values))
	for i, v := range values {
		if !finite(v) {
			return nil, fmt.Errorf("%w: values[%d] = %v", ErrNotFinite, i, v)
		}
		t[i] = v
	}
	n := len(t)
	for k := 1; k <= n; k++ {
		if parent := k + k&-k; parent <= n {
			s, ok := addChecked(t[parent-1], t[k-1])
			if !ok {
				return nil, fmt.Errorf("%w: the sum of values[%d:%d]", ErrOverflow, parent-parent&-parent, parent)
			}
			t[parent-1] = s
		}
	}
	return &FenwickTree[T]{t: t}, nil
}

// Len returns the number of elements n (0 for a nil or zero tree).
//
// ID:      ALG-TREE-027
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Fenwick, "A new data structure for cumulative frequency tables", Software: Practice and Experience, 1994, 24(3):327–336, doi:10.1002/spe.4380240306
// Tags:    fenwick-tree, binary-indexed-tree, length
func (f *FenwickTree[T]) Len() int {
	if f == nil {
		return 0
	}
	return len(f.t)
}

// Add adds delta to element i (0-based).
// Errors: ErrRange unless 0 <= i < Len(); ErrNotFinite when a float delta is NaN or ±Inf; ErrOverflow
// when a node's new sum would not fit T. On error the tree is unchanged.
//
// ID:      ALG-TREE-024
// Time:    O(log n)   n = Len(); two passes over at most ⌊log2 n⌋ + 1 nodes (check, then write)
// Space:   O(1) extra, no allocation
// Source:  Fenwick, "A new data structure for cumulative frequency tables", Software: Practice and Experience, 1994, 24(3):327–336, doi:10.1002/spe.4380240306
// Source:  Laaksonen, Competitive Programmer's Handbook, draft of 3 July 2018, §9.2 "Binary indexed tree", "Implementation", p. 88 (add: k += k & −k), https://cses.fi/book/book.pdf (read 2026-11-12)
// Tags:    fenwick-tree, binary-indexed-tree, point-update, add
func (f *FenwickTree[T]) Add(i int, delta T) error {
	if f == nil || i < 0 || i >= len(f.t) {
		return fmt.Errorf("%w: i = %d, Len = %d", ErrRange, i, f.Len())
	}
	if !finite(delta) {
		return fmt.Errorf("%w: delta = %v", ErrNotFinite, delta)
	}
	return f.add(i, delta, false)
}

// add adds delta (or subtracts it, when subtract is true) at index i, atomically.
func (f *FenwickTree[T]) add(i int, delta T, subtract bool) error {
	n := len(f.t)
	op := addChecked[T]
	if subtract {
		op = subChecked[T]
	}
	for k := i + 1; k <= n; k += k & -k {
		if _, ok := op(f.t[k-1], delta); !ok {
			return fmt.Errorf("%w: the sum of elements [%d, %d) after changing element %d by %v", ErrOverflow, k-k&-k, k, i, delta)
		}
	}
	for k := i + 1; k <= n; k += k & -k {
		f.t[k-1], _ = op(f.t[k-1], delta)
	}
	return nil
}

// PrefixSum returns the sum of the first k elements, a[0] + … + a[k−1]; k = 0 gives 0.
// Errors: ErrRange unless 0 <= k <= Len(); ErrOverflow when a partial sum does not fit T.
//
// ID:      ALG-TREE-025
// Time:    O(log n)   n = Len(); at most ⌊log2 n⌋ + 1 terms
// Space:   O(1) extra, no allocation
// Source:  Fenwick, "A new data structure for cumulative frequency tables", Software: Practice and Experience, 1994, 24(3):327–336, doi:10.1002/spe.4380240306
// Source:  Laaksonen, Competitive Programmer's Handbook, draft of 3 July 2018, §9.2 "Binary indexed tree", p. 87 (sum of [1, 7] = 16 + 7 + 4 = 27), p. 88 (sum: k −= k & −k), https://cses.fi/book/book.pdf (read 2026-11-12)
// Tags:    fenwick-tree, binary-indexed-tree, prefix-sum, cumulative-frequency
func (f *FenwickTree[T]) PrefixSum(k int) (T, error) {
	return f.RangeSum(0, k)
}

// RangeSum returns the sum of the elements in the half-open range [lo, hi): a[lo] + … + a[hi−1]; lo == hi
// gives 0. It equals PrefixSum(hi) − PrefixSum(lo), evaluated by walking the two index chains only until
// they meet, so the shared part is never added and subtracted.
// Errors: ErrRange unless 0 <= lo <= hi <= Len(); ErrOverflow when a partial sum does not fit T.
//
// ID:      ALG-TREE-026
// Time:    O(log n)   n = Len(); at most 2(⌊log2 n⌋ + 1) terms
// Space:   O(1) extra, no allocation
// Source:  Laaksonen, Competitive Programmer's Handbook, draft of 3 July 2018, §9.2 "Binary indexed tree", p. 87 (sum(a, b) = sum(1, b) − sum(1, a − 1)), https://cses.fi/book/book.pdf (read 2026-11-12)
// Source:  Fenwick, "A new data structure for cumulative frequency tables", Software: Practice and Experience, 1994, 24(3):327–336, doi:10.1002/spe.4380240306
// Tags:    fenwick-tree, binary-indexed-tree, range-sum, range-query
func (f *FenwickTree[T]) RangeSum(lo, hi int) (T, error) {
	var s T
	if lo < 0 || hi < lo || hi > f.Len() {
		return s, fmt.Errorf("%w: [%d, %d), Len = %d", ErrRange, lo, hi, f.Len())
	}
	// Prefix chains: k → k − p(k). The chains of lo and hi reach a common value; the nodes below it
	// belong to both prefix sums and cancel.
	lo0, hi0 := lo, hi
	var ok bool
	for hi != lo {
		if hi > lo {
			if s, ok = addChecked(s, f.t[hi-1]); !ok {
				return 0, fmt.Errorf("%w: a partial sum of the range [%d, %d)", ErrOverflow, lo0, hi0)
			}
			hi -= hi & -hi
		} else {
			if s, ok = subChecked(s, f.t[lo-1]); !ok {
				return 0, fmt.Errorf("%w: a partial sum of the range [%d, %d)", ErrOverflow, lo0, hi0)
			}
			lo -= lo & -lo
		}
	}
	return s, nil
}
