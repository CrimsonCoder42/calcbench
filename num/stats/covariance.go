package stats

import (
	"fmt"
	"math"

	"github.com/CrimsonCoder42/calcbench/algo/reduce"
)

// SampleCovariance returns the sample covariance of the paired samples x and y,
// s_xy = Σ(x − x̄)(y − ȳ)/(n − 1), by the corrected two-pass formula with compensated sums (as
// PearsonCorrelation). x and y are read, never modified.
// Errors: ErrEmpty when x is empty; ErrLengthMismatch when len(y) ≠ len(x); ErrTooFew when there is one
// pair; ErrNotFinite when a value is NaN or ±Inf; ErrOverflow when the result is beyond float64.
//
// ID:      NUM-STATS-038
// Time:    O(n)   n = len(x); three passes (means, then the centred sums)
// Space:   O(1) extra, no allocation
// Source:  NIST/SEMATECH, e-Handbook of Statistical Methods, §6.5.4.1 "Mean Vector and Covariance Matrix", https://www.itl.nist.gov/div898/handbook/pmc/section5/pmc541.htm (read 2026-11-13) (s_jk = Σ(x_ij − x̄_j)(x_ik − x̄_k)/(n − 1))
// Tags:    covariance, sample, paired, association, descriptive
// Units:   x in unit X, y in unit Y; result in X·Y
// Example: e-Handbook §6.5.4.1, X1 = 4.0, 4.2, 3.9, 4.3, 4.1 and X2 = 2.0, 2.1, 2.0, 2.1, 2.2 → 0.0075 (as printed)
func SampleCovariance(x, y []float64) (float64, error) {
	if err := checkPaired(x, y); err != nil {
		return 0, err
	}
	return covarianceOf(x, y, float64(len(x)-1))
}

// PopulationCovariance returns the population covariance of the paired values x and y,
// σ_xy = Σ(x − μ_x)(y − μ_y)/N, by the corrected two-pass formula with compensated sums. One pair gives 0.
// x and y are read, never modified.
// Errors: ErrEmpty when x is empty; ErrLengthMismatch when len(y) ≠ len(x); ErrNotFinite when a value
// is NaN or ±Inf; ErrOverflow when the result is beyond float64.
//
// ID:      NUM-STATS-039
// Time:    O(n)   n = len(x); three passes (means, then the centred sums)
// Space:   O(1) extra, no allocation
// Source:  NIST/SEMATECH, e-Handbook of Statistical Methods, §6.5.4.1 "Mean Vector and Covariance Matrix", https://www.itl.nist.gov/div898/handbook/pmc/section5/pmc541.htm (read 2026-11-13) (the covariance with divisor N for the whole population, as §1.3.5.6 for the variance)
// Tags:    covariance, population, paired, association, descriptive
// Units:   x in unit X, y in unit Y; result in X·Y
// Example: x = 1, 2, 3, 4 and y = 2, 1, 4, 3 → 0.75 (Σ cross-products 3, N = 4; hand-checked)
func PopulationCovariance(x, y []float64) (float64, error) {
	if len(x) == 1 && len(y) == 1 {
		if err := checkFinite("x", x); err != nil {
			return 0, err
		}
		if err := checkFinite("y", y); err != nil {
			return 0, err
		}
		return 0, nil
	}
	if err := checkPaired(x, y); err != nil {
		return 0, err
	}
	return covarianceOf(x, y, float64(len(x)))
}

// covarianceOf returns Σ(x − x̄)(y − ȳ)/divisor for checked x and y, or ErrOverflow.
func covarianceOf(x, y []float64, divisor float64) (float64, error) {
	_, _, sxy, ex, ey := centredProducts(x, y, meanOf(x), meanOf(y))
	c := math.Ldexp(sxy/divisor, ex+ey)
	if !reduce.AllFinite(c) {
		return 0, fmt.Errorf("%w: covariance", ErrOverflow)
	}
	return c, nil
}

// CovarianceMatrix returns the k×k sample covariance matrix of k variables observed together:
// variables[i][t] is variable i at observation t, and entry (i, j) of the result, at index i·k + j
// (row-major, the layout num/linalg uses), is SampleCovariance(variables[i], variables[j]). The matrix is
// symmetric: each pair is computed once. variables is read, never modified.
// Errors: ErrEmpty when there are no variables or no observations; ErrLengthMismatch when the variables
// have different lengths; ErrTooFew when there is one observation; ErrNotFinite when a value is NaN or
// ±Inf; ErrOverflow when an entry is beyond float64.
//
// ID:      NUM-STATS-040
// Time:    O(k²·n)   k = len(variables), n = observations per variable; three passes per pair
// Space:   O(k²) for the result (one allocation); O(k) extra for the means (one allocation)
// Source:  NIST/SEMATECH, e-Handbook of Statistical Methods, §6.5.4.1 "Mean Vector and Covariance Matrix", https://www.itl.nist.gov/div898/handbook/pmc/section5/pmc541.htm (read 2026-11-13) (the variance-covariance matrix S and its worked example)
// Tags:    covariance, covariance-matrix, multivariate, sample, descriptive
// Units:   variable i in unit Uᵢ; entry (i, j) in Uᵢ·Uⱼ
// Example: e-Handbook §6.5.4.1, five observations of three variables → S = [0.025 0.0075 0.00175; 0.0075 0.007 0.00135; 0.00175 0.00135 0.00043] (as printed)
func CovarianceMatrix(variables [][]float64) ([]float64, error) {
	k, err := checkVariables(variables)
	if err != nil {
		return nil, err
	}
	means := make([]float64, k)
	for i, v := range variables {
		means[i] = meanOf(v)
	}
	div := float64(len(variables[0]) - 1)
	out := make([]float64, k*k)
	for i := range k {
		for j := i; j < k; j++ {
			_, _, sxy, ex, ey := centredProducts(variables[i], variables[j], means[i], means[j])
			c := math.Ldexp(sxy/div, ex+ey)
			if !reduce.AllFinite(c) {
				return nil, fmt.Errorf("%w: covariance (%d, %d)", ErrOverflow, i, j)
			}
			out[i*k+j], out[j*k+i] = c, c
		}
	}
	return out, nil
}

// checkVariables checks k ≥ 1 variables of one length n ≥ 2, all finite, and returns k.
func checkVariables(variables [][]float64) (int, error) {
	k := len(variables)
	if k == 0 {
		return 0, fmt.Errorf("%w: variables", ErrEmpty)
	}
	n := len(variables[0])
	if n == 0 {
		return 0, fmt.Errorf("%w: variables[0]", ErrEmpty)
	}
	for i, v := range variables {
		if len(v) != n {
			return 0, fmt.Errorf("%w: len(variables[%d]) = %d, len(variables[0]) = %d", ErrLengthMismatch, i, len(v), n)
		}
	}
	if n < 2 {
		return 0, fmt.Errorf("%w: %d observation, want ≥ 2", ErrTooFew, n)
	}
	for i, v := range variables {
		if err := checkFinite(fmt.Sprintf("variables[%d]", i), v); err != nil {
			return 0, err
		}
	}
	return k, nil
}
