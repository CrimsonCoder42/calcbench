package prob

import "fmt"

// UnionProbability returns P(A ∪ B) = P(A) + P(B) − P(A ∩ B). Mutually exclusive events pass pAAndB = 0.
// The consistency checks are exact comparisons; no tolerance is invented.
// Errors: ErrNotFinite when an input is NaN or ±Inf; ErrProbability when one is outside [0, 1];
// ErrInconsistent when pAAndB > min(pA, pB) or the union would exceed 1.
//
// ID:      NUM-PROB-005
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Ross, A First Course in Probability, 10th ed., Pearson, 2019, §2.4 "Some Simple Propositions"
// Source:  Casella, Berger, Statistical Inference, 2nd ed., Duxbury, 2002, §1.2 "Basics of Probability Theory"
// Tags:    probability, union, inclusion-exclusion, addition-rule, rules
// Units:   pA, pB, pAAndB and result: probabilities in [0, 1]
// Example: one card, A = a heart (13/52), B = a king (4/52), A∩B = 1/52 → 16/52 = 4/13
func UnionProbability(pA, pB, pAAndB float64) (float64, error) {
	if err := checkProbability("pA", pA); err != nil {
		return 0, err
	}
	if err := checkProbability("pB", pB); err != nil {
		return 0, err
	}
	if err := checkProbability("pAAndB", pAAndB); err != nil {
		return 0, err
	}
	if pAAndB > min(pA, pB) {
		return 0, fmt.Errorf("%w: pAAndB = %v > min(pA, pB) = %v", ErrInconsistent, pAAndB, min(pA, pB))
	}
	// pA − pAAndB is exact or nearly so (both in [0, 1], the larger first), then one addition.
	u := (pA - pAAndB) + pB
	if u > 1 {
		return 0, fmt.Errorf("%w: pA + pB − pAAndB = %v > 1", ErrInconsistent, u)
	}
	return u, nil
}
