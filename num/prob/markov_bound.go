package prob

import "fmt"

// MarkovBound returns min(1, E[X]/a), Markov's upper bound on P(X ≥ a) for a non-negative X and a > 0.
// That X is non-negative is the caller's assumption. It is a bound, not a probability.
// Errors: ErrNotFinite when an input is NaN or ±Inf; ErrOutOfRange when mean < 0 or threshold ≤ 0.
//
// ID:      NUM-PROB-022
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Ross, A First Course in Probability, 10th ed., Pearson, 2019, §8.2 "Chebyshev's Inequality and the Weak Law of Large Numbers" (Markov's inequality)
// Source:  Casella, Berger, Statistical Inference, 2nd ed., Duxbury, 2002, §3.6 "Inequalities and Identities" (P(g(X) ≥ r) ≤ E g(X)/r for g ≥ 0)
// Tags:    probability, markov, inequality, tail-bound
// Units:   mean and threshold: the same unit; result: a probability bound in [0, 1]
// Example: mean 50, threshold 75 → 2/3
func MarkovBound(mean, threshold float64) (float64, error) {
	if err := checkFinite("mean", mean); err != nil {
		return 0, err
	}
	if mean < 0 {
		return 0, fmt.Errorf("%w: mean = %v, want >= 0", ErrOutOfRange, mean)
	}
	if err := checkPositive("threshold", threshold); err != nil {
		return 0, err
	}
	return min(1, mean/threshold), nil
}
