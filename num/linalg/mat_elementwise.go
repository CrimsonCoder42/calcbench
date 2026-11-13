package linalg

import (
	"fmt"

	"github.com/CrimsonCoder42/calcbench/algo/reduce"
)

// elementwise checks that a and b are rows×cols matrices with finite elements.
func elementwise(a []float64, b []float64, rows, cols int) error {
	if err := checkMatrix("a", a, rows, cols); err != nil {
		return err
	}
	if err := checkMatrix("b", b, rows, cols); err != nil {
		return err
	}
	if i := firstNotFinite(a); i >= 0 {
		return notFiniteAt("a", i)
	}
	if i := firstNotFinite(b); i >= 0 {
		return notFiniteAt("b", i)
	}
	return nil
}

// finishElementwise reports ErrOverflow for the first element of dst outside the float64 range.
func finishElementwise(dst []float64, op string) ([]float64, error) {
	if i := firstNotFinite(dst); i >= 0 {
		return nil, fmt.Errorf("%w: %s element %d", ErrOverflow, op, i)
	}
	return dst, nil
}

// MatAdd returns the elementwise sum A + B of two rows×cols matrices in dst (dst may be a or b).
// Errors: ErrDimension when a length disagrees with rows×cols; ErrNotFinite for a NaN or ±Inf element;
// ErrOverflow when a sum is outside the float64 range.
//
// ID:      NUM-LINALG-116
// Time:    O(r·c)   r = rows, c = cols
// Space:   O(1) extra; O(r·c) for the result, one allocation only if dst lacks capacity
// Source:  Strang, Introduction to Linear Algebra, 5th ed., Wellesley-Cambridge Press, 2016, §2.4 "Rules for Matrix Operations" (A + B entry by entry)
// Tags:    matrix, n-dimensional, add, elementwise, linear-algebra
// Units:   a and b in one unit; the result in that unit
// Example: [1 2; 3 4] + [5 6; 7 8] → [6 8; 10 12] (hand-checked)
func MatAdd(dst, a, b []float64, rows, cols int) ([]float64, error) {
	if err := elementwise(a, b, rows, cols); err != nil {
		return nil, err
	}
	dst = grow(dst, rows*cols)
	for i := range a {
		dst[i] = a[i] + b[i]
	}
	return finishElementwise(dst, "sum")
}

// MatSub returns the elementwise difference A − B of two rows×cols matrices in dst (dst may be a or b).
// Errors: ErrDimension when a length disagrees with rows×cols; ErrNotFinite for a NaN or ±Inf element;
// ErrOverflow when a difference is outside the float64 range.
//
// ID:      NUM-LINALG-117
// Time:    O(r·c)   r = rows, c = cols
// Space:   O(1) extra; O(r·c) for the result, one allocation only if dst lacks capacity
// Source:  Strang, Introduction to Linear Algebra, 5th ed., Wellesley-Cambridge Press, 2016, §2.4 "Rules for Matrix Operations" (A − B entry by entry)
// Tags:    matrix, n-dimensional, subtract, elementwise, linear-algebra
// Units:   a and b in one unit; the result in that unit
// Example: [5 6; 7 8] − [1 2; 3 4] → [4 4; 4 4] (hand-checked)
func MatSub(dst, a, b []float64, rows, cols int) ([]float64, error) {
	if err := elementwise(a, b, rows, cols); err != nil {
		return nil, err
	}
	dst = grow(dst, rows*cols)
	for i := range a {
		dst[i] = a[i] - b[i]
	}
	return finishElementwise(dst, "difference")
}

// MatScale returns k·A for a rows×cols matrix A in dst (dst may be a).
// Errors: ErrDimension when len(a) disagrees with rows×cols; ErrNotFinite for a NaN or ±Inf in a or k;
// ErrOverflow when a product is outside the float64 range.
//
// ID:      NUM-LINALG-118
// Time:    O(r·c)   r = rows, c = cols
// Space:   O(1) extra; O(r·c) for the result, one allocation only if dst lacks capacity
// Source:  Strang, Introduction to Linear Algebra, 5th ed., Wellesley-Cambridge Press, 2016, §2.4 "Rules for Matrix Operations" (cA multiplies every entry by c)
// Tags:    matrix, n-dimensional, scale, scalar-multiply, linear-algebra
// Units:   a in unit A, k in unit K; the result in A·K
// Example: 3·[1 2; 3 4] → [3 6; 9 12] (hand-checked)
func MatScale(dst, a []float64, rows, cols int, k float64) ([]float64, error) {
	if err := checkMatrix("a", a, rows, cols); err != nil {
		return nil, err
	}
	if i := firstNotFinite(a); i >= 0 {
		return nil, notFiniteAt("a", i)
	}
	if !reduce.AllFinite(k) {
		return nil, errNotFinite("k")
	}
	dst = grow(dst, rows*cols)
	for i := range a {
		dst[i] = k * a[i]
	}
	return finishElementwise(dst, "product")
}

// MatHadamard returns the Hadamard (entrywise) product A ∘ B, (A ∘ B)ᵢⱼ = aᵢⱼ·bᵢⱼ, of two rows×cols matrices
// in dst (dst may be a or b).
// Errors: ErrDimension when a length disagrees with rows×cols; ErrNotFinite for a NaN or ±Inf element;
// ErrOverflow when a product is outside the float64 range.
//
// ID:      NUM-LINALG-119
// Time:    O(r·c)   r = rows, c = cols
// Space:   O(1) extra; O(r·c) for the result, one allocation only if dst lacks capacity
// Source:  Styan, "Hadamard products and multivariate statistical analysis", Linear Algebra and its Applications, 1973, 6:217–240, doi:10.1016/0024-3795(73)90023-2 (definition: the elementwise product A∘B)
// Tags:    matrix, n-dimensional, hadamard, schur-product, elementwise, linear-algebra
// Units:   a in unit A, b in unit B; the result in A·B
// Example: [1 2; 3 4] ∘ [5 6; 7 8] → [5 12; 21 32] (hand-checked)
func MatHadamard(dst, a, b []float64, rows, cols int) ([]float64, error) {
	if err := elementwise(a, b, rows, cols); err != nil {
		return nil, err
	}
	dst = grow(dst, rows*cols)
	for i := range a {
		dst[i] = a[i] * b[i]
	}
	return finishElementwise(dst, "product")
}
