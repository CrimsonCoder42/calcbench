package prob

import "fmt"

// PostTestProbability returns the probability of a condition after a test result with likelihood ratio
// LR, from the probability before it: post-test odds = pre-test odds × LR, converted back to a
// probability. It is computed as p·LR/(p·LR + 1 − p) (as p/(p + (1 − p)/LR) when LR ≥ 1), which is the
// same quantity without forming odds, so p = 1 and large LR stay finite. Several results applied one after
// another (calling this again with the previous result) assume the tests are independent given the
// condition; that assumption is the caller's.
// Errors: ErrNotFinite when an input is NaN or ±Inf; ErrProbability when pretestProbability is outside
// [0, 1]; ErrOutOfRange when likelihoodRatio < 0, or when pretestProbability = 1 and likelihoodRatio = 0
// (0/0).
//
// ID:      NUM-PROB-068
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Deeks JJ, Altman DG, "Diagnostic tests 4: likelihood ratios", BMJ, 2004, 329(7458):168–169, doi:10.1136/bmj.329.7458.168, PMID 15258077, PMC478236 (pre-test odds p/(1 − p) × likelihood ratio = post-test odds; probability = odds/(1 + odds); worked example 0.1 → odds 0.11 × 20.43 = 2.27 → 0.69)
// Tags:    probability, diagnostic-test, likelihood-ratio, post-test, pre-test, odds, bayes
// Units:   pretestProbability a probability in [0, 1]; likelihoodRatio a ratio ≥ 0; result a probability in [0, 1]
// Example: Deeks & Altman 2004: pre-test 0.1, LR 20.43 → 2.043/2.943 ≈ 0.6942 (printed 0.69)
func PostTestProbability(pretestProbability, likelihoodRatio float64) (float64, error) {
	if err := checkProbability("pretestProbability", pretestProbability); err != nil {
		return 0, err
	}
	if err := checkFinite("likelihoodRatio", likelihoodRatio); err != nil {
		return 0, err
	}
	if likelihoodRatio < 0 {
		return 0, fmt.Errorf("%w: likelihoodRatio = %v, want ≥ 0", ErrOutOfRange, likelihoodRatio)
	}
	p := pretestProbability
	if likelihoodRatio >= 1 {
		if p == 0 {
			return 0, nil
		}
		return p / (p + (1-p)/likelihoodRatio), nil
	}
	if p == 1 {
		if likelihoodRatio == 0 {
			return 0, fmt.Errorf("%w: pretestProbability = 1 with likelihoodRatio = 0 is 0/0", ErrOutOfRange)
		}
		return 1, nil
	}
	num := p * likelihoodRatio
	return num / (num + (1 - p)), nil
}
