package prob

import (
	"fmt"
	"slices"
)

// growTo returns dst resized to length n, reusing its storage when the capacity allows (num/linalg's dst
// convention: the contents are overwritten, not appended to).
func growTo(dst []float64, n int) []float64 {
	return slices.Grow(dst[:0], n)[:n]
}

// MultinomialMeanVariance returns the mean vector and the covariance matrix of the multinomial distribution
// of n trials over m categories with probabilities p: E[Xᵢ] = n·pᵢ, Var(Xᵢ) = n·pᵢ(1 − pᵢ) and
// Cov(Xᵢ, Xⱼ) = −n·pᵢ·pⱼ for i ≠ j. The mean is written to meanDst (length m) and the covariance, m×m
// row-major, to covDst (length m²), each resized as needed (num/linalg's dst convention); the returned
// slices are the results. n = 0 gives zeros.
// Errors: ErrOutOfRange when n < 0, n > 2⁵³ or probs is empty; ErrNotFinite, ErrProbability,
// ErrNotDistribution as the package's distribution check on probs. On error both returns are nil.
//
// ID:      NUM-PROB-149
// Time:    O(m²)   m = len(probs)
// Space:   O(1) extra; O(m + m²) for the results, at most one allocation each if dst lacks capacity
// Source:  Johnson, Kotz, Balakrishnan, Discrete Multivariate Distributions, Wiley, 1997, ch. 35 "Multinomial Distributions", §2 (moments: E[Nᵢ] = n·pᵢ, var(Nᵢ) = n·pᵢ(1 − pᵢ), cov(Nᵢ, Nⱼ) = −n·pᵢ·pⱼ)
// Tags:    probability, distribution, multinomial, mean, variance, covariance, discrete, multivariate, categorical
// Units:   n: trials; probs: probabilities summing to 1; mean: counts; covariance: counts²
// Example: n = 10, p = (0.2, 0.3, 0.5) → mean (2, 3, 5); covariance [1.6 −0.6 −1; −0.6 2.1 −1.5; −1 −1.5 2.5]
func MultinomialMeanVariance(meanDst, covDst []float64, n int, probs []float64) (mean, cov []float64, err error) {
	if n < 0 || n > 1<<53 {
		return nil, nil, fmt.Errorf("%w: n = %d, want 0..2^53", ErrOutOfRange, n)
	}
	m := len(probs)
	if m == 0 {
		return nil, nil, fmt.Errorf("%w: probs is empty", ErrOutOfRange)
	}
	if err := checkDistribution("probs", probs); err != nil {
		return nil, nil, err
	}
	nf := float64(n)
	mean = growTo(meanDst, m)
	cov = growTo(covDst, m*m)
	for i, pi := range probs {
		mean[i] = nf * pi
		row := cov[i*m : (i+1)*m]
		for j, pj := range probs {
			if i == j {
				row[j] = nf * pi * (1 - pi)
			} else {
				row[j] = -nf * pi * pj
			}
		}
	}
	return mean, cov, nil
}
