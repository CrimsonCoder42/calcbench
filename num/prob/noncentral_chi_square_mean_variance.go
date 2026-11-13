package prob

// NoncentralChiSquareMeanVariance returns the mean df + λ and the variance 2(df + 2λ) of the noncentral
// chi-square distribution with df degrees of freedom and noncentrality λ.
// Errors: ErrNotFinite when df or λ is NaN or ±Inf; ErrOutOfRange when df ≤ 0, df > 10⁶, λ < 0 or λ > 10⁵.
//
// ID:      NUM-PROB-153
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Johnson, Kotz, Balakrishnan, Continuous Univariate Distributions, vol. 2, 2nd ed., Wiley, 1995, ch. 29 "Noncentral χ² Distributions" (moments: mean ν + λ, variance 2(ν + 2λ))
// Tags:    probability, distribution, chi-square, noncentral, mean, variance, continuous
// Units:   df, lambda dimensionless; mean and variance dimensionless
// Range:   0 < df ≤ 10⁶, 0 ≤ λ ≤ 10⁵
// Example: df = 3, λ = 2 → mean 5, variance 14
func NoncentralChiSquareMeanVariance(df, lambda float64) (mean, variance float64, err error) {
	if err := checkChiSquareDF(df); err != nil {
		return 0, 0, err
	}
	if err := checkNoncentrality(lambda); err != nil {
		return 0, 0, err
	}
	return df + lambda, 2 * (df + 2*lambda), nil
}
