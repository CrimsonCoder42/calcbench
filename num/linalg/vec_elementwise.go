package linalg

import (
	"fmt"

	"github.com/CrimsonCoder42/calcbench/algo/reduce"
)

// VecHadamard returns the elementwise product a ∘ b (numpy's multiply on two vectors), in dst (reused when
// its capacity suffices; dst may be a or b).
// Errors: ErrDimension when len(a) ≠ len(b); ErrNotFinite for a NaN or ±Inf element; ErrOverflow when a
// product leaves the float64 range.
//
// ID:      NUM-LINALG-187
// Time:    O(n)   n = len(a)
// Space:   O(1) extra; O(n) for the result, one allocation only if dst lacks capacity
// Source:  Horn, Johnson, Topics in Matrix Analysis, Cambridge University Press, 1991, §5.0 (the Hadamard product, entrywise)
// Tags:    vector, n-dimensional, hadamard, elementwise, multiply, linear-algebra
// Units:   a in A, b in B; the result in A·B
// Example: (1, 2, 3) ∘ (4, 5, 6) → (4, 10, 18) (hand-checked)
func VecHadamard(dst, a, b []float64) ([]float64, error) {
	if len(a) != len(b) {
		return nil, lengthError("a", len(a), "b", len(b))
	}
	dst = grow(dst, len(a))
	for i := range a {
		x := a[i] * b[i]
		if !reduce.AllFinite(x) {
			return nil, elementError(i, "a", a[i], "b", b[i])
		}
		dst[i] = x
	}
	return dst, nil
}

// VecDivElementwise returns a ⊘ b, aᵢ/bᵢ element by element (numpy's divide on two vectors), in dst (reused
// when its capacity suffices; dst may be a or b).
// Errors: ErrDimension when len(a) ≠ len(b); ErrNotFinite for a NaN or ±Inf element; ErrZeroDivisor when
// some bᵢ = 0; ErrOverflow when a quotient leaves the float64 range.
//
// ID:      NUM-LINALG-188
// Time:    O(n)   n = len(a)
// Space:   O(1) extra; O(n) for the result, one allocation only if dst lacks capacity
// Source:  Horn, Johnson, Topics in Matrix Analysis, Cambridge University Press, 1991, §5.0 (entrywise operations; the quotient is the Hadamard product with the entrywise inverse of b)
// Tags:    vector, n-dimensional, elementwise, divide, hadamard-division, linear-algebra
// Units:   a in A, b in B; the result in A/B
// Example: (4, 10, 18) ⊘ (4, 5, 6) → (1, 2, 3) (hand-checked)
func VecDivElementwise(dst, a, b []float64) ([]float64, error) {
	if len(a) != len(b) {
		return nil, lengthError("a", len(a), "b", len(b))
	}
	dst = grow(dst, len(a))
	for i := range a {
		if !reduce.AllFinite(a[i]) || !reduce.AllFinite(b[i]) {
			return nil, elementError(i, "a", a[i], "b", b[i])
		}
		if b[i] == 0 {
			return nil, fmt.Errorf("%w: b[%d] = 0", ErrZeroDivisor, i)
		}
		x := a[i] / b[i]
		if !reduce.AllFinite(x) {
			return nil, fmt.Errorf("%w: element %d", ErrOverflow, i)
		}
		dst[i] = x
	}
	return dst, nil
}
