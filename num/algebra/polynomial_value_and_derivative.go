package algebra

import (
	"fmt"
	"math"

	"github.com/CrimsonCoder42/calcbench/algo/reduce"
)

// PolynomialValueAndDerivative returns p(x) and its first derivative p′(x) in one Horner pass, where
// p(x) = coeffs[0] + coeffs[1]·x + … + coeffs[n]·xⁿ (ascending order: coeffs[i] multiplies xⁱ).
// Each step is two fused multiply-adds: d ← d·x + s, then s ← s·x + coeffs[i]. No derivative
// coefficients are formed. coeffs is read, never modified. This is the Newton step's numerator and
// denominator, and the slope of the tangent line at x.
// Errors: ErrEmpty when coeffs is empty; ErrNotFinite when x or a coefficient is NaN or ±Inf;
// ErrOverflow when the value or the derivative is not finite. Both results are 0 with any error.
//
// ID:      NUM-ALGEB-002
// Time:    O(n)   n = degree = len(coeffs)−1; 2n fused multiply-adds
// Space:   O(1), no allocation
// Source:  NIST Digital Library of Mathematical Functions, §1.11(i) (extended Horner scheme, the derivative), https://dlmf.nist.gov/1.11 (read 2026-11-12)
// Source:  Higham, Accuracy and Stability of Numerical Algorithms, 2nd ed., SIAM, 2002, §5.2 "Evaluating Derivatives"
// Tags:    polynomial, horner, derivative, newton, algebra
// Units:   coeffs, x and value are dimensionless; derivative is per unit of x
// Example: coeffs = [−1, 2, −6, 2], x = 3 → (5, 20) (hand-checked: p′(x) = 6x² − 12x + 2)
func PolynomialValueAndDerivative(coeffs []float64, x float64) (value, derivative float64, err error) {
	n := len(coeffs) - 1
	if n < 0 {
		return 0, 0, fmt.Errorf("%w: coeffs", ErrEmpty)
	}
	if !reduce.AllFinite(x) {
		return 0, 0, fmt.Errorf("%w: x = %v", ErrNotFinite, x)
	}
	s := coeffs[n]
	if !reduce.AllFinite(s) {
		return 0, 0, fmt.Errorf("%w: coeffs[%d] = %v", ErrNotFinite, n, s)
	}
	d := 0.0
	for i := n - 1; i >= 0; i-- {
		a := coeffs[i]
		if !reduce.AllFinite(a) {
			return 0, 0, fmt.Errorf("%w: coeffs[%d] = %v", ErrNotFinite, i, a)
		}
		d = math.FMA(d, x, s)
		s = math.FMA(s, x, a)
	}
	if !reduce.AllFinite(s) {
		return 0, 0, fmt.Errorf("%w: p(x)", ErrOverflow)
	}
	if !reduce.AllFinite(d) {
		return 0, 0, fmt.Errorf("%w: p′(x)", ErrOverflow)
	}
	return s, d, nil
}
