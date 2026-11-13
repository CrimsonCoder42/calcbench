package trig

import (
	"errors"
	"fmt"
	"math"
)

// ErrNotTriangle: the lengths cannot be the sides of a triangle (one exceeds the sum of the other two),
// or a right triangle's leg is not shorter than its hypotenuse.
var ErrNotTriangle = errors.New("trig: lengths do not form a triangle")

// LawOfCosinesAngle returns the angle C opposite side c of the triangle with sides a, b and c, in [0, π].
// Mathematically cos C = (a² + b² − c²)/(2ab) (the law of cosines), but arccos of that quotient loses
// digits for needle-like triangles. Kahan's formula is used instead: with a ≥ b (swapped if needed),
// μ = c − (a−b) if b ≥ c, else b − (a−c), and C = 2 arctan(√(((a−b)+c)·μ / ((a+(b+c))·((a−c)+b)))),
// with the parentheses kept as written. Sides beyond 2^±500 are first scaled by an exact power of two,
// which leaves C unchanged and keeps the sums from overflowing. A flat triangle (one side equal to the
// sum of the others) is allowed and gives 0 or π.
// Errors: ErrNotFinite when an input is NaN or ±Inf; ErrNonPositive when a side ≤ 0; ErrNotTriangle
// when one side exceeds the sum of the other two.
//
// ID:      NUM-TRIG-021
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  NIST Digital Library of Mathematical Functions, Release 1.2.8 (2026-09-15), eq. 4.42.5, https://dlmf.nist.gov/4.42 (read 2026-11-12)
// Source:  Kahan, "Miscalculating Area and Angles of a Needle-like Triangle", lecture notes, UC Berkeley, 2014-09-04, §3 (the formula for C with μ), p. 3, and §5 Table 1, p. 4, https://people.eecs.berkeley.edu/~wkahan/Triangle.pdf (read 2026-11-12)
// Tags:    trig, triangle, law-of-cosines, needle
// Units:   sideA, sideB, sideC in one length unit; result in radians, in [0, π]
// Example: (10, 10, 10) → 60°; (100000, 99999.99979, 0.00029) → 1.145915591e-7°; (10000, 5000.000001, 15000) → 179.9985965° (Kahan 2014 Table 1)
func LawOfCosinesAngle(sideA, sideB, sideC float64) (float64, error) {
	if err := checkSides(sideA, sideB, sideC); err != nil {
		return 0, err
	}
	angle, ok := cosinesAngle(sideA, sideB, sideC, false)
	if !ok {
		return 0, fmt.Errorf("%w: sideA, sideB, sideC", ErrNotTriangle)
	}
	return angle, nil
}

// checkSides validates three lengths.
func checkSides(sideA, sideB, sideC float64) error {
	if err := checkLength(sideA, "sideA"); err != nil {
		return err
	}
	if err := checkLength(sideB, "sideB"); err != nil {
		return err
	}
	return checkLength(sideC, "sideC")
}

// sideShift returns the power-of-two exponent shift that brings the longest side m to [0.5, 1) when m is
// beyond 2^±500, and 0 otherwise. Scaling a side by math.Ldexp(side, shift) is exact (no side becomes
// subnormal unless it is 2^1000 times shorter than m), angles are unchanged by it, and no product of three
// scaled sides can overflow or underflow.
func sideShift(m float64) int {
	if m > 0x1p500 || m < 0x1p-500 {
		_, e := math.Frexp(m)
		return -e
	}
	return 0
}

// cosinesAngle is Kahan 2014 §3 on positive, finite sides: the angle opposite c. ok is false when the
// sides fail the triangle inequality, unless clamp is set, in which case a factor that rounded below zero
// is taken as zero (for sides that form a triangle by construction).
func cosinesAngle(a, b, c float64, clamp bool) (angle float64, ok bool) {
	if a < b {
		a, b = b, a
	}
	if k := sideShift(max(a, c)); k != 0 {
		a, b, c = math.Ldexp(a, k), math.Ldexp(b, k), math.Ldexp(c, k)
	}
	var mu float64
	if b >= c {
		mu = c - (a - b)
	} else {
		mu = b - (a - c)
	}
	den := (a - c) + b
	if mu < 0 || den < 0 {
		if !clamp {
			return 0, false
		}
		mu, den = max(mu, 0), max(den, 0)
	}
	if mu == 0 {
		return 0, true // flat: c = a − b (or b = a − c)
	}
	// ((a−b)+c)/(a+(b+c)) is in (0, 1]; μ/den is +Inf only for the flat c = a + b, where C = π.
	t := (((a - b) + c) / (a + (b + c))) * (mu / den)
	return 2 * math.Atan(math.Sqrt(t)), true
}
