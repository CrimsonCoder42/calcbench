package intmath

import "math/bits"

// ISqrt returns ⌊√n⌋, the greatest integer whose square is ≤ n, exactly for every uint64.
// Integer Newton (Heron) iteration x ← ⌊(x + ⌊n/x⌋)/2⌋ from a start 2^⌈b/2⌉ ≥ √n (b = bit length of
// n), stopped at the first step that does not decrease x. Termination: while x > ⌊√n⌋ the step
// strictly decreases x, and no step goes below ⌊√n⌋ (AM–GM), so the sequence of integers falls to
// ⌊√n⌋ and stops there. No intermediate sum overflows: x ≤ 2^32 and ⌊n/x⌋ ≤ 2^32 + 2.
//
// ID:      ALG-INTMTH-002
// Time:    O(log b)   b = bit length of n (≤ 64); quadratic convergence from within a factor 2, at most 7 divisions
// Space:   O(1), no allocation
// Source:  Cohen, A Course in Computational Algebraic Number Theory, Springer GTM 138, 1993, §1.7.1, Algorithm 1.7.1 (Integer Square Root) and its proof of correctness
// Source:  Crandall, Pomerance, Prime Numbers: A Computational Perspective, 2nd ed., Springer, 2005, §9.2.2, Algorithm 9.2.11 (integer part of square root)
// Tags:    integer, square-root, isqrt, newton, heron, exact
func ISqrt(n uint64) uint64 {
	if n < 2 {
		return n
	}
	x := uint64(1) << ((bits.Len64(n) + 1) / 2)
	for {
		y := (x + n/x) >> 1
		if y >= x {
			return x
		}
		x = y
	}
}
