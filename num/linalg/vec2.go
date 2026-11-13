package linalg

import (
	"errors"
	"fmt"
	"math"

	"github.com/CrimsonCoder42/calcbench/algo/reduce"
)

// Errors returned by this package. Each is wrapped with the offending parameter's name; test with
// errors.Is.
var (
	// ErrNotFinite is returned when an input is NaN or ±Inf.
	ErrNotFinite = errors.New("linalg: value is NaN or infinite")
	// ErrOverflow is returned when finite inputs give a result outside the float64 range.
	ErrOverflow = errors.New("linalg: result is outside the float64 range")
	// ErrZeroVector is returned when a direction is needed from a zero vector or a zero quaternion.
	ErrZeroVector = errors.New("linalg: zero vector or quaternion has no direction")
	// ErrZeroDivisor is returned when a vector is divided by zero.
	ErrZeroDivisor = errors.New("linalg: division by zero")
	// ErrDimension is returned when lengths or matrix dimensions do not agree.
	ErrDimension = errors.New("linalg: lengths or dimensions do not match")
	// ErrSingular is returned when a matrix has no inverse.
	ErrSingular = errors.New("linalg: matrix is singular")
	// ErrNotSymmetric is returned by the symmetric eigenvalue functions for a matrix that is not
	// exactly symmetric.
	ErrNotSymmetric = errors.New("linalg: matrix is not symmetric")
	// ErrOutOfRange is returned when a parameter is outside its domain.
	ErrOutOfRange = errors.New("linalg: parameter outside its domain")
)

// Vec2 is a vector in the plane, [x, y]. It is a value: functions take and return copies and never
// allocate.
type Vec2 [2]float64

func finite2(v Vec2) bool { return reduce.AllFinite(v[:]...) }

func errNotFinite(name string) error { return fmt.Errorf("%w: %s", ErrNotFinite, name) }

func errOverflow(name string) error { return fmt.Errorf("%w: %s", ErrOverflow, name) }

// pow2Scale returns s = 2⁻ᵉ, an exact power of two that brings m (> 0) into [1, 2) — or, for a subnormal
// m, as close as a finite s allows (e is clamped to [−1022, 1023], so s itself never overflows).
func pow2Scale(m float64) (s float64, e int) {
	e = max(-1022, min(1023, math.Ilogb(m)))
	return math.Ldexp(1, -e), e
}

// det2 returns a·d − b·c by Kahan's algorithm with a fused multiply-add, which keeps the relative
// error within 2 units in the last place even when the two products nearly cancel (Jeannerod, Louvet,
// Muller 2013). When a product overflows, the inputs are rescaled by exact powers of two and the
// determinant recomputed, so a non-finite result means the true value is outside the float64 range.
func det2(a, b, c, d float64) float64 {
	r := det2raw(a, b, c, d)
	if reduce.AllFinite(r) {
		return r
	}
	// A product overflowed although the difference may be finite: scale (a, b) and (c, d) by exact powers
	// of two, which scales the result by their product, and undo it at the end.
	m1, m2 := max(math.Abs(a), math.Abs(b)), max(math.Abs(c), math.Abs(d))
	if m1 == 0 || m2 == 0 || !reduce.AllFinite(m1) || !reduce.AllFinite(m2) {
		return r
	}
	_, e1 := pow2Scale(m1)
	_, e2 := pow2Scale(m2)
	r = det2raw(math.Ldexp(a, -e1), math.Ldexp(b, -e1), math.Ldexp(c, -e2), math.Ldexp(d, -e2))
	return math.Ldexp(r, e1+e2)
}

func det2raw(a, b, c, d float64) float64 {
	w := float64(b * c) // the explicit conversion forbids fusing this product into a later operation
	e := math.FMA(-b, c, w)
	f := math.FMA(a, d, -w)
	return f + e
}

