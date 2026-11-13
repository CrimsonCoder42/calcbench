package prob

import "fmt"

// BayesTheorem returns P(A | B) = P(B | A)·P(A) / P(B), defined only when P(B) > 0.
// Errors: ErrNotFinite when an input is NaN or ±Inf; ErrProbability when one is outside [0, 1];
// ErrOutOfRange when pB = 0; ErrInconsistent when P(B | A)·P(A) > P(B) (a posterior above 1).
//
// ID:      NUM-PROB-008
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Ross, A First Course in Probability, 10th ed., Pearson, 2019, §3.3 "Bayes's Formula"
// Source:  Casella, Berger, Statistical Inference, 2nd ed., Duxbury, 2002, §1.3 "Conditional Probability and Independence" (Bayes' Rule)
// Tags:    probability, bayes, conditional, posterior, rules
// Units:   pBGivenA, pA, pB and result: probabilities in [0, 1]
// Example: P(B|A) = 0.9, P(A) = 0.01, P(B) = 0.0585 → 0.009/0.0585 = 2/13 ≈ 0.15384615384615385
func BayesTheorem(pBGivenA, pA, pB float64) (float64, error) {
	if err := checkProbability("pBGivenA", pBGivenA); err != nil {
		return 0, err
	}
	if err := checkProbability("pA", pA); err != nil {
		return 0, err
	}
	if err := checkProbability("pB", pB); err != nil {
		return 0, err
	}
	if pB == 0 {
		return 0, fmt.Errorf("%w: pB = 0, posterior undefined", ErrOutOfRange)
	}
	joint := pBGivenA * pA
	if joint > pB {
		return 0, fmt.Errorf("%w: pBGivenA·pA = %v > pB = %v", ErrInconsistent, joint, pB)
	}
	return joint / pB, nil
}
