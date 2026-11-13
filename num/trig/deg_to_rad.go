package trig

import (
	"errors"
	"fmt"
	"math"

	"github.com/CrimsonCoder42/calcbench/algo/reduce"
)

// Errors returned by this package's angle and reciprocal functions. They are wrapped with the offending
// parameter; test with errors.Is. ErrDomain, ErrNonPositive and ErrNotTriangle are declared beside the
// functions that first need them (csc.go, law_of_sines_side.go, law_of_cosines_angle.go).
var (
	// ErrNotFinite: an input is NaN or ±Inf.
	ErrNotFinite = errors.New("trig: value is NaN or infinite")
	// ErrOverflow: the true result of a finite input is beyond the float64 range.
	ErrOverflow = errors.New("trig: result exceeds the float64 range")
)

// DegToRad returns angleDeg, an angle in degrees, in radians: (angleDeg/180)·π.
// Dividing by 180 first keeps the multiples of 45° exact: 180 gives math.Pi and 90 gives math.Pi/2.
// The angle is not reduced; NormalizeAngle does that.
// Errors: ErrNotFinite when angleDeg is NaN or ±Inf.
//
// ID:      NUM-TRIG-001
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Thompson and Taylor, Guide for the Use of the International System of Units (SI), NIST Special Publication 811, 2008 ed., §5.1.1, Table 6 ("1° = (π/180) rad"), https://www.nist.gov/pml/special-publication-811 (read 2026-11-12)
// Tags:    trig, angle, unit, degree, radian
// Units:   angleDeg in degrees; result in radians
// Example: 180 → math.Pi; 90 → math.Pi/2; −45 → −math.Pi/4 (from Table 6's 1° = (π/180) rad)
func DegToRad(angleDeg float64) (float64, error) {
	if !reduce.AllFinite(angleDeg) {
		return 0, fmt.Errorf("%w: angleDeg", ErrNotFinite)
	}
	return (angleDeg / 180) * math.Pi, nil
}