// Vec2Add returns a + b, component by component.
// Errors: ErrNotFinite for a NaN or ±Inf component; ErrOverflow when a component of the sum is
// outside the float64 range. The result is the zero vector with any error.
//
// ID:      NUM-LINALG-001
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Strang, Introduction to Linear Algebra, 5th ed., Wellesley-Cambridge Press, 2016, §1.1 "Vectors and Linear Combinations"
// Tags:    vector, 2d, add, linear-algebra
// Units:   a and b in one unit; the result in that unit
// Example: (1, 2) + (3, −5) → (4, −3) (hand-checked)
func Vec2Add(a, b Vec2) (Vec2, error) {
	if !finite2(a) {
		return Vec2{}, errNotFinite("a")
	}
	if !finite2(b) {
		return Vec2{}, errNotFinite("b")
	}
	r := Vec2{a[0] + b[0], a[1] + b[1]}
	if !finite2(r) {
		return Vec2{}, errOverflow("a + b")
	}
	return r, nil
}

// Vec2Sub returns a − b, component by component.
// Errors: as Vec2Add.
//
// ID:      NUM-LINALG-002
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Strang, Introduction to Linear Algebra, 5th ed., Wellesley-Cambridge Press, 2016, §1.1 "Vectors and Linear Combinations"
// Tags:    vector, 2d, subtract, difference, linear-algebra
// Units:   a and b in one unit; the result in that unit
// Example: (1, 2) − (3, −5) → (−2, 7) (hand-checked)
func Vec2Sub(a, b Vec2) (Vec2, error) {
	if !finite2(a) {
		return Vec2{}, errNotFinite("a")
	}
	if !finite2(b) {
		return Vec2{}, errNotFinite("b")
	}
	r := Vec2{a[0] - b[0], a[1] - b[1]}
	if !finite2(r) {
		return Vec2{}, errOverflow("a − b")
	}
	return r, nil
}

// Vec2Scale returns k·v.
// Errors: ErrNotFinite for a NaN or ±Inf input; ErrOverflow when a component of the product is outside
// the float64 range.
//
// ID:      NUM-LINALG-003
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Strang, Introduction to Linear Algebra, 5th ed., Wellesley-Cambridge Press, 2016, §1.1 "Vectors and Linear Combinations"
// Tags:    vector, 2d, scale, scalar-multiplication, linear-algebra
// Units:   v in any unit, k dimensionless; the result in v's unit
// Example: 3·(1, −2) → (3, −6) (hand-checked)
func Vec2Scale(v Vec2, k float64) (Vec2, error) {
	if !finite2(v) {
		return Vec2{}, errNotFinite("v")
	}
	if !reduce.AllFinite(k) {
		return Vec2{}, errNotFinite("k")
	}
	r := Vec2{k * v[0], k * v[1]}
	if !finite2(r) {
		return Vec2{}, errOverflow("k·v")
	}
	return r, nil
}

// Vec2Div returns v/k.
// Errors: ErrZeroDivisor when k = 0; ErrNotFinite for a NaN or ±Inf input; ErrOverflow when a
// component of the quotient is outside the float64 range.
//
// ID:      NUM-LINALG-004
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Strang, Introduction to Linear Algebra, 5th ed., Wellesley-Cambridge Press, 2016, §1.1 "Vectors and Linear Combinations"
// Tags:    vector, 2d, divide, scalar-division, linear-algebra
// Units:   v in any unit, k dimensionless; the result in v's unit
// Example: (3, −6)/3 → (1, −2) (hand-checked)
func Vec2Div(v Vec2, k float64) (Vec2, error) {
	if !finite2(v) {
		return Vec2{}, errNotFinite("v")
	}
	if !reduce.AllFinite(k) {
		return Vec2{}, errNotFinite("k")
	}
	if k == 0 {
		return Vec2{}, fmt.Errorf("%w: k", ErrZeroDivisor)
	}
	r := Vec2{v[0] / k, v[1] / k}
	if !finite2(r) {
		return Vec2{}, errOverflow("v/k")
	}
	return r, nil
}

