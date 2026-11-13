package trig

import (
	"fmt"
	"math"

	"github.com/CrimsonCoder42/calcbench/algo/reduce"
)

// Acoth returns the inverse hyperbolic cotangent of x, arctanh(1/x) = ½·ln((x+1)/(x−1)), for |x| > 1.
// It is odd and nonzero. It computes ½·log1p(2/(|x|−1)) with the sign of x: for |x| in (1, 2] the
// subtraction |x|−1 is exact, so the only roundings are the division and log1p, which keeps full accuracy
// as x → 1⁺, where arctanh(1/x) would amplify the rounding of 1/x.
// Errors: ErrNotFinite when x is NaN or ±Inf; ErrDomain when |x| ≤ 1 (±1 is a pole; inside is not real).
//
// ID:      NUM-TRIG-017
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  NIST Digital Library of Mathematical Functions, Release 1.2.8 (2026-09-15), eq. 4.37.9 and 4.37.15, https://dlmf.nist.gov/4.37 (read 2026-11-12)
// Source:  Abramowitz and Stegun (eds.), Handbook of Mathematical Functions, NBS Applied Mathematics Series 55, 10th printing, 1972, eq. 4.6.6, p. 86, and 4.6.25 (the logarithmic form, x² > 1), p. 87
// Tags:    trig, hyperbolic, inverse
// Units:   x dimensionless, |x| > 1; result dimensionless
// Example: Acoth(2) = ½ ln 3 = 0.5493061443340549 (A&S 4.6.25 at x = 2)
func Acoth(x float64) (float64, error) {
	if !reduce.AllFinite(x) {
		return 0, fmt.Errorf("%w: x", ErrNotFinite)
	}
	t := math.Abs(x)
	if t <= 1 {
		return 0, fmt.Errorf("%w: |x| ≤ 1", ErrDomain)
	}
	return math.Copysign(0.5*math.Log1p(2/(t-1)), x), nil
}
