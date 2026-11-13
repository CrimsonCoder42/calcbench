package combinatorics

import (
	"fmt"
	"math/big"
)

// derangements holds d(n) for n = 0…20; d(21) = 18 795 307 255 050 944 540 exceeds 2⁶⁴ − 1.
// Values: OEIS A000166, checked in the tests against DLMF 26.13.4 in math/big.
var derangements = [...]uint64{
	1,
	0,
	1,
	2,
	9,
	44,
	265,
	1854,
	14833,
	133496,
	1334961,
	14684570,
	176214841,
	2290792932,
	32071101049,
	481066515734,
	7697064251745,
	130850092279664,
	2355301661033953,
	44750731559645106,
	895014631192902121,
}

// Derangements returns d(n), the number of permutations of n objects that leave no object in its place
// (the subfactorial !n): d(n) = n!·Σⱼ₌₀..ₙ (−1)ʲ/j!. d(0) = 1, d(1) = 0. d(20) = 895 014 631 192 902 121 is
// the largest that fits a uint64; past it use DerangementsBig.
// Errors: ErrNegative when n < 0; ErrOverflow when n > 20.
//
// ID:      NUM-COMBI-011
// Time:    O(1)   (a 21-entry table)
// Space:   O(1), no allocation
// Source:  NIST DLMF, https://dlmf.nist.gov (read 2026-11-12), §26.13 eq. 26.13.4
// Source:  OEIS A000166, https://oeis.org/A000166 (read 2026-11-12)
// Tags:    combinatorics, derangements, subfactorial, counting
// Units:   n: a count of objects; result: a count of permutations
// Example: Derangements(4) = 9, Derangements(6) = 265 (OEIS A000166)
func Derangements(n int) (uint64, error) {
	if n < 0 {
		return 0, fmt.Errorf("%w: n", ErrNegative)
	}
	if n >= len(derangements) {
		return 0, fmt.Errorf("%w: n", ErrOverflow)
	}
	return derangements[n], nil
}

// DerangementsBig returns d(n), the number of permutations of n objects with no fixed point, as a
// *big.Int: Derangements without the uint64 limit. It uses d(n) = n·d(n−1) + (−1)ⁿ, d(0) = 1, which
// follows from DLMF 26.13.4 by splitting off the j = n term of the sum.
// Errors: ErrNegative when n < 0.
//
// ID:      NUM-COMBI-018
// Time:    O(n·b) bit operations   b = bit length of the result, O(n log n)
// Space:   O(b) bits
// Source:  NIST DLMF, https://dlmf.nist.gov (read 2026-11-12), §26.13 eq. 26.13.4
// Source:  OEIS A000166, https://oeis.org/A000166 (read 2026-11-12)
// Tags:    combinatorics, derangements, subfactorial, counting, big-int
// Units:   n: a count of objects; result: a count of permutations
// Example: DerangementsBig(22) = 413 496 759 611 120 779 881 (OEIS A000166)
func DerangementsBig(n int) (*big.Int, error) {
	if n < 0 {
		return nil, fmt.Errorf("%w: n", ErrNegative)
	}
	d := big.NewInt(1)
	var f big.Int
	one := big.NewInt(1)
	for i := 1; i <= n; i++ {
		d.Mul(d, f.SetInt64(int64(i)))
		if i%2 == 0 {
			d.Add(d, one)
		} else {
			d.Sub(d, one)
		}
	}
	return d, nil
}