// Vec2Dot returns the dot product a·b = a₀b₀ + a₁b₁, computed as a 2×2 determinant by Kahan's
// fused multiply-add algorithm so that it is within 2 ulps even when the two products cancel.
// Errors: ErrNotFinite for a NaN or ±Inf component; ErrOverflow when the product is outside the
// float64 range.
//
// ID:      NUM-LINALG-005
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Strang, Introduction to Linear Algebra, 5th ed., Wellesley-Cambridge Press, 2016, §1.2 "Lengths and Dot Products"
// Source:  Jeannerod, Louvet, Muller, "Further analysis of Kahan's algorithm for the accurate computation of 2×2 determinants", Mathematics of Computation, 2013, 82(284):2245–2264, doi:10.1090/S0025-5718-2013-02679-8
// Tags:    vector, 2d, dot-product, inner-product, linear-algebra
// Units:   a in unit A, b in unit B; the result in A·B
// Example: (1, 2)·(3, 4) → 11 (hand-checked)
func Vec2Dot(a, b Vec2) (float64, error) {
	if !finite2(a) {
		return 0, errNotFinite("a")
	}
	if !finite2(b) {
		return 0, errNotFinite("b")
	}
	d := det2(a[0], -a[1], b[1], b[0])
	if !reduce.AllFinite(d) {
		return 0, errOverflow("a·b")
	}
	return d, nil
}

// Vec2Cross returns the scalar cross product a × b = a₀b₁ − a₁b₀, the determinant of the matrix with
// rows a and b. It is positive when b is counterclockwise from a (right-hand rule: it is the z
// component of the 3D cross product), negative when clockwise, zero when parallel; its absolute value
// is the area of the parallelogram the two vectors span. Computed by Kahan's fused multiply-add
// algorithm, within 2 ulps.
// Errors: as Vec2Dot.
//
// ID:      NUM-LINALG-006
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Strang, Introduction to Linear Algebra, 5th ed., Wellesley-Cambridge Press, 2016, §5.3 "Cramer's Rule, Inverses, and Volumes"
// Source:  Jeannerod, Louvet, Muller, "Further analysis of Kahan's algorithm for the accurate computation of 2×2 determinants", Mathematics of Computation, 2013, 82(284):2245–2264, doi:10.1090/S0025-5718-2013-02679-8
// Tags:    vector, 2d, cross-product, determinant, right-hand-rule, orientation, area, linear-algebra
// Units:   a in unit A, b in unit B; the result in A·B
// Example: (1, 0) × (0, 1) → 1; (3, 1) × (2, 4) → 10 (hand-checked)
func Vec2Cross(a, b Vec2) (float64, error) {
	if !finite2(a) {
		return 0, errNotFinite("a")
	}
	if !finite2(b) {
		return 0, errNotFinite("b")
	}
	d := det2(a[0], a[1], b[0], b[1])
	if !reduce.AllFinite(d) {
		return 0, errOverflow("a × b")
	}
	return d, nil
}

// Vec2Norm returns the Euclidean length |v| = √(x² + y²), without overflow or underflow in the squares.
// Errors: ErrNotFinite for a NaN or ±Inf component; ErrOverflow when the length exceeds the float64
// range.
//
// ID:      NUM-LINALG-007
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Strang, Introduction to Linear Algebra, 5th ed., Wellesley-Cambridge Press, 2016, §1.2 "Lengths and Dot Products"
// Source:  Blue, "A Portable Fortran Program to Find the Euclidean Norm of a Vector", ACM Transactions on Mathematical Software, 1978, 4(1):15–23, doi:10.1145/355769.355771
// Tags:    vector, 2d, norm, length, magnitude, linear-algebra
// Units:   v in any unit; the result in that unit
// Example: (3, 4) → 5 (hand-checked)
func Vec2Norm(v Vec2) (float64, error) {
	if !finite2(v) {
		return 0, errNotFinite("v")
	}
	n := math.Hypot(v[0], v[1])
	if !reduce.AllFinite(n) {
		return 0, errOverflow("|v|")
	}
	return n, nil
}

// Vec2Normalize returns the unit vector v/|v|.
// Errors: ErrZeroVector when v = 0; ErrNotFinite for a NaN or ±Inf component.
//
// ID:      NUM-LINALG-008
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Strang, Introduction to Linear Algebra, 5th ed., Wellesley-Cambridge Press, 2016, §1.2 "Lengths and Dot Products"
// Tags:    vector, 2d, normalize, unit-vector, direction, linear-algebra
// Units:   v in any unit; the result is dimensionless
// Example: (3, 4) → (0.6, 0.8) (hand-checked)
func Vec2Normalize(v Vec2) (Vec2, error) {
	if !finite2(v) {
		return Vec2{}, errNotFinite("v")
	}
	u, ok := unit2(v)
	if !ok {
		return Vec2{}, fmt.Errorf("%w: v", ErrZeroVector)
	}
	return u, nil
}

