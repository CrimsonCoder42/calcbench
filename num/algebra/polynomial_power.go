package algebra

import "fmt"

// maxPolynomialPower caps the exponent of PolynomialPower (a policy: the result has k·(len(a)−1)+1
// coefficients, so a larger k is an input error).
const maxPolynomialPower = 1 << 20

// PolynomialPower returns the coefficients of pᵏ (ascending order) for an integer k ≥ 0 by binary
// powering with PolynomialMultiply (NUM-ALGEB-003); p⁰ = [1].
// Errors: ErrEmpty when a is empty; ErrNotFinite for a NaN or ±Inf coefficient; ErrDegenerate when k < 0
// or k > 2²⁰ (a policy cap); ErrOverflow when a coefficient leaves the float64 range.
//
// ID:      NUM-ALGEB-014
// Time:    O((k·n)²)   n = len(a); the last squaring dominates (schoolbook products)
// Space:   O(k·n) for the result and the intermediate powers
// Source:  Knuth, The Art of Computer Programming, Vol. 2: Seminumerical Algorithms, 3rd ed., Addison-Wesley, 1997, §4.6.3 "Evaluation of Powers", Algorithm A (right-to-left binary method)
// Tags:    polynomial, power, binary-powering, binomial-expansion, algebra
// Units:   coefficients are dimensionless
// Example: (1 + x)³ → [1, 3, 3, 1] (hand-checked, binomial coefficients)
func PolynomialPower(a []float64, k int) ([]float64, error) {
	if err := checkCoeffs("a", a); err != nil {
		return nil, err
	}
	if k < 0 || k > maxPolynomialPower {
		return nil, fmt.Errorf("%w: k = %d, want 0 ≤ k ≤ 2^20", ErrDegenerate, k)
	}
	result := []float64{1}
	base := append([]float64(nil), a...)
	var err error
	for k > 0 {
		if k&1 == 1 {
			if result, err = PolynomialMultiply(result, base); err != nil {
				return nil, err
			}
		}
		k >>= 1
		if k > 0 {
			if base, err = PolynomialMultiply(base, base); err != nil {
				return nil, err
			}
		}
	}
	return result, nil
}
