package combinatorics

import (
	"fmt"
	"math/bits"
)

// CircularPermutations returns the number of ways to place k of n distinct objects around a circle, two
// placements being equal when one is a rotation of the other: P(n,k)/k = C(n,k)·(k−1)!. With k = n it is
// (n−1)!. Equivalently, the number of k-cycles in the symmetric group on n letters. The result is formed
// as C(n,k)·(k−1)!, a product of two factors ≤ the result, so an overflow is never false.
// Errors: ErrNegative when n < 0 or k < 0; ErrOutOfRange when k = 0 or k > n; ErrOverflow when the
// count > 2⁶⁴ − 1.
//
// ID:      NUM-COMBI-009
// Time:    O(m log n)   m = min(k, n−k, 64), as Binomial
// Space:   O(1), no allocation
// Source:  OEIS A111492, https://oeis.org/A111492 (read 2026-11-12), a(n,k) = (k−1)!·C(n,k) = P(n,k)/k
// Tags:    combinatorics, permutations, circular, cycles
// Units:   n: objects available; k: objects placed; result: a count of circular arrangements
// Example: OEIS A111492 row n = 5, k = 1…5: 5, 10, 20, 30, 24
func CircularPermutations(n, k int) (uint64, error) {
	c, err := circularCheck(n, k)
	if err != nil {
		return 0, err
	}
	r, ok := mulFactorial(c, uint64(k-1), 1)
	if !ok {
		return 0, fmt.Errorf("%w: P(n,k)/k", ErrOverflow)
	}
	return r, nil
}

// CircularPermutationsUndirected returns the number of ways to place k of n distinct objects around a
// circle that may also be turned over (a necklace or key ring): placements equal under rotation or
// reflection. For k ≥ 3 that is C(n,k)·(k−1)!/2; for k ≤ 2 every placement is its own mirror image and the
// count equals CircularPermutations. Equivalently, the number of undirected k-cycles on n labelled nodes.
// It is formed as C(n,k)·((k−1)!/2), so an overflow is never false.
// Errors: ErrNegative when n < 0 or k < 0; ErrOutOfRange when k = 0 or k > n; ErrOverflow when the
// count > 2⁶⁴ − 1.
//
// ID:      NUM-COMBI-010
// Time:    O(m log n)   m = min(k, n−k, 64), as Binomial
// Space:   O(1), no allocation
// Source:  OEIS A144151, https://oeis.org/A144151 (read 2026-11-12), T(n,k) = C(n,k) for k ≤ 2, else C(n,k)·(k−1)!/2
// Tags:    combinatorics, permutations, circular, necklace, cycles
// Units:   n: objects available; k: objects placed; result: a count of circular arrangements
// Example: OEIS A144151 row n = 5, k = 1…5: 5, 10, 10, 15, 12
func CircularPermutationsUndirected(n, k int) (uint64, error) {
	c, err := circularCheck(n, k)
	if err != nil {
		return 0, err
	}
	halve := uint64(1)
	if k >= 3 {
		halve = 2 // (k−1)! is even for k ≥ 3
	}
	r, ok := mulFactorial(c, uint64(k-1), halve)
	if !ok {
		return 0, fmt.Errorf("%w: P(n,k)/(2k)", ErrOverflow)
	}
	return r, nil
}

// circularCheck validates (n, k) for the circular counts and returns C(n,k).
func circularCheck(n, k int) (uint64, error) {
	if n < 0 {
		return 0, fmt.Errorf("%w: n", ErrNegative)
	}
	if k < 0 {
		return 0, fmt.Errorf("%w: k", ErrNegative)
	}
	if k == 0 || k > n {
		return 0, fmt.Errorf("%w: k must be in 1..n", ErrOutOfRange)
	}
	c, ok := binomialU64(uint64(n), uint64(k))
	if !ok {
		return 0, fmt.Errorf("%w: C(n,k)", ErrOverflow)
	}
	return c, nil
}

// mulFactorial returns c·(f!/div) and true, or false when it exceeds 2⁶⁴ − 1. div divides f!. c ≥ 1, so
// f!/div > 2⁶⁴ − 1 means the product does too.
func mulFactorial(c, f, div uint64) (uint64, bool) {
	if f >= uint64(len(factorials)) {
		// f ≥ 21: f!/2 ≥ 21!/2 > 2⁶⁴ − 1.
		return 0, false
	}
	hi, lo := bits.Mul64(c, factorials[f]/div)
	return lo, hi == 0
}