// unit2 returns v/|v| and false for the zero vector. |v| may overflow for components near the float64
// maximum, so v is first scaled by an exact power of two.
func unit2(v Vec2) (Vec2, bool) {
	m := max(math.Abs(v[0]), math.Abs(v[1]))
	if m == 0 {
		return Vec2{}, false
	}
	s, _ := pow2Scale(m)
	x, y := v[0]*s, v[1]*s
	n := math.Hypot(x, y)
	return Vec2{x / n, y / n}, true
}

// Vec2AngleBetween returns the unoriented angle between a and b, in [0, π] radians, by Kahan's formula
// θ = 2·atan2(|â − b̂|, |â + b̂|) with â = a/|a|, b̂ = b/|b|. Unlike arccos of the normalised dot
// product, it keeps full accuracy when the vectors are nearly parallel or antiparallel, and unlike
// arcsin of the cross product it keeps it near π/2.
// Errors: ErrZeroVector when a or b is zero; ErrNotFinite for a NaN or ±Inf component.
//
// ID:      NUM-LINALG-009
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Kahan, "How Futile are Mindless Assessments of Roundoff in Floating-Point Computation?", University of California, Berkeley, 11 January 2006, §12 "Mangled Angles", pp. 46–47, https://people.eecs.berkeley.edu/~wkahan/Mindless.pdf (read 2026-11-12)
// Tags:    vector, 2d, angle, angle-between, linear-algebra
// Units:   a and b in any units; the result in radians
// Example: (1, 0) and (1, 1) → π/4 (hand-checked)
func Vec2AngleBetween(a, b Vec2) (float64, error) {
	if !finite2(a) {
		return 0, errNotFinite("a")
	}
	if !finite2(b) {
		return 0, errNotFinite("b")
	}
	ua, ok := unit2(a)
	if !ok {
		return 0, fmt.Errorf("%w: a", ErrZeroVector)
	}
	ub, ok := unit2(b)
	if !ok {
		return 0, fmt.Errorf("%w: b", ErrZeroVector)
	}
	return 2 * math.Atan2(math.Hypot(ua[0]-ub[0], ua[1]-ub[1]), math.Hypot(ua[0]+ub[0], ua[1]+ub[1])), nil
}

// Vec2Rotate returns v turned counterclockwise by angleRad about the origin: the active rotation
// [cos −sin; sin cos]·v. Each component is a 2×2 determinant computed by Kahan's algorithm.
// Errors: ErrNotFinite for a NaN or ±Inf input; ErrOverflow when a component of the result is outside
// the float64 range (possible only for |v| near the float64 maximum).
//
// ID:      NUM-LINALG-010
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Diebel, "Representing Attitude: Euler Angles, Unit Quaternions, and Rotation Vectors", Stanford University, 20 October 2006, §2.4, eq. (16) (the upper-left 2×2 block of R₃, transposed: Diebel's matrices are frame rotations)
// Tags:    vector, 2d, rotate, rotation, linear-algebra
// Units:   v in any unit, angleRad in radians; the result in v's unit
// Example: (1, 0) by π/2 → (0, 1) (hand-checked; the computed x is cos(π/2) ≈ 6.1e−17)
func Vec2Rotate(v Vec2, angleRad float64) (Vec2, error) {
	if !finite2(v) {
		return Vec2{}, errNotFinite("v")
	}
	if !reduce.AllFinite(angleRad) {
		return Vec2{}, errNotFinite("angleRad")
	}
	s, c := math.Sincos(angleRad)
	r := Vec2{det2(c, s, v[1], v[0]), det2(s, -c, v[1], v[0])}
	if !finite2(r) {
		return Vec2{}, errOverflow("rotated v")
	}
	return r, nil
}
