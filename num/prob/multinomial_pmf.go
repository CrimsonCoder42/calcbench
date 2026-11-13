package prob

import (
	"fmt"
	"math"

	"github.com/CrimsonCoder42/calcbench/algo/reduce"
)

// MultinomialPMF returns P(X₁ = k₁, …, X_m = k_m) = n!/(k₁!⋯k_m!) · Π pᵢ^kᵢ, n = Σ kᵢ, for n trials each
// falling in category i with probability pᵢ. It is formed in log space, ln n! − Σ ln kᵢ! + Σ kᵢ ln pᵢ
// (math.Lgamma, the sum compensated), then exponentiated once, so no factorial overflows; a result below
// the smallest subnormal is 0. A category with pᵢ = 0 and kᵢ > 0 gives 0; kᵢ = 0 contributes nothing.
// Errors: ErrOutOfRange when the lengths differ or are 0, or a count is negative, or n exceeds 2⁵³;
// ErrNotFinite, ErrProbability, ErrNotDistribution as the package's distribution check on probs.
//
// ID:      NUM-PROB-131
// Time:    O(m)   m = len(counts)
// Space:   O(1), no allocation
// Source:  Johnson, Kotz, Balakrishnan, Discrete Multivariate Distributions, Wiley, 1997, ch. 35 "Multinomial Distributions" (the PMF)
// Tags:    probability, distribution, multinomial, pmf, discrete, multivariate, categorical
// Units:   counts: counts of trials per category; probs: probabilities summing to 1; result: a probability
// Example: counts (1, 2, 1), probs (0.2, 0.5, 0.3) → 12·0.2·0.25·0.3 = 0.18
func MultinomialPMF(counts []int, probs []float64) (float64, error) {
	if err := checkPairedLengths("counts", len(counts), "probs", len(probs)); err != nil {
		return 0, err
	}
	if err := checkDistribution("probs", probs); err != nil {
		return 0, err
	}
	n := 0
	for i, k := range counts {
		if k < 0 {
			return 0, fmt.Errorf("%w: counts[%d] = %d, want >= 0", ErrOutOfRange, i, k)
		}
		if k > 1<<53-n {
			return 0, fmt.Errorf("%w: total count above 2^53", ErrOutOfRange)
		}
		n += k
	}
	lgN, _ := math.Lgamma(float64(n) + 1)
	sum, comp := lgN, 0.0
	for i, k := range counts {
		if k == 0 {
			continue
		}
		if probs[i] == 0 {
			return 0, nil
		}
		lg, _ := math.Lgamma(float64(k) + 1)
		sum, comp = reduce.NeumaierAdd(sum, comp, -lg)
		sum, comp = reduce.NeumaierAdd(sum, comp, float64(k)*math.Log(probs[i]))
	}
	return min(1, math.Exp(sum+comp)), nil
}
