package trig

import (
	"fmt"
	"math"

	"github.com/CrimsonCoder42/calcbench/algo/reduce"
)

// Sec returns the secant of xRad, 1/cos xRad.
// No float64 makes cos exactly 0 (π/2 is not representable), so there is no pole to report: at
// math.Pi/2 the result is the large finite ≈ 1.633123935319537e16, the true secant of that float64 (Go's
// math.Sincos/Cos/Sin reduce the argument with a three-part π/4, so near a multiple of π/2 the small
// cos or sin, and this result, is accurate to about 2e-15 relative, not to the last ulp).
// Errors: ErrNotFinite when xRad is NaN or ±Inf.
//
// ID:      NUM-TRIG-006
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  NIST Digital Library of Mathematical Functions, Release 1.2.8 (2026-09-15), eq. 4.14.6, https://dlmf.nist.gov/4.14 (read 2026-11-12)
// Source:  Abramowitz and Stegun (eds.), Handbook of Mathematical Functions, NBS Applied Mathematics Series 55, 10th printing, 1972, eq. 4.3.5, p. 72
// Tags:    trig, secant, reciprocal
// Units:   xRad in radians; result dimensionless, |result| ≥ 1
// Example: Sec(0) = 1; Sec(math.Pi) = −1 (hand-derived from 4.14.6)
func Sec(xRad float64) (float64, error) {
	if !reduce.AllFinite(xRad) {
		return 0, fmt.Errorf("%w: xRad", ErrNotFinite)
	}
	return 1 / math.Cos(xRad), nil
}
