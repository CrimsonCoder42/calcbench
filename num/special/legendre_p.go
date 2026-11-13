package special

import (
	"fmt"
	"math"
)

// LegendreP returns the Legendre polynomial P_n(x), by the three-term recurrence
// (k+1)P_(k+1) = (2k+1)x P_k − k P_(k−1), P_0 = 1, P_1 = x (DLMF 18.9.1 with Table 18.9.1). The
// recurrence is forward-stable on and off [−1, 1].
//
// ID:      NUM-SPEC-009
// Time:    O(n)   n = degree
// Space:   O(1), no allocation
// Source:  NIST DLMF, Release 1.2.1, §18.9 (eq. 18.9.1, Table 18.9.1), §18.3 (Table 18.3.1 normalization)
// Source:  Abramowitz & Stegun, Handbook of Mathematical Functions, 1964, eq. 22.7.10; Table 22.11
// Tags:    special-function, orthogonal-polynomial, legendre, recurrence
// Units:   dimensionless
// Range:   n ≥ 0; x finite. Errors: ErrOutOfRange (n < 0), ErrNotFinite, ErrOverflow (|P_n(x)| > MaxFloat64).
// Example: P_2(x) = (3x² − 1)/2, so P_2(0.5) = −0.125
func LegendreP(n int, x float64) (float64, error) {
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
		prev, cur = cur, ((2*fk+1)*x*cur-fk*prev)/(fk+1)
	}
	if math.IsInf(cur, 0) || math.IsNaN(cur) {
		return 0, fmt.Errorf("%w: LegendreP(%d, %v)", ErrOverflow, n, x)
	}
	return cur, nil
}
