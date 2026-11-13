package tree

import "fmt"

// LazyOps defines a lazy segment tree: a monoid of aggregates T and an update type U acting on it.
//
//   - Combine(left, right) is associative (need not commute); Identity is the aggregate of an empty range.
//   - Apply(u, agg, length) is the aggregate of a range of length elements whose old aggregate was agg,
//     after update u was applied to every element of it.
//   - Compose(newer, older) is the single update equal to applying older and then newer.
//   - NoUpdate is the update that changes nothing: Apply(NoUpdate, agg, length) = agg and
//     Compose(u, NoUpdate) = Compose(NoUpdate, u) = u.
//
// Example (range add, range sum): Combine = +, Identity = 0, Apply(u, s, len) = s + u·len,
// Compose(a, b) = a + b, NoUpdate = 0.
type LazyOps[T, U any] struct {
	Combine  func(left, right T) T
	Identity T
	Apply    func(update U, aggregate T, length int) T
	Compose  func(newer, older U) U
	NoUpdate U
}

// LazySegmentTree answers range queries and applies range updates over an array of n elements, each in
// O(log n), by keeping at every internal node its aggregate and one pending update for its children
// that is pushed down only when a later update needs to split the node (lazy propagation).
//
// Layout and recursion as SegmentTree (2n − 1 nodes, depth ≤ ⌈log2 n⌉ + 1 ≤ 64). Query does not modify
// the tree: pending updates met on the way down are composed and applied to the nodes it reads.
// RangeUpdate is atomic: it records every node it writes and restores them if an operation reports
// overflow (ready-made trees). The zero value is an empty tree. Methods on a nil *LazySegmentTree never
// panic. Not safe for concurrent use with RangeUpdate.
type LazySegmentTree[T, U any] struct {
	n           int
	node        []T
	lazy        []U // pending update for the children of an internal node; unused at leaves
	combine     func(left, right T) (T, bool)
	apply       func(u U, agg T, length int) (T, bool)
	compose     func(newer, older U) (U, bool)
	noUpdate    U
	validU      func(u U) error
	identity    T
	hasIdentity bool
	journal     []lazyUndo[T, U]
	// opp holds each node's opposite extreme for the range-add min/max trees (the max for a min tree,
	// the min for a max tree), so an add is checked against both ends of every range it touches;
	// oppCombine is nil (and opp unused) for every other tree.
	opp        []T
	oppCombine func(left, right T) (T, bool)
}

// lazyUndo is one node overwritten by RangeUpdate, kept to roll back on overflow.
type lazyUndo[T, U any] struct {
	i int
	v T
	z U
	w T // opp[i], when the tree keeps opposite extremes
}

// NewLazySegmentTree returns a lazy segment tree over a copy of values defined by ops (see LazyOps).
// values is read, never modified. The functions in ops must not panic.
// Errors: ErrNilFunc when ops.Combine, ops.Apply or ops.Compose is nil.
//
// ID:      ALG-TREE-039
// Time:    O(n)   n = len(values); n − 1 calls to Combine
// Space:   O(n) for the result (2n − 1 aggregates and 2n − 1 pending updates, one allocation each, plus an O(log n) undo buffer)
// Source:  Laaksonen, Competitive Programmer's Handbook, draft of 3 July 2018, §28.1 "Lazy propagation", pp. 258–260 (node value plus pending update; updates pushed down only when needed; pending updates combine), https://cses.fi/book/book.pdf (read 2026-11-12)
// Tags:    segment-tree, lazy-propagation, range-update, range-query, monoid, generic
func NewLazySegmentTree[T, U any](values []T, ops LazyOps[T, U]) (*LazySegmentTree[T, U], error) {
	switch {
	case ops.Combine == nil:
		return nil, fmt.Errorf("%w: ops.Combine", ErrNilFunc)
	case ops.Apply == nil:
		return nil, fmt.Errorf("%w: ops.Apply", ErrNilFunc)
	case ops.Compose == nil:
		return nil, fmt.Errorf("%w: ops.Compose", ErrNilFunc)
	}
	s := &LazySegmentTree[T, U]{
		combine:     func(l, r T) (T, bool) { return ops.Combine(l, r), true },
		apply:       func(u U, agg T, length int) (T, bool) { return ops.Apply(u, agg, length), true },
		compose:     func(newer, older U) (U, bool) { return ops.Compose(newer, older), true },
		noUpdate:    ops.NoUpdate,
		identity:    ops.Identity,
		hasIdentity: true,
	}
	if err := s.build(values); err != nil {
		return nil, err
	}
	return s, nil
}

