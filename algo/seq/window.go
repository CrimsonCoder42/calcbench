package seq

import (
	"fmt"
	"iter"
)

// Window returns an iterator over every run of width consecutive elements of seq, oldest element first:
// m elements give m-width+1 windows, none when m < width. Every window is the SAME slice, a view into
// a reused ring buffer: it is valid until the loop body returns, so copy it (slices.Clone) to keep it
// and do not modify it. When the consumer stops early, seq is stopped at once. A nil seq is the empty
// sequence. Until the first window fills, the buffer grows with the elements seen, so a huge width
// costs nothing on a short seq.
// Errors: ErrSize when width < 1.
//
// Ring: after the first window, the buffer has 2·width slots and element i is written to slot i mod
// width AND to slot i mod width + width, so the window that starts at slot s is the contiguous slice
// [s, s+width): each step is two writes, never a copy of the window.
//
// ID:      ALG-SEQ-008
// Time:    O(1) amortized per element (two writes once the first window is full) plus the cost of seq
// Space:   O(width) extra (2·width slots, allocated once the first window is full)
// Source:  The Go Programming Language Specification, Go 1.23+, "For statements with range clause" (range over a function: yield returns false when the loop body stops and must not be called again), https://go.dev/ref/spec#For_range (read 2026-11-12)
// Source:  Knuth, The Art of Computer Programming, Vol. 1, 3rd ed., Addison-Wesley, 1997, §2.2.2 (a queue in a circular array, indices taken modulo M)
// Tags:    iterator, iter-seq, range-over-func, sliding-window, window, ring-buffer, lazy
func Window[T any](seq iter.Seq[T], width int) (iter.Seq[[]T], error) {
	if width < 1 {
		return nil, fmt.Errorf("%w: width = %d", ErrSize, width)
	}
	return func(yield func([]T) bool) {
		if seq == nil {
			return
		}
		var fill []T // the first width elements, grown as they arrive
		var ring []T // 2·width slots once the first window is full
		slot := 0    // slot of the next element, in [0, width)
		for v := range seq {
			if ring == nil {
				fill = append(fill, v)
				if len(fill) < width {
					continue
				}
				ring = make([]T, 2*width) // width ≤ len(fill), which fits in memory, so 2·width fits in an int
				copy(ring, fill)
				copy(ring[width:], fill)
				fill = nil
				if !yield(ring[:width]) {
					return
				}
				continue
			}
			ring[slot] = v
			ring[slot+width] = v
			slot++
			if slot == width {
				slot = 0
			}
			if !yield(ring[slot : slot+width]) {
				return
			}
		}
	}, nil
}
