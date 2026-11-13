package combinatorics

import (
	"fmt"
	"math/bits"

	"github.com/CrimsonCoder42/calcbench/num/integer"
)

// Binomial returns the binomial coefficient C(n,k) = n!/(k!(n−k)!) exactly: the number of k-element
// subsets of an n-set. C(n,k) = 0 when k > n. Every C(n,k) with n ≤ 67 fits a uint64; C(68,34) does not.
// Past uint64 use new(big.Int).Binomial(n, k), or LogBinomial for the log.
// Errors: ErrNegative when n < 0 or k < 0; ErrOverflow when C(n,k) > 2⁶⁴ − 1.
//
// ID:      NUM-COMBI-003
// Time:    O(m log n)   m = min(k, n−k, 64) steps (the partial result C(n−k+i, i) ≥ 2ⁱ, so a result that fits needs fewer than 64), each with one O(log n) binary gcd
// Space:   O(1), no allocation
// Source:  NIST DLMF, https://dlmf.nist.gov (read 2026-11-12), §26.3 eqs. 26.3.1, 26.3.2, 26.3.6, Table 26.3.1
// Source:  OEIS A007318, https://oeis.org/A007318 (read 2026-11-12)
// Tags:    combinatorics, binomial, combinations, nCr, choose, counting
// Units:   n, k: counts of objects; result: a count of subsets
// Example: Binomial(6, 3) = 20, Binomial(8, 2) = 28, Binomial(10, 5) = 252 (DLMF Table 26.3.1)
func Binomial(n, k int) (uint64, error) {
	if n < 0 {
		return 0, fmt.Errorf("%w: n", ErrNegative)
	}
	if k < 0 {
		return 0, fmt.Errorf("%w: k", ErrNegative)
	}
	r, ok := binomialU64(uint64(n), uint64(k))
	if !ok {
		return 0, fmt.Errorf("%w: C(n,k)", ErrOverflow)
	}
	return r, nil
}

// binomialU64 returns C(n,k) and true, or 0 and false when it exceeds 2⁶⁴ − 1. Step i turns
// r = C(n−k+i−1, i−1) into C(n−k+i, i) = r·(n−k+i)/i (DLMF 26.3.6). With g = gcd(r, i), i/g divides
// n−k+i because i divides r·(n−k+i) and gcd(r/g, i/g) = 1; so (r/g)·((n−k+i)/(i/g)) is exact, and the
// only product formed is the next partial result, which is ≤ C(n,k): an overflow is never false.
func binomialU64(n, k uint64) (uint64, bool) {
	if k > n {
		return 0, true
	}
	if k > n-k {
		k = n - k
	}
	base := n - k
	r := uint64(1)
	for i := uint64(1); i <= k; i++ {
		g := integer.GCD(r, i)
		hi, lo := bits.Mul64(r/g, (base+i)/(i/g))
		if hi != 0 {
			return 0, false
		}
		r = lo
	}
	return r, true
}
