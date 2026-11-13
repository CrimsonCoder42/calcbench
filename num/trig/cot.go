package trig

import (
	"fmt"
	"math"

	"github.com/CrimsonCoder42/calcbench/algo/reduce"
)

// Cot returns the cotangent of xRad, cos xRad / sin xRad, from one math.Sincos (not 1/math.Tan, which is
// huge near π/2). At math.Pi/2 the result is ≈ 6.123233995736766e-17, not 0, because math.Pi/2 is not π/2
// (Go's math.Sincos/Cos/Sin reduce the argument with a three-part π/4, so near a multiple of π/2 the small
// cos or sin, and this result, is accurate to about 2e-15 relative, not to the last ulp).
// Errors: ErrNotFinite when xRad is NaN or ±Inf; ErrDomain when xRad is ±0; ErrOverflow when
// 0 < |xRad| < about 5.6e-309, where the result is beyond float64.
//
// ID:      NUM-TRIG-008
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  NIST Digital Library of Mathematical Functions, Release 1.2.8 (2026-09-15), eq. 4.14.7, https://dlmf.nist.gov/4.14 (read 2026-11-12)
// Source:  Abramowitz and Stegun (eds.), Handbook of Mathematical Functions, NBS Applied Mathematics Series 55, 10th printing, 1972, eq. 4.3.6, p. 72
// Tags:    trig, cotangent, reciprocal
// Units:   xRad in radians; result dimensionless
// Example: Cot(math.Pi/4) = 1.0000000000000002 (hand-derived from 4.14.7)
func Cot(xRad float64) (float64, error) {
	if !reduce.AllFinite(xRad) {
		return 0, fmt.Errorf("%w: xRad", ErrNotFinite)
	}
	if xRad == 0 {
		return 0, fmt.Errorf("%w: xRad is 0", ErrDomain)
	}
	s, c := math.Sincos(xRad)
	r := c / s
	if math.IsInf(r, 0) {
		return 0, fmt.Errorf("%w: xRad", ErrOverflow)
	}
	return r, nil
}
