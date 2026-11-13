package linalg

import (
	"fmt"

	"github.com/CrimsonCoder42/calcbench/algo/reduce"
)

// Vandermonde returns the len(x)×cols Vandermonde matrix V[i][j] = x[i]ʲ (ascending powers, j = 0…cols−1)
// row-major in dst (reused when its capacity suffices). Each row is built by repeated multiplication.
// Errors: ErrDimension when cols < 1 or len(x)·cols overflows; ErrNotFinite for a NaN or ±Inf x[i];
// ErrOverflow when a power leaves the float64 range.
//
// ID:      NUM-LINALG-170
// Time:    O(m·cols)   m = len(x)
// Space:   O(1) extra; O(m·cols) for the result, one allocation only if dst lacks capacity
// Source:  Golub, Van Loan, Matrix Computations, 4th ed., Johns Hopkins University Press, 2013, §4.6 "Vandermonde Systems" (V with entries xᵢʲ)
// Tags:    matrix, n-dimensional, vandermonde, polynomial, interpolation, least-squares, linear-algebra
// Units:   x in unit X; column j in Xʲ
// Example: x = (1, 2, 3), cols 3 → [1 1 1; 1 2 4; 1 3 9] (hand-checked)
func Vandermonde(dst, x []float64, cols int) ([]float64, error) {
	if cols < 1 {
		return nil, fmt.Errorf("%w: cols = %d", ErrDimension, cols)
	}
	m := len(x)
	if m > 0 && cols > int(^uint(0)>>1)/m {
		return nil, fmt.Errorf("%w: %d×%d too large", ErrDimension, m, cols)
	}
	if i := firstNotFinite(x); i >= 0 {
		return nil, notFiniteAt("x", i)
	}
	dst = grow(dst, m*cols)
	for i, xi := range x {
		row := dst[i*cols : i*cols+cols]
		p := 1.0
		for j := range row {
			if !reduce.AllFinite(p) {
				return nil, fmt.Errorf("%w: x[%d]^%d", ErrOverflow, i, j)
			}
			row[j] = p
			p *= xi
		}
	}
	return dst, nil
}
