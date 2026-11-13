package prob

// PowerFunctionMeanVariance returns the mean a/(a + 1) and the variance a/((a + 1)²·(a + 2)) of the
// power-function distribution with shape a on [0, 1].
// Errors: ErrNotFinite when shape is NaN or ±Inf; ErrOutOfRange when shape ≤ 0.
//
// ID:      NUM-PROB-130
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Johnson, Kotz, Balakrishnan, Continuous Univariate Distributions, vol. 2, 2nd ed., Wiley, 1995, ch. 25 "Beta Distributions", §1 (the power-function distribution; beta moments with β = 1)
// Tags:    probability, distribution, power-function, beta, mean, variance
// Units:   dimensionless
// Example: a = 1 → (½, 1/12) (the uniform case); a = 2 → (2/3, 1/18)
func PowerFunctionMeanVariance(shape float64) (mean, variance float64, err error) {
	if err := checkPositive("shape", shape); err != nil {
		return 0, 0, err
	}
	mean = shape / (shape + 1)
	variance = mean / (shape + 1) / (shape + 2)
	return mean, variance, nil
}
