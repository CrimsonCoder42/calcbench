package prob

// BinomialMeanVariance returns the mean np and the variance np(1 − p) of a binomial variable with n
// trials and success probability p.
// Errors: as BinomialPMF.
//
// ID:      NUM-PROB-028
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  NIST/SEMATECH, e-Handbook of Statistical Methods, §1.3.6.6.18 "Binomial Distribution", https://www.itl.nist.gov/div898/handbook/eda/section3/eda366i.htm (read 2026-11-12)
// Source:  Ross, A First Course in Probability, 10th ed., Pearson, 2019, §4.6 (Bernoulli and binomial random variables)
// Tags:    probability, distribution, binomial, mean, variance, expected-value
// Units:   n: a count of trials; p: a probability; mean: a count of successes; variance: that count squared
// Example: n = 10, p = 0.5 → (5, 2.5)
func BinomialMeanVariance(n int, p float64) (mean, variance float64, err error) {
	if err := checkBinomial(n, p); err != nil {
		return 0, 0, err
	}
	mean = float64(n) * p
	return mean, mean * (1 - p), nil
}
