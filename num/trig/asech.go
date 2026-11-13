package trig

import (
	"fmt"
	"math"

	"github.com/CrimsonCoder42/calcbench/algo/reduce"
)

// Asech returns the principal inverse hyperbolic secant of x, arccosh(1/x), for x in (0, 1]; the result
// is ≥ 0. It evaluates the logarithmic form ln[1/x + (1/x² − 1)^½] without forming 1/x:
// log1p(((1−x) + √((1−x)(1+x)))/x) for x > ½, where 1−x is exact, and ln(1 + √((1−x)(1+x))) − ln x
// otherwise, which cannot overflow even for subnormal x. Near x = 1 the function itself is
// ill-conditioned (it behaves like √(2(1−x))), as arccosh is near 1; that is inherent.
// Errors: ErrNotFinite when x is NaN or ±Inf; ErrDomain when x ≤ 0 or x > 1 (no real principal value).
//
// ID:      NUM-TRIG-015
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  NIST Digital Library of Mathematical Functions, Release 1.2.8 (2026-09-15), eq. 4.37.8, https://dlmf.nist.gov/4.37 (read 2026-11-12)
// Source:  Abramowitz and Stegun (eds.), Handbook of Mathematical Functions, NBS Applied Mathematics Series 55, 10th printing, 1972, eq. 4.6.5, p. 86, and 4.6.24 (the logarithmic form, 0 < x ≤ 1), p. 87
// Tags:    trig, hyperbolic, inverse
// Units:   x dimensionless, in (0, 1]; result dimensionless, ≥ 0
// Example: Asech(1) = 0; Asech(0.5) = arccosh 2 = 1.3169578969248166 (hand-derived from 4.37.8)
func Asech(x float64) (float64, error) {
	if !reduce.AllFinite(x) {
		return 0, fmt.Errorf("%w: x", ErrNotFinite)
	}
	if x <= 0 || x > 1 {
		return 0, fmt.Errorf("%w: x outside (0, 1]", ErrDomain)
	}
	root := math.Sqrt((1 - x) * (1 + x))
	if x > 0.5 {
		return math.Log1p(((1 - x) + root) / x), nil
	}
	return math.Log(1+root) - math.Log(x), nil
}
