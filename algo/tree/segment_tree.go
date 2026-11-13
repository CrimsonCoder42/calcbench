package tree

import (
	"cmp"
	"fmt"
	"math/bits"
)

// SegmentTree answers range queries over an array of n elements under an associative operation (sum,
// min, max, gcd, matrix product, …) and supports changing one element, each in O(log n).
//
// The operation is a monoid: combine(left, right) must be associative; it need not be commutative (the
// tree always combines left-to-right). An identity, when the tree has one, is returned for an empty
// range. Nodes are laid out in 2n − 1 slots (a node for [lo, hi) at i has its left child at i + 1 and
// its right child at i + 2·(mid − lo), mid = lo + (hi − lo)/2), and the recursion is over the index
// range, so its depth is at most ⌈log2 n⌉ + 1 (≤ 64).
//
// The zero value is an empty tree (Len 0). Methods on a nil *SegmentTree never panic. Set is atomic: on
// error the tree is unchanged. Query does not modify the tree. Not safe for concurrent use with Set.
type SegmentTree[T any] struct {
	n           int
	node        []T
	combine     func(left, right T) (T, bool) // false: the result does not fit (ready-made sums)
	valid       func(v T) error               // nil or a check on values given to Set
	identity    T
	hasIdentity bool
	journal     []segUndo[T]
}

// segUndo is one node value overwritten by Set, kept to roll back on overflow.
type segUndo[T any] struct {
	i int
	v T
}

// NewSegmentTree returns a segment tree over a copy of values combined by combine, with identity as
// the value of an empty range (combine(identity, x) = combine(x, identity) = x). values is read, never
// modified. combine must be associative and must not panic; it is called O(n) times here.
// Errors: ErrNilFunc when combine is nil.
//
// ID:      ALG-TREE-032
// Time:    O(n)   n = len(values); n − 1 calls to combine
// Space:   O(n) for the result (2n − 1 nodes, one allocation, plus an O(log n) undo buffer)
// Source:  Laaksonen, Competitive Programmer's Handbook, draft of 3 July 2018, §9.3 "Segment tree", pp. 89–93 (range queries under any associative operation, point updates, O(log n) each), https://cses.fi/book/book.pdf (read 2026-11-12)
// Tags:    segment-tree, range-query, monoid, point-update, generic
func NewSegmentTree[T any](values []T, combine func(left, right T) T, identity T) (*SegmentTree[T], error) {
	if combine == nil {
		return nil, fmt.Errorf("%w: combine", ErrNilFunc)
	}
	s := &SegmentTree[T]{
		combine:     func(l, r T) (T, bool) { return combine(l, r), true },
		identity:    identity,
		hasIdentity: true,
	}
	if err := s.build(values); err != nil {
		return nil, err
	}
	return s, nil
}

// NewSumSegmentTree returns a segment tree of sums over a copy of values; the empty range sums to 0.
// Integer sums are exact or ErrOverflow (from here, Set or Query) when a sum the tree forms does not fit
// T; float values must be finite and sums are formed pairwise up the tree. values is read, never modified.
// Errors: ErrNotFinite when a float value is NaN or ±Inf; ErrOverflow when a node's sum does not fit T.
//
// ID:      ALG-TREE-033
// Time:    O(n)   n = len(values)
// Space:   O(n) for the result (2n − 1 nodes, one allocation)
// Source:  Laaksonen, Competitive Programmer's Handbook, draft of 3 July 2018, §9.3 "Segment tree", pp. 89–91 (sum queries), https://cses.fi/book/book.pdf (read 2026-11-12)
// Tags:    segment-tree, range-sum, point-update, sum
func NewSumSegmentTree[T number](values []T) (*SegmentTree[T], error) {
	s := &SegmentTree[T]{combine: addChecked[T], valid: checkFinite[T], hasIdentity: true}
	for i, v := range values {
		if !finite(v) {
			return nil, fmt.Errorf("%w: values[%d] = %v", ErrNotFinite, i, v)
		}
	}
	if err := s.build(values); err != nil {
		return nil, err
	}
	return s, nil
}

