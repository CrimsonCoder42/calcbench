package seq

import "iter"

// Map returns an iterator over fn(v) for each element v of seq, in order. A nil seq or a nil fn gives
// the empty sequence. When the consumer stops early, seq is stopped at once.
//
// ID:      ALG-SEQ-001
// Time:    O(1) per element plus the cost of fn and of seq
// Space:   O(1) extra; one closure allocated when Map is called, none per element
// Source:  The Go Programming Language Specification, Go 1.23+, "For statements with range clause" (range over a function: yield returns false when the loop body stops and must not be called again), https://go.dev/ref/spec#For_range (read 2026-11-12)
// Source:  Go standard library, package iter documentation, https://pkg.go.dev/iter (read 2026-11-12)
// Tags:    iterator, iter-seq, range-over-func, map, transform, lazy
func Map[T, U any](seq iter.Seq[T], fn func(T) U) iter.Seq[U] {
	return func(yield func(U) bool) {
		if seq == nil || fn == nil {
			return
		}
		for v := range seq {
			if !yield(fn(v)) {
				return
			}
		}
	}
}
