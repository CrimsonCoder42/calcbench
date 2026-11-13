package intmath

import (
	"fmt"
	"math"
	"math/bits"
)

// hammingCountUint64 is the number of Hamming numbers ≤ math.MaxUint64; the largest of them,
// 18_432_000_000_000_000_000 = 2^26·3^2·5^15, is NthHammingNumber(hammingCountUint64). Counted exactly
// with math/big (every 2^i·3^j·5^k ≤ 2^64−1); the test recounts it.
const hammingCountUint64 = 13282

// IsHammingNumber reports whether n = 2^i·3^j·5^k for some i, j, k ≥ 0 (a Hamming, regular or 5-smooth
// number: 1, 2, 3, 4, 5, 6, 8, 9, 10, 12, …). By unique factorisation it suffices to divide out every 2,
// 3 and 5 and see whether 1 remains. 0 is not a Hamming number.
//
// ID:      ALG-INTMTH-009
// Time:    O(log n)   at most log2 n + log3 n + log5 n divisions (the 2s by one shift)
// Space:   O(1), no allocation
// Source:  Dijkstra, A Discipline of Programming, Prentice-Hall, 1976, ch. 17 "An exercise attributed to R. W. Hamming" (the set: 1 is in it, and 2x, 3x, 5x are whenever x is, and nothing else)
// Tags:    integer, hamming-number, regular-number, 5-smooth, ugly-number, smooth-number
func IsHammingNumber(n uint64) bool {
	if n == 0 {
		return false
	}
	n >>= uint(bits.TrailingZeros64(n))
	for n%3 == 0 {
		n /= 3
	}
	for n%5 == 0 {
		n /= 5
	}
	return n == 1
}

// timesSaturated returns h·k, or math.MaxUint64 when the product overflows. MaxUint64 is odd, so it is
// never a Hamming number reached by doubling; NthHammingNumber only calls this for n ≤ 13282, where the
// minimum of the three candidates is always a true product below 2^64.
func timesSaturated(h, k uint64) uint64 {
	hi, lo := bits.Mul64(h, k)
	if hi != 0 {
		return math.MaxUint64
	}
	return lo
}

// NthHammingNumber returns the n-th Hamming number in ascending order, counting from NthHammingNumber(1)
// = 1 (then 2, 3, 4, 5, 6, 8, …). Dijkstra's method: the sequence is merged from three copies of itself
// multiplied by 2, 3 and 5, each read by its own index, so every number is produced once, in order.
// Errors: ErrOutOfRange when n < 1; ErrOverflow when n > 13282, as the 13282nd Hamming number
// (18_432_000_000_000_000_000) is the largest that fits in a uint64 (checked before any allocation).
//
// ID:      ALG-INTMTH-010
// Time:    O(n)   three candidate products and one comparison per number
// Space:   O(n) extra (the sequence so far; at most 13282 values, 104 KiB), one allocation
// Source:  Dijkstra, A Discipline of Programming, Prentice-Hall, 1976, ch. 17 "An exercise attributed to R. W. Hamming" (generation in increasing order by three advancing indexes into the sequence)
// Tags:    integer, hamming-number, regular-number, 5-smooth, ugly-number, merge, sequence
func NthHammingNumber(n int) (uint64, error) {
	if n < 1 {
		return 0, fmt.Errorf("%w: n = %d, want n ≥ 1", ErrOutOfRange, n)
	}
	if n > hammingCountUint64 {
		return 0, fmt.Errorf("%w: n = %d exceeds the %d Hamming numbers below 2^64", ErrOverflow, n, hammingCountUint64)
	}
	seq := make([]uint64, n)
	seq[0] = 1
	i2, i3, i5 := 0, 0, 0
	for k := 1; k < n; k++ {
		c2, c3, c5 := timesSaturated(seq[i2], 2), timesSaturated(seq[i3], 3), timesSaturated(seq[i5], 5)
		next := min(c2, c3, c5)
		seq[k] = next
		// Advance every index that produced next, so equal products (6 = 2·3 = 3·2) appear once.
		if c2 == next {
			i2++
		}
		if c3 == next {
			i3++
		}
		if c5 == next {
			i5++
		}
	}
	return seq[n-1], nil
}
