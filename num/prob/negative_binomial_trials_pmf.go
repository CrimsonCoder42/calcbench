package prob

import "fmt"

// checkNegativeBinomial validates r ≥ 1 and p a probability in (0, 1].
func checkNegativeBinomial(r int, p float64) error {
	if r < 1 {
		return fmt.Errorf("%w: r = %d, want >= 1", ErrOutOfRange, r)
	}
	return checkGeometric(p)
}

// NegativeBinomialTrialsPMF returns P(X = x) = C(x−1, r−1)·pʳ·(1 − p)^(x−r), where X counts the TRIALS up
// to and including the r-th success (the Pascal form, support r, r+1, …, matching GeometricTrialsPMF; the
// failures count is X − r); 0 for x < r. It is evaluated as (r/x)·b(r; x, p) through the identity
// C(x−1, r−1) = (r/x)·C(x, r), with b the saddle-point binomial term of BinomialPMF.
// Errors: ErrOutOfRange when r < 1 or p = 0; ErrNotFinite when p is NaN or ±Inf; ErrProbability when p
// is outside [0, 1].
//
// ID:      NUM-PROB-032
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Ross, A First Course in Probability, 10th ed., Pearson, 2019, §4.8 (other discrete distributions: the negative binomial random variable)
// Source:  Casella, Berger, Statistical Inference, 2nd ed., Duxbury, 2002, §3.2 "Discrete Distributions" (negative binomial distribution)
// Source:  Loader, "Fast and Accurate Computation of Binomial Probabilities", technical report, 2000 (rev. 2002), https://www.r-project.org/doc/reports/CLoader-dbinom-2002.pdf (read 2026-11-12), §2 eqs. (3)–(7) and Appendix A (stirlerr, bd0)
// Tags:    probability, distribution, negative-binomial, pascal, pmf, discrete, trials, waiting-time
// Units:   x: a count of trials; r: a count of successes; p and result: probabilities in [0, 1]
// Example: r = 2, p = 0.5, x = 4 → C(3,1)·0.25·0.25 = 0.1875
func NegativeBinomialTrialsPMF(x, r int, p float64) (float64, error) {
	if err := checkNegativeBinomial(r, p); err != nil {
		return 0, err
	}
	if x < r {
		return 0, nil
	}
	if p == 1 {
		return boolProb(x == r), nil
	}
	return float64(r) / float64(x) * binomialTerm(r, x, p), nil
}
