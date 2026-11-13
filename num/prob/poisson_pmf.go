package prob

import (
	"fmt"
	"math"
)

// checkPoisson validates λ finite and ≥ 0.
func checkPoisson(lambda float64) error {
	if err := checkFinite("lambda", lambda); err != nil {
		return err
	}
	if lambda < 0 {
		return fmt.Errorf("%w: lambda = %v, want >= 0", ErrOutOfRange, lambda)
	}
	return nil
}

// poissonTerm returns e^(−λ)·λᵏ/k! for k ≥ 0 and λ > 0 in Loader's form: e^(−λ) at k = 0, else
// exp(−δ(k) − D₀(k, λ)) / √(2πk), with δ the Stirling remainder and D₀ the deviance.
func poissonTerm(k float64, lambda float64) float64 {
	if k == 0 {
		return math.Exp(-lambda)
	}
	return min(1, math.Exp(-stirlingRemainder(k)-deviance(k, lambda))/math.Sqrt(2*math.Pi*k))
}

// PoissonPMF returns P(X = k) = e^(−λ)·λᵏ/k! for a Poisson variable with mean λ; 0 for k < 0, and at
// λ = 0 the point mass at 0. It never forms λᵏ or k!: it uses Loader's saddle-point form
// exp(−δ(k) − D₀(k, λ))/√(2πk): a few ulps near the mode, and in the tails a relative error of about
// 2⁻⁵³ times D₀(k, λ).
// Errors: ErrNotFinite when lambda is NaN or ±Inf; ErrOutOfRange when lambda < 0.
//
// ID:      NUM-PROB-041
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  NIST/SEMATECH, e-Handbook of Statistical Methods, §1.3.6.6.19 "Poisson Distribution", https://www.itl.nist.gov/div898/handbook/eda/section3/eda366j.htm (read 2026-11-12)
// Source:  Loader, "Fast and Accurate Computation of Binomial Probabilities", technical report, 2000 (rev. 2002), https://www.r-project.org/doc/reports/CLoader-dbinom-2002.pdf (read 2026-11-12), §2 eqs. (3)–(7) and Appendix A (stirlerr, bd0)
// Source:  Ross, A First Course in Probability, 10th ed., Pearson, 2019, §4.7 (Poisson random variable)
// Tags:    probability, distribution, poisson, pmf, discrete, counts, saddle-point
// Units:   k: a count of events; lambda: the mean count per interval; result: a probability in [0, 1]
// Example: λ = 2, k = 0 → e⁻² ≈ 0.1353352832366127; λ = 1000, k = 1000 → 0.012614611348719664 (λᵏ/k! overflows)
func PoissonPMF(k int, lambda float64) (float64, error) {
	if err := checkPoisson(lambda); err != nil {
		return 0, err
	}
	if k < 0 {
		return 0, nil
	}
	if lambda == 0 {
		return boolProb(k == 0), nil
	}
	return poissonTerm(float64(k), lambda), nil
}
