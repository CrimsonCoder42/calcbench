package prob

import (
	"fmt"
	"math"
)

// checkGeometric validates p as a probability in (0, 1].
func checkGeometric(p float64) error {
	if err := checkProbability("p", p); err != nil {
		return err
	}
	if p == 0 {
		return fmt.Errorf("%w: p = 0, want > 0", ErrOutOfRange)
	}
	return nil
}

// GeometricTrialsPMF returns P(X = k) = (1 − p)^(k−1)·p, where X counts the TRIALS up to and including
// the first success (support 1, 2, …; doc/DECISIONS.md D-013); 0 for k < 1. The failures count is X − 1.
// It is evaluated as exp((k − 1)·log1p(−p))·p, accurate for a small p.
// Errors: ErrNotFinite when p is NaN or ±Inf; ErrProbability when p is outside [0, 1]; ErrOutOfRange
// when p = 0.
//
// ID:      NUM-PROB-029
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Ross, A First Course in Probability, 10th ed., Pearson, 2019, §4.8 (other discrete distributions: the geometric random variable)
// Source:  Casella, Berger, Statistical Inference, 2nd ed., Duxbury, 2002, §3.2 "Discrete Distributions" (geometric distribution)
// Tags:    probability, distribution, geometric, pmf, discrete, trials, waiting-time
// Units:   k: a count of trials; p and result: probabilities in [0, 1]
// Example: p = 0.2, k = 3 → 0.8²·0.2 = 0.128
func GeometricTrialsPMF(k int, p float64) (float64, error) {
	if err := checkGeometric(p); err != nil {
		return 0, err
	}
	if k < 1 {
		return 0, nil
	}
	if p == 1 {
		return boolProb(k == 1), nil
	}
	return math.Exp(float64(k-1)*math.Log1p(-p)) * p, nil
}