// NewMinSegmentTree returns a segment tree of minimums over a copy of values in cmp.Compare order (a NaN
// is less than every other float, D-002). An empty range has no minimum: Query reports ErrEmptyRange.
// values is read, never modified.
//
// ID:      ALG-TREE-034
// Time:    O(n)   n = len(values)
// Space:   O(n) for the result (2n − 1 nodes, one allocation)
// Source:  Laaksonen, Competitive Programmer's Handbook, draft of 3 July 2018, §9.3 "Segment tree", "Other queries", p. 92 (minimum queries), https://cses.fi/book/book.pdf (read 2026-11-12)
// Tags:    segment-tree, range-minimum, rmq, point-update, min
func NewMinSegmentTree[T cmp.Ordered](values []T) (*SegmentTree[T], error) {
	s := &SegmentTree[T]{combine: func(l, r T) (T, bool) {
		if cmp.Less(r, l) {
			return r, true
		}
		return l, true
	}}
	if err := s.build(values); err != nil {
		return nil, err
	}
	return s, nil
}

// NewMaxSegmentTree returns a segment tree of maximums over a copy of values in cmp.Compare order (a NaN
// is less than every other float, D-002). An empty range has no maximum: Query reports ErrEmptyRange.
// values is read, never modified.
//
// ID:      ALG-TREE-035
// Time:    O(n)   n = len(values)
// Space:   O(n) for the result (2n − 1 nodes, one allocation)
// Source:  Laaksonen, Competitive Programmer's Handbook, draft of 3 July 2018, §9.3 "Segment tree", "Other queries", p. 92 (maximum queries), https://cses.fi/book/book.pdf (read 2026-11-12)
// Tags:    segment-tree, range-maximum, point-update, max
func NewMaxSegmentTree[T cmp.Ordered](values []T) (*SegmentTree[T], error) {
	s := &SegmentTree[T]{combine: func(l, r T) (T, bool) {
		if cmp.Less(l, r) {
			return r, true
		}
		return l, true
	}}
	if err := s.build(values); err != nil {
		return nil, err
	}
	return s, nil
}

// checkFinite reports ErrNotFinite for a NaN or infinite float.
func checkFinite[T number](v T) error {
	if !finite(v) {
		return fmt.Errorf("%w: %v", ErrNotFinite, v)
	}
	return nil
}

// treeDepth is the most nodes on a root-to-leaf path of a 2n−1 layout: ⌈log2 n⌉ + 1.
func treeDepth(n int) int {
	if n <= 1 {
		return 1
	}
	return bits.Len(uint(n-1)) + 1
}

// build fills the nodes from values. An error leaves s unusable; the constructors then return nil.
func (s *SegmentTree[T]) build(values []T) error {
	s.n = len(values)
	if s.n == 0 {
		return nil
	}
	s.node = make([]T, 2*s.n-1)
	s.journal = make([]segUndo[T], 0, treeDepth(s.n))
	return s.buildNode(values, 0, 0, s.n)
}

func (s *SegmentTree[T]) buildNode(values []T, i, lo, hi int) error {
	if hi-lo == 1 {
		s.node[i] = values[lo]
		return nil
	}
	mid := lo + (hi-lo)/2
	l, r := i+1, i+2*(mid-lo)
	if err := s.buildNode(values, l, lo, mid); err != nil {
		return err
	}
	if err := s.buildNode(values, r, mid, hi); err != nil {
		return err
	}
	v, ok := s.combine(s.node[l], s.node[r])
	if !ok {
		return fmt.Errorf("%w: the combined value of values[%d:%d]", ErrOverflow, lo, hi)
	}
	s.node[i] = v
	return nil
}

// Len returns the number of elements n (0 for a nil or zero tree).
//
// ID:      ALG-TREE-038
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Laaksonen, Competitive Programmer's Handbook, draft of 3 July 2018, §9.3 "Segment tree", p. 89, https://cses.fi/book/book.pdf (read 2026-11-12)
// Tags:    segment-tree, length
func (s *SegmentTree[T]) Len() int {
	if s == nil {
		return 0
	}
	return s.n
}

