package trig

import (
	"fmt"
	"math"

	"github.com/CrimsonCoder42/calcbench/algo/reduce"
)

// Sech returns the hyperbolic secant of x, 1/cosh x, in (0, 1]. For |x| above about 710, cosh x is
// beyond float64 and the result is 0, the correctly rounded value there.
// Errors: ErrNotFinite when x is NaN or ±Inf.
//
// ID:      NUM-TRIG-012
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  NIST Digital Library of Mathematical Functions, Release 1.2.8 (2026-09-15), eq. 4.28.6, https://dlmf.nist.gov/4.28 (read 2026-11-12)
// Source:  Abramowitz and Stegun (eds.), Handbook of Mathematical Functions, NBS Applied Mathematics Series 55, 10th printing, 1972, eq. 4.5.5, p. 83
// Tags:    trig, hyperbolic, secant
// Units:   x dimensionless; result dimensionless, in (0, 1]
// Example: Sech(0) = 1; Sech(1) = 0.6480542736638855 (hand-derived from 4.28.6)
func Sech(x float64) (float64, error) {
	if !reduce.AllFinite(x) {
		return 0, fmt.Errorf("%w: x", ErrNotFinite)
	}
	return 1 / math.Cosh(x), nil
}
