package trig

import (
	"errors"
	"fmt"
	"math"

	"github.com/CrimsonCoder42/calcbench/algo/reduce"
)

// ErrNonPositive: a length is zero or negative.
var ErrNonPositive = errors.New("trig: length must be positive")

// checkLength returns ErrNotFinite or ErrNonPositive, wrapped with name, unless v is a finite length > 0.
func checkLength(v float64, name string) error {
	if !reduce.AllFinite(v) {
		return fmt.Errorf("%w: %s", ErrNotFinite, name)
	}
	if v <= 0 {
		return fmt.Errorf("%w: %s", ErrNonPositive, name)
	}
	return nil
}

// checkAngle returns ErrNotFinite or ErrDomain, wrapped with name, unless v is a finite angle in
// (0, maxRad). maxRad is math.Pi for a triangle's angle and math.Pi/2 for a right triangle's acute angle.
func checkAngle(v, maxRad float64, name string) error {
	if !reduce.AllFinite(v) {
		return fmt.Errorf("%w: %s", ErrNotFinite, name)
	}
	if v <= 0 || v >= maxRad {
		return fmt.Errorf("%w: %s outside (0, %v)", ErrDomain, name, maxRad)
	}
	return nil
}

// checkTwoAngles validates two triangle angles: each in (0, π) and their sum below π.
func checkTwoAngles(angleARad, angleBRad float64) error {
	if err := checkAngle(angleARad, math.Pi, "angleARad"); err != nil {
		return err
	}
	if err := checkAngle(angleBRad, math.Pi, "angleBRad"); err != nil {
		return err
	}
	if angleARad+angleBRad >= math.Pi {
		return fmt.Errorf("%w: angleARad + angleBRad ≥ π", ErrDomain)
	}
	return nil
}

// LawOfSinesSide returns side b of a triangle from side a and the angles A (opposite a) and B
// (opposite b), by the law of sines a/sin A = b/sin B: b = a·(sin B / sin A).
// Errors: ErrNotFinite when an input is NaN or ±Inf; ErrNonPositive when sideA ≤ 0; ErrDomain when an
// angle is outside (0, π) or angleARad + angleBRad ≥ π (no triangle); ErrOverflow when b is beyond float64.
//
// ID:      NUM-TRIG-018
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  NIST Digital Library of Mathematical Functions, Release 1.2.8 (2026-09-15), eq. 4.42.4, https://dlmf.nist.gov/4.42 (read 2026-11-12)
// Source:  Abramowitz and Stegun (eds.), Handbook of Mathematical Functions, NBS Applied Mathematics Series 55, 10th printing, 1972, eq. 4.3.148, p. 79
// Tags:    trig, triangle, law-of-sines
// Units:   sideA in any length unit; angleARad, angleBRad in radians; result in the unit of sideA
// Example: a = 10, A = B = 60° → b = 10; a = 3, A = atan2(3, 4), B = π/2 → b = 5 (hand-derived from 4.42.4)
func LawOfSinesSide(sideA, angleARad, angleBRad float64) (float64, error) {
	if err := checkLength(sideA, "sideA"); err != nil {
		return 0, err
	}
	if err := checkTwoAngles(angleARad, angleBRad); err != nil {
		return 0, err
	}
	b := sideA * (math.Sin(angleBRad) / math.Sin(angleARad))
	if math.IsInf(b, 0) {
		return 0, fmt.Errorf("%w: sideB", ErrOverflow)
	}
	return b, nil
}
