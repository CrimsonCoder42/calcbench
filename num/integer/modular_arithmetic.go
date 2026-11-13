package integer

import "math/bits"

// mulMod returns a·b mod n through the 128-bit product. It requires n > 0; with a, b < n the high word is
// below n, and bits.Rem64 never panics for n > 0 in any case.
func mulMod(a, b, n uint64) uint64 {
	hi, lo := bits.Mul64(a, b)
	return bits.Rem64(hi, lo, n)
}

// addMod returns (a+b) mod n for a, b < n without overflow, even when n is close to 2⁶⁴.
func addMod(a, b, n uint64) uint64 {
	s, carry := bits.Add64(a, b, 0)
	if carry != 0 || s >= n {
		s -= n
	}
	return s
}

// powMod returns base^exp mod n by right-to-left binary exponentiation (Knuth, TAOCP Vol. 2, 3rd ed.,
// §4.6.3, Algorithm A). It requires n > 1. O(log exp) multiplications, O(1) space.
func powMod(base, exp, n uint64) uint64 {
	result := uint64(1)
	base %= n
	for exp > 0 {
		if exp&1 == 1 {
			result = mulMod(result, base, n)
		}
		base = mulMod(base, base, n)
		exp >>= 1
	}
	return result
}
