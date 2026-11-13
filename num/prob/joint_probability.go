package prob

// JointProbability returns P(A ∩ B) = P(A | B)·P(B), the multiplication rule for events that need not be
// independent.
// Errors: ErrNotFinite when an input is NaN or ±Inf; ErrProbability when one is outside [0, 1].
//
// ID:      NUM-PROB-004
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Ross, A First Course in Probability, 10th ed., Pearson, 2019, §3.2 "Conditional Probabilities" (the multiplication rule)
// Tags:    probability, joint, multiplication-rule, conditional, rules
// Units:   pAGivenB, pB and result: probabilities in [0, 1]
// Example: two cards without replacement, both aces: (3/51)·(4/52) = 1/221 ≈ 0.004524886877828055
func JointProbability(pAGivenB, pB float64) (float64, error) {
	if err := checkProbability("pAGivenB", pAGivenB); err != nil {
		return 0, err
	}
	if err := checkProbability("pB", pB); err != nil {
		return 0, err
	}
	return pAGivenB * pB, nil
}
