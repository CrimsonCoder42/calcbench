package seq

import (
	"fmt"
	"iter"
)

// Chunk returns an iterator over consecutive groups of size elements of seq; the last group holds the
// remaining 1…size elements and an empty seq gives no groups. Every group is the SAME slice, rewritten
// between steps: it is valid until the loop body returns, so copy it (slices.Clone) to keep it and do not
// modify it. When the consumer stops early, seq is stopped at once. A nil seq is the empty sequence.
// The buffer grows with the elements actually seen, so a huge size costs nothing until it fills.
// Errors: ErrSize when size < 1.
//
// ID:      ALG-SEQ-004
// Time:    O(1) amortized per element plus the cost of seq
// Space:   O(min(size, m)) extra   m = elements seq produces; no allocation per group once grown
// Source:  The Go Programming Language Specification, Go 1.23+, "For statements with range clause" (range over a function: yield returns false when the loop body stops and must not be called again), https://go.dev/ref/spec#For_range (read 2026-11-12)
// Source:  Go standard library, package iter documentation, https://pkg.go.dev/iter (read 2026-11-12)
// Tags:    iterator, iter-seq, range-over-func, chunk, group, partition, fixed-size, lazy
func Chunk[T any](seq iter.Seq[T], size int) (iter.Seq[[]T], error) {
	if size < 1 {
		return nil, fmt.Errorf("%w: size = %d", ErrSize, size)
	}
	return func(yield func([]T) bool) {
		if seq == nil {
			return
		}
		var buf []T
		for v := range seq {
			buf = append(buf, v)
			if len(buf) == size {
				if !yield(buf) {
					return
				}
				clear(buf) // release what the old elements refer to
				buf = buf[:0]
			}
		}
		if len(buf) > 0 {
			yield(buf)
		}
	}, nil
}
