package trig

import (
	"fmt"
	"math"

	"github.com/CrimsonCoder42/calcbench/algo/reduce"
)

// Acsc returns the principal inverse cosecant of x, arcsin(1/x), in radians in [−π/2, π/2] excluding 0.
// It is odd: Acsc(−x) = −Acsc(x).
// Errors: ErrNotFinite when x is NaN or ±Inf; ErrDomain when |x| < 1 (no real value).
//
// ID:      NUM-TRIG-010
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  NIST Digital Library of Mathematical Functions, Release 1.2.8 (2026-09-15), eq. 4.23.7, 4.23.13 and Table 4.23.1, https://dlmf.nist.gov/4.23 (read 2026-11-12)
// Source:  Abramowitz and Stegun (eds.), Handbook of Mathematical Functions, NBS Applied Mathematics Series 55, 10th printing, 1972, eq. 4.4.6 and 4.4.4 (−π/2 ≤ arcsin x ≤ π/2), p. 79
// Tags:    trig, inverse, arccosecant
// Units:   x dimensionless, |x| ≥ 1; result in radians, in [−π/2, π/2]
// Example: Acsc(1) = π/2; Acsc(−1) = −π/2 (DLMF Table 4.23.1); Acsc(2) = π/6 (hand-derived)
func Acsc(x float64) (float64, error) {
	if !reduce.AllFinite(x) {
		return 0, fmt.Errorf("%w: x", ErrNotFinite)
	}
	if math.Abs(x) < 1 {
		return 0, fmt.Errorf("%w: |x| < 1", ErrDomain)
	}
	return math.Asin(1 / x), nil
}
