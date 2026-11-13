package prob

import "fmt"

// ConditionalProbability returns P(A | B) = P(A ∩ B) / P(B), defined only when P(B) > 0.
// Errors: ErrNotFinite when an input is NaN or ±Inf; ErrProbability when one is outside [0, 1];
// ErrOutOfRange when pB = 0 (undefined, never NaN); ErrInconsistent when pAAndB > pB.
//
// ID:      NUM-PROB-006
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Ross, A First Course in Probability, 10th ed., Pearson, 2019, §3.2 "Conditional Probabilities"
// Source:  Casella, Berger, Statistical Inference, 2nd ed., Duxbury, 2002, §1.3 "Conditional Probability and Independence"
// Tags:    probability, conditional, rules
// Units:   pAAndB, pB and result: probabilities in [0, 1]
// Example: one die, A = {6}, B = {even}: (1/6)/(1/2) = 1/3
func ConditionalProbability(pAAndB, pB float64) (float64, error) {
	if err := checkProbability("pAAndB", pAAndB); err != nil {
		return 0, err
	}
	if err := checkProbability("pB", pB); err != nil {
		return 0, err
	}
	if pB == 0 {
		return 0, fmt.Errorf("%w: pB = 0, conditional probability undefined", ErrOutOfRange)
	}
	if pAAndB > pB {
		return 0, fmt.Errorf("%w: pAAndB = %v > pB = %v", ErrInconsistent, pAAndB, pB)
	}
	return pAAndB / pB, nil
}
