package prob

// ChebyshevBound returns min(1, 1/k²), Chebyshev's upper bound on P(|X − μ| ≥ kσ) for any X with a
// finite variance. It is a bound, not a probability; below k = 1 it is the trivial bound 1, and for a
// huge k it underflows to 0.
// Errors: ErrNotFinite when k is NaN or ±Inf; ErrOutOfRange when k ≤ 0.
//
// ID:      NUM-PROB-020
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Ross, A First Course in Probability, 10th ed., Pearson, 2019, §8.2 "Chebyshev's Inequality and the Weak Law of Large Numbers"
// Source:  Casella, Berger, Statistical Inference, 2nd ed., Duxbury, 2002, §3.6 "Inequalities and Identities" (P(g(X) ≥ r) ≤ E g(X)/r for g ≥ 0)
// Tags:    probability, chebyshev, inequality, tail-bound
// Units:   k: a number of standard deviations, dimensionless; result: a probability bound in [0, 1]
// Example: k = 2 → 0.25; k = 0.5 → 1
func ChebyshevBound(k float64) (float64, error) {
	if err := checkPositive("k", k); err != nil {
		return 0, err
	}
	if k <= 1 {
		return 1, nil
	}
	return (1 / k) / k, nil
}
