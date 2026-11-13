package special

import (
	"fmt"
	"math"
)

// ChebyshevU returns the Chebyshev polynomial of the second kind U_n(x), by U_(k+1) = 2x U_k − U_(k−1),
// U_0 = 1, U_1 = 2x (DLMF 18.9.1 with Table 18.9.1). On (−1, 1), U_n(cos θ) = sin((n+1)θ)/sin θ (DLMF 18.5.2).
//
// ID:      NUM-SPEC-011
// Time:    O(n)   n = degree
// Space:   O(1), no allocation
// Source:  NIST DLMF, Release 1.2.1, §18.9 (eq. 18.9.1, Table 18.9.1); §18.5 (eq. 18.5.2)
// Source:  Abramowitz & Stegun, Handbook of Mathematical Functions, 1964, eq. 22.7.5
// Tags:    special-function, orthogonal-polynomial, chebyshev, recurrence
// Units:   dimensionless
// Range:   n ≥ 0; x finite. Errors: ErrOutOfRange (n < 0), ErrNotFinite, ErrOverflow.
// Example: U_2(x) = 4x² − 1, so U_2(0.5) = 0
func ChebyshevU(n int, x float64) (float64, error) {
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
		prev, cur = cur, 2*x*cur-prev
	}
	if math.IsInf(cur, 0) || math.IsNaN(cur) {
		return 0, fmt.Errorf("%w: ChebyshevU(%d, %v)", ErrOverflow, n, x)
	}
	return cur, nil
}
