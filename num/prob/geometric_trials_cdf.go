package prob

import "math"

// GeometricTrialsCDF returns P(X ≤ k) = 1 − (1 − p)^k, X the number of TRIALS up to and including the
// first success (D-013); 0 for k < 1. It is evaluated as −expm1(k·log1p(−p)), with no cancellation for a
// small p.
// Errors: as GeometricTrialsPMF.
//
// ID:      NUM-PROB-030
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Ross, A First Course in Probability, 10th ed., Pearson, 2019, §4.8 (other discrete distributions: the geometric random variable)
// Source:  Casella, Berger, Statistical Inference, 2nd ed., Duxbury, 2002, §3.2 "Discrete Distributions" (geometric distribution)
// Tags:    probability, distribution, geometric, cdf, discrete, trials, waiting-time
// Units:   k: a count of trials; p and result: probabilities in [0, 1]
// Example: p = 0.2, k = 3 → 1 − 0.512 = 0.488
func GeometricTrialsCDF(k int, p float64) (float64, error) {
	if err := checkGeometric(p); err != nil {
		return 0, err
	}
	if k < 1 {
		return 0, nil
	}
	return geometricTrialsCDF(k, p), nil
}

// geometricTrialsCDF is GeometricTrialsCDF for a validated p in (0, 1] and k ≥ 1.
func geometricTrialsCDF(k int, p float64) float64 {
	switch {
	case p == 1:
		return 1
	case k == 1:
		return p // F(1) = p exactly, so the quantile's u ≤ p shortcut agrees with the CDF
	}
	return -math.Expm1(float64(k) * math.Log1p(-p))
}
