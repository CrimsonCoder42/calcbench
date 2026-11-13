package linalg

import (
	"fmt"
	"math"

	"github.com/CrimsonCoder42/calcbench/algo/reduce"
)

// Vec3 is a vector in space, [x, y, z]. It is a value: functions take and return copies and never
// allocate.
type Vec3 [3]float64

func finite3(v Vec3) bool { return reduce.AllFinite(v[:]...) }

// Vec3Add returns a + b, component by component.
// Errors: ErrNotFinite for a NaN or ±Inf component; ErrOverflow when a component of the sum is
// outside the float64 range. The result is the zero vector with any error.
//
// ID:      NUM-LINALG-011
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Strang, Introduction to Linear Algebra, 5th ed., Wellesley-Cambridge Press, 2016, §1.1 "Vectors and Linear Combinations"
// Tags:    vector, 3d, add, linear-algebra
// Units:   a and b in one unit; the result in that unit
// Example: (1, 2, 3) + (4, −5, 6) → (5, −3, 9) (hand-checked)
func Vec3Add(a, b Vec3) (Vec3, error) {
	if !finite3(a) {
		return Vec3{}, errNotFinite("a")
	}
	if !finite3(b) {
		return Vec3{}, errNotFinite("b")
	}
	r := Vec3{a[0] + b[0], a[1] + b[1], a[2] + b[2]}
	if !finite3(r) {
		return Vec3{}, errOverflow("a + b")
	}
	return r, nil
}

// Vec3Sub returns a − b, component by component.
// Errors: as Vec3Add.
//
// ID:      NUM-LINALG-012
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Strang, Introduction to Linear Algebra, 5th ed., Wellesley-Cambridge Press, 2016, §1.1 "Vectors and Linear Combinations"
// Tags:    vector, 3d, subtract, difference, linear-algebra
// Units:   a and b in one unit; the result in that unit
// Example: (1, 2, 3) − (4, −5, 6) → (−3, 7, −3) (hand-checked)
func Vec3Sub(a, b Vec3) (Vec3, error) {
	if !finite3(a) {
		return Vec3{}, errNotFinite("a")
	}
	if !finite3(b) {
		return Vec3{}, errNotFinite("b")
	}
	r := Vec3{a[0] - b[0], a[1] - b[1], a[2] - b[2]}
	if !finite3(r) {
		return Vec3{}, errOverflow("a − b")
	}
	return r, nil
}

// Vec3Scale returns k·v.
// Errors: ErrNotFinite for a NaN or ±Inf input; ErrOverflow when a component of the product is outside
// the float64 range.
//
// ID:      NUM-LINALG-013
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Strang, Introduction to Linear Algebra, 5th ed., Wellesley-Cambridge Press, 2016, §1.1 "Vectors and Linear Combinations"
// Tags:    vector, 3d, scale, scalar-multiplication, linear-algebra
// Units:   v in any unit, k dimensionless; the result in v's unit
// Example: 2·(1, −2, 3) → (2, −4, 6) (hand-checked)
func Vec3Scale(v Vec3, k float64) (Vec3, error) {
	if !finite3(v) {
		return Vec3{}, errNotFinite("v")
	}
	if !reduce.AllFinite(k) {
		return Vec3{}, errNotFinite("k")
	}
	r := Vec3{k * v[0], k * v[1], k * v[2]}
	if !finite3(r) {
		return Vec3{}, errOverflow("k·v")
	}
	return r, nil
}

// Vec3Div returns v/k.
// Errors: ErrZeroDivisor when k = 0; ErrNotFinite for a NaN or ±Inf input; ErrOverflow when a
// component of the quotient is outside the float64 range.
//
// ID:      NUM-LINALG-014
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Strang, Introduction to Linear Algebra, 5th ed., Wellesley-Cambridge Press, 2016, §1.1 "Vectors and Linear Combinations"
// Tags:    vector, 3d, divide, scalar-division, linear-algebra
// Units:   v in any unit, k dimensionless; the result in v's unit
// Example: (2, −4, 6)/2 → (1, −2, 3) (hand-checked)
func Vec3Div(v Vec3, k float64) (Vec3, error) {
	if !finite3(v) {
		return Vec3{}, errNotFinite("v")
	}
	if !reduce.AllFinite(k) {
		return Vec3{}, errNotFinite("k")
	}
	if k == 0 {
		return Vec3{}, fmt.Errorf("%w: k", ErrZeroDivisor)
	}
	r := Vec3{v[0] / k, v[1] / k, v[2] / k}
	if !finite3(r) {
		return Vec3{}, errOverflow("v/k")
	}
	return r, nil
}

