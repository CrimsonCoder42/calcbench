package linalg

import (
	"fmt"
	"math"

	"github.com/CrimsonCoder42/calcbench/algo/reduce"
)

// CovarianceToCorrelation returns the correlation matrix R = D^(−1/2)·Σ·D^(−1/2), rᵢⱼ = σᵢⱼ/(σᵢσⱼ) with
// σᵢ = √σᵢᵢ, of an n×n covariance matrix Σ, in corrDst (the diagonal is exactly 1); stdDst receives the
// standard deviations σᵢ (length n), which CorrelationToCovariance needs to go back. Symmetry is not
// checked (MatIsSymmetric does that); entries are computed as σᵢⱼ/σᵢ/σⱼ so no product overflows.
// Errors: ErrDimension when n < 1 or len(cov) ≠ n²; ErrNotFinite for a NaN or ±Inf entry; ErrOutOfRange when
// a variance σᵢᵢ is ≤ 0; ErrOverflow when an entry σᵢⱼ/(σᵢσⱼ) leaves the float64 range (cov is then not a
// covariance matrix).
//
// ID:      NUM-LINALG-167
// Time:    O(n²)
// Space:   O(1) extra; O(n²) and O(n) for the results, one allocation each only if a dst lacks capacity
// Source:  NIST/SEMATECH e-Handbook of Statistical Methods, §6.5.4.1 "Mean Vector and Covariance Matrix", https://www.itl.nist.gov/div898/handbook/pmc/section5/pmc541.htm (read 2026-11-12) (the variance-covariance matrix S)
// Source:  NIST, Dataplot Reference Manual, "CORRELATION", https://www.itl.nist.gov/div898/software/dataplot/refman2/auxillar/correlat.htm (read 2026-11-12) (correlation = covariance/(s_x·s_y))
// Tags:    matrix, n-dimensional, covariance, correlation, statistics, scaling, linear-algebra
// Units:   cov in X·Y units; the correlation dimensionless; std in X
// Example: Σ = [4 2; 2 9] → σ = (2, 3), R = [1 1/3; 1/3 1] (hand-checked)
func CovarianceToCorrelation(corrDst, stdDst, cov []float64, n int) ([]float64, []float64, error) {
	if err := checkSquare(cov, n); err != nil {
		return nil, nil, err
	}
	stdDst = grow(stdDst, n)
	for i := range n {
		v := cov[i*n+i]
		if !(v > 0) {
			return nil, nil, fmt.Errorf("%w: variance cov[%d,%d] = %g, want > 0", ErrOutOfRange, i, i, v)
		}
		stdDst[i] = math.Sqrt(v)
	}
	corrDst = grow(corrDst, n*n)
	for i := range n {
		for j := range n {
			if i == j {
				corrDst[i*n+j] = 1
				continue
			}
			corrDst[i*n+j] = cov[i*n+j] / stdDst[i] / stdDst[j]
		}
	}
	if i := firstNotFinite(corrDst); i >= 0 {
		return nil, nil, fmt.Errorf("%w: correlation entry %d (|σᵢⱼ| far exceeds σᵢσⱼ: not a covariance matrix)", ErrOverflow, i)
	}
	return corrDst, stdDst, nil
}

// CorrelationToCovariance returns the covariance matrix Σ = D^(1/2)·R·D^(1/2), σᵢⱼ = rᵢⱼ·σᵢ·σⱼ, from an n×n
// correlation matrix R and the standard deviations σᵢ (each > 0), in dst.
// Errors: ErrDimension when n < 1, len(corr) ≠ n² or len(std) ≠ n; ErrNotFinite for a NaN or ±Inf entry;
// ErrOutOfRange when a standard deviation is ≤ 0; ErrOverflow when an entry leaves the float64 range.
//
// ID:      NUM-LINALG-168
// Time:    O(n²)
// Space:   O(1) extra; O(n²) for the result, one allocation only if dst lacks capacity
// Source:  NIST, Dataplot Reference Manual, "CORRELATION", https://www.itl.nist.gov/div898/software/dataplot/refman2/auxillar/correlat.htm (read 2026-11-12) (covariance = correlation·s_x·s_y)
// Source:  NIST/SEMATECH e-Handbook of Statistical Methods, §6.5.4.1 "Mean Vector and Covariance Matrix", https://www.itl.nist.gov/div898/handbook/pmc/section5/pmc541.htm (read 2026-11-12)
// Tags:    matrix, n-dimensional, covariance, correlation, statistics, scaling, linear-algebra
// Units:   corr dimensionless, std in X; the result in X²
// Example: R = [1 1/3; 1/3 1], σ = (2, 3) → Σ = [4 2; 2 9] (hand-checked)
func CorrelationToCovariance(dst, corr []float64, n int, std []float64) ([]float64, error) {
	if err := checkSquare(corr, n); err != nil {
		return nil, err
	}
	if len(std) != n {
		return nil, lengthError("std", len(std), "n", n)
	}
	for i, s := range std {
		if !reduce.AllFinite(s) {
			return nil, notFiniteAt("std", i)
		}
		if !(s > 0) {
			return nil, fmt.Errorf("%w: std[%d] = %g, want > 0", ErrOutOfRange, i, s)
		}
	}
	dst = grow(dst, n*n)
	for i := range n {
		for j := range n {
			dst[i*n+j] = corr[i*n+j] * std[i] * std[j]
		}
	}
	if i := firstNotFinite(dst); i >= 0 {
		return nil, fmt.Errorf("%w: entry %d", ErrOverflow, i)
	}
	return dst, nil
}
