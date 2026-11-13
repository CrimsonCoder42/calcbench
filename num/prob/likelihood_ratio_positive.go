package prob

import "fmt"

// checkSensSpec validates a sensitivity and a specificity: finite and in [0, 1].
func checkSensSpec(sensitivity, specificity float64) error {
	if err := checkProbability("sensitivity", sensitivity); err != nil {
		return err
	}
	return checkProbability("specificity", specificity)
}

// LikelihoodRatioPositive returns the likelihood ratio of a positive test result, LR+ = sensitivity /
// (1 − specificity): how many times more likely a positive result is with the condition than without it.
// Errors: ErrNotFinite when an input is NaN or ±Inf; ErrProbability when an input is outside [0, 1];
// ErrOutOfRange when specificity = 1 (LR+ infinite, or undefined when sensitivity is also 0).
//
// ID:      NUM-PROB-066
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Deeks JJ, Altman DG, "Diagnostic tests 4: likelihood ratios", BMJ, 2004, 329(7458):168–169, doi:10.1136/bmj.329.7458.168, PMID 15258077, PMC478236 (LR+ = sensitivity/(1 − specificity); worked example (42/148)/(2/144) = 20.4)
// Tags:    probability, diagnostic-test, likelihood-ratio, sensitivity, specificity, bayes, evidence
// Units:   sensitivity and specificity probabilities in [0, 1]; result a ratio ≥ 0
// Example: Deeks & Altman 2004: sensitivity 42/148, specificity 142/144 → (42/148)/(2/144) ≈ 20.43
func LikelihoodRatioPositive(sensitivity, specificity float64) (float64, error) {
	if err := checkSensSpec(sensitivity, specificity); err != nil {
		return 0, err
	}
	if specificity == 1 {
		return 0, fmt.Errorf("%w: specificity = 1, LR+ infinite or undefined", ErrOutOfRange)
	}
	return sensitivity / (1 - specificity), nil
}
