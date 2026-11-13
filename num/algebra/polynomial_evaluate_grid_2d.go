package algebra

import (
	"fmt"

	"github.com/CrimsonCoder42/calcbench/algo/reduce"
)

// PolynomialEvaluateGrid2D returns p(xs[i], ys[j]) for every pair of the Cartesian product of xs and ys,
// in dst (row-major, dst[i·len(ys)+j]), where p(x, y) = Σₖ Σₗ c[k·cols+l]·xᵏ·yˡ is given by a row-major
// rows×cols coefficient grid (degreeX = rows−1, degreeY = cols−1), as numpy's polygrid2d. Nested Horner:
// for each y the row polynomials are evaluated once (Horner in y), then Horner in x over those row values
// for every x, so the work per grid point is O(rows), not O(rows·cols). An empty xs or ys gives an
// empty result. dst may be nil; it is reused when its capacity suffices.
// Errors: ErrEmpty when rows or cols < 1; ErrDegenerate when len(coeffs) ≠ rows·cols; ErrNotFinite for a
// NaN or ±Inf coefficient, x or y; ErrOverflow when a value leaves the float64 range.
//
// ID:      NUM-ALGEB-027
// Time:    O(ny·rows·cols + nx·ny·rows)   nx = len(xs), ny = len(ys)
// Space:   O(rows) extra (one allocation); O(nx·ny) for the result
// Source:  Knuth, The Art of Computer Programming, Vol. 2: Seminumerical Algorithms, 3rd ed., Addison-Wesley, 1997, §4.6.4 "Evaluation of Polynomials" (Horner's rule; polynomials in several variables)
// Tags:    polynomial, evaluate, bivariate, two-dimensional, grid, cartesian-product, horner, algebra
// Units:   coefficients, x, y and the results are dimensionless
// Example: c = [1 2; 3 4] (1 + 2y + 3x + 4xy), xs = (0, 2), ys = (1, 3) → [3 7; 17 37] (hand-checked)
func PolynomialEvaluateGrid2D(dst, coeffs []float64, rows, cols int, xs, ys []float64) ([]float64, error) {
	if rows < 1 || cols < 1 {
		return nil, fmt.Errorf("%w: grid is %d×%d", ErrEmpty, rows, cols)
	}
	if rows > len(coeffs) || len(coeffs) != rows*cols {
		return nil, fmt.Errorf("%w: len(coeffs) = %d, want %d×%d", ErrDegenerate, len(coeffs), rows, cols)
	}
	if err := checkCoeffs("coeffs", coeffs); err != nil {
		return nil, err
	}
	for i, v := range xs {
		if !reduce.AllFinite(v) {
			return nil, fmt.Errorf("%w: xs[%d] = %v", ErrNotFinite, i, v)
		}
	}
	for j, v := range ys {
		if !reduce.AllFinite(v) {
			return nil, fmt.Errorf("%w: ys[%d] = %v", ErrNotFinite, j, v)
		}
	}
	nx, ny := len(xs), len(ys)
	if cap(dst) >= nx*ny {
		dst = dst[:nx*ny]
	} else {
		dst = make([]float64, nx*ny)
	}
	if nx == 0 || ny == 0 {
		return dst, nil
	}
	rowVal := make([]float64, rows)
	for j, y := range ys {
		for k := range rows {
			row := coeffs[k*cols : k*cols+cols]
			s := 0.0
			for l := cols - 1; l >= 0; l-- {
				s = s*y + row[l]
			}
			rowVal[k] = s
		}
		for i, x := range xs {
			s := 0.0
			for k := rows - 1; k >= 0; k-- {
				s = s*x + rowVal[k]
			}
			if !reduce.AllFinite(s) {
				return nil, fmt.Errorf("%w: p(xs[%d], ys[%d])", ErrOverflow, i, j)
			}
			dst[i*ny+j] = s
		}
	}
	return dst, nil
}