// NewRangeAddSumTree returns a lazy segment tree over a copy of values whose update adds a constant to
// every element of a range and whose query returns a range's sum (the empty range sums to 0).
// Integer arithmetic is exact or ErrOverflow: a sum, a pending total of adds, or delta·length that does
// not fit T is reported, never wrapped (a pending total can overflow while every element fits; use a
// wider T). Float values and deltas must be finite. values is read, never modified.
// Errors: ErrNotFinite when a float value is NaN or ±Inf; ErrOverflow when a node's sum does not fit T.
//
// ID:      ALG-TREE-040
// Time:    O(n)   n = len(values)
// Space:   O(n) for the result
// Source:  Laaksonen, Competitive Programmer's Handbook, draft of 3 July 2018, §28.1 "Lazy propagation", "Lazy segment trees", pp. 258–260 (increase each value in [a, b] by a constant; sum of [a, b]; a node's sum grows by h·u), https://cses.fi/book/book.pdf (read 2026-11-12)
// Tags:    segment-tree, lazy-propagation, range-add, range-sum, range-update
func NewRangeAddSumTree[T number](values []T) (*LazySegmentTree[T, T], error) {
	s := &LazySegmentTree[T, T]{
		combine: addChecked[T],
		apply: func(u, agg T, length int) (T, bool) {
			p, ok := mulCount(u, length)
			if !ok {
				return p, false
			}
			return addChecked(agg, p)
		},
		compose:     addChecked[T],
		validU:      checkFinite[T],
		hasIdentity: true,
	}
	return newNumberLazy(s, values)
}

// NewRangeAddMinTree returns a lazy segment tree over a copy of values whose update adds a constant to
// every element of a range and whose query returns a range's minimum. An empty range has no minimum:
// Query reports ErrEmptyRange. Each node also keeps its range's maximum, so RangeUpdate checks the add
// against both extremes and rejects it (ErrOverflow, tree unchanged) when any element would leave T;
// a pending total of adds that does not fit T is also ErrOverflow (as NewRangeAddSumTree). Float values
// and deltas must be finite. values is read, never modified.
// Errors: ErrNotFinite when a float value is NaN or ±Inf.
//
// ID:      ALG-TREE-041
// Time:    O(n)   n = len(values)
// Space:   O(n) for the result (2n − 1 minimums, 2n − 1 maximums, 2n − 1 pending adds)
// Source:  Laaksonen, Competitive Programmer's Handbook, draft of 3 July 2018, §28.1 "Lazy propagation", p. 258 (a node holds "the sum or some other value" of its subarray plus a pending update; an added constant shifts a range's minimum by the same constant), https://cses.fi/book/book.pdf (read 2026-11-12)
// Tags:    segment-tree, lazy-propagation, range-add, range-minimum, range-update
func NewRangeAddMinTree[T number](values []T) (*LazySegmentTree[T, T], error) {
	s := &LazySegmentTree[T, T]{
		combine:    func(l, r T) (T, bool) { return min(l, r), true },
		apply:      func(u, agg T, _ int) (T, bool) { return addChecked(agg, u) },
		compose:    addChecked[T],
		validU:     checkFinite[T],
		oppCombine: func(l, r T) (T, bool) { return max(l, r), true },
	}
	return newNumberLazy(s, values)
}

// NewRangeAddMaxTree returns a lazy segment tree over a copy of values whose update adds a constant to
// every element of a range and whose query returns a range's maximum. An empty range has no maximum:
// Query reports ErrEmptyRange. Each node also keeps its range's minimum, so RangeUpdate checks the add
// against both extremes and rejects it (ErrOverflow, tree unchanged) when any element would leave T;
// a pending total of adds that does not fit T is also ErrOverflow (as NewRangeAddSumTree). Float values
// and deltas must be finite. values is read, never modified.
// Errors: ErrNotFinite when a float value is NaN or ±Inf.
//
// ID:      ALG-TREE-042
// Time:    O(n)   n = len(values)
// Space:   O(n) for the result (2n − 1 maximums, 2n − 1 minimums, 2n − 1 pending adds)
// Source:  Laaksonen, Competitive Programmer's Handbook, draft of 3 July 2018, §28.1 "Lazy propagation", p. 258 (a node holds "the sum or some other value" of its subarray plus a pending update; an added constant shifts a range's maximum by the same constant), https://cses.fi/book/book.pdf (read 2026-11-12)
// Tags:    segment-tree, lazy-propagation, range-add, range-maximum, range-update
func NewRangeAddMaxTree[T number](values []T) (*LazySegmentTree[T, T], error) {
	s := &LazySegmentTree[T, T]{
		combine:    func(l, r T) (T, bool) { return max(l, r), true },
		apply:      func(u, agg T, _ int) (T, bool) { return addChecked(agg, u) },
		compose:    addChecked[T],
		validU:     checkFinite[T],
		oppCombine: func(l, r T) (T, bool) { return min(l, r), true },
	}
	return newNumberLazy(s, values)
}

