package prob

import "fmt"

// NegativePredictiveValue returns NPV = Sp·(1 − π) / (Sp·(1 − π) + (1 − Se)·π), the probability of no
// condition given a negative test, from sensitivity Se, specificity Sp and prevalence π. It states a
// number only.
// Errors: ErrNotFinite when an input is NaN or ±Inf; ErrProbability when one is outside [0, 1];
// ErrOutOfRange when the denominator is 0 (no negative test is possible).
//
// ID:      NUM-PROB-011
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Ross, A First Course in Probability, 10th ed., Pearson, 2019, §3.3 "Bayes's Formula"
// Tags:    probability, bayes, diagnostic-test, predictive-value, npv, medical
// Units:   sensitivity, specificity, prevalence and result: probabilities in [0, 1]
// Example: a 2×2 table TP 231, FP 32, FN 27, TN 54: Se = 231/258, Sp = 54/86, π = 258/344 → 54/81 = 2/3
func NegativePredictiveValue(sensitivity, specificity, prevalence float64) (float64, error) {
	if err := checkPredictiveInputs(sensitivity, specificity, prevalence); err != nil {
		return 0, err
	}
	trueNeg := specificity * (1 - prevalence)
	falseNeg := (1 - sensitivity) * prevalence
	denom := trueNeg + falseNeg
	if denom == 0 {
		return 0, fmt.Errorf("%w: specificity·(1−prevalence) + (1−sensitivity)·prevalence = 0", ErrOutOfRange)
	}
	return min(1, trueNeg/denom), nil
}
