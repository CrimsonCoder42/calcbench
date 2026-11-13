package algebra

import (
	"fmt"

	"github.com/CrimsonCoder42/calcbench/algo/reduce"
)

// PolynomialDivide returns the quotient and remainder of p ÷ q (ascending coefficients): p = q·quotient +
// remainder with deg remainder < deg q. The quotient has len(a) − len(b) + 1 coefficients ([0] when
// len(a) < len(b)); the remainder has len(b) − 1 coefficients ([0] when q is a constant).
// Errors: ErrEmpty when a or b is empty; ErrNotFinite for a NaN or ±Inf coefficient; ErrDegenerate when
// the leading coefficient b[len(b)−1] is 0 (trim it first); ErrOverflow when a coefficient leaves the
// float64 range.
//
// ID:      NUM-ALGEB-013
// Time:    O((n − m + 1)·m)   n = len(a), m = len(b)
// Space:   O(n) for the two results (one allocation; the remainder is a prefix of the working copy), O(1) extra
// Source:  Knuth, The Art of Computer Programming, Vol. 2: Seminumerical Algorithms, 3rd ed., Addison-Wesley, 1997, §4.6.1 "Division of Polynomials", Algorithm D
// Tags:    polynomial, divide, long-division, quotient, remainder, algebra
// Units:   coefficients are dimensionless
// Example: (x³ − 2x² − 4) ÷ (x − 3) → quotient x² + x + 3 = [3, 1, 1], remainder 5 = [5] (hand-checked)
func PolynomialDivide(a, b []float64) (quotient, remainder []float64, err error) {
	if err := checkCoeffs("a", a); err != nil {
		return nil, nil, err
	}
	if err := checkCoeffs("b", b); err != nil {
		return nil, nil, err
	}
	m := len(b)
	lead := b[m-1]
	if lead == 0 {
		return nil, nil, fmt.Errorf("%w: leading coefficient b[%d] is 0", ErrDegenerate, m-1)
	}
	n := len(a)
	if n < m {
		buf := make([]float64, 1+max(m-1, 1))
		r := buf[1:]
		copy(r, a)
		return buf[:1], r[:max(m-1, 1)], nil
	}
	buf := make([]float64, n+n-m+1)
	r, q := buf[:n], buf[n:]
	copy(r, a)
	for k := n - m; k >= 0; k-- {
		c := r[k+m-1] / lead
		if !reduce.AllFinite(c) {
			return nil, nil, fmt.Errorf("%w: quotient coefficient %d", ErrOverflow, k)
		}
		q[k] = c
		for j := range m {
			r[k+j] -= c * b[j]
			if !reduce.AllFinite(r[k+j]) {
				return nil, nil, fmt.Errorf("%w: remainder coefficient %d", ErrOverflow, k+j)
			}
		}
	}
	if m == 1 {
		return q, []float64{0}, nil
	}
	return q, r[:m-1], nil
}
