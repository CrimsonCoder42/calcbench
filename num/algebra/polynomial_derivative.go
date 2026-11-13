package algebra

import (
	"fmt"

	"github.com/CrimsonCoder42/calcbench/algo/reduce"
)

// PolynomialDerivative returns the coefficients of the order-th derivative of p (ascending order):
// c'[i] = c[i+order]·(i+1)(i+2)…(i+order). order = 0 returns a copy; when order ≥ len(a) the result is
// [0].
// Errors: ErrEmpty when a is empty; ErrNotFinite for a NaN or ±Inf coefficient; ErrDegenerate when
// order < 0; ErrOverflow when a coefficient leaves the float64 range.
//
// ID:      NUM-ALGEB-015
// Time:    O(n·order)   n = len(a)
// Space:   O(n) for the result (one allocation), O(1) extra
// Source:  NIST Digital Library of Mathematical Functions, §1.4(iii) "Derivatives" (term-by-term differentiation of xⁿ); Knuth, TAOCP Vol. 2, 3rd ed., 1997, §4.6 "Polynomial Arithmetic"
// Tags:    polynomial, derivative, differentiate, coefficients, calculus, algebra
// Units:   coefficients are dimensionless
// Example: d/dx (1 + 2x + 3x²) → [2, 6]; second derivative → [6] (hand-checked)
func PolynomialDerivative(a []float64, order int) ([]float64, error) {
	if err := checkCoeffs("a", a); err != nil {
		return nil, err
	}
	if order < 0 {
		return nil, fmt.Errorf("%w: order = %d", ErrDegenerate, order)
	}
	if order >= len(a) {
		return []float64{0}, nil
	}
	c := make([]float64, len(a)-order)
	for i := range c {
		v := a[i+order]
		for f := i + 1; f <= i+order; f++ {
			v *= float64(f)
		}
		if !reduce.AllFinite(v) {
			return nil, fmt.Errorf("%w: coefficient %d", ErrOverflow, i)
		}
		c[i] = v
	}
	return c, nil
}
