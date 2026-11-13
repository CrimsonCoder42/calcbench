package prob

// HypergeometricCDF returns P(X ≤ k) for X hypergeometric (k successes in draws draws without
// replacement from population items, successes of them successes): 0 below the support, 1 at or above
// its top. It sums the tail on the far side of k from the mean nK/N, from one PMF term outward with the
// ratio recurrences P(i+1)/P(i) = (K−i)(n−i)/((i+1)(N−K−n+i+1)) and its inverse, compensated, and stops
// when the rest is below 2⁻⁵³ of the sum; the other side is 1 minus that sum.
// Errors: as HypergeometricPMF; ErrOutOfRange also when the sum would need more than 2²⁸ terms (σ above
// roughly 10⁷).
//
// ID:      NUM-PROB-036
// Time:    O(min(k − lo + 1, hi − k, c·(σ + 1)))   lo, hi = the support bounds, σ = the standard deviation (HypergeometricMeanVariance), c ≈ 10
// Space:   O(1), no allocation
// Source:  Ross, A First Course in Probability, 10th ed., Pearson, 2019, §4.8 (other discrete distributions: the hypergeometric random variable)
// Source:  Casella, Berger, Statistical Inference, 2nd ed., Duxbury, 2002, §3.2 "Discrete Distributions" (hypergeometric distribution)
// Tags:    probability, distribution, hypergeometric, cdf, discrete, sampling-without-replacement, cards
// Units:   k, population, successes, draws: counts of items; result: a probability in [0, 1]
// Example: at most one ace in a 5-card hand: N = 52, K = 4, n = 5, k = 1 → (C(48,5) + 4·C(48,4))/C(52,5) ≈ 0.958315633945886
func HypergeometricCDF(k, population, successes, draws int) (float64, error) {
	if err := checkHypergeometric(population, successes, draws); err != nil {
		return 0, err
	}
	lo, hi := hypergeometricSupport(population, successes, draws)
	switch {
	case k < lo:
		return 0, nil
	case k >= hi:
		return 1, nil
	}
	nf, kf, df := float64(population), float64(successes), float64(draws)
	var s compensatedSum
	if float64(k) < df*(kf/nf) {
		t := hypergeometricTerm(k, population, successes, draws)
		s.add(t)
		for i, steps := k, 0; i > lo; i, steps = i-1, steps+1 {
			if steps == maxTailTerms {
				return 0, errTailTooLong("hypergeometric tail")
			}
			// P(i−1)/P(i) = i·(N−K−n+i) / ((K−i+1)·(n−i+1))
			r := float64(i) * float64(population-successes-draws+i) /
				(float64(successes-i+1) * float64(draws-i+1))
			t *= r
			s.add(t)
			if tailDone(t, r, &s) {
				break
			}
		}
		return min(1, s.total()), nil
	}
	t := hypergeometricTerm(k+1, population, successes, draws)
	s.add(t)
	for i, steps := k+1, 0; i < hi; i, steps = i+1, steps+1 {
		if steps == maxTailTerms {
			return 0, errTailTooLong("hypergeometric tail")
		}
		r := float64(successes-i) * float64(draws-i) /
			(float64(i+1) * float64(population-successes-draws+i+1))
		t *= r
		s.add(t)
		if tailDone(t, r, &s) {
			break
		}
	}
	return 1 - min(1, s.total()), nil
}
