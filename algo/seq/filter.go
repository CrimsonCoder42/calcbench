package seq

import "iter"

// Filter returns an iterator over the elements v of seq for which keep(v) is true, in order. A nil seq
// or a nil keep gives the empty sequence. When the consumer stops early, seq is stopped at once.
//
// ID:      ALG-SEQ-002
// Time:    O(1) per upstream element plus the cost of keep and of seq
// Space:   O(1) extra; one closure allocated when Filter is called, none per element
// Source:  The Go Programming Language Specification, Go 1.23+, "For statements with range clause" (range over a function: yield returns false when the loop body stops and must not be called again), https://go.dev/ref/spec#For_range (read 2026-11-12)
// Source:  Go standard library, package iter documentation, https://pkg.go.dev/iter (read 2026-11-12)
// Tags:    iterator, iter-seq, range-over-func, filter, predicate, lazy
func Filter[T any](seq iter.Seq[T], keep func(T) bool) iter.Seq[T] {
	return func(yield func(T) bool) {
		if seq == nil || keep == nil {
			return
		}
		for v := range seq {
			if keep(v) && !yield(v) {
				return
			}
		}
	}
}
