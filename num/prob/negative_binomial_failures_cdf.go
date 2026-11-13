package prob

// NegativeBinomialFailuresCDF returns P(Y ≤ k), Y the number of FAILURES before the r-th success; 0 for
// k < 0. It is NegativeBinomialTrialsCDF(k + r, r, p) (NUM-PROB-033): the binomial upper tail
// P(Bin(k + r, p) ≥ r).
// Errors: as NegativeBinomialFailuresPMF; ErrOutOfRange also when the sum would need more than 2²⁸ terms.
//
// ID:      NUM-PROB-100
// Time:    O(min(r, k + 1, c·(σ + 1)))   σ = √((k + r)p(1 − p)); c ≈ 10 (as BinomialCDF on k + r trials)
// Space:   O(1), no allocation
// Source:  Johnson, Kemp, Kotz, Univariate Discrete Distributions, 3rd ed., Wiley, 2005, ch. 5 "Negative Binomial Distribution", §5.1 (definition), §5.6 (relation to the binomial tail)
// Tags:    probability, distribution, negative-binomial, cdf, discrete, failures
// Units:   k: a count of failures; r: a count of successes; p and result: probabilities in [0, 1]
// Example: r = 2, p = 0.5, k = 2 → 0.25 + 0.25 + 0.1875 = 0.6875
func NegativeBinomialFailuresCDF(k, r int, p float64) (float64, error) {
	if err := checkNegativeBinomial(r, p); err != nil {
		return 0, err
	}
	if k < 0 {
		return 0, nil
	}
	x, err := negativeBinomialTrialsFor(k, r)
	if err != nil {
		return 0, err
	}
	return NegativeBinomialTrialsCDF(x, r, p)
}
