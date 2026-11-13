package prob

import "fmt"

// HypergeometricMeanVariance returns the mean nK/N and the variance n·(K/N)·((N−K)/N)·((N−n)/(N−1)) of X
// hypergeometric with N = population, K = successes, n = draws; the variance is 0 when N = 1. Every
// factor but n is a ratio in [0, 1], so nothing overflows.
// Errors: as HypergeometricPMF; ErrOutOfRange also when population = 0.
//
// ID:      NUM-PROB-037
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Ross, A First Course in Probability, 10th ed., Pearson, 2019, §4.8 (other discrete distributions: the hypergeometric random variable)
// Source:  Casella, Berger, Statistical Inference, 2nd ed., Duxbury, 2002, §3.2 "Discrete Distributions" (hypergeometric distribution)
// Tags:    probability, distribution, hypergeometric, mean, variance, expected-value
// Units:   population, successes, draws: counts of items; mean: a count of successes; variance: that count squared
// Example: N = 52, K = 4, n = 5 → (5/13 ≈ 0.38461538461538464, 940/2873 ≈ 0.32718412808910546)
func HypergeometricMeanVariance(population, successes, draws int) (mean, variance float64, err error) {
	if err := checkHypergeometric(population, successes, draws); err != nil {
		return 0, 0, err
	}
	if population == 0 {
		return 0, 0, fmt.Errorf("%w: population = 0", ErrOutOfRange)
	}
	nf := float64(population)
	pK := float64(successes) / nf
	mean = float64(draws) * pK
	if population == 1 {
		return mean, 0, nil
	}
	variance = mean * (float64(population-successes) / nf) * (float64(population-draws) / float64(population-1))
	return mean, variance, nil
}
