package linalg

import (
	"fmt"
	"math"

	"github.com/CrimsonCoder42/calcbench/algo/reduce"
)

// VecAngleBetween returns the unoriented angle between a and b, in [0, π] radians, in any dimension, by
// Kahan's formula θ = 2·atan2(‖â − b̂‖, ‖â + b̂‖) with â = a/‖a‖, b̂ = b/‖b‖. Its absolute error stays
// near one ulp of π for nearly parallel, perpendicular and antiparallel vectors alike.
// Errors: ErrDimension when len(a) ≠ len(b); ErrZeroVector when a or b is all zeros (or empty);
// ErrNotFinite for a NaN or ±Inf element.
//
// ID:      NUM-LINALG-047
// Time:    O(n)   n = len(a); five passes
// Space:   O(1), no allocation
// Source:  Kahan, "How Futile are Mindless Assessments of Roundoff in Floating-Point Computation?", University of California, Berkeley, 11 January 2006, §12 "Mangled Angles", pp. 46–47, https://people.eecs.berkeley.edu/~wkahan/Mindless.pdf (read 2026-11-12)
// Tags:    vector, n-dimensional, angle, angle-between, linear-algebra
// Units:   a and b in any units; the result in radians
// Example: (1, 0, 0, 0) and (1, 1, 1, 1) → π/3 (hand-checked: cos θ = 1/2)
func VecAngleBetween(a, b []float64) (float64, error) {
	if len(a) != len(b) {
		return 0, lengthError("a", len(a), "b", len(b))
	}
	ma, bad := maxAbs(a)
	if bad >= 0 {
		return 0, notFiniteAt("a", bad)
	}
	mb, bad := maxAbs(b)
	if bad >= 0 {
		return 0, notFiniteAt("b", bad)
	}
	if ma == 0 {
		return 0, fmt.Errorf("%w: a", ErrZeroVector)
	}
	if mb == 0 {
		return 0, fmt.Errorf("%w: b", ErrZeroVector)
	}
	na, ea := scaledNorm(a, ma)
	nb, eb := scaledNorm(b, mb)
	var diff, sum dot2
	for i := range a {
		// Scale by the power of two first, then divide: 2^(−e)/n alone can overflow for subnormal input.
		x, y := math.Ldexp(a[i], -ea)/na, math.Ldexp(b[i], -eb)/nb
		d, s := x-y, x+y
		diff.add(d, d)
		sum.add(s, s)
	}
	return 2 * math.Atan2(math.Sqrt(diff.result()), math.Sqrt(sum.result())), nil
}

// VecProject returns the projection of v onto the line of onto, (v·n)n with n = onto/‖onto‖₂, written
// into dst (reused when its capacity suffices; dst may be v or onto). The dot product is taken by Dot2.
// Errors: ErrDimension when len(v) ≠ len(onto); ErrZeroVector when onto is all zeros (or empty);
// ErrNotFinite for a NaN or ±Inf element; ErrOverflow when an element of the result is outside the
// float64 range.
//
// ID:      NUM-LINALG-048
// Time:    O(n)   n = len(v); four passes
// Space:   O(1) extra; O(n) for the result, one allocation only if dst lacks capacity
// Source:  Strang, Introduction to Linear Algebra, 5th ed., Wellesley-Cambridge Press, 2016, §4.2 "Projections" (projection onto a line: p = a(aᵀb)/(aᵀa))
// Tags:    vector, n-dimensional, projection, vector-projection, linear-algebra
// Units:   v in any unit, onto in any unit; the result in v's unit
// Example: v = (1, 2, 3, 4), onto = (0, 0, 0, 2) → (0, 0, 0, 4) (hand-checked)
func VecProject(dst, v, onto []float64) ([]float64, error) {
	if len(v) != len(onto) {
		return nil, lengthError("v", len(v), "onto", len(onto))
	}
	if _, bad := maxAbs(v); bad >= 0 {
		return nil, notFiniteAt("v", bad)
	}
	mo, bad := maxAbs(onto)
	if bad >= 0 {
		return nil, notFiniteAt("onto", bad)
	}
	if mo == 0 {
		return nil, fmt.Errorf("%w: onto", ErrZeroVector)
	}
	no, eo := scaledNorm(onto, mo)
	unit := func(i int) float64 { return math.Ldexp(onto[i], -eo) / no }
	var d dot2
	for i := range v {
		d.add(v[i], unit(i))
	}
	c := d.result()
	if !reduce.AllFinite(c) {
		return nil, errOverflow("v·n")
	}
	dst = grow(dst, len(v))
	for i := range onto {
		y := c * unit(i)
		if !reduce.AllFinite(y) {
			return nil, fmt.Errorf("%w: element %d", ErrOverflow, i)
		}
		dst[i] = y
	}
	return dst, nil
}
