package algebra

import (
	"fmt"

	"github.com/CrimsonCoder42/calcbench/algo/reduce"
)

// PolynomialEvaluate2D returns p(x, y) = Σᵢ Σⱼ c[i·cols+j]·xⁱ·yʲ for a degreeX = rows−1, degreeY = cols−1
// coefficient grid (row-major) by nested Horner's rule: Horner in y along each row, then Horner in x over
// the row values.
// Errors: ErrEmpty when rows or cols < 1; ErrDegenerate when len(coeffs) ≠ rows·cols; ErrNotFinite for a
// NaN or ±Inf coefficient, x or y; ErrOverflow when the value leaves the float64 range.
//
// ID:      NUM-ALGEB-021
// Time:    O(rows·cols)
// Space:   O(1)
// Source:  Knuth, The Art of Computer Programming, Vol. 2: Seminumerical Algorithms, 3rd ed., Addison-Wesley, 1997, §4.6.4 "Evaluation of Polynomials" (Horner's rule; polynomials in several variables)
// Tags:    polynomial, evaluate, bivariate, two-dimensional, horner, algebra
// Units:   coefficients, x, y and the result are dimensionless
// Example: c = [1 2; 3 4] (1 + 2y + 3x + 4xy) at x = 2, y = 3 → 1 + 6 + 6 + 24 = 37 (hand-checked)
func PolynomialEvaluate2D(coeffs []float64, rows, cols int, x, y float64) (float64, error) {
	if rows < 1 || cols < 1 {
		return 0, fmt.Errorf("%w: grid is %d×%d", ErrEmpty, rows, cols)
	}
	if rows > len(coeffs) || len(coeffs) != rows*cols {
		return 0, fmt.Errorf("%w: len(coeffs) = %d, want %d×%d", ErrDegenerate, len(coeffs), rows, cols)
	}
	if !reduce.AllFinite(x) || !reduce.AllFinite(y) {
		return 0, fmt.Errorf("%w: x = %v, y = %v", ErrNotFinite, x, y)
	}
	if err := checkCoeffs("coeffs", coeffs); err != nil {
		return 0, err
	}
	outer := 0.0
	for i := rows - 1; i >= 0; i-- {
		row := coeffs[i*cols : i*cols+cols]
		inner := 0.0
		for j := cols - 1; j >= 0; j-- {
			inner = inner*y + row[j]
		}
		outer = outer*x + inner
	}
	if !reduce.AllFinite(outer) {
		return 0, fmt.Errorf("%w: p(x, y)", ErrOverflow)
	}
	return outer, nil
}
