package algebra

import (
	"errors"
	"fmt"
	"math"

	"github.com/CrimsonCoder42/calcbench/algo/reduce"
)

// Errors returned by this package. They are wrapped with the offending parameter; test with errors.Is.
var (
	// ErrEmpty: a coefficient slice has too few elements.
	ErrEmpty = errors.New("algebra: empty input")
	// ErrNotFinite: an input is NaN or ±Inf.
	ErrNotFinite = errors.New("algebra: input is NaN or infinite")
	// ErrOverflow: the inputs are finite but a result (or an intermediate the result needs) is not.
	ErrOverflow = errors.New("algebra: result outside the float64 range")
	// ErrDegenerate: a leading coefficient or a slope is zero, so the equation is of lower degree.
	ErrDegenerate = errors.New("algebra: degenerate equation")
	// ErrNoConvergence: an iteration did not meet its stopping rule within its sweep limit.
	ErrNoConvergence = errors.New("algebra: iteration did not converge")
)

// unitRoundoff is u = 2⁻⁵³, the unit roundoff of float64 with round-to-nearest.
const unitRoundoff = 0x1p-53

// PolynomialEvaluate returns p(x) = coeffs[0] + coeffs[1]·x + … + coeffs[n]·xⁿ, evaluated by Horner's
// rule with one fused multiply-add (one rounding) per step. Coefficients are in ascending order:
// coeffs[i] multiplies xⁱ. Trailing zero coefficients are allowed. coeffs is read, never modified.
// Horner's rule is backward stable: the result is the exact value of a polynomial whose coefficients
// differ from coeffs by a relative amount of at most γ₂ₙ = 2nu/(1−2nu) (Higham §5.1).
// Errors: ErrEmpty when coeffs is empty; ErrNotFinite when x or a coefficient is NaN or ±Inf;
// ErrOverflow when the value is not finite. The value is 0 with any error.
//
// ID:      NUM-ALGEB-001
// Time:    O(n)   n = degree = len(coeffs)−1; n fused multiply-adds
// Space:   O(1), no allocation
// Source:  Knuth, The Art of Computer Programming, Vol. 2: Seminumerical Algorithms, 3rd ed., Addison-Wesley, 1997, §4.6.4, pp. 486–488 (Horner's rule)
// Source:  Higham, Accuracy and Stability of Numerical Algorithms, 2nd ed., SIAM, 2002, §5.1 "Horner's Method", p. 94
// Source:  NIST Digital Library of Mathematical Functions, §1.11(i) (Horner scheme), https://dlmf.nist.gov/1.11 (read 2026-11-12)
// Tags:    polynomial, horner, evaluate, algebra
// Units:   coeffs, x and the result are dimensionless (a caller that attaches units keeps them consistent)
// Example: coeffs = [−1, 2, −6, 2] (2x³ − 6x² + 2x − 1), x = 3 → 5 (hand-checked)
func PolynomialEvaluate(coeffs []float64, x float64) (float64, error) {
	n := len(coeffs) - 1
	if n < 0 {
		return 0, fmt.Errorf("%w: coeffs", ErrEmpty)
	}
	if !reduce.AllFinite(x) {
		return 0, fmt.Errorf("%w: x = %v", ErrNotFinite, x)
	}
	s := coeffs[n]
	if !reduce.AllFinite(s) {
		return 0, fmt.Errorf("%w: coeffs[%d] = %v", ErrNotFinite, n, s)
	}
	for i := n - 1; i >= 0; i-- {
		a := coeffs[i]
		if !reduce.AllFinite(a) {
			return 0, fmt.Errorf("%w: coeffs[%d] = %v", ErrNotFinite, i, a)
		}
		s = math.FMA(s, x, a)
	}
	if !reduce.AllFinite(s) {
		return 0, fmt.Errorf("%w: p(x)", ErrOverflow)
	}
	return s, nil
}
