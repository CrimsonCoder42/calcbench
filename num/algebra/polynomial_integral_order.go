package algebra

import (
	"fmt"

	"github.com/CrimsonCoder42/calcbench/algo/reduce"
)

// PolynomialIntegralOrder returns the coefficients (ascending order) of the order-th repeated integral of
// p, as numpy's polyint(c, m, k, lbnd): order antiderivative steps are taken, and after step i the
// integration constant is chosen so that the new polynomial equals constants[i] at x = lowerBound
// (constants shorter than order are padded with zeros; nil means all zero). order = 0 returns a copy of
// a. The result has len(a) + order coefficients. PolynomialIntegral (NUM-ALGEB-016) is the one-step case
// with lowerBound = 0.
// Errors: ErrEmpty when a is empty; ErrDegenerate when order < 0 or len(constants) > order; ErrNotFinite
// for a NaN or ±Inf coefficient, constant or lowerBound; ErrOverflow when a coefficient or a value at
// lowerBound leaves the float64 range.
//
// ID:      NUM-ALGEB-028
// Time:    O(m·(n + m))   n = len(a), m = order (each step divides n + i coefficients and evaluates P(lowerBound) by Horner)
// Space:   O(n + m) for the result (one allocation), O(1) extra
// Source:  NIST Digital Library of Mathematical Functions, §1.4(iv) "Indefinite Integrals" (∫xⁿ dx = xⁿ⁺¹/(n+1)); Knuth, TAOCP Vol. 2, 3rd ed., 1997, §4.6.4 (Horner's rule for P(lowerBound))
// Tags:    polynomial, integral, antiderivative, integrate, repeated, order, lower-bound, coefficients, calculus, algebra
// Units:   coefficients, lowerBound and constants are dimensionless
// Example: ∫∫ 6 dx dx with lowerBound 1 and constants (0, 0) → P₁ = −6 + 6x, P₂ = 3 − 6x + 3x² = 3(x − 1)² (hand-checked)
func PolynomialIntegralOrder(a []float64, order int, lowerBound float64, constants []float64) ([]float64, error) {
	if err := checkCoeffs("a", a); err != nil {
		return nil, err
	}
	if order < 0 || len(constants) > order {
		return nil, fmt.Errorf("%w: order = %d, len(constants) = %d", ErrDegenerate, order, len(constants))
	}
	if !reduce.AllFinite(lowerBound) {
		return nil, fmt.Errorf("%w: lowerBound = %v", ErrNotFinite, lowerBound)
	}
	for i, k := range constants {
		if !reduce.AllFinite(k) {
			return nil, fmt.Errorf("%w: constants[%d] = %v", ErrNotFinite, i, k)
		}
	}
	n := len(a)
	c := make([]float64, n+order)
	copy(c, a)
	for step := range order {
		size := n + step // current length; the integral has size+1 coefficients
		for i := size - 1; i >= 0; i-- {
			c[i+1] = c[i] / float64(i+1)
		}
		c[0] = 0
		atLower := 0.0
		for i := size; i >= 0; i-- {
			atLower = atLower*lowerBound + c[i]
		}
		k := 0.0
		if step < len(constants) {
			k = constants[step]
		}
		c[0] = k - atLower
		if !reduce.AllFinite(atLower) || !reduce.AllFinite(c[0]) {
			return nil, fmt.Errorf("%w: integral %d at lowerBound", ErrOverflow, step+1)
		}
	}
	return c, nil
}
