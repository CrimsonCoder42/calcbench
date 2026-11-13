package special

import (
	"fmt"
	"math"
)

// HermiteH returns the physicists' Hermite polynomial H_n(x) (weight e^(−x²)), by
// H_(k+1) = 2x H_k − 2k H_(k−1), H_0 = 1, H_1 = 2x (DLMF 18.9.1 with Table 18.9.1).
//
// ID:      NUM-SPEC-012
// Time:    O(n)   n = degree
// Space:   O(1), no allocation
// Source:  NIST DLMF, Release 1.2.1, §18.9 (eq. 18.9.1, Table 18.9.1), §18.3 (Table 18.3.1)
// Source:  Abramowitz & Stegun, Handbook of Mathematical Functions, 1964, eq. 22.7.13; Table 22.12
// Tags:    special-function, orthogonal-polynomial, hermite, recurrence, physics
// Units:   dimensionless
// Range:   n ≥ 0; x finite. Errors: ErrOutOfRange (n < 0), ErrNotFinite, ErrOverflow.
// Example: H_3(x) = 8x³ − 12x, so H_3(1) = −4
func HermiteH(n int, x float64) (float64, error) {
	if n < 0 {
		return 0, fmt.Errorf("%w: n = %d, want ≥ 0", ErrOutOfRange, n)
	}
	if math.IsNaN(x) || math.IsInf(x, 0) {
		return 0, fmt.Errorf("%w: x = %v", ErrNotFinite, x)
	}
	prev, cur := 1.0, 2*x
	if n == 0 {
		return prev, nil
	}
	for k := 1; k < n; k++ {
		fk := float64(k)
		prev, cur = cur, 2*x*cur-2*fk*prev
	}
	if math.IsInf(cur, 0) || math.IsNaN(cur) {
		return 0, fmt.Errorf("%w: HermiteH(%d, %v)", ErrOverflow, n, x)
	}
	return cur, nil
}
