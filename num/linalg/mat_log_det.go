package linalg

import (
	"fmt"
	"math"

	"github.com/CrimsonCoder42/calcbench/algo/reduce"
)

// MatLogDet returns the sign (±1) and the natural logarithm of |det A| of an n×n matrix from its LU
// factorisation (det A = det P·∏uᵢᵢ, log|det A| = Σ log|uᵢᵢ|), so a determinant far outside the float64
// range (a 1000×1000 matrix with entries near 10) still has a finite answer; det A = sign·exp(logAbsDet).
// Errors: ErrSingular when a pivot is exactly zero (det A = 0); otherwise as LUDecompose.
//
// ID:      NUM-LINALG-174
// Time:    O(n³)   n = order of the matrix
// Space:   O(n²) for the factor
// Source:  Higham, Accuracy and Stability of Numerical Algorithms, 2nd ed., SIAM, 2002, §14.6 "Computing the Determinant" (det from the LU factors)
// Source:  Golub, Van Loan, Matrix Computations, 4th ed., Johns Hopkins University Press, 2013, §3.2 "The LU Factorization"
// Tags:    matrix, n-dimensional, determinant, log-determinant, slogdet, lu, linear-algebra
// Units:   a in unit A; logAbsDet is log of Aⁿ
// Example: [1 2; 3 4] → sign −1, log 2 ≈ 0.693147 (det = −2; hand-checked)
func MatLogDet(a []float64, n int) (sign, logAbsDet float64, err error) {
	f, err := LUDecompose(a, n)
	if err != nil {
		return 0, 0, err
	}
	if f.singular {
		return 0, 0, fmt.Errorf("%w: zero pivot", ErrSingular)
	}
	sign = f.sign
	s, c := 0.0, 0.0
	for i := range n {
		u := f.lu[i*n+i]
		if u < 0 {
			sign = -sign
		}
		s, c = reduce.NeumaierAdd(s, c, math.Log(math.Abs(u)))
	}
	return sign, s + c, nil
}
