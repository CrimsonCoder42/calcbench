package intmath

// maxPowerOfThreeUint64 is 3^40, the largest power of 3 that fits in a uint64 (3^41 > 2^64).
const maxPowerOfThreeUint64 = 12_157_665_459_056_928_801

// IsPowerOfThree reports whether x is 3^k for some k ≥ 0 (1, 3, 9, 27, …). Because 3 is prime, the
// positive divisors of 3^40 are exactly 3^0 … 3^40, so x > 0 is a power of three exactly when it divides
// 3^40: one remainder, no loop. Zero and negative values return false.
//
// ID:      ALG-INTMTH-012
// Time:    O(1)   one 64-bit remainder
// Space:   O(1), no allocation
// Source:  Hardy, Wright, An Introduction to the Theory of Numbers, 6th ed., Oxford University Press, 2008, Theorem 2 (the fundamental theorem of arithmetic: factorisation into primes is unique; stated in §1.3, proved in §2.10), so every divisor of p^k is p^j, j ≤ k
// Tags:    integer, power-of-three, divisibility, prime-power
func IsPowerOfThree[T integer](x T) bool {
	return x > 0 && maxPowerOfThreeUint64%uint64(x) == 0
}
