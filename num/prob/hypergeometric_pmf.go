package prob

import (
	"fmt"
	"math"
)

// checkHypergeometric validates population ≥ 0 and successes, draws in [0, population].
func checkHypergeometric(population, successes, draws int) error {
	if population < 0 {
		return fmt.Errorf("%w: population = %d, want >= 0", ErrOutOfRange, population)
	}
	if successes < 0 || successes > population {
		return fmt.Errorf("%w: successes = %d, want 0..%d", ErrOutOfRange, successes, population)
	}
	if draws < 0 || draws > population {
		return fmt.Errorf("%w: draws = %d, want 0..%d", ErrOutOfRange, draws, population)
	}
	return nil
}

// hypergeometricSupport returns the least and greatest k with a positive probability.
func hypergeometricSupport(population, successes, draws int) (lo, hi int) {
	return max(0, draws-(population-successes)), min(draws, successes)
}

// hypergeometricTerm returns P(X = k) for validated parameters and lo ≤ k ≤ hi.
// With N = population, K = successes, n = draws and any q in (0, 1),
// C(K,k)·C(N−K,n−k)/C(N,n) = b(k; K, q)·b(n−k; N−K, q) / b(n; N, q), because the powers of q and of 1 − q
// cancel. q = n/N puts the denominator at its mode, so it never underflows; the three saddle-point
// terms are combined in their exponents before exp, so the numerator cannot underflow early either.
func hypergeometricTerm(k, population, successes, draws int) float64 {
	if draws == 0 || draws == population || successes == 0 || successes == population {
		return 1 // a single support point
	}
	q := float64(draws) / float64(population)
	lc1, s1 := binomialTermParts(k, successes, q)
	lc2, s2 := binomialTermParts(draws-k, population-successes, q)
	lc3, s3 := binomialTermParts(draws, population, q)
	return min(1, math.Exp(lc1+lc2-lc3)*(s1*s2/s3))
}

// HypergeometricPMF returns P(X = k) = C(K, k)·C(N−K, n−k) / C(N, n): the probability of k successes in
// n = draws draws without replacement from N = population items of which K = successes are successes;
// 0 for k outside max(0, n+K−N)…min(n, K). It is evaluated as three saddle-point binomial terms,
// b(k; K, q)·b(n−k; N−K, q) / b(n; N, q) with q = n/N, with no cancellation and no factorial.
// Errors: ErrOutOfRange when population < 0, or successes or draws is outside 0…population.
//
// ID:      NUM-PROB-035
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Ross, A First Course in Probability, 10th ed., Pearson, 2019, §4.8 (other discrete distributions: the hypergeometric random variable)
// Source:  Casella, Berger, Statistical Inference, 2nd ed., Duxbury, 2002, §3.2 "Discrete Distributions" (hypergeometric distribution)
// Source:  Loader, "Fast and Accurate Computation of Binomial Probabilities", technical report, 2000 (rev. 2002), https://www.r-project.org/doc/reports/CLoader-dbinom-2002.pdf (read 2026-11-12), §2 eqs. (3)–(7) and Appendix A (stirlerr, bd0)
// Tags:    probability, distribution, hypergeometric, pmf, discrete, sampling-without-replacement, cards
// Units:   k, population, successes, draws: counts of items; result: a probability in [0, 1]
// Example: one ace in a 5-card hand: N = 52, K = 4, n = 5, k = 1 → 4·194580/2598960 = 3243/10829 ≈ 0.2994736356080894
func HypergeometricPMF(k, population, successes, draws int) (float64, error) {
	if err := checkHypergeometric(population, successes, draws); err != nil {
		return 0, err
	}
	lo, hi := hypergeometricSupport(population, successes, draws)
	if k < lo || k > hi {
		return 0, nil
	}
	return hypergeometricTerm(k, population, successes, draws), nil
}
