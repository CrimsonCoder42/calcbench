package prob

import (
	"fmt"
	"math"
)

// NoncentralFMeanVariance returns the mean and variance of the noncentral F distribution with df1 = ν₁,
// df2 = ν₂ degrees of freedom and noncentrality λ: mean ν₂(ν₁ + λ)/(ν₁(ν₂ − 2)) and variance
// 2(ν₂/ν₁)²·((ν₁ + λ)² + (ν₁ + 2λ)(ν₂ − 2))/((ν₂ − 2)²(ν₂ − 4)). Both are returned only when both exist
// (ν₂ > 4), as FDistributionMeanVariance does.
// Errors: ErrNotFinite when an input is NaN or ±Inf; ErrOutOfRange when df1 ≤ 0, df2 ≤ 4, λ < 0 or λ > 10⁵;
// ErrOverflow when the variance exceeds the float64 range.
//
// ID:      NUM-PROB-155
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Johnson, Kotz, Balakrishnan, Continuous Univariate Distributions, vol. 2, 2nd ed., Wiley, 1995, ch. 30 "Noncentral F-Distributions" (moments of the noncentral F)
// Tags:    probability, distribution, f-distribution, noncentral, mean, variance, continuous
// Units:   df1, df2, lambda dimensionless; mean and variance dimensionless
// Range:   df1 > 0, df2 > 4, 0 ≤ λ ≤ 10⁵
// Example: df1 = 4, df2 = 10, λ = 2 → mean 10·6/(4·8) = 1.875, variance 2·6.25·(36 + 8·8)/(64·6) = 3.2552083…
func NoncentralFMeanVariance(df1, df2, lambda float64) (mean, variance float64, err error) {
	if err := checkFDistribution(df1, df2); err != nil {
		return 0, 0, err
	}
	if err := checkNoncentrality(lambda); err != nil {
		return 0, 0, err
	}
	if df2 <= 4 {
		return 0, 0, fmt.Errorf("%w: df2 = %v, want > 4 (variance undefined)", ErrOutOfRange, df2)
	}
	r := df2 / df1
	mean = r * (df1 + lambda) / (df2 - 2)
	variance = 2 * r * r * ((df1+lambda)*(df1+lambda) + (df1+2*lambda)*(df2-2)) / ((df2 - 2) * (df2 - 2) * (df2 - 4))
	if math.IsInf(mean, 0) || math.IsInf(variance, 0) {
		return 0, 0, fmt.Errorf("%w: noncentral F moments for df1 = %v, df2 = %v, lambda = %v", ErrOverflow, df1, df2, lambda)
	}
	return mean, variance, nil
}
