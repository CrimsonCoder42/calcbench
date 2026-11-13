package special

import (
	"fmt"
	"math"
)

// LaguerreL returns the Laguerre polynomial L_n(x) (weight e^(−x) on [0, ∞)), by
// (k+1)L_(k+1) = (2k+1−x)L_k − k L_(k−1), L_0 = 1, L_1 = 1 − x (DLMF 18.9.1 with Table 18.9.1).
//
// ID:      NUM-SPEC-014
// Time:    O(n)   n = degree
// Space:   O(1), no allocation
// Source:  NIST DLMF, Release 1.2.1, §18.9 (eq. 18.9.1, Table 18.9.1 with α = 0), §18.3 (Table 18.3.1)
// Source:  Abramowitz & Stegun, Handbook of Mathematical Functions, 1964, eq. 22.7.12; Table 22.10
// Tags:    special-function, orthogonal-polynomial, laguerre, recurrence
// Units:   dimensionless
// Range:   n ≥ 0; x finite. Errors: ErrOutOfRange (n < 0), ErrNotFinite, ErrOverflow.
// Example: L_2(x) = (x² − 4x + 2)/2, so L_2(1) = −0.5
func LaguerreL(n int, x float64) (float64, error) {
	if n < 0 {
		return 0, fmt.Errorf("%w: n = %d, want ≥ 0", ErrOutOfRange, n)
	}
	if math.IsNaN(x) || math.IsInf(x, 0) {
		return 0, fmt.Errorf("%w: x = %v", ErrNotFinite, x)
	}
	prev, cur := 1.0, 1-x
	if n == 0 {
		return prev, nil
	}
	for k := 1; k < n; k++ {
		fk := float64(k)
		prev, cur = cur, ((2*fk+1-x)*cur-fk*prev)/(fk+1)
	}
	if math.IsInf(cur, 0) || math.IsNaN(cur) {
		return 0, fmt.Errorf("%w: LaguerreL(%d, %v)", ErrOverflow, n, x)
	}
	return cur, nil
}
