package linalg

import (
	"fmt"

	"github.com/CrimsonCoder42/calcbench/algo/reduce"
)

// grow returns dst resliced to length n when its capacity allows, else a new slice of length n (the one
// allocation an n-dimensional function may make for its result).
func grow(dst []float64, n int) []float64 {
	if cap(dst) >= n {
		return dst[:n]
	}
	return make([]float64, n)
}

// elementError classifies a non-finite element result at index i: ErrNotFinite naming the input that
// was already NaN or ±Inf, otherwise ErrOverflow.
func elementError(i int, aName string, a float64, bName string, b float64) error {
	if !reduce.AllFinite(a) {
		return fmt.Errorf("%w: %s[%d]", ErrNotFinite, aName, i)
	}
	if !reduce.AllFinite(b) {
		return fmt.Errorf("%w: %s[%d]", ErrNotFinite, bName, i)
	}
	return fmt.Errorf("%w: element %d", ErrOverflow, i)
}

func lengthError(aName string, na int, bName string, nb int) error {
	return fmt.Errorf("%w: len(%s) = %d, len(%s) = %d", ErrDimension, aName, na, bName, nb)
}

// VecAdd returns a + b, element by element, written into dst (reused when its capacity suffices, else
// one allocation; dst may be a or b).
// Errors: ErrDimension when len(a) ≠ len(b); ErrNotFinite for a NaN or ±Inf element; ErrOverflow when an
// element of the sum is outside the float64 range. On error the result is nil and dst's contents are
// unspecified.
//
// ID:      NUM-LINALG-035
// Time:    O(n)   n = len(a)
// Space:   O(1) extra; O(n) for the result, one allocation only if dst lacks capacity
// Source:  Strang, Introduction to Linear Algebra, 5th ed., Wellesley-Cambridge Press, 2016, §1.1 "Vectors and Linear Combinations"
// Tags:    vector, n-dimensional, add, linear-algebra
// Units:   a and b in one unit; the result in that unit
// Example: (1, 2, 3, 4) + (4, 3, 2, 1) → (5, 5, 5, 5) (hand-checked)
func VecAdd(dst, a, b []float64) ([]float64, error) {
	if len(a) != len(b) {
		return nil, lengthError("a", len(a), "b", len(b))
	}
	dst = grow(dst, len(a))
	for i := range a {
		x := a[i] + b[i]
		if !reduce.AllFinite(x) {
			return nil, elementError(i, "a", a[i], "b", b[i])
		}
		dst[i] = x
	}
	return dst, nil
}

// VecSub returns a − b, element by element, written into dst (as VecAdd).
// Errors: as VecAdd.
//
// ID:      NUM-LINALG-036
// Time:    O(n)   n = len(a)
// Space:   O(1) extra; O(n) for the result, one allocation only if dst lacks capacity
// Source:  Strang, Introduction to Linear Algebra, 5th ed., Wellesley-Cambridge Press, 2016, §1.1 "Vectors and Linear Combinations"
// Tags:    vector, n-dimensional, subtract, difference, linear-algebra
// Units:   a and b in one unit; the result in that unit
// Example: (1, 2, 3, 4) − (4, 3, 2, 1) → (−3, −1, 1, 3) (hand-checked)
func VecSub(dst, a, b []float64) ([]float64, error) {
	if len(a) != len(b) {
		return nil, lengthError("a", len(a), "b", len(b))
	}
	dst = grow(dst, len(a))
	for i := range a {
		x := a[i] - b[i]
		if !reduce.AllFinite(x) {
			return nil, elementError(i, "a", a[i], "b", b[i])
		}
		dst[i] = x
	}
	return dst, nil
}

// VecScale returns k·v written into dst (reused when its capacity suffices; dst may be v).
// Errors: ErrNotFinite for a NaN or ±Inf input; ErrOverflow when an element of the product is outside
// the float64 range. On error the result is nil.
//
// ID:      NUM-LINALG-037
// Time:    O(n)   n = len(v)
// Space:   O(1) extra; O(n) for the result, one allocation only if dst lacks capacity
// Source:  Strang, Introduction to Linear Algebra, 5th ed., Wellesley-Cambridge Press, 2016, §1.1 "Vectors and Linear Combinations"
// Tags:    vector, n-dimensional, scale, scalar-multiplication, linear-algebra
// Units:   v in any unit, k dimensionless; the result in v's unit
// Example: 2·(1, −2, 3, 0) → (2, −4, 6, 0) (hand-checked)
func VecScale(dst, v []float64, k float64) ([]float64, error) {
	if !reduce.AllFinite(k) {
		return nil, errNotFinite("k")
	}
	dst = grow(dst, len(v))
	for i, x := range v {
		y := k * x
		if !reduce.AllFinite(y) {
			return nil, elementError(i, "v", x, "k", k)
		}
		dst[i] = y
	}
	return dst, nil
}

