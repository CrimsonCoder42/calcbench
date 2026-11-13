package prob

// BernoulliMeanVariance returns the mean p and the variance p(1 − p) of a Bernoulli variable.
// Errors: ErrNotFinite when p is NaN or ±Inf; ErrProbability when p is outside [0, 1].
//
// ID:      NUM-PROB-025
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Ross, A First Course in Probability, 10th ed., Pearson, 2019, §4.6 (Bernoulli and binomial random variables)
// Source:  Casella, Berger, Statistical Inference, 2nd ed., Duxbury, 2002, §3.2 "Discrete Distributions" (Bernoulli distribution)
// Tags:    probability, distribution, bernoulli, mean, variance, expected-value
// Units:   p: a probability; mean: a probability; variance: dimensionless
// Example: p = 0.3 → (0.3, 0.21)
func BernoulliMeanVariance(p float64) (mean, variance float64, err error) {
	if err := checkProbability("p", p); err != nil {
		return 0, 0, err
	}
	return p, p * (1 - p), nil
}
