package seq

import (
	"fmt"
	"iter"
)

// Batch returns an iterator over consecutive groups of seq's elements, each closed when it reaches
// maxSize elements OR when endsBatch(v) is true for the element v just added (v is the group's last
// element); the remaining elements form a final group, and an empty seq gives no groups. With a nil
// endsBatch it groups exactly as Chunk; it differs from Chunk only by the predicate flush (for example,
// records up to a terminator, or messages up to a commit marker, but never more than maxSize).
// Every group is the SAME slice, rewritten between steps: it is valid until the loop body returns, so
// copy it (slices.Clone) to keep it and do not modify it. When the consumer stops early, seq is stopped
// at once. A nil seq is the empty sequence.
// Errors: ErrSize when maxSize < 1.
//
// ID:      ALG-SEQ-007
// Time:    O(1) amortized per element plus the cost of endsBatch and of seq
// Space:   O(min(maxSize, m)) extra   m = elements seq produces; no allocation per group once grown
// Source:  The Go Programming Language Specification, Go 1.23+, "For statements with range clause" (range over a function: yield returns false when the loop body stops and must not be called again), https://go.dev/ref/spec#For_range (read 2026-11-12)
// Source:  Go standard library, package iter documentation, https://pkg.go.dev/iter (read 2026-11-12)
// Tags:    iterator, iter-seq, range-over-func, batch, group, delimiter, flush, size-limit, lazy
func Batch[T any](seq iter.Seq[T], maxSize int, endsBatch func(T) bool) (iter.Seq[[]T], error) {
	if maxSize < 1 {
		return nil, fmt.Errorf("%w: maxSize = %d", ErrSize, maxSize)
	}
	return func(yield func([]T) bool) {
		if seq == nil {
			return
		}
		var buf []T
		for v := range seq {
			buf = append(buf, v)
			if len(buf) == maxSize || (endsBatch != nil && endsBatch(v)) {
				if !yield(buf) {
					return
				}
				clear(buf)
				buf = buf[:0]
			}
		}
		if len(buf) > 0 {
			yield(buf)
		}
	}, nil
}
