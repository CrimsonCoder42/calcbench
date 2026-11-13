package trig

import (
	"fmt"
	"math"

	"github.com/CrimsonCoder42/calcbench/algo/reduce"
)

// Asec returns the principal inverse secant of x, arccos(1/x), in radians in [0, π]: in [0, π/2) for
// x ≥ 1 and in (π/2, π] for x ≤ −1.
// Errors: ErrNotFinite when x is NaN or ±Inf; ErrDomain when |x| < 1 (no real value).
//
// ID:      NUM-TRIG-009
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  NIST Digital Library of Mathematical Functions, Release 1.2.8 (2026-09-15), eq. 4.23.8 and Table 4.23.1, https://dlmf.nist.gov/4.23 (read 2026-11-12)
// Source:  Abramowitz and Stegun (eds.), Handbook of Mathematical Functions, NBS Applied Mathematics Series 55, 10th printing, 1972, eq. 4.4.7 and 4.4.4 (0 ≤ arccos x ≤ π), p. 79
// Tags:    trig, inverse, arcsecant
// Units:   x dimensionless, |x| ≥ 1; result in radians, in [0, π]
// Example: Asec(1) = 0; Asec(−1) = π (DLMF Table 4.23.1); Asec(2) = π/3 (hand-derived)
func Asec(x float64) (float64, error) {
	if !reduce.AllFinite(x) {
		return 0, fmt.Errorf("%w: x", ErrNotFinite)
	}
	if math.Abs(x) < 1 {
		return 0, fmt.Errorf("%w: |x| < 1", ErrDomain)
	}
	return math.Acos(1 / x), nil
}
