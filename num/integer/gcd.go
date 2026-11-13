package integer

import "math/bits"

// GCD returns the greatest common divisor of a and b, by the binary algorithm: shifts and subtractions,
// no division. GCD(0, b) = b, GCD(a, 0) = a and GCD(0, 0) = 0.
//
// ID:      NUM-INTEG-004
// Time:    O(log a + log b) steps (each step removes at least one bit)
// Space:   O(1) (no allocation)
// Source:  Knuth, The Art of Computer Programming, Vol. 2: Seminumerical Algorithms, 3rd ed., Addison-Wesley, 1998, §4.5.2, Algorithm B (binary gcd algorithm)
// Tags:    gcd, greatest-common-divisor, binary-gcd, stein, number-theory, no-alloc
func GCD(a, b uint64) uint64 {
	if a == 0 {
		return b
	}
	if b == 0 {
		return a
	}
	// B1: the common power of two; B2-B6: strip factors of 2 and subtract the smaller odd from the larger.
	shift := bits.TrailingZeros64(a | b)
	a >>= bits.TrailingZeros64(a)
	for b != 0 {
		b >>= bits.TrailingZeros64(b)
		if a > b {
			a, b = b, a
		}
		b -= a
	}
	return a << shift
}
