package prob

import "fmt"

// PositivePredictiveValue returns PPV = Se·π / (Se·π + (1 − Sp)·(1 − π)), the probability of the
// condition given a positive test, from the test's sensitivity Se, specificity Sp and the prevalence
// (pre-test probability) π. It is Bayes' theorem with the law of total probability in the denominator.
// It states a number only.
// Errors: ErrNotFinite when an input is NaN or ±Inf; ErrProbability when one is outside [0, 1];
// ErrOutOfRange when the denominator is 0 (no positive test is possible).
//
// ID:      NUM-PROB-010
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Ross, A First Course in Probability, 10th ed., Pearson, 2019, §3.3 "Bayes's Formula"
// Tags:    probability, bayes, diagnostic-test, predictive-value, ppv, medical
// Units:   sensitivity, specificity, prevalence and result: probabilities in [0, 1]
// Example: a 2×2 table TP 231, FP 32, FN 27, TN 54: Se = 231/258, Sp = 54/86, π = 258/344 → 231/263 ≈ 0.878
func PositivePredictiveValue(sensitivity, specificity, prevalence float64) (float64, error) {
	if err := checkPredictiveInputs(sensitivity, specificity, prevalence); err != nil {
		return 0, err
	}
	truePos := sensitivity * prevalence
	falsePos := (1 - specificity) * (1 - prevalence)
	denom := truePos + falsePos
	if denom == 0 {
		return 0, fmt.Errorf("%w: sensitivity·prevalence + (1−specificity)·(1−prevalence) = 0", ErrOutOfRange)
	}
	return min(1, truePos/denom), nil
}

// checkPredictiveInputs validates the three probabilities shared by the predictive-value functions.
func checkPredictiveInputs(sensitivity, specificity, prevalence float64) error {
	if err := checkProbability("sensitivity", sensitivity); err != nil {
		return err
	}
	if err := checkProbability("specificity", specificity); err != nil {
		return err
	}
	return checkProbability("prevalence", prevalence)
}
