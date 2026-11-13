package prob

// IndependentJointProbability returns P(A ∩ B) = P(A)·P(B), the probability that two independent events
// both occur. Independence is the caller's assumption; it cannot be checked from two numbers.
// Errors: ErrNotFinite when an input is NaN or ±Inf; ErrProbability when one is outside [0, 1].
//
// ID:      NUM-PROB-003
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Ross, A First Course in Probability, 10th ed., Pearson, 2019, §3.4 "Independent Events"
// Source:  Casella, Berger, Statistical Inference, 2nd ed., Duxbury, 2002, §1.3 "Conditional Probability and Independence"
// Tags:    probability, joint, independence, intersection, rules
// Units:   pA, pB and result: probabilities in [0, 1]
// Example: two fair coins both heads: 0.5 · 0.5 = 0.25
func IndependentJointProbability(pA, pB float64) (float64, error) {
	if err := checkProbability("pA", pA); err != nil {
		return 0, err
	}
	if err := checkProbability("pB", pB); err != nil {
		return 0, err
	}
	return pA * pB, nil
}
