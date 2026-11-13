package prob

// HypergeometricQuantile returns the smallest k with P(X ≤ k) ≥ u for X hypergeometric (successes in
// draws draws without replacement from population items, successes of them successes). It bisects on
// HypergeometricCDF (NUM-PROB-036) over the support max(0, n+K−N)…min(n, K), so the result agrees with the
// CDF as this package computes it. u = 0 returns the bottom of the support, u = 1 the top.
// Errors: ErrNotFinite when u is NaN or ±Inf; ErrProbability when u is outside [0, 1]; as
// HypergeometricCDF otherwise.
//
// ID:      NUM-PROB-098
// Time:    O(log(hi − lo + 1) · C)   lo, hi = the support bounds, C = the cost of one HypergeometricCDF
// Space:   O(1), no allocation
// Source:  Johnson, Kemp, Kotz, Univariate Discrete Distributions, 3rd ed., Wiley, 2005, ch. 6 "Hypergeometric Distributions", §6.1–6.2 (definition, cdf)
// Tags:    probability, distribution, hypergeometric, quantile, inverse-cdf, discrete, sampling-without-replacement
// Units:   u: a probability; population, successes, draws, result: counts of items
// Example: N = 52, K = 4, n = 5: u = 0.5 → 0 (F(0) ≈ 0.6588); u = 0.9 → 1 (F(1) ≈ 0.9583)
func HypergeometricQuantile(u float64, population, successes, draws int) (int, error) {
	if err := checkProbability("u", u); err != nil {
		return 0, err
	}
	if err := checkHypergeometric(population, successes, draws); err != nil {
		return 0, err
	}
	lo, hi := hypergeometricSupport(population, successes, draws)
	if u == 0 {
		return lo, nil
	}
	l := lo - 1 // F(l) = 0 < u
	for hi-l > 1 {
		mid := l + (hi-l)/2
		c, err := HypergeometricCDF(mid, population, successes, draws)
		if err != nil {
			return 0, err
		}
		if c >= u {
			hi = mid
		} else {
			l = mid
		}
	}
	return hi, nil
}
