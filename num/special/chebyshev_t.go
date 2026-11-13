package special

import (
	"fmt"
	"math"
)

// ChebyshevT returns the Chebyshev polynomial of the first kind T_n(x), by T_(k+1) = 2x T_k − T_(k−1),
// T_0 = 1, T_1 = x (DLMF 18.9.1 with Table 18.9.1). On [−1, 1], T_n(cos θ) = cos nθ (DLMF 18.5.1).
//
// ID:      NUM-SPEC-010
// Time:    O(n)   n = degree
// Space:   O(1), no allocation
// Source:  NIST DLMF, Release 1.2.1, §18.9 (eq. 18.9.1, Table 18.9.1); §18.5 (eq. 18.5.1 trigonometric form, used by the tests)
// Source:  Abramowitz & Stegun, Handbook of Mathematical Functions, 1964, eq. 22.7.4
// Tags:    special-function, orthogonal-polynomial, chebyshev, recurrence
// Units:   dimensionless
// Range:   n ≥ 0; x finite. Errors: ErrOutOfRange (n < 0), ErrNotFinite, ErrOverflow.
// Example: T_3(x) = 4x³ − 3x, so T_3(0.5) = −1
func ChebyshevT(n int, x float64) (float64, error) {
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
		prev, cur = cur, 2*x*cur-prev
	}
	if math.IsInf(cur, 0) || math.IsNaN(cur) {
		return 0, fmt.Errorf("%w: ChebyshevT(%d, %v)", ErrOverflow, n, x)
	}
	return cur, nil
}
