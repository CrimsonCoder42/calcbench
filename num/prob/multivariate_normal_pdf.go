package prob

import (
	"errors"
	"fmt"
	"math"

	"github.com/CrimsonCoder42/calcbench/algo/reduce"
)

// ErrDimension: the lengths of the vectors and the matrix passed together do not match.
var ErrDimension = errors.New("prob: dimensions do not match")

// mvnStackDim is the largest dimension whose forward-substitution scratch lives on the stack.
const mvnStackDim = 16

// multivariateNormalLog returns ln f(x) for the N(μ, Σ) density with Σ = L·Lᵀ, L the caller's lower-triangular
// Cholesky factor (row-major n×n; the strict upper triangle is not read). With d = x − μ and z the solution of
// L·z = d by forward substitution, (x − μ)ᵀΣ⁻¹(x − μ) = zᵀz and ln det Σ = 2·Σᵢ ln Lᵢᵢ, so
// ln f = −½·(n·ln 2π + ln det Σ + zᵀz).
func multivariateNormalLog(x, mean, cholL []float64) (float64, error) {
	n := len(x)
	if n == 0 {
		return 0, fmt.Errorf("%w: x is empty", ErrDimension)
	}
	if len(mean) != n {
		return 0, fmt.Errorf("%w: len(mean) = %d, len(x) = %d", ErrDimension, len(mean), n)
	}
	if len(cholL) != n*n {
		return 0, fmt.Errorf("%w: len(cholL) = %d, want %d (n×n, row-major)", ErrDimension, len(cholL), n*n)
	}
	for i := range n {
		if !reduce.AllFinite(x[i]) {
			return 0, fmt.Errorf("%w: x[%d] = %v", ErrNotFinite, i, x[i])
		}
		if !reduce.AllFinite(mean[i]) {
			return 0, fmt.Errorf("%w: mean[%d] = %v", ErrNotFinite, i, mean[i])
		}
		for j := 0; j <= i; j++ {
			if v := cholL[i*n+j]; !reduce.AllFinite(v) {
				return 0, fmt.Errorf("%w: cholL[%d][%d] = %v", ErrNotFinite, i, j, v)
			}
		}
		if d := cholL[i*n+i]; d <= 0 {
			return 0, fmt.Errorf("%w: cholL[%d][%d] = %v, want > 0 (Σ positive definite)", ErrOutOfRange, i, i, d)
		}
	}
	var buf [mvnStackDim]float64
	var z []float64
	if n <= mvnStackDim {
		z = buf[:n]
	} else {
		z = make([]float64, n)
	}
	logDet, quad := 0.0, 0.0
	for i := range n {
		s := x[i] - mean[i]
		row := cholL[i*n : i*n+i]
		for j, l := range row {
			s -= l * z[j]
		}
		z[i] = s / cholL[i*n+i]
		quad += z[i] * z[i]
		logDet += math.Log(cholL[i*n+i])
	}
	if !reduce.AllFinite(quad) {
		return 0, fmt.Errorf("%w: quadratic form (x − μ)ᵀΣ⁻¹(x − μ) exceeds the float64 range", ErrOverflow)
	}
	return -0.5*(float64(n)*math.Log(2*math.Pi)+quad) - logDet, nil
}

// MultivariateNormalPDF returns the n-variate normal density
// f(x) = (2π)^(−n/2)·det(Σ)^(−½)·exp(−½·(x − μ)ᵀΣ⁻¹(x − μ)), with Σ given by its lower-triangular Cholesky
// factor L (Σ = L·Lᵀ), row-major n×n, exactly as NUM-LINALG-152 Cholesky returns it; the strict upper
// triangle of cholL is not read. num/prob does not import num/linalg (D-009, D-083), so the caller factors
// Σ once and reuses L for every x. det Σ = Πᵢ Lᵢᵢ²; the quadratic form is zᵀz with L·z = x − μ solved by
// forward substitution. A density below the float64 range underflows to 0 (the true value rounded).
// Errors: ErrDimension when x is empty, len(mean) ≠ len(x) or len(cholL) ≠ n²; ErrNotFinite when an
// entry of x, mean or the lower triangle of cholL is NaN or ±Inf; ErrOutOfRange when a diagonal entry
// of cholL is ≤ 0; ErrOverflow when the quadratic form or the density exceeds the float64 range.
//
// ID:      NUM-PROB-094
// Time:    O(n²)   n = len(x) (forward substitution)
// Space:   O(1) extra for n ≤ 16 (scratch on the stack, no allocation); O(n), one allocation, above
// Source:  Kotz, Balakrishnan, Johnson, Continuous Multivariate Distributions, vol. 1: Models and Applications, 2nd ed., Wiley, 2000, ch. 45 "Multivariate Normal Distributions", §1 (density)
// Source:  Golub, Van Loan, Matrix Computations, 4th ed., Johns Hopkins, 2013, §3.1.1 (forward substitution), §4.2 (Cholesky, det Σ = Π Lᵢᵢ²)
// Tags:    probability, distribution, multivariate-normal, gaussian, pdf, continuous, multivariate, cholesky
// Units:   x and mean: one unit per coordinate; cholL: those units; result: probability per unit volume
// Example: n = 2, μ = 0, Σ = [[4,2],[2,3]] (L = [[2,0],[1,√2]]), x = (1,1) → e^(−3/16)/(2π·√8) ≈ 0.04664
func MultivariateNormalPDF(x, mean, cholL []float64) (float64, error) {
	lf, err := multivariateNormalLog(x, mean, cholL)
	if err != nil {
		return 0, err
	}
	if lf > math.Log(math.MaxFloat64) {
		return 0, fmt.Errorf("%w: density e^%v", ErrOverflow, lf)
	}
	return math.Exp(lf), nil
}
