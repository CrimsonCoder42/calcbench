package trig

import (
	"fmt"
	"math"

	"github.com/CrimsonCoder42/calcbench/algo/reduce"
)

// Acsch returns the inverse hyperbolic cosecant of x, arcsinh(1/x), for x ≠ 0. It is odd and nonzero.
// For |x| > 1 it is arcsinh(1/|x|); for |x| ≤ 1 it is the logarithmic form ln(1 + √(1 + x²)) − ln|x|,
// which never forms 1/x and so cannot overflow for subnormal x. The sign of x is restored at the end.
// Errors: ErrNotFinite when x is NaN or ±Inf; ErrDomain when x is ±0 (the pole).
//
// ID:      NUM-TRIG-016
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  NIST Digital Library of Mathematical Functions, Release 1.2.8 (2026-09-15), eq. 4.37.7 and 4.37.13, https://dlmf.nist.gov/4.37 (read 2026-11-12)
// Source:  Abramowitz and Stegun (eds.), Handbook of Mathematical Functions, NBS Applied Mathematics Series 55, 10th printing, 1972, eq. 4.6.4, p. 86, and 4.6.23 (the logarithmic form, x ≠ 0), p. 87
// Tags:    trig, hyperbolic, inverse
// Units:   x dimensionless; result dimensionless
// Example: Acsch(1) = arcsinh 1 = 0.881373587019543 (hand-derived from 4.37.7)
func Acsch(x float64) (float64, error) {
	if !reduce.AllFinite(x) {
		return 0, fmt.Errorf("%w: x", ErrNotFinite)
	}
	if x == 0 {
		return 0, fmt.Errorf("%w: x is 0", ErrDomain)
	}
	t := math.Abs(x)
	var r float64
	if t > 1 {
		r = math.Asinh(1 / t)
	} else {
		r = math.Log(1+math.Sqrt(1+t*t)) - math.Log(t)
	}
	return math.Copysign(r, x), nil
}
