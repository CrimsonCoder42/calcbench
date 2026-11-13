package algebra

import (
	"fmt"

	"github.com/CrimsonCoder42/calcbench/algo/reduce"
)

// PolynomialCompanion returns the n×n companion matrix of p, n = len(a)−1 = deg p, row-major in dst
// (reused when its capacity suffices): ones on the subdiagonal and last column −a[i]/a[n], so its
// characteristic polynomial is p/a[n] and its eigenvalues are the roots of p.
// Errors: ErrEmpty when len(a) < 2 (a constant has no companion); ErrNotFinite for a NaN or ±Inf
// coefficient; ErrDegenerate when the leading coefficient a[n] is 0; ErrOverflow when −a[i]/a[n] leaves
// the float64 range.
//
// ID:      NUM-ALGEB-019
// Time:    O(n²)   n = deg p (writing the matrix)
// Space:   O(n²) for the result, one allocation only if dst lacks capacity; O(1) extra
// Source:  Edelman, Murakami, "Polynomial Roots from Companion Matrix Eigenvalues", Mathematics of Computation, 1995, 64(210):763–776, doi:10.1090/S0025-5718-1995-1262279-2
// Source:  Golub, Van Loan, Matrix Computations, 4th ed., Johns Hopkins University Press, 2013, §7.4.6 "Companion Matrix Form"
// Tags:    polynomial, companion-matrix, roots, eigenvalues, algebra
// Units:   coefficients are dimensionless
// Example: x² − 3x + 2 = [2, −3, 1] → [0 −2; 1 3] (eigenvalues 1 and 2; hand-checked)
func PolynomialCompanion(dst, a []float64) ([]float64, error) {
	if len(a) < 2 {
		return nil, fmt.Errorf("%w: len(a) = %d, need degree ≥ 1", ErrEmpty, len(a))
	}
	if err := checkCoeffs("a", a); err != nil {
		return nil, err
	}
	n := len(a) - 1
	lead := a[n]
	if lead == 0 {
		return nil, fmt.Errorf("%w: leading coefficient a[%d] is 0", ErrDegenerate, n)
	}
	if cap(dst) >= n*n {
		dst = dst[:n*n]
		clear(dst)
	} else {
		dst = make([]float64, n*n)
	}
	for i := 1; i < n; i++ {
		dst[i*n+i-1] = 1
	}
	for i := range n {
		v := -a[i] / lead
		if !reduce.AllFinite(v) {
			return nil, fmt.Errorf("%w: entry (%d, %d)", ErrOverflow, i, n-1)
		}
		dst[i*n+n-1] = v
	}
	return dst, nil
}
