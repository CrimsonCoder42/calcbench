package algebra

import (
	"fmt"

	"github.com/CrimsonCoder42/calcbench/algo/reduce"
)

// PolynomialIntegral returns the coefficients of the antiderivative P of p with P(0) = constant
// (ascending order): P[0] = constant, P[i+1] = c[i]/(i+1). Integrate again for a higher order.
// Errors: ErrEmpty when a is empty; ErrNotFinite for a NaN or ±Inf coefficient or constant; ErrOverflow
// never in practice (division by i+1 ≥ 1 cannot grow a value).
//
// ID:      NUM-ALGEB-016
// Time:    O(n)   n = len(a)
// Space:   O(n) for the result (one allocation), O(1) extra
// Source:  NIST Digital Library of Mathematical Functions, §1.4(iv) "Indefinite Integrals" (∫xⁿ dx = xⁿ⁺¹/(n+1)); Knuth, TAOCP Vol. 2, 3rd ed., 1997, §4.6
// Tags:    polynomial, integral, antiderivative, integrate, coefficients, calculus, algebra
// Units:   coefficients are dimensionless
// Example: ∫(2 + 6x) dx with P(0) = 1 → [1, 2, 3] (hand-checked)
func PolynomialIntegral(a []float64, constant float64) ([]float64, error) {
	if err := checkCoeffs("a", a); err != nil {
		return nil, err
	}
	if !reduce.AllFinite(constant) {
		return nil, fmt.Errorf("%w: constant = %v", ErrNotFinite, constant)
	}
	c := make([]float64, len(a)+1)
	c[0] = constant
	for i, v := range a {
		c[i+1] = v / float64(i+1)
	}
	return c, nil
}
