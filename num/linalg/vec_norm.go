package linalg

import (
	"fmt"
	"math"

	"github.com/CrimsonCoder42/calcbench/algo/reduce"
)

// maxAbs returns the largest |vᵢ| and the index of the first NaN or ±Inf element (−1 when all are
// finite).
func maxAbs(v []float64) (float64, int) {
	m := 0.0
	for i, x := range v {
		if !reduce.AllFinite(x) {
			return 0, i
		}
		m = max(m, math.Abs(x))
	}
	return m, -1
}

// scaledNorm returns ‖v‖₂ as ns·2ᵉ: v is scaled by the exact power of two 2⁻ᵉ that brings its largest
// element into [1, 2), and the sum of squares of the scaled elements is accumulated by Dot2, so
// neither overflow nor underflow can occur (Blue 1978) and the sum is as accurate as in twice the
// working precision. ns is 0 for the zero vector.
func scaledNorm(v []float64, m float64) (ns float64, e int) {
	if m == 0 {
		return 0, 0
	}
	s, e := pow2Scale(m)
	var d dot2
	for _, x := range v {
		y := x * s
		d.add(y, y)
	}
	return math.Sqrt(d.result()), e
}

func notFiniteAt(name string, i int) error { return fmt.Errorf("%w: %s[%d]", ErrNotFinite, name, i) }

// VecNorm returns the Euclidean length ‖v‖₂ = √(Σ vᵢ²). The vector is scaled by an exact power of two
// so that no square overflows or underflows, and the squares are summed as if in twice the working
// precision (Dot2), so the result is within about one ulp for any length of vector.
// Errors: ErrNotFinite for a NaN or ±Inf element; ErrOverflow when the length exceeds the float64
// range. The norm of an empty vector is 0.
//
// ID:      NUM-LINALG-041
// Time:    O(n)   n = len(v); two passes
// Space:   O(1), no allocation
// Source:  Blue, "A Portable Fortran Program to Find the Euclidean Norm of a Vector", ACM Transactions on Mathematical Software, 1978, 4(1):15–23, doi:10.1145/355769.355771
// Source:  Ogita, Rump, Oishi, "Accurate Sum and Dot Product", SIAM Journal on Scientific Computing, 2005, 26(6):1955–1988, doi:10.1137/030601818, Algorithm 5.3 (Dot2)
// Tags:    vector, n-dimensional, norm, length, euclidean-norm, l2-norm, compensated, linear-algebra
// Units:   v in any unit; the result in that unit
// Example: (1, 2, 2, 4) → 5 (hand-checked)
func VecNorm(v []float64) (float64, error) {
	m, bad := maxAbs(v)
	if bad >= 0 {
		return 0, notFiniteAt("v", bad)
	}
	ns, e := scaledNorm(v, m)
	n := math.Ldexp(ns, e)
	if !reduce.AllFinite(n) {
		return 0, errOverflow("‖v‖₂")
	}
	return n, nil
}

// VecNorm1 returns the 1-norm ‖v‖₁ = Σ|vᵢ|, summed by Ogita, Rump and Oishi's cascaded summation (Sum2),
// as accurate as in twice the working precision.
// Errors: ErrNotFinite for a NaN or ±Inf element; ErrOverflow when the sum exceeds the float64 range.
//
// ID:      NUM-LINALG-042
// Time:    O(n)   n = len(v)
// Space:   O(1), no allocation
// Source:  Higham, Accuracy and Stability of Numerical Algorithms, 2nd ed., SIAM, 2002, §6.1 "Vector Norms"
// Source:  Ogita, Rump, Oishi, "Accurate Sum and Dot Product", SIAM Journal on Scientific Computing, 2005, 26(6):1955–1988, doi:10.1137/030601818, Algorithm 4.4 (Sum2)
// Tags:    vector, n-dimensional, norm, l1-norm, taxicab, manhattan, compensated, linear-algebra
// Units:   v in any unit; the result in that unit
// Example: (1, −2, 3, −4) → 10 (hand-checked)
func VecNorm1(v []float64) (float64, error) {
	var p, s float64
	for i, x := range v {
		if !reduce.AllFinite(x) {
			return 0, notFiniteAt("v", i)
		}
		var q float64
		p, q = reduce.TwoSum(p, math.Abs(x))
		s += q
	}
	r := p + s
	if !reduce.AllFinite(r) {
		return 0, errOverflow("‖v‖₁")
	}
	return r, nil
}

// VecNormInf returns the infinity norm ‖v‖∞ = max|vᵢ| (0 for an empty vector).
// Errors: ErrNotFinite for a NaN or ±Inf element.
//
// ID:      NUM-LINALG-043
// Time:    O(n)   n = len(v)
// Space:   O(1), no allocation
// Source:  Higham, Accuracy and Stability of Numerical Algorithms, 2nd ed., SIAM, 2002, §6.1 "Vector Norms"
// Tags:    vector, n-dimensional, norm, infinity-norm, max-norm, chebyshev, linear-algebra
// Units:   v in any unit; the result in that unit
// Example: (1, −7, 3) → 7 (hand-checked)
func VecNormInf(v []float64) (float64, error) {
	m, bad := maxAbs(v)
	if bad >= 0 {
		return 0, notFiniteAt("v", bad)
	}
	return m, nil
}

