package prob

// PoissonMeanVariance returns the mean λ and the variance λ of a Poisson variable.
// Errors: as PoissonPMF.
//
// ID:      NUM-PROB-043
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  NIST/SEMATECH, e-Handbook of Statistical Methods, §1.3.6.6.19 "Poisson Distribution", https://www.itl.nist.gov/div898/handbook/eda/section3/eda366j.htm (read 2026-11-12)
// Source:  Ross, A First Course in Probability, 10th ed., Pearson, 2019, §4.7 (Poisson random variable)
// Tags:    probability, distribution, poisson, mean, variance, expected-value
// Units:   lambda: the mean count per interval; mean: a count; variance: that count squared
// Example: λ = 2 → (2, 2)
func PoissonMeanVariance(lambda float64) (mean, variance float64, err error) {
	if err := checkPoisson(lambda); err != nil {
		return 0, 0, err
	}
	return lambda, lambda, nil
}
