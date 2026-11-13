package prob

import (
	"fmt"
	"slices"
)

// BayesPosterior appends to dst the posteriors P(Aⱼ | B) = P(B | Aⱼ)·P(Aⱼ) / Σᵢ P(B | Aᵢ)·P(Aᵢ) over a
// partition A₁…Aₙ and returns the extended slice. priors must be a distribution (sum 1 within the
// package tolerance n·2⁻⁵²). The denominator is one compensated sum, so the posteriors sum to 1 up to
// rounding. likelihoods and priors are read, never modified.
// Errors: ErrOutOfRange when the lengths differ or are 0, or when every product is 0 (B is impossible
// under every Aᵢ); ErrNotFinite, ErrProbability, ErrNotDistribution as TotalProbability. On error dst is
// returned with its length unchanged.
//
// ID:      NUM-PROB-009
// Time:    O(n)   n = len(priors)
// Space:   O(1) extra; O(n) for the result, one allocation only if dst lacks capacity
// Source:  Ross, A First Course in Probability, 10th ed., Pearson, 2019, §3.3 "Bayes's Formula" (over a partition)
// Source:  Casella, Berger, Statistical Inference, 2nd ed., Duxbury, 2002, §1.3 "Conditional Probability and Independence" (Bayes' Rule)
// Tags:    probability, bayes, posterior, partition, compensated-summation, rules
// Units:   every input and output: probabilities in [0, 1]
// Example: priors [0.5, 0.5], likelihoods [0.3, 0.6] → [1/3, 2/3]
func BayesPosterior(dst, likelihoods, priors []float64) ([]float64, error) {
	if err := checkPairedLengths("likelihoods", len(likelihoods), "priors", len(priors)); err != nil {
		return dst, err
	}
	if err := checkProbabilities("likelihoods", likelihoods); err != nil {
		return dst, err
	}
	if err := checkDistribution("priors", priors); err != nil {
		return dst, err
	}
	var s compensatedSum
	for i, l := range likelihoods {
		s.addProduct(l, priors[i])
	}
	denom := s.total()
	if denom <= 0 {
		return dst, fmt.Errorf("%w: every likelihood·prior is 0, posterior undefined", ErrOutOfRange)
	}
	out := slices.Grow(dst, len(priors))
	for i, l := range likelihoods {
		out = append(out, min(1, l*priors[i]/denom))
	}
	return out, nil
}
