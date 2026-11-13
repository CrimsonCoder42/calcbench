package linalg

import (
	"fmt"

	"github.com/CrimsonCoder42/calcbench/num/algebra"
)

// PolynomialRootsCompanion returns the deg p complex roots of the real polynomial p (ascending
// coefficients, numpy's polynomial.polyroots) as the eigenvalues (MatEigenvalues, NUM-LINALG-169) of its
// companion matrix (algebra.PolynomialCompanion, NUM-ALGEB-019), in rootsDst. This is backward stable for
// the companion matrix (Edelman–Murakami); NUM-ALGEB-010 PolynomialRoots (Aberth–Ehrlich) is the
// alternative. Leading zero coefficients are not trimmed (use algebra.PolynomialTrim).
// Errors: algebra.ErrEmpty when len(coeffs) < 2; algebra.ErrNotFinite; algebra.ErrDegenerate for a zero
// leading coefficient; algebra.ErrOverflow; ErrNoConvergence from the QR iteration.
//
// ID:      NUM-LINALG-172
// Time:    O(n³)   n = deg p
// Space:   O(n²) for the companion matrix and the eigen-solver's working copy; O(n) for the result
// Source:  Edelman, Murakami, "Polynomial Roots from Companion Matrix Eigenvalues", Mathematics of Computation, 1995, 64(210):763–776, doi:10.1090/S0025-5718-1995-1262279-2
// Source:  Golub, Van Loan, Matrix Computations, 4th ed., Johns Hopkins University Press, 2013, §7.4.6 "Companion Matrix Form" and §7.5
// Tags:    polynomial, roots, companion-matrix, eigenvalues, zeros, linear-algebra
// Units:   coefficients and roots are dimensionless
// Example: x² − 3x + 2 = [2, −3, 1] → roots 2, 1; x² + 1 = [1, 0, 1] → ±i (hand-checked)
func PolynomialRootsCompanion(rootsDst []complex128, coeffs []float64) ([]complex128, error) {
	c, err := algebra.PolynomialCompanion(nil, coeffs)
	if err != nil {
		return nil, err
	}
	roots, err := MatEigenvalues(rootsDst, c, len(coeffs)-1)
	if err != nil {
		return nil, fmt.Errorf("companion eigenvalues: %w", err)
	}
	return roots, nil
}
