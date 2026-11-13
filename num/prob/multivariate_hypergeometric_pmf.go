package prob

import (
	"fmt"
	"math"

	"github.com/CrimsonCoder42/calcbench/algo/reduce"
)

// logChoose returns ln C(n, k) for 0 ≤ k ≤ n through math.Lgamma.
func logChoose(n, k int) float64 {
	a, _ := math.Lgamma(float64(n) + 1)
	b, _ := math.Lgamma(float64(k) + 1)
	c, _ := math.Lgamma(float64(n-k) + 1)
	return a - b - c
}

// MultivariateHypergeometricPMF returns P(X₁ = k₁, …, X_m = k_m) = Π C(Kᵢ, kᵢ) / C(N, n) for n = Σ kᵢ draws
// without replacement from an urn of N = Σ Kᵢ items, Kᵢ of colour i. It is formed in log space
// (math.Lgamma, compensated) and exponentiated once. A kᵢ > Kᵢ gives 0.
// Errors: ErrOutOfRange when the lengths differ or are 0, an element is negative, or N exceeds 2⁵³.
//
// ID:      NUM-PROB-133
// Time:    O(m)   m = len(draws)
// Space:   O(1), no allocation
// Source:  Johnson, Kotz, Balakrishnan, Discrete Multivariate Distributions, Wiley, 1997, ch. 39 "Multivariate Hypergeometric and Related Distributions" (the PMF)
// Tags:    probability, distribution, hypergeometric, multivariate, pmf, discrete, sampling-without-replacement
// Units:   draws, colors: counts of items; result: a probability in [0, 1]
// Example: colors (5, 10, 15), draws (2, 2, 2) → C(5,2)·C(10,2)·C(15,2)/C(30,6) = 47250/593775 ≈ 0.07957559681697612
func MultivariateHypergeometricPMF(draws, colors []int) (float64, error) {
	if err := checkPairedLengths("draws", len(draws), "colors", len(colors)); err != nil {
		return 0, err
	}
	total, n := 0, 0
	for i := range draws {
		if draws[i] < 0 || colors[i] < 0 {
			return 0, fmt.Errorf("%w: draws[%d] = %d, colors[%d] = %d, want >= 0", ErrOutOfRange, i, draws[i], i, colors[i])
		}
		if colors[i] > 1<<53-total {
			return 0, fmt.Errorf("%w: total items above 2^53", ErrOutOfRange)
		}
		total += colors[i]
		n += min(draws[i], colors[i])
	}
	sum, comp := -logChoose(total, n), 0.0
	for i, k := range draws {
		if k > colors[i] {
			return 0, nil
		}
		sum, comp = reduce.NeumaierAdd(sum, comp, logChoose(colors[i], k))
	}
	return min(1, math.Exp(sum+comp)), nil
}
