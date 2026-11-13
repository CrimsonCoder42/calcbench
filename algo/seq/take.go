package seq

import (
	"errors"
	"fmt"
	"iter"
)

// Errors returned by this package. They are wrapped with the offending parameter; test with errors.Is.
var (
	// ErrCount: a count is negative.
	ErrCount = errors.New("seq: count must not be negative")
	// ErrSize: a group size or window width is less than 1.
	ErrSize = errors.New("seq: size must be at least 1")
)

// Take returns an iterator over the first n elements of seq (all of them when seq has fewer). After the
// n-th element it stops seq without asking it for the (n+1)-th, so Take bounds an infinite iterator.
// n = 0 gives the empty sequence and never starts seq. A nil seq is the empty sequence.
// Errors: ErrCount when n < 0.
//
// ID:      ALG-SEQ-003
// Time:    O(min(n, m))   m = elements seq produces; plus the cost of seq
// Space:   O(1) extra; one closure allocated when Take is called, none per element
// Source:  The Go Programming Language Specification, Go 1.23+, "For statements with range clause" (range over a function: yield returns false when the loop body stops and must not be called again), https://go.dev/ref/spec#For_range (read 2026-11-12)
// Source:  Go standard library, package iter documentation, https://pkg.go.dev/iter (read 2026-11-12)
// Tags:    iterator, iter-seq, range-over-func, take, limit, prefix, lazy, infinite
func Take[T any](seq iter.Seq[T], n int) (iter.Seq[T], error) {
	if n < 0 {
		return nil, fmt.Errorf("%w: n = %d", ErrCount, n)
	}
	return func(yield func(T) bool) {
		if seq == nil || n == 0 {
			return
		}
		taken := 0
		for v := range seq {
			if !yield(v) {
				return
			}
			taken++
			if taken == n {
				return // stop seq before it produces element n+1
			}
		}
	}, nil
}
