package stream

import (
	"slices"

	"github.com/CrimsonCoder42/calcbench/algo/interval"
	"github.com/CrimsonCoder42/calcbench/algo/search"
)

// integer is the set of types whose consecutive values DisjointIntervals can join.
type integer interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 | ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr
}

// DisjointIntervals is the set of integers added so far, kept as sorted, disjoint, non-adjacent closed
// ranges: adding 1, 3, 7, 2, 6 gives [1, 3] and [6, 7]. Two ranges are joined when the gap between them
// is filled; a value already covered changes nothing. Works up to the extremes of T with no overflow.
// The ranges are kept in two parallel sorted slices (starts, ends); Add finds its place by binary search
// (search.UpperBound, ALG-SEARCH-002) and inserts or removes at most one range. The zero value is an
// empty set, ready to use. Not safe for concurrent use.
//
// Methods: Add(v) adds v. Contains(v) reports whether v was added. Len is the number of ranges.
// Intervals(dst) appends the ranges, ascending, as interval.Interval values (ALG-INTVL-001's type).
//
// ID:      ALG-STREAM-003
// Time:    Add O(log k + k) worst case (the binary search, then at most one insertion or removal, a memmove of up to k elements); Contains O(log k); Intervals O(k); Len O(1)   k = number of ranges
// Space:   O(k); Add allocates only when a new range outgrows the slices (amortized O(1)); Intervals O(k) for the result
// Source:  Knuth, The Art of Computer Programming, Vol. 3: Sorting and Searching, 2nd ed., Addison-Wesley, 1998, §6.2.1, Algorithm B (binary search in an ordered table), via ALG-SEARCH-002
// Source:  Cormen, Leiserson, Rivest, Stein, Introduction to Algorithms, 3rd ed., MIT Press, 2009, ch. 13 (red-black trees: the O(log k) Add that lost, see the catalog)
// Tags:    stream, intervals, disjoint-intervals, ranges, merge, ordered-set, binary-search, leetcode-352
type DisjointIntervals[T integer] struct {
	starts []T // ascending; starts[i] ≤ ends[i] < starts[i+1] − 1
	ends   []T
}

// Add adds v to the set: joins it to the range ending at v−1 and/or the range starting at v+1, or
// starts a new range [v, v]. O(log k + k).
func (d *DisjointIntervals[T]) Add(v T) {
	i := search.UpperBound(d.starts, v) // ranges [0, i) start at or before v
	if i > 0 && d.ends[i-1] >= v {
		return // already covered
	}
	// ends[i-1] < v, so v−1 does not underflow; v < starts[i], so starts[i]−1 does not underflow.
	joinsLeft := i > 0 && d.ends[i-1] == v-1
	joinsRight := i < len(d.starts) && d.starts[i]-1 == v
	switch {
	case joinsLeft && joinsRight:
		d.ends[i-1] = d.ends[i]
		d.starts = slices.Delete(d.starts, i, i+1)
		d.ends = slices.Delete(d.ends, i, i+1)
	case joinsLeft:
		d.ends[i-1] = v
	case joinsRight:
		d.starts[i] = v
	default:
		d.starts = slices.Insert(d.starts, i, v)
		d.ends = slices.Insert(d.ends, i, v)
	}
}

// Contains reports whether v has been added. O(log k).
func (d *DisjointIntervals[T]) Contains(v T) bool {
	i := search.UpperBound(d.starts, v)
	return i > 0 && d.ends[i-1] >= v
}

// Len returns the number of disjoint ranges. O(1).
func (d *DisjointIntervals[T]) Len() int {
	return len(d.starts)
}

// Intervals appends the ranges to dst, ascending by Start, and returns the extended slice. Pass nil to
// allocate once, sized up front, or buf[:0] to reuse a buffer. O(k).
func (d *DisjointIntervals[T]) Intervals(dst []interval.Interval[T]) []interval.Interval[T] {
	dst = slices.Grow(dst, len(d.starts))
	for i, s := range d.starts {
		dst = append(dst, interval.Interval[T]{Start: s, End: d.ends[i]})
	}
	return dst
}
