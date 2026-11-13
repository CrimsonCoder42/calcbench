package prob

// BernoulliPMF returns P(X = k) for a Bernoulli variable with success probability p: p at k = 1,
// 1 − p at k = 0, and 0 at every other k (outside the support is a value, not an error).
// Errors: ErrNotFinite when p is NaN or ±Inf; ErrProbability when p is outside [0, 1].
//
// ID:      NUM-PROB-023
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Ross, A First Course in Probability, 10th ed., Pearson, 2019, §4.6 (Bernoulli and binomial random variables)
// Source:  Casella, Berger, Statistical Inference, 2nd ed., Duxbury, 2002, §3.2 "Discrete Distributions" (Bernoulli distribution)
// Tags:    probability, distribution, bernoulli, pmf, discrete
// Units:   k: an outcome, 0 or 1; p and result: probabilities in [0, 1]
// Example: p = 0.3: k = 1 → 0.3; k = 0 → 0.7; k = 2 → 0
func BernoulliPMF(k int, p float64) (float64, error) {
	if err := checkProbability("p", p); err != nil {
		return 0, err
	}
	switch k {
	case 0:
		return 1 - p, nil
	case 1:
		return p, nil
	}
	return 0, nil
}