// newNumberLazy checks that every value is finite, then builds s.
func newNumberLazy[T number](s *LazySegmentTree[T, T], values []T) (*LazySegmentTree[T, T], error) {
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

func (s *LazySegmentTree[T, U]) build(values []T) error {
	s.n = len(values)
	if s.n == 0 {
		return nil
	}
	s.node = make([]T, 2*s.n-1)
	if s.oppCombine != nil {
		s.opp = make([]T, 2*s.n-1)
	}
	s.lazy = make([]U, 2*s.n-1)
	for i := range s.lazy {
		s.lazy[i] = s.noUpdate
	}
	// Each level of an update writes at most 4 visited nodes and 2 children per push: 12 entries a level.
	s.journal = make([]lazyUndo[T, U], 0, 12*treeDepth(s.n))
	return s.buildNode(values, 0, 0, s.n)
}

func (s *LazySegmentTree[T, U]) buildNode(values []T, i, lo, hi int) error {
	if hi-lo == 1 {
		s.node[i] = values[lo]
		if s.opp != nil {
			s.opp[i] = values[lo]
		}
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
	if s.opp != nil {
		s.opp[i], _ = s.oppCombine(s.opp[l], s.opp[r])
	}
	return nil
}

// Len returns the number of elements n (0 for a nil or zero tree).
//
// ID:      ALG-TREE-045
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Laaksonen, Competitive Programmer's Handbook, draft of 3 July 2018, §28.1 "Lazy propagation", p. 258, https://cses.fi/book/book.pdf (read 2026-11-12)
// Tags:    segment-tree, lazy-propagation, length
func (s *LazySegmentTree[T, U]) Len() int {
	if s == nil {
		return 0
	}
	return s.n
}

// RangeUpdate applies update to every element of the half-open range [lo, hi); lo == hi is a no-op.
// Errors: ErrRange unless 0 <= lo <= hi <= Len(); for the ready-made trees, ErrNotFinite for a NaN or
// ±Inf float update and ErrOverflow when an aggregate or a pending total would not fit T. On error the
// tree is unchanged.
//
// ID:      ALG-TREE-043
// Time:    O(log n)   n = Len(); at most 4 nodes per level visited, each push O(1) calls to Apply and Compose
// Space:   O(1) extra, no allocation (the undo buffer is sized at construction; recursion depth ≤ ⌈log2 n⌉ + 1)
// Source:  Laaksonen, Competitive Programmer's Handbook, draft of 3 July 2018, §28.1 "Lazy propagation", "Lazy segment trees", pp. 258–260 (a node inside the range takes the update and stops; a partly covered node pushes its pending update down first), https://cses.fi/book/book.pdf (read 2026-11-12)
// Tags:    segment-tree, lazy-propagation, range-update, range-add
func (s *LazySegmentTree[T, U]) RangeUpdate(lo, hi int, update U) error {
	n := s.Len()
	if lo < 0 || hi < lo || hi > n {
		return fmt.Errorf("%w: [%d, %d), Len = %d", ErrRange, lo, hi, n)
	}
	if s != nil && s.validU != nil {
		if err := s.validU(update); err != nil {
			return err
		}
	}
	if lo == hi {
		return nil
	}
	s.journal = s.journal[:0]
	if !s.updateNode(0, 0, n, lo, hi, update) {
		for j := len(s.journal) - 1; j >= 0; j-- {
			e := s.journal[j]
			s.node[e.i], s.lazy[e.i] = e.v, e.z
			if s.opp != nil {
				s.opp[e.i] = e.w
			}
		}
		return fmt.Errorf("%w: updating [%d, %d) by %v", ErrOverflow, lo, hi, update)
	}
	return nil
}

// record saves node i before it is written.
func (s *LazySegmentTree[T, U]) record(i int) {
	e := lazyUndo[T, U]{i: i, v: s.node[i], z: s.lazy[i]}
	if s.opp != nil {
		e.w = s.opp[i]
	}
	s.journal = append(s.journal, e)
}

// applyTo applies u to node i covering length elements: its aggregate, and its pending update unless leaf.
func (s *LazySegmentTree[T, U]) applyTo(i, length int, u U) bool {
	s.record(i)
	v, ok := s.apply(u, s.node[i], length)
	if !ok {
		return false
	}
	s.node[i] = v
	if s.opp != nil {
		w, ok := s.apply(u, s.opp[i], length)
		if !ok {
			return false
		}
		s.opp[i] = w
	}
	if length > 1 {
		z, ok := s.compose(u, s.lazy[i])
		if !ok {
			return false
		}
		s.lazy[i] = z
	}
	return true
}

func (s *LazySegmentTree[T, U]) updateNode(i, nlo, nhi, lo, hi int, u U) bool {
	if lo <= nlo && nhi <= hi {
		return s.applyTo(i, nhi-nlo, u)
	}
	mid := nlo + (nhi-nlo)/2
	l, r := i+1, i+2*(mid-nlo)
	// Push node i's pending update to its children before splitting it.
	s.record(i)
	pending := s.lazy[i]
	if !s.applyTo(l, mid-nlo, pending) || !s.applyTo(r, nhi-mid, pending) {
		return false
	}
	s.lazy[i] = s.noUpdate
	if lo < mid && !s.updateNode(l, nlo, mid, lo, hi, u) {
		return false
	}
	if hi > mid && !s.updateNode(r, mid, nhi, lo, hi, u) {
		return false
	}
	v, ok := s.combine(s.node[l], s.node[r])
	if !ok {
		return false
	}
	s.node[i] = v
	if s.opp != nil {
		if s.opp[i], ok = s.oppCombine(s.opp[l], s.opp[r]); !ok {
			return false
		}
	}
	return true
}

// Query returns the combined value of the elements in the half-open range [lo, hi), combined
// left-to-right, with every update applied. An empty range gives the identity, or ErrEmptyRange for
// the min and max trees. The tree is not modified.
// Errors: ErrRange unless 0 <= lo <= hi <= Len(); ErrEmptyRange as above; for the ready-made trees,
// ErrOverflow when a partial result does not fit T.
//
// ID:      ALG-TREE-044
// Time:    O(log n)   n = Len(); at most 4 nodes per level visited
// Space:   O(1) extra, no allocation (recursion depth ≤ ⌈log2 n⌉ + 1)
// Source:  Laaksonen, Competitive Programmer's Handbook, draft of 3 July 2018, §28.1 "Lazy propagation", "Lazy segment trees", p. 259 (a node fully inside the range contributes its value; pending updates are accounted for on the way down), https://cses.fi/book/book.pdf (read 2026-11-12)
// Tags:    segment-tree, lazy-propagation, range-query, range-sum, range-minimum, range-maximum
func (s *LazySegmentTree[T, U]) Query(lo, hi int) (T, error) {
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
	v, ok := s.queryNode(0, 0, n, lo, hi, s.noUpdate)
	if !ok {
		return zero, fmt.Errorf("%w: combining [%d, %d)", ErrOverflow, lo, hi)
	}
	return v, nil
}

// queryNode combines the part of [lo, hi) inside node i's range [nlo, nhi) (they intersect), with
// above, the composition of the pending updates of i's ancestors, applied to what it reads.
func (s *LazySegmentTree[T, U]) queryNode(i, nlo, nhi, lo, hi int, above U) (T, bool) {
	if lo <= nlo && nhi <= hi {
		return s.apply(above, s.node[i], nhi-nlo)
	}
	// Ancestors' pending updates are newer than node i's own.
	acc, ok := s.compose(above, s.lazy[i])
	if !ok {
		var zero T
		return zero, false
	}
	mid := nlo + (nhi-nlo)/2
	l, r := i+1, i+2*(mid-nlo)
	if hi <= mid {
		return s.queryNode(l, nlo, mid, lo, hi, acc)
	}
	if lo >= mid {
		return s.queryNode(r, mid, nhi, lo, hi, acc)
	}
	a, ok := s.queryNode(l, nlo, mid, lo, hi, acc)
	if !ok {
		return a, false
	}
	b, ok := s.queryNode(r, mid, nhi, lo, hi, acc)
	if !ok {
		return b, false
	}
	return s.combine(a, b)
}
