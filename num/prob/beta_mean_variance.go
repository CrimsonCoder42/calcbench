package prob

// BetaMeanVariance returns the mean α/(α + β) and the variance αβ/((α + β)²·(α + β + 1)) of a beta variable.
// The variance is formed as (m·(1 − m))/(α + β + 1), m the mean, so it does not overflow for large shapes.
// Errors: ErrNotFinite when an input is NaN or ±Inf; ErrOutOfRange when α ≤ 0 or β ≤ 0.
//
// ID:      NUM-PROB-076
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Johnson, Kotz, Balakrishnan, Continuous Univariate Distributions, vol. 2, 2nd ed., Wiley, 1995, ch. 25 "Beta Distributions", §2 (moments)
// Tags:    probability, distribution, beta, mean, variance, expected-value
// Units:   dimensionless
// Example: α = 2, β = 3 → (0.4, 0.04)
func BetaMeanVariance(alpha, beta float64) (mean, variance float64, err error) {
	if err := checkBetaShapes(alpha, beta); err != nil {
		return 0, 0, err
	}
	half := alpha/2 + beta/2 // (α + β)/2 never overflows
	mean = (alpha / 2) / half
	variance = mean * ((beta / 2) / half) / (alpha + beta + 1) // an infinite α + β gives the limit 0
	return mean, variance, nil
}
