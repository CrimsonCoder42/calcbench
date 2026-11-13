package trig

import "math"

// LawOfSinesAngle returns the angles B (opposite side b) of the triangles with sides a and b and the
// non-included angle A (opposite a): the solutions of sin B / b = sin A / a in (0, π) that close a
// triangle. count is 0 (no triangle: a < b·sin A, or a ≤ b with A ≥ π/2), 1 (a > b, or the two
// solutions coincide at B = π/2) or 2 (a < b and A < π/2: B and π − B). angleB1Rad holds the smaller
// (acute) B, angleB2Rad the other; a slot that holds no solution is 0. count 0 is not an error.
// B is computed by Kahan's procedure B(a, A, b), which picks, from the data, the least inaccurate of
// arcsin((b/a) sin A), A + 2 arctan((b−a)/((a/tan A) + √((a/tan A)² − (b−a)(b+a)))) and
// arccos(√((b cos A)² + (a−b)(a+b))/a); the bare arcsin can lose half the digits near B = π/2.
// Errors: ErrNotFinite when an input is NaN or ±Inf; ErrNonPositive when a side ≤ 0; ErrDomain when
// angleARad is outside (0, π).
//
// ID:      NUM-TRIG-019
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  NIST Digital Library of Mathematical Functions, Release 1.2.8 (2026-09-15), eq. 4.42.4, https://dlmf.nist.gov/4.42 (read 2026-11-12)
// Source:  Kahan, "Miscalculating Area and Angles of a Needle-like Triangle", lecture notes, UC Berkeley, 2014-09-04, §8 (formulas B_S, B_T, B_C and the selection B(a, A, b)), p. 6, and §9 Table 2, p. 8, https://people.eecs.berkeley.edu/~wkahan/Triangle.pdf (read 2026-11-12)
// Tags:    trig, triangle, law-of-sines, ssa, ambiguous-case
// Units:   sideA, sideB in one length unit; angleARad in radians; results in radians
// Example: a = 3, b = 17, A = 10° → count 2, B = 79.73894317° or 100.26105683° (Kahan 2014 Table 2, row 1)
func LawOfSinesAngle(sideA, sideB, angleARad float64) (angleB1Rad, angleB2Rad float64, count int, err error) {
	if err := checkLength(sideA, "sideA"); err != nil {
		return 0, 0, 0, err
	}
	if err := checkLength(sideB, "sideB"); err != nil {
		return 0, 0, 0, err
	}
	if err := checkAngle(angleARad, math.Pi, "angleARad"); err != nil {
		return 0, 0, 0, err
	}
	b1, b2, n := sinesAngle(sideA, sideB, angleARad)
	return b1, b2, n, nil
}

// sinesAngle is LawOfSinesAngle on validated inputs: Kahan 2014 §8, in radians. A NaN from arcsin(>1)
// or √(<0) is Kahan's own signal that no B exists.
func sinesAngle(a, b, angleA float64) (b1, b2 float64, count int) {
	const right = math.Pi / 2
	// B depends only on b/a; an exact power-of-two scaling keeps the squares below from overflowing.
	if k := sideShift(max(a, b)); k != 0 {
		a, b = math.Ldexp(a, k), math.Ldexp(b, k)
	}
	if a < b {
		// B2(A, X): no B at all when A is not acute; otherwise the pair {X, π − X}.
		if angleA >= right {
			return 0, 0, 0
		}
		var x float64
		if b-a < a {
			x = kahanBT(a, b, angleA)
		} else {
			x = kahanBS(a, b, angleA)
		}
		if math.IsNaN(x) {
			return 0, 0, 0
		}
		// The second triangle exists when its third angle, B1 − A, is positive.
		if x <= angleA || x >= right {
			return x, 0, 1
		}
		return x, math.Pi - x, 2
	}
	// a ≥ b: at most one B, and it is acute. a = b closes a triangle only when A is acute (B = A).
	if a == b && angleA >= right {
		return 0, 0, 0
	}
	ratio := b / a
	var x float64
	switch {
	case angleA >= right && ratio*math.Sin(angleA) >= 0.756:
		x = kahanBC(a, b, angleA)
	case angleA < right && ratio >= 0.756:
		x = kahanBT(a, b, angleA)
	default:
		x = kahanBS(a, b, angleA)
	}
	if math.IsNaN(x) {
		return 0, 0, 0
	}
	return x, 0, 1
}

// kahanBS is B_S(a, A, b) = arcsin((b/a) sin A).
func kahanBS(a, b, angleA float64) float64 {
	return math.Asin((b / a) * math.Sin(angleA))
}

// kahanBT is B_T(a, A, b) = A + 2 arctan((b−a)/((a/tan A) + √((a/tan A)² − (b−a)(b+a)))), the parentheses
// kept as Kahan writes them. When q = a/tan A is beyond 2^500 (A tiny), q² would overflow, so the root is
// taken in the factored forms hypot(q, √((a−b)(a+b))) for b ≤ a and √(q−r)·√(q+r), r = √((b−a)(b+a)), for
// b > a, which equal it exactly in real arithmetic.
func kahanBT(a, b, angleA float64) float64 {
	q := a / math.Tan(angleA)
	var root float64
	switch d := (b - a) * (b + a); {
	case q <= 0x1p500:
		root = math.Sqrt(q*q - d)
	case d <= 0:
		root = math.Hypot(q, math.Sqrt(-d))
	default:
		r := math.Sqrt(d)
		root = math.Sqrt(q-r) * math.Sqrt(q+r)
	}
	return angleA + 2*math.Atan((b-a)/(q+root))
}

// kahanBC is B_C(a, A, b) = arccos(√((b cos A)² + (a−b)(a+b))/a).
func kahanBC(a, b, angleA float64) float64 {
	bc := b * math.Cos(angleA)
	return math.Acos(math.Sqrt(bc*bc+(a-b)*(a+b)) / a)
}
