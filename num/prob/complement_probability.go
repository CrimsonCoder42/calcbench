package prob

// ComplementProbability returns P(Aᶜ) = 1 − P(A), the probability that A does not occur.
// The subtraction is exact for p in [0.5, 1] (Sterbenz); a tiny complement of a p close to 1 is only as
// accurate as p itself, so upper tails of a distribution come from its survival function instead.
// Errors: ErrNotFinite when p is NaN or ±Inf; ErrProbability when p is outside [0, 1].
//
// ID:      NUM-PROB-002
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Ross, A First Course in Probability, 10th ed., Pearson, 2019, §2.4 "Some Simple Propositions"
// Source:  Casella, Berger, Statistical Inference, 2nd ed., Duxbury, 2002, §1.2 "Basics of Probability Theory"
// Tags:    probability, complement, rules
// Units:   p and result: probabilities in [0, 1]
// Example: 0.25 → 0.75
func ComplementProbability(p float64) (float64, error) {
	if err := checkProbability("p", p); err != nil {
		return 0, err
	}
	return 1 - p, nil
}
