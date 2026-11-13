package prob

// BernoulliCDF returns P(X ≤ k) for a Bernoulli variable with success probability p: 0 for k < 0,
// 1 − p for k = 0, and 1 for k ≥ 1.
// Errors: ErrNotFinite when p is NaN or ±Inf; ErrProbability when p is outside [0, 1].
//
// ID:      NUM-PROB-024
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Ross, A First Course in Probability, 10th ed., Pearson, 2019, §4.6 (Bernoulli and binomial random variables)
// Source:  Casella, Berger, Statistical Inference, 2nd ed., Duxbury, 2002, §3.2 "Discrete Distributions" (Bernoulli distribution)
// Tags:    probability, distribution, bernoulli, cdf, discrete
// Units:   k: an outcome; p and result: probabilities in [0, 1]
// Example: p = 0.3, k = 0 → 0.7
func BernoulliCDF(k int, p float64) (float64, error) {
	if err := checkProbability("p", p); err != nil {
		return 0, err
	}
	switch {
	case k < 0:
		return 0, nil
	case k == 0:
		return 1 - p, nil
	}
	return 1, nil
}