// VecNormP returns the p-norm ‖v‖ₚ = (Σ|vᵢ|ᵖ)^(1/p) for p ≥ 1, computed as m·(Σ(|vᵢ|/m)ᵖ)^(1/p) with
// m = max|vᵢ| so that no power overflows, and the sum taken by Sum2. p = 1 and p = 2 agree with VecNorm1
// and VecNorm to within rounding; use VecNormInf for p = ∞.
// Errors: ErrOutOfRange when p < 1; ErrNotFinite for a NaN or ±Inf element or p; ErrOverflow when the
// norm exceeds the float64 range.
//
// ID:      NUM-LINALG-044
// Time:    O(n)   n = len(v); two passes and one math.Pow per element
// Space:   O(1), no allocation
// Source:  Higham, Accuracy and Stability of Numerical Algorithms, 2nd ed., SIAM, 2002, §6.1 "Vector Norms" (Hölder p-norm)
// Source:  Ogita, Rump, Oishi, "Accurate Sum and Dot Product", SIAM Journal on Scientific Computing, 2005, 26(6):1955–1988, doi:10.1137/030601818, Algorithm 4.4 (Sum2)
// Tags:    vector, n-dimensional, norm, p-norm, holder, linear-algebra
// Units:   v in any unit, p dimensionless; the result in v's unit
// Example: (3, 4), p = 2 → 5; (1, 1, 1, 1), p = 3 → 4^(1/3) ≈ 1.5874010519681994 (hand-checked)
func VecNormP(v []float64, p float64) (float64, error) {
	if !reduce.AllFinite(p) {
		return 0, errNotFinite("p")
	}
	if p < 1 {
		return 0, fmt.Errorf("%w: p = %v (need p ≥ 1)", ErrOutOfRange, p)
	}
	m, bad := maxAbs(v)
	if bad >= 0 {
		return 0, notFiniteAt("v", bad)
	}
	if m == 0 {
		return 0, nil
	}
	var sum, comp float64
	for _, x := range v {
		var q float64
		sum, q = reduce.TwoSum(sum, math.Pow(math.Abs(x)/m, p))
		comp += q
	}
	r := m * math.Pow(sum+comp, 1/p)
	if !reduce.AllFinite(r) {
		return 0, errOverflow("‖v‖ₚ")
	}
	return r, nil
}

// VecNormalize returns v/‖v‖₂ written into dst (reused when its capacity suffices; dst may be v).
// The division uses the scaled vector, so it is exact in range even when ‖v‖₂ itself would overflow.
// Errors: ErrZeroVector when v is all zeros (or empty); ErrNotFinite for a NaN or ±Inf element.
//
// ID:      NUM-LINALG-045
// Time:    O(n)   n = len(v); three passes
// Space:   O(1) extra; O(n) for the result, one allocation only if dst lacks capacity
// Source:  Strang, Introduction to Linear Algebra, 5th ed., Wellesley-Cambridge Press, 2016, §1.2 "Lengths and Dot Products"
// Source:  Blue, "A Portable Fortran Program to Find the Euclidean Norm of a Vector", ACM Transactions on Mathematical Software, 1978, 4(1):15–23, doi:10.1145/355769.355771
// Tags:    vector, n-dimensional, normalize, unit-vector, direction, linear-algebra
// Units:   v in any unit; the result is dimensionless
// Example: (1, 2, 2, 4) → (0.2, 0.4, 0.4, 0.8) (hand-checked)
func VecNormalize(dst, v []float64) ([]float64, error) {
	m, bad := maxAbs(v)
	if bad >= 0 {
		return nil, notFiniteAt("v", bad)
	}
	if m == 0 {
		return nil, fmt.Errorf("%w: v", ErrZeroVector)
	}
	ns, e := scaledNorm(v, m)
	s := math.Ldexp(1, -e) // e is within [−1022, 1023]
	dst = grow(dst, len(v))
	for i, x := range v {
		dst[i] = x * s / ns
	}
	return dst, nil
}

// VecDistance returns ‖a − b‖₂, the Euclidean distance between the points a and b, without forming
// a − b in memory, scaled against overflow and summed by Dot2.
// Errors: ErrDimension when len(a) ≠ len(b); ErrNotFinite for a NaN or ±Inf element; ErrOverflow when an
// element difference or the distance exceeds the float64 range.
//
// ID:      NUM-LINALG-046
// Time:    O(n)   n = len(a); two passes
// Space:   O(1), no allocation
// Source:  Strang, Introduction to Linear Algebra, 5th ed., Wellesley-Cambridge Press, 2016, §1.2 "Lengths and Dot Products" (length of v − w)
// Source:  Blue, "A Portable Fortran Program to Find the Euclidean Norm of a Vector", ACM Transactions on Mathematical Software, 1978, 4(1):15–23, doi:10.1145/355769.355771
// Tags:    vector, n-dimensional, distance, euclidean-distance, linear-algebra
// Units:   a and b in one unit; the result in that unit
// Example: (1, 1, 1, 1) and (2, 3, 3, 5) → 5 (hand-checked)
func VecDistance(a, b []float64) (float64, error) {
	if len(a) != len(b) {
		return 0, lengthError("a", len(a), "b", len(b))
	}
	m := 0.0
	for i := range a {
		d := a[i] - b[i]
		if !reduce.AllFinite(d) {
			return 0, elementError(i, "a", a[i], "b", b[i])
		}
		m = max(m, math.Abs(d))
	}
	if m == 0 {
		return 0, nil
	}
	s, e := pow2Scale(m)
	var acc dot2
	for i := range a {
		y := (a[i] - b[i]) * s
		acc.add(y, y)
	}
	r := math.Ldexp(math.Sqrt(acc.result()), e)
	if !reduce.AllFinite(r) {
		return 0, errOverflow("‖a − b‖₂")
	}
	return r, nil
}
