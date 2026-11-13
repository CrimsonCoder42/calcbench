package window

import (
	"fmt"
	"math"
	"slices"

	"github.com/CrimsonCoder42/calcbench/algo/heap"
)

// medianEntry is one value in a MovingMedian heap with its position in values. Entries are ordered by
// (value, position), a total order in which no two entries are equal, so which heap holds a value is
// decided by one comparison with the lower heap's top.
type medianEntry struct {
	v float64
	i int
}

// entryLess orders entries by value, then by position.
func entryLess(a, b medianEntry) bool {
	return a.v < b.v || (a.v == b.v && a.i < b.i)
}

// MovingMedian appends to dst the median of each window of windowLen consecutive values and returns
// the extended slice: n-windowLen+1 medians for n values, none when windowLen > n.
// For an odd windowLen the median is the middle value of the window; for an even windowLen it is the
// arithmetic mean of the two middle values (computed without overflow), so every result is exact for
// odd windows and correctly rounded for even ones. -0 and +0 compare equal.
// values is read, never modified.
// Errors: ErrWindowLen when windowLen < 1; ErrNotFinite when a value is NaN or ±Inf (a NaN has no
// place in the order and the mean of +Inf and -Inf is undefined). On error dst is returned with its
// length unchanged.
//
// Method: a max-heap holds the lower half of the window and a min-heap the upper half, the lower half
// keeping the extra value for an odd count. A leaving value is not searched for: it is counted out of
// the half it belongs to (one comparison with the lower top) and its entry is discarded when it
// reaches a heap's top ("lazy deletion"). A heap whose stale entries make it larger than
// 2·windowLen+2 is rebuilt from its live entries, so memory stays O(windowLen).
//
// ID:      ALG-WINDOW-015
// Time:    O(n log w) amortized   n = len(values), w = windowLen (each value is pushed and popped O(1)
// times amortized, each heap holds O(w) entries)
// Space:   O(w) extra (two heaps of at most 2w+3 entries and a rebuild buffer; a fixed number of
// allocations, independent of n); O(n-w+1) for the result,
// one allocation only if dst lacks capacity
// Source:  Härdle, Steiger, "Algorithm AS 296: Optimal Median Smoothing", Journal of the Royal Statistical Society Series C (Applied Statistics), 1995, 44(2):258–264, doi:10.2307/2986349 (running median of span w by a double heap, O(log w) per step)
// Source:  Sedgewick, Wayne, Algorithms, 4th ed., Addison-Wesley, 2011, §2.4 Creative Problems, "Dynamic median-finding" (max-heap below the median, min-heap above)
// Tags:    window, sliding-window, moving-median, running-median, median-filter, two-heaps, lazy-deletion, time-series
func MovingMedian(dst, values []float64, windowLen int) ([]float64, error) {
	if windowLen < 1 {
		return dst, fmt.Errorf("%w: windowLen = %d", ErrWindowLen, windowLen)
	}
	for i, v := range values {
		if math.IsNaN(v) || math.IsInf(v, 0) {
			return dst, fmt.Errorf("%w: values[%d] = %v", ErrNotFinite, i, v)
		}
	}
	n := len(values)
	if windowLen > n {
		return dst, nil
	}
	out := slices.Grow(dst, n-windowLen+1)

	limit := 2*windowLen + 2 // heap size that triggers a rebuild
	// The less functions are non-nil, so NewHeap cannot fail.
	lower, _ := heap.NewHeap(func(a, b medianEntry) bool { return entryLess(b, a) }, limit+1) // max-heap
	upper, _ := heap.NewHeap(entryLess, limit+1)                                              // min-heap
	var scratch []medianEntry
	liveLower, liveUpper := 0, 0 // entries of the current window in each heap

	for end, v := range values {
		stale := end - windowLen // entries at positions <= stale are outside the window
		// 1. Count the leaving value out of its half. The tops are live for the previous window, and
		//    every live lower entry orders at or before the lower top while every live upper entry
		//    orders after it, so one comparison names the half.
		if stale >= 0 {
			leaving := medianEntry{values[stale], stale}
			if top, ok := lower.Peek(); ok && !entryLess(top, leaving) {
				liveLower--
			} else {
				liveUpper--
			}
		}
		// 2. Place the entering value. Comparing with a top that has just left is still correct: that
		//    entry ordered after every live lower entry and before every live upper entry.
		entering := medianEntry{v, end}
		if top, ok := lower.Peek(); !ok || entryLess(entering, top) {
			lower.Push(entering)
			liveLower++
		} else {
			upper.Push(entering)
			liveUpper++
		}
		// 3. Rebalance so liveLower - liveUpper is 0 or 1, moving only live tops.
		for liveLower > liveUpper+1 {
			dropStale(lower, stale)
			x, _ := lower.Pop()
			upper.Push(x)
			liveLower--
			liveUpper++
		}
		for liveLower < liveUpper {
			dropStale(upper, stale)
			x, _ := upper.Pop()
			lower.Push(x)
			liveUpper--
			liveLower++
		}
		dropStale(lower, stale)
		dropStale(upper, stale)
		// 4. Bound memory: rebuild a heap that stale entries have grown past the limit.
		if lower.Len() > limit {
			scratch = rebuildLive(lower, scratch, stale)
		}
		if upper.Len() > limit {
			scratch = rebuildLive(upper, scratch, stale)
		}
		if end < windowLen-1 {
			continue
		}
		lo, _ := lower.Peek()
		if windowLen%2 == 1 {
			out = append(out, lo.v)
			continue
		}
		hi, _ := upper.Peek()
		out = append(out, midpoint(lo.v, hi.v))
	}
	return out, nil
}

// dropStale pops entries at positions <= stale from the top of h, so its top is live.
func dropStale(h *heap.Heap[medianEntry], stale int) {
	for {
		top, ok := h.Peek()
		if !ok || top.i > stale {
			return
		}
		h.Pop()
	}
}

// rebuildLive empties h and pushes back only its live entries (positions > stale), using scratch as the
// buffer, and returns scratch for reuse. O(L log L) for L = h.Len(); called only when L > 2w+2 and
// leaving at most w+1 live entries, so its cost is amortized over at least w+1 later pushes.
func rebuildLive(h *heap.Heap[medianEntry], scratch []medianEntry, stale int) []medianEntry {
	if cap(scratch) < h.Len() {
		scratch = make([]medianEntry, 0, h.Len())
	}
	scratch = scratch[:0]
	for {
		x, ok := h.Pop()
		if !ok {
			break
		}
		if x.i > stale {
			scratch = append(scratch, x)
		}
	}
	for _, x := range scratch {
		h.Push(x)
	}
	return scratch
}

// midpoint returns (a+b)/2 correctly rounded, without overflow for finite a and b: the sum rounds once
// and the halving is exact unless the result is subnormal; on overflow each half is taken first.
func midpoint(a, b float64) float64 {
	if m := (a + b) / 2; !math.IsInf(m, 0) {
		return m
	}
	return a/2 + b/2
}
