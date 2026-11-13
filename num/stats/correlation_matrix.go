package stats

import (
	"fmt"
	"math"
)

// CorrelationMatrix returns the k×k Pearson correlation matrix of k variables observed together:
// variables[i][t] is variable i at observation t, and entry (i, j), at index i·k + j (row-major, the layout
// num/linalg uses), is PearsonCorrelation(variables[i], variables[j]); the diagonal is exactly 1 and each
// entry is clamped to [−1, 1] against rounding. variables is read, never modified.
// Errors: ErrEmpty when there are no variables or no observations; ErrLengthMismatch when the variables
// have different lengths; ErrTooFew when there is one observation; ErrNotFinite when a value is NaN or
// ±Inf; ErrZeroVariance when every value of some variable is equal.
//
// ID:      NUM-STATS-041
// Time:    O(k²·n)   k = len(variables), n = observations per variable; three passes per pair
// Space:   O(k²) for the result (one allocation); O(k) extra for the means (one allocation)
// Source:  Press, Teukolsky, Vetterling, Flannery, Numerical Recipes, 3rd ed., Cambridge, 2007, §14.5, eq. 14.5.1 (linear correlation coefficient r)
// Source:  NIST, Dataplot Reference Manual, "CORRELATION MATRIX" and "CORRELATION", https://www.itl.nist.gov/div898/software/dataplot/refman2/auxillar/correlat.htm (read 2026-11-13)
// Tags:    correlation, correlation-matrix, pearson, multivariate, descriptive
// Units:   variables in any units; result dimensionless in [−1, 1]
// Example: x = 1, 2, 3, 4; y = 2, 1, 4, 3; z = 4, 3, 2, 1 → r_xy = 0.6, r_xz = −1, r_yz = −0.6 (hand-checked)
func CorrelationMatrix(variables [][]float64) ([]float64, error) {
	k, err := checkVariables(variables)
	if err != nil {
		return nil, err
	}
	means := make([]float64, k)
	for i, v := range variables {
		means[i] = meanOf(v)
	}
	out := make([]float64, k*k)
	for i := range k {
		out[i*k+i] = 1
		for j := i; j < k; j++ {
			sxx, syy, sxy, _, _ := centredProducts(variables[i], variables[j], means[i], means[j])
			if sxx == 0 {
				return nil, fmt.Errorf("%w: every variables[%d] equals %v", ErrZeroVariance, i, variables[i][0])
			}
			if syy == 0 {
				return nil, fmt.Errorf("%w: every variables[%d] equals %v", ErrZeroVariance, j, variables[j][0])
			}
			if j == i {
				continue
			}
			r := min(max(sxy/(math.Sqrt(sxx)*math.Sqrt(syy)), -1), 1)
			out[i*k+j], out[j*k+i] = r, r
		}
	}
	return out, nil
}
