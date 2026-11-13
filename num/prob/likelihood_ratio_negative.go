package prob

import "fmt"

// LikelihoodRatioNegative returns the likelihood ratio of a negative test result, LR− = (1 − sensitivity) /
// specificity: how many times more (or less) likely a negative result is with the condition than without it.
// Errors: ErrNotFinite when an input is NaN or ±Inf; ErrProbability when an input is outside [0, 1];
// ErrOutOfRange when specificity = 0 (LR− infinite, or undefined when sensitivity is 1).
//
// ID:      NUM-PROB-067
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Deeks JJ, Altman DG, "Diagnostic tests 4: likelihood ratios", BMJ, 2004, 329(7458):168–169, doi:10.1136/bmj.329.7458.168, PMID 15258077, PMC478236 (negative likelihood ratio = (1 − sensitivity)/specificity)
// Tags:    probability, diagnostic-test, likelihood-ratio, sensitivity, specificity, bayes, evidence
// Units:   sensitivity and specificity probabilities in [0, 1]; result a ratio ≥ 0
// Example: Deeks & Altman 2004 table, cut-off 40 pack-years: sensitivity 42/148, specificity 142/144 → (106/148)/(142/144) ≈ 0.7263
func LikelihoodRatioNegative(sensitivity, specificity float64) (float64, error) {
	if err := checkSensSpec(sensitivity, specificity); err != nil {
		return 0, err
	}
	if specificity == 0 {
		return 0, fmt.Errorf("%w: specificity = 0, LR− infinite or undefined", ErrOutOfRange)
	}
	return (1 - sensitivity) / specificity, nil
}
