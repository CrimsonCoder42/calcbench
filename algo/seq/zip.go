package seq

import "iter"

// Zip returns an iterator over the pairs (a₀, b₀), (a₁, b₁), … of the two iterators, stopping at the
// shorter: when either runs out, or when the consumer stops, both are stopped. first is ranged over
// directly; second is pulled with iter.Pull, and its stop function is ALWAYS called (deferred): on early
// stop, when first runs out, when second runs out, and if the loop body panics. second is not started
// until first has produced its first element. A nil iterator is the empty sequence.
//
// ID:      ALG-SEQ-005
// Time:    O(1) per pair plus the cost of the two iterators (iter.Pull switches coroutines, no goroutine)
// Space:   O(1) extra; iter.Pull allocates its coroutine state once per range loop, none per pair
// Source:  The Go Programming Language Specification, Go 1.23+, "For statements with range clause" (range over a function: yield returns false when the loop body stops and must not be called again), https://go.dev/ref/spec#For_range (read 2026-11-12)
// Source:  Go standard library, package iter documentation, "Pulling Values" (stop must be called when the caller is no longer interested in next values and next has not yet signaled that the sequence is over), https://pkg.go.dev/iter#Pull (read 2026-11-12)
// Tags:    iterator, iter-seq, iter-seq2, range-over-func, zip, pair, pull, lazy
func Zip[A, B any](first iter.Seq[A], second iter.Seq[B]) iter.Seq2[A, B] {
	return func(yield func(A, B) bool) {
		if first == nil || second == nil {
			return
		}
		var next func() (B, bool)
		var stop func()
		defer func() {
			if stop != nil {
				stop()
			}
		}()
		for a := range first {
			if next == nil {
				next, stop = iter.Pull(second)
			}
			b, ok := next()
			if !ok || !yield(a, b) {
				return
			}
		}
	}
}
