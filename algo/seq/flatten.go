package seq

import "iter"

// Flatten returns an iterator over the elements of each inner iterator of seqs, inner by inner, in order.
// A nil seqs or a nil inner iterator is empty. When the consumer stops early, the current inner iterator
// and seqs are both stopped at once and no further inner iterator is started.
//
// ID:      ALG-SEQ-006
// Time:    O(1) per element and per inner iterator, plus the cost of the iterators
// Space:   O(1) extra; one closure allocated when Flatten is called, none per element
// Source:  The Go Programming Language Specification, Go 1.23+, "For statements with range clause" (range over a function: yield returns false when the loop body stops and must not be called again), https://go.dev/ref/spec#For_range (read 2026-11-12)
// Source:  Go standard library, package iter documentation, https://pkg.go.dev/iter (read 2026-11-12)
// Tags:    iterator, iter-seq, range-over-func, flatten, concat, chain, lazy
func Flatten[T any](seqs iter.Seq[iter.Seq[T]]) iter.Seq[T] {
	return func(yield func(T) bool) {
		if seqs == nil {
			return
		}
		for inner := range seqs {
			if inner == nil {
				continue
			}
			for v := range inner {
				if !yield(v) {
					return
				}
			}
		}
	}
}
