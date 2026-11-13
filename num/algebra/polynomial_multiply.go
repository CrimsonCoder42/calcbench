package algebra

import (
	"fmt"
	"math"

	"github.com/CrimsonCoder42/calcbench/algo/reduce"
)

// PolynomialMultiply returns the coefficients of the product p·q, where p has coefficients a and q has
// coefficients b, all in ascending order (a[i] multiplies xⁱ). The result has len(a)+len(b)−1
// coefficients: c[k] = Σ a[i]·b[k−i] over 0 ≤ i < len(a), 0 ≤ k−i < len(b). Each term is added with a
// fused multiply-add (one rounding per term). Trailing zeros are kept, never trimmed. The product of two
// binomials, (a₁x + a₀)(b₁x + b₀), is the case len(a) = len(b) = 2. a and b are read, never modified.
// Errors: ErrEmpty when a or b is empty; ErrNotFinite when a coefficient is NaN or ±Inf; ErrOverflow
// when a product coefficient is not finite. The result is nil with any error.
//
// ID:      NUM-ALGEB-003
// Time:    O(n·m)   n = len(a), m = len(b); n·m fused multiply-adds
// Space:   O(n+m) for the result (one allocation), O(1) extra
// Source:  Cormen, Leiserson, Rivest, Stein, Introduction to Algorithms, 3rd ed., MIT Press, 2009, §30.1, eqs. (30.1)–(30.2) (product in coefficient form)
// Tags:    polynomial, multiply, binomial, foil, convolution, algebra
// Units:   coefficients are dimensionless
// Example: a = [3, 2] (2x + 3), b = [−4, 1] (x − 4) → [−12, −5, 2] (2x² − 5x − 12, hand-checked)
func PolynomialMultiply(a, b []float64) ([]float64, error) {
	if len(a) == 0 {
		return nil, fmt.Errorf("%w: a", ErrEmpty)
	}
	if len(b) == 0 {
		return nil, fmt.Errorf("%w: b", ErrEmpty)
	}
	for i, v := range a {
		if !reduce.AllFinite(v) {
			return nil, fmt.Errorf("%w: a[%d] = %v", ErrNotFinite, i, v)
		}
	}
	for j, v := range b {
		if !reduce.AllFinite(v) {
			return nil, fmt.Errorf("%w: b[%d] = %v", ErrNotFinite, j, v)
		}
	}
	c := make([]float64, len(a)+len(b)-1)
	for i, ai := range a {
		row := c[i : i+len(b)]
		for j, bj := range b {
			row[j] = math.FMA(ai, bj, row[j])
		}
	}
	for k, v := range c {
		if !reduce.AllFinite(v) {
			return nil, fmt.Errorf("%w: coefficient %d of the product", ErrOverflow, k)
		}
	}
	return c, nil
}
