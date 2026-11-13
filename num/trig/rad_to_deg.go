package trig

import (
	"fmt"
	"math"

	"github.com/CrimsonCoder42/calcbench/algo/reduce"
)

// RadToDeg returns angleRad, an angle in radians, in degrees: (angleRad/π)·180.
// Dividing by π first makes math.Pi give exactly 180 and math.Pi/2 exactly 90.
// Errors: ErrNotFinite when angleRad is NaN or ±Inf; ErrOverflow when |angleRad| is so large (above about
// 3.1e306) that the result is beyond float64.
//
// ID:      NUM-TRIG-002
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Thompson and Taylor, Guide for the Use of the International System of Units (SI), NIST Special Publication 811, 2008 ed., §5.1.1, Table 6 ("1° = (π/180) rad"), https://www.nist.gov/pml/special-publication-811 (read 2026-11-12)
// Tags:    trig, angle, unit, degree, radian
// Units:   angleRad in radians; result in degrees
// Example: math.Pi → 180; math.Pi/2 → 90 (the inverse of Table 6's 1° = (π/180) rad)
func RadToDeg(angleRad float64) (float64, error) {
	if !reduce.AllFinite(angleRad) {
		return 0, fmt.Errorf("%w: angleRad", ErrNotFinite)
	}
	deg := (angleRad / math.Pi) * 180
	if math.IsInf(deg, 0) {
		return 0, fmt.Errorf("%w: angleRad", ErrOverflow)
	}
	return deg, nil
}