// Set replaces element i (0-based) with v and recomputes the O(log n) nodes above it.
// Errors: ErrRange unless 0 <= i < Len(); for a sum tree, ErrNotFinite for a NaN or ±Inf float and
// ErrOverflow when a node's new sum would not fit T. On error the tree is unchanged.
//
// ID:      ALG-TREE-036
// Time:    O(log n)   n = Len(); ⌈log2 n⌉ calls to combine
// Space:   O(1) extra, no allocation (the undo buffer is sized at construction; recursion depth ≤ ⌈log2 n⌉ + 1)
// Source:  Laaksonen, Competitive Programmer's Handbook, draft of 3 July 2018, §9.3 "Segment tree", p. 90 (after an update, every node on the path to the root changes), https://cses.fi/book/book.pdf (read 2026-11-12)
// Tags:    segment-tree, point-update, set
func (s *SegmentTree[T]) Set(i int, v T) error {
	if s == nil || i < 0 || i >= s.n {
		return fmt.Errorf("%w: i = %d, Len = %d", ErrRange, i, s.Len())
	}
	if s.valid != nil {
		if err := s.valid(v); err != nil {
			return err
		}
	}
	s.journal = s.journal[:0]
	if !s.setNode(0, 0, s.n, i, v) {
		for j := len(s.journal) - 1; j >= 0; j-- {
			s.node[s.journal[j].i] = s.journal[j].v
		}
		return fmt.Errorf("%w: setting element %d to %v", ErrOverflow, i, v)
	}
	return nil
}

func (s *SegmentTree[T]) setNode(i, lo, hi, at int, v T) bool {
	s.journal = append(s.journal, segUndo[T]{i, s.node[i]})
	if hi-lo == 1 {
		s.node[i] = v
		return true
	}
	mid := lo + (hi-lo)/2
	l, r := i+1, i+2*(mid-lo)
	if at < mid {
		if !s.setNode(l, lo, mid, at, v) {
			return false
		}
	} else if !s.setNode(r, mid, hi, at, v) {
		return false
	}
	c, ok := s.combine(s.node[l], s.node[r])
	if !ok {
		return false
	}
	s.node[i] = c
	return true
}

// Query returns the combined value of the elements in the half-open range [lo, hi), combined
// left-to-right. An empty range (lo == hi) gives the identity, or ErrEmptyRange for a min or max tree.
// Errors: ErrRange unless 0 <= lo <= hi <= Len(); ErrEmptyRange as above; for a sum tree, ErrOverflow
// when a partial sum does not fit T.
//
// ID:      ALG-TREE-037
// Time:    O(log n)   n = Len(); at most 4 nodes per level are visited
// Space:   O(1) extra, no allocation (recursion depth ≤ ⌈log2 n⌉ + 1)
// Source:  Laaksonen, Competitive Programmer's Handbook, draft of 3 July 2018, §9.3 "Segment tree", pp. 90–91 (a range is the union of O(log n) node ranges), https://cses.fi/book/book.pdf (read 2026-11-12)
// Tags:    segment-tree, range-query, range-sum, range-minimum, range-maximum
func (s *SegmentTree[T]) Query(lo, hi int) (T, error) {
	var zero T
	n := s.Len()
	if lo < 0 || hi < lo || hi > n {
		return zero, fmt.Errorf("%w: [%d, %d), Len = %d", ErrRange, lo, hi, n)
	}
	if lo == hi {
		if s != nil && s.hasIdentity {
			return s.identity, nil
		}
		return zero, fmt.Errorf("%w: [%d, %d)", ErrEmptyRange, lo, hi)
	}
	v, ok := s.queryNode(0, 0, n, lo, hi)
	if !ok {
		return zero, fmt.Errorf("%w: combining [%d, %d)", ErrOverflow, lo, hi)
	}
	return v, nil
}

// queryNode combines the part of [lo, hi) inside node i's range [nlo, nhi); the two must intersect.
func (s *SegmentTree[T]) queryNode(i, nlo, nhi, lo, hi int) (T, bool) {
	if lo <= nlo && nhi <= hi {
		return s.node[i], true
	}
	mid := nlo + (nhi-nlo)/2
	l, r := i+1, i+2*(mid-nlo)
	if hi <= mid {
		return s.queryNode(l, nlo, mid, lo, hi)
	}
	if lo >= mid {
		return s.queryNode(r, mid, nhi, lo, hi)
	}
	a, ok := s.queryNode(l, nlo, mid, lo, hi)
	if !ok {
		return a, false
	}
	b, ok := s.queryNode(r, mid, nhi, lo, hi)
	if !ok {
		return b, false
	}
	return s.combine(a, b)
}
