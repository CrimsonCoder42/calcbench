package prob

// BinomialQuantile returns the smallest k in 0…n with P(X ≤ k) ≥ u for X ~ Binomial(n, p), searching by
// bisection on BinomialCDF (NUM-PROB-027), so the result agrees with the CDF as this package computes it.
// u = 0 returns 0; u = 1 returns the smallest k with BinomialCDF(k) = 1.
// Errors: ErrNotFinite when u or p is NaN or ±Inf; ErrProbability when u or p is outside [0, 1];
// ErrOutOfRange when n < 0, or as BinomialCDF.
//
// ID:      NUM-PROB-085
// Time:    O(log n · C)   C = the cost of one BinomialCDF, O(√(np(1−p)) + 1)
// Space:   O(1), no allocation
// Source:  Johnson, Kemp, Kotz, Univariate Discrete Distributions, 3rd ed., Wiley, 2005, ch. 3 "Binomial Distribution", §3.2 (cdf), inverted as the least k with F(k) ≥ u
// Source:  Devroye, Non-Uniform Random Variate Generation, Springer, 1986, §III.2.1 (inversion by search for discrete distributions)
// Tags:    probability, distribution, binomial, quantile, inverse-cdf, discrete
// Units:   u, p: probabilities; n and result: counts
// Example: n = 10, p = 0.5, u = 0.5 → 5 (F(4) = 0.376953125 < 0.5 ≤ F(5) = 0.623046875)
func BinomialQuantile(u float64, n int, p float64) (int, error) {
	if err := checkProbability("u", u); err != nil {
		return 0, err
	}
	if err := checkBinomial(n, p); err != nil {
		return 0, err
	}
	lo, hi := -1, n // invariant: F(lo) < u ≤ F(hi), with F(−1) = 0 and F(n) = 1
	for hi-lo > 1 {
		mid := lo + (hi-lo)/2
		c, err := BinomialCDF(mid, n, p)
		if err != nil {
			return 0, err
		}
		if c >= u {
			hi = mid
		} else {
			lo = mid
		}
	}
	return hi, nil
}
