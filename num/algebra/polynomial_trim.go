package algebra

import (
	"fmt"
	"math"

	"github.com/CrimsonCoder42/calcbench/algo/reduce"
)

// PolynomialTrim returns a with its trailing (highest-degree) coefficients of magnitude ≤ tol removed,
// keeping at least one coefficient (the zero polynomial trims to [a[0]] when |a[0]| ≤ tol, as given). The
// result is a subslice of a (no copy, no allocation). tol = 0 removes exact zeros only.
// Errors: ErrEmpty when a is empty; ErrNotFinite for a NaN or ±Inf coefficient or tol; ErrDegenerate when
// tol < 0.
//
// ID:      NUM-ALGEB-020
// Time:    O(n)   n = len(a)
// Space:   O(1); the result shares a's storage
// Source:  Knuth, The Art of Computer Programming, Vol. 2: Seminumerical Algorithms, 3rd ed., Addison-Wesley, 1997, §4.6 (the degree is the index of the highest nonzero coefficient); tolerance policy per calcbench .claude/rules/numeric.md (caller-supplied)
// Tags:    polynomial, trim, degree, coefficients, algebra
// Units:   coefficients and tol are dimensionless
// Example: [1, 2, 0, 1e−20] with tol 1e−15 → [1, 2] (hand-checked)
func PolynomialTrim(a []float64, tol float64) ([]float64, error) {
	if err := checkCoeffs("a", a); err != nil {
		return nil, err
	}
	if !reduce.AllFinite(tol) {
		return nil, fmt.Errorf("%w: tol = %v", ErrNotFinite, tol)
	}
	if tol < 0 {
		return nil, fmt.Errorf("%w: tol = %v", ErrDegenerate, tol)
	}
	n := len(a)
	for n > 1 && math.Abs(a[n-1]) <= tol {
		n--
	}
	return a[:n:n], nil
}
