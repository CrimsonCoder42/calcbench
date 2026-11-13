package trig

import (
	"fmt"
	"math"

	"github.com/CrimsonCoder42/calcbench/algo/reduce"
)

// Csch returns the hyperbolic cosecant of x, 1/sinh x. It is odd and nonzero.
// Errors: ErrNotFinite when x is NaN or ±Inf; ErrDomain when x is ±0 (the pole); ErrOverflow when
// 0 < |x| < about 5.6e-309, where the result is beyond float64.
//
// ID:      NUM-TRIG-013
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  NIST Digital Library of Mathematical Functions, Release 1.2.8 (2026-09-15), eq. 4.28.5, https://dlmf.nist.gov/4.28 (read 2026-11-12)
// Source:  Abramowitz and Stegun (eds.), Handbook of Mathematical Functions, NBS Applied Mathematics Series 55, 10th printing, 1972, eq. 4.5.4, p. 83
// Tags:    trig, hyperbolic, cosecant
// Units:   x dimensionless; result dimensionless
// Example: Csch(1) = 0.8509181282393216; Csch(−1) = −0.8509181282393216 (hand-derived from 4.28.5)
func Csch(x float64) (float64, error) {
	if !reduce.AllFinite(x) {
		return 0, fmt.Errorf("%w: x", ErrNotFinite)
	}
	if x == 0 {
		return 0, fmt.Errorf("%w: x is 0", ErrDomain)
	}
	r := 1 / math.Sinh(x)
	if math.IsInf(r, 0) {
		return 0, fmt.Errorf("%w: x", ErrOverflow)
	}
	return r, nil
}
