package prob

import (
	"fmt"

	"github.com/CrimsonCoder42/calcbench/algo/reduce"
)

// DirichletMeanVariance returns the mean vector and the covariance matrix of the Dirichlet distribution with
// concentrations α (m ≥ 2): with α₀ = Σ αᵢ (compensated sum) and α̃ᵢ = αᵢ/α₀, E[Xᵢ] = α̃ᵢ,
// Var(Xᵢ) = α̃ᵢ(1 − α̃ᵢ)/(α₀ + 1) and Cov(Xᵢ, Xⱼ) = −α̃ᵢ·α̃ⱼ/(α₀ + 1) for i ≠ j. The mean is written to
// meanDst (length m) and the covariance, m×m row-major, to covDst (length m²), each resized as needed
// (num/linalg's dst convention); the returned slices are the results.
// Errors: ErrOutOfRange when m < 2 or an αᵢ ≤ 0; ErrNotFinite for a NaN or ±Inf αᵢ; ErrOverflow when α₀
// overflows. On error both returns are nil.
//
// ID:      NUM-PROB-150
// Time:    O(m²)   m = len(alpha)
// Space:   O(1) extra; O(m + m²) for the results, at most one allocation each if dst lacks capacity
// Source:  Kotz, Balakrishnan, Johnson, Continuous Multivariate Distributions, vol. 1, 2nd ed., Wiley, 2000, ch. 49 "Dirichlet and Inverted Dirichlet Distributions", §2 (moments: E[Xᵢ] = αᵢ/α₀, var and cov with the factor 1/(α₀²(α₀ + 1)))
// Tags:    probability, distribution, dirichlet, mean, variance, covariance, continuous, multivariate, simplex, bayesian
// Units:   alpha: dimensionless concentrations; mean: proportions; covariance: proportions²
// Example: α = (1, 2, 2) → α₀ = 5, mean (0.2, 0.4, 0.4); Var(X₁) = 0.2·0.8/6 = 0.02666…, Cov(X₁, X₂) = −0.08/6 = −0.01333…
func DirichletMeanVariance(meanDst, covDst []float64, alpha []float64) (mean, cov []float64, err error) {
	m := len(alpha)
	if m < 2 {
		return nil, nil, fmt.Errorf("%w: dimension %d, want >= 2", ErrOutOfRange, m)
	}
	for i, a := range alpha {
		if !reduce.AllFinite(a) {
			return nil, nil, fmt.Errorf("%w: alpha[%d] = %v", ErrNotFinite, i, a)
		}
		if a <= 0 {
			return nil, nil, fmt.Errorf("%w: alpha[%d] = %v, want > 0", ErrOutOfRange, i, a)
		}
	}
	a0, err := reduce.Sum(alpha)
	if err != nil {
		return nil, nil, fmt.Errorf("%w: sum of alpha: %w", ErrOverflow, err)
	}
	mean = growTo(meanDst, m)
	for i, a := range alpha {
		mean[i] = a / a0
	}
	inv := 1 / (a0 + 1)
	cov = growTo(covDst, m*m)
	for i, mi := range mean {
		row := cov[i*m : (i+1)*m]
		for j, mj := range mean {
			if i == j {
				row[j] = mi * (1 - mi) * inv
			} else {
				row[j] = -mi * mj * inv
			}
		}
	}
	return mean, cov, nil
}