// VecDiv returns v/k written into dst (reused when its capacity suffices; dst may be v).
// Errors: ErrZeroDivisor when k = 0; otherwise as VecScale.
//
// ID:      NUM-LINALG-038
// Time:    O(n)   n = len(v)
// Space:   O(1) extra; O(n) for the result, one allocation only if dst lacks capacity
// Source:  Strang, Introduction to Linear Algebra, 5th ed., Wellesley-Cambridge Press, 2016, §1.1 "Vectors and Linear Combinations"
// Tags:    vector, n-dimensional, divide, scalar-division, linear-algebra
// Units:   v in any unit, k dimensionless; the result in v's unit
// Example: (2, −4, 6, 0)/2 → (1, −2, 3, 0) (hand-checked)
func VecDiv(dst, v []float64, k float64) ([]float64, error) {
	if !reduce.AllFinite(k) {
		return nil, errNotFinite("k")
	}
	if k == 0 {
		return nil, fmt.Errorf("%w: k", ErrZeroDivisor)
	}
	dst = grow(dst, len(v))
	for i, x := range v {
		y := x / k
		if !reduce.AllFinite(y) {
			return nil, elementError(i, "v", x, "k", k)
		}
		dst[i] = y
	}
	return dst, nil
}

// VecLerp returns (1−t)·a + t·b written into dst: a at t = 0 and b at t = 1, exactly; t outside [0, 1]
// extrapolates. dst may be a or b.
// Errors: ErrDimension when len(a) ≠ len(b); ErrNotFinite for a NaN or ±Inf input; ErrOverflow when an
// element is outside the float64 range.
//
// ID:      NUM-LINALG-049
// Time:    O(n)   n = len(a)
// Space:   O(1) extra; O(n) for the result, one allocation only if dst lacks capacity
// Source:  Strang, Introduction to Linear Algebra, 5th ed., Wellesley-Cambridge Press, 2016, §1.1 "Vectors and Linear Combinations" (the combination (1−t)a + tb)
// Tags:    vector, n-dimensional, lerp, interpolation, linear-combination, linear-algebra
// Units:   a and b in one unit, t dimensionless; the result in that unit
// Example: a = (0, 10), b = (4, 20), t = 0.5 → (2, 15) (hand-checked)
func VecLerp(dst, a, b []float64, t float64) ([]float64, error) {
	if len(a) != len(b) {
		return nil, lengthError("a", len(a), "b", len(b))
	}
	if !reduce.AllFinite(t) {
		return nil, errNotFinite("t")
	}
	dst = grow(dst, len(a))
	for i := range a {
		x := lerp(a[i], b[i], t)
		if !reduce.AllFinite(x) {
			return nil, elementError(i, "a", a[i], "b", b[i])
		}
		dst[i] = x
	}
	return dst, nil
}

// VecMin returns the element-by-element minimum of a and b written into dst (dst may be a or b).
// Errors: ErrDimension when len(a) ≠ len(b); ErrNotFinite for a NaN or ±Inf element.
//
// ID:      NUM-LINALG-050
// Time:    O(n)   n = len(a)
// Space:   O(1) extra; O(n) for the result, one allocation only if dst lacks capacity
// Source:  Strang, Introduction to Linear Algebra, 5th ed., Wellesley-Cambridge Press, 2016, §1.1 "Vectors and Linear Combinations" (a vector as its list of components; the minimum is taken component by component)
// Tags:    vector, n-dimensional, min, componentwise, linear-algebra
// Units:   a and b in one unit; the result in that unit
// Example: (1, 5, −2) and (3, 0, −4) → (1, 0, −4) (hand-checked)
func VecMin(dst, a, b []float64) ([]float64, error) {
	if len(a) != len(b) {
		return nil, lengthError("a", len(a), "b", len(b))
	}
	dst = grow(dst, len(a))
	for i := range a {
		x, y := a[i], b[i]
		if !reduce.AllFinite(x) || !reduce.AllFinite(y) {
			return nil, elementError(i, "a", x, "b", y)
		}
		dst[i] = min(x, y)
	}
	return dst, nil
}

// VecMax returns the element-by-element maximum of a and b written into dst (dst may be a or b).
// Errors: as VecMin.
//
// ID:      NUM-LINALG-051
// Time:    O(n)   n = len(a)
// Space:   O(1) extra; O(n) for the result, one allocation only if dst lacks capacity
// Source:  Strang, Introduction to Linear Algebra, 5th ed., Wellesley-Cambridge Press, 2016, §1.1 "Vectors and Linear Combinations" (a vector as its list of components; the maximum is taken component by component)
// Tags:    vector, n-dimensional, max, componentwise, linear-algebra
// Units:   a and b in one unit; the result in that unit
// Example: (1, 5, −2) and (3, 0, −4) → (3, 5, −2) (hand-checked)
func VecMax(dst, a, b []float64) ([]float64, error) {
	if len(a) != len(b) {
		return nil, lengthError("a", len(a), "b", len(b))
	}
	dst = grow(dst, len(a))
	for i := range a {
		x, y := a[i], b[i]
		if !reduce.AllFinite(x) || !reduce.AllFinite(y) {
			return nil, elementError(i, "a", x, "b", y)
		}
		dst[i] = max(x, y)
	}
	return dst, nil
}
