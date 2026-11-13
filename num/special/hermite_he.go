package special

import (
	"fmt"
	"math"
)

// HermiteHe returns the probabilists' Hermite polynomial He_n(x) (weight e^(−x²/2)), by
// He_(k+1) = x He_k − k He_(k−1), He_0 = 1, He_1 = x (DLMF 18.9.1 with Table 18.9.1).
//
// ID:      NUM-SPEC-013
// Time:    O(n)   n = degree
// Space:   O(1), no allocation
// Source:  NIST DLMF, Release 1.2.1, §18.9 (eq. 18.9.1, Table 18.9.1), §18.3 (Table 18.3.1)
// Source:  Abramowitz & Stegun, Handbook of Mathematical Functions, 1964, eq. 22.7.14
// Tags:    special-function, orthogonal-polynomial, hermite, recurrence, probability
// Units:   dimensionless
// Range:   n ≥ 0; x finite. Errors: ErrOutOfRange (n < 0), ErrNotFinite, ErrOverflow.
// Example: He_3(x) = x³ − 3x, so He_3(2) = 2
func HermiteHe(n int, x float64) (float64, error) {
	if n < 0 {
		return 0, fmt.Errorf("%w: n = %d, want ≥ 0", ErrOutOfRange, n)
	}
	if math.IsNaN(x) || math.IsInf(x, 0) {
		return 0, fmt.Errorf("%w: x = %v", ErrNotFinite, x)
	}
	prev, cur := 1.0, x
	if n == 0 {
		return prev, nil
	}
	for k := 1; k < n; k++ {
		fk := float64(k)
		prev, cur = cur, x*cur-fk*prev
	}
	if math.IsInf(cur, 0) || math.IsNaN(cur) {
		return 0, fmt.Errorf("%w: HermiteHe(%d, %v)", ErrOverflow, n, x)
	}
	return cur, nil
}
