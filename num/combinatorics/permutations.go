package combinatorics

import (
	"fmt"
	"math/bits"
)

// Permutations returns P(n,k) = n!/(n−k)! = n(n−1)…(n−k+1) exactly, the falling factorial: the number of
// ordered selections of k of n distinct objects. P(n,k) = 0 when k > n. Past uint64 use
// new(big.Int).MulRange(n−k+1, n).
// Errors: ErrNegative when n < 0 or k < 0; ErrOverflow when P(n,k) > 2⁶⁴ − 1.
//
// ID:      NUM-COMBI-005
// Time:    O(min(k, 64))   (every factor but the last is ≥ 2, so overflow stops the loop within 64 multiplies)
// Space:   O(1), no allocation
// Source:  NIST DLMF, https://dlmf.nist.gov (read 2026-11-12), §26.1 (falling factorial), §26.17 Table 26.17.1 (one-to-one maps)
// Tags:    combinatorics, permutations, nPr, falling-factorial, counting
// Units:   n, k: counts of objects; result: a count of ordered selections
// Example: Permutations(10, 3) = 720
func Permutations(n, k int) (uint64, error) {
	if n < 0 {
		return 0, fmt.Errorf("%w: n", ErrNegative)
	}
	if k < 0 {
		return 0, fmt.Errorf("%w: k", ErrNegative)
	}
	if k > n {
		return 0, nil
	}
	r := uint64(1)
	for i := 0; i < k; i++ {
		hi, lo := bits.Mul64(r, uint64(n-i))
		if hi != 0 {
			return 0, fmt.Errorf("%w: P(n,k)", ErrOverflow)
		}
		r = lo
	}
	return r, nil
}