// lerp returns (1−t)·a + t·b, exact at t = 0 (a) and t = 1 (b): below t = ½ it is a + t(b−a), above it
// b − (1−t)(b−a), each with one rounding (fused multiply-add). When b − a overflows it falls back to
// the weighted sum, which then cannot overflow for t in [0, 1].
func lerp(a, b, t float64) float64 {
	d := b - a
	if !reduce.AllFinite(d) {
		return (1-t)*a + t*b
	}
	if t <= 0.5 {
		return math.FMA(t, d, a)
	}
	return math.FMA(t-1, d, b)
}

// Vec3Lerp returns the linear interpolation (1−t)·a + t·b: a at t = 0 and b at t = 1, exactly. t
// outside [0, 1] extrapolates along the line through a and b.
// Errors: ErrNotFinite for a NaN or ±Inf input; ErrOverflow when a component of the result is outside
// the float64 range.
//
// ID:      NUM-LINALG-032
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Strang, Introduction to Linear Algebra, 5th ed., Wellesley-Cambridge Press, 2016, §1.1 "Vectors and Linear Combinations" (the combination (1−t)a + tb)
// Tags:    vector, 3d, lerp, interpolation, linear-combination, linear-algebra
// Units:   a and b in one unit, t dimensionless; the result in that unit
// Example: a = (0, 0, 0), b = (2, 4, 6), t = 0.25 → (0.5, 1, 1.5) (hand-checked)
func Vec3Lerp(a, b Vec3, t float64) (Vec3, error) {
	if !finite3(a) {
		return Vec3{}, errNotFinite("a")
	}
	if !finite3(b) {
		return Vec3{}, errNotFinite("b")
	}
	if !reduce.AllFinite(t) {
		return Vec3{}, errNotFinite("t")
	}
	r := Vec3{lerp(a[0], b[0], t), lerp(a[1], b[1], t), lerp(a[2], b[2], t)}
	if !finite3(r) {
		return Vec3{}, errOverflow("lerp")
	}
	return r, nil
}

// Vec3Min returns the componentwise minimum of a and b.
// Errors: ErrNotFinite for a NaN or ±Inf component.
//
// ID:      NUM-LINALG-033
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Strang, Introduction to Linear Algebra, 5th ed., Wellesley-Cambridge Press, 2016, §1.1 "Vectors and Linear Combinations" (a vector as its list of components; the minimum is taken component by component)
// Tags:    vector, 3d, min, componentwise, bounding-box, linear-algebra
// Units:   a and b in one unit; the result in that unit
// Example: (1, 5, −2) and (3, 0, −4) → (1, 0, −4) (hand-checked)
func Vec3Min(a, b Vec3) (Vec3, error) {
	if !finite3(a) {
		return Vec3{}, errNotFinite("a")
	}
	if !finite3(b) {
		return Vec3{}, errNotFinite("b")
	}
	return Vec3{min(a[0], b[0]), min(a[1], b[1]), min(a[2], b[2])}, nil
}

// Vec3Max returns the componentwise maximum of a and b.
// Errors: ErrNotFinite for a NaN or ±Inf component.
//
// ID:      NUM-LINALG-034
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Strang, Introduction to Linear Algebra, 5th ed., Wellesley-Cambridge Press, 2016, §1.1 "Vectors and Linear Combinations" (a vector as its list of components; the maximum is taken component by component)
// Tags:    vector, 3d, max, componentwise, bounding-box, linear-algebra
// Units:   a and b in one unit; the result in that unit
// Example: (1, 5, −2) and (3, 0, −4) → (3, 5, −2) (hand-checked)
func Vec3Max(a, b Vec3) (Vec3, error) {
	if !finite3(a) {
		return Vec3{}, errNotFinite("a")
	}
	if !finite3(b) {
		return Vec3{}, errNotFinite("b")
	}
	return Vec3{max(a[0], b[0]), max(a[1], b[1]), max(a[2], b[2])}, nil
}
