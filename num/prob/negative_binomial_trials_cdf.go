package prob

// NegativeBinomialTrialsCDF returns P(X ≤ x), X the number of TRIALS up to and including the r-th
// success; 0 for x < r. The r-th success has come by trial x exactly when x trials hold at least r
// successes, so this is the binomial upper tail P(Bin(x, p) ≥ r) = Σ_{i=r..x} C(x, i)·pⁱ·(1 − p)^(x−i),
// summed as in BinomialCDF.
// Errors: as NegativeBinomialTrialsPMF; ErrOutOfRange also when the sum would need more than 2²⁸ terms
// (as BinomialCDF).
//
// ID:      NUM-PROB-033
// Time:    O(min(r, x − r + 1, c·(σ + 1)))   σ = √(xp(1 − p)); c ≈ 10 (as BinomialCDF on x trials)
// Space:   O(1), no allocation
// Source:  Ross, A First Course in Probability, 10th ed., Pearson, 2019, §4.8 (other discrete distributions: the negative binomial random variable)
// Source:  Casella, Berger, Statistical Inference, 2nd ed., Duxbury, 2002, §3.2 "Discrete Distributions" (negative binomial distribution)
// Source:  Loader, "Fast and Accurate Computation of Binomial Probabilities", technical report, 2000 (rev. 2002), https://www.r-project.org/doc/reports/CLoader-dbinom-2002.pdf (read 2026-11-12), §2 eqs. (3)–(7) and Appendix A (stirlerr, bd0)
// Tags:    probability, distribution, negative-binomial, pascal, cdf, discrete, trials
// Units:   x: a count of trials; r: a count of successes; p and result: probabilities in [0, 1]
// Example: r = 2, p = 0.5, x = 4 → 0.25 + 0.25 + 0.1875 = 0.6875
func NegativeBinomialTrialsCDF(x, r int, p float64) (float64, error) {
	if err := checkNegativeBinomial(r, p); err != nil {
		return 0, err
	}
	if x < r {
		return 0, nil
	}
	if p == 1 {
		return 1, nil
	}
	_, upper, err := binomialTails(r-1, x, p)
	return upper, err
}
