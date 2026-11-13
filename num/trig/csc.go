package trig

import (
	"errors"
	"fmt"
	"math"

	"github.com/CrimsonCoder42/calcbench/algo/reduce"
)

// ErrDomain: an input is outside the domain of the function (a pole, or a value with no real result).
var ErrDomain = errors.New("trig: input outside the function's domain")

// Csc returns the cosecant of xRad, 1/sin xRad.
// ±0 is the only float64 at which sin is exactly 0; elsewhere the result is finite (at math.Pi it is
// ≈ 8.165619676597685e15, the true cosecant of that float64; Go's
// math.Sincos/Cos/Sin reduce the argument with a three-part π/4, so near a multiple of π/2 the small
// cos or sin, and this result, is accurate to about 2e-15 relative, not to the last ulp).
// Errors: ErrNotFinite when xRad is NaN or ±Inf; ErrDomain when xRad is ±0; ErrOverflow when
// 0 < |xRad| < about 5.6e-309, where 1/sin xRad is beyond float64.
//
// ID:      NUM-TRIG-007
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  NIST Digital Library of Mathematical Functions, Release 1.2.8 (2026-09-15), eq. 4.14.5, https://dlmf.nist.gov/4.14 (read 2026-11-12)
// Source:  Abramowitz and Stegun (eds.), Handbook of Mathematical Functions, NBS Applied Mathematics Series 55, 10th printing, 1972, eq. 4.3.4, p. 72
// Tags:    trig, cosecant, reciprocal
// Units:   xRad in radians; result dimensionless, |result| ≥ 1
// Example: Csc(math.Pi/2) = 1; Csc(math.Pi/6) = 2.0000000000000004 (hand-derived from 4.14.5)
func Csc(xRad float64) (float64, error) {
	if !reduce.AllFinite(xRad) {
		return 0, fmt.Errorf("%w: xRad", ErrNotFinite)
	}
	if xRad == 0 {
		return 0, fmt.Errorf("%w: xRad is 0", ErrDomain)
	}
	r := 1 / math.Sin(xRad)
	if math.IsInf(r, 0) {
		return 0, fmt.Errorf("%w: xRad", ErrOverflow)
	}
	return r, nil
}
