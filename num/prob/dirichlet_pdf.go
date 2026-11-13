package prob

import (
	"fmt"
	"math"

	"github.com/CrimsonCoder42/calcbench/algo/reduce"
)

// DirichletPDF returns the Dirichlet density f(x) = Γ(α₀)/Π Γ(αᵢ) · Π xᵢ^(αᵢ−1), α₀ = Σ αᵢ, at a point x of
// the simplex (xᵢ ≥ 0, Σ xᵢ = 1 within the package tolerance m·2⁻⁵²). It is formed in log space (math.Lgamma,
// the sum compensated) and exponentiated once. On a face (some xᵢ = 0): 0 when αᵢ > 1, the αᵢ = 1 factor is
// 1, and αᵢ < 1 makes the density unbounded (ErrOverflow).
// Errors: ErrOutOfRange when the lengths differ, m < 2, or an αᵢ ≤ 0; ErrNotFinite for a NaN or ±Inf
// element; ErrProbability when an xᵢ is outside [0, 1]; ErrNotDistribution when x does not sum to 1;
// ErrOverflow when the density exceeds the float64 range.
//
// ID:      NUM-PROB-132
// Time:    O(m)   m = len(x)
// Space:   O(1), no allocation
// Source:  Kotz, Balakrishnan, Johnson, Continuous Multivariate Distributions, vol. 1, 2nd ed., Wiley, 2000, ch. 49 "Dirichlet and Inverted Dirichlet Distributions" (the density)
// Tags:    probability, distribution, dirichlet, pdf, continuous, multivariate, simplex, bayesian
// Units:   x: proportions summing to 1; alpha: dimensionless concentrations; result: density on the simplex
// Example: α = (2, 2, 2), x = (0.2, 0.3, 0.5) → Γ(6)/Γ(2)³ · 0.2·0.3·0.5 = 120·0.03 = 3.6
func DirichletPDF(x, alpha []float64) (float64, error) {
	if err := checkPairedLengths("x", len(x), "alpha", len(alpha)); err != nil {
		return 0, err
	}
	if len(x) < 2 {
		return 0, fmt.Errorf("%w: dimension %d, want >= 2", ErrOutOfRange, len(x))
	}
	for i, a := range alpha {
		if !(a > 0) || math.IsInf(a, 1) {
			// The element name is formatted only here, so the success path does not allocate.
			return 0, checkPositive(fmt.Sprintf("alpha[%d]", i), a)
		}
	}
	if err := checkDistribution("x", x); err != nil {
		return 0, err
	}
	var a0, c0 float64
	for _, a := range alpha {
		a0, c0 = reduce.NeumaierAdd(a0, c0, a)
	}
	lgA0, _ := math.Lgamma(a0 + c0)
	sum, comp := lgA0, 0.0
	for i, a := range alpha {
		lg, _ := math.Lgamma(a)
		sum, comp = reduce.NeumaierAdd(sum, comp, -lg)
		switch {
		case a == 1:
		case x[i] == 0 && a > 1:
			return 0, nil
		case x[i] == 0:
			return 0, fmt.Errorf("%w: x[%d] = 0 with alpha < 1 (unbounded density)", ErrOverflow, i)
		default:
			sum, comp = reduce.NeumaierAdd(sum, comp, (a-1)*math.Log(x[i]))
		}
	}
	d := math.Exp(sum + comp)
	if math.IsInf(d, 0) {
		return 0, fmt.Errorf("%w: Dirichlet density", ErrOverflow)
	}
	return d, nil
}
