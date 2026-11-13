package invest

import (
	"fmt"

	"github.com/CrimsonCoder42/calcbench/algo/reduce"
)

// EWMAVariance returns the exponentially weighted moving-average variance forecast after the last return,
// by the recursion σ²ₜ₊₁ = λ·σ²ₜ + (1 − λ)·rₜ², started from initialVariance (σ² before the first
// return; the source leaves the seed to the user, so it is a parameter). The decay factor λ is the
// caller's (RiskMetrics uses 0.94 for daily and 0.97 for monthly data). An empty series returns
// initialVariance. returns is read, never modified.
// Errors: ErrNotFinite for NaN or ±Inf in returns or initialVariance; ErrDomain when lambda is not in
// (0, 1); ErrNegative when initialVariance < 0; ErrOverflow when a squared return exceeds float64.
//
// ID:      FIN-INVEST-039
// Time:    O(n)   n = len(returns)
// Space:   O(1), no allocation
// Source:  J.P. Morgan/Reuters, RiskMetrics — Technical Document, 4th ed., 1996, §5.2.1 "Estimating volatility and correlation", eq. [5.3] (recursive EWMA variance, zero mean)
// Tags:    invest, risk, volatility, ewma, riskmetrics, exponential-smoothing, variance-forecast, no-alloc
// Units:   returns decimal per period; lambda dimensionless in (0, 1); initialVariance and result decimal² per period
// Example: hand-checked: λ = 0.94 (RiskMetrics daily), σ²₀ = 0.0001, returns 0.01, −0.02 → 0.0001, then 0.94·0.0001 + 0.06·0.0004 = 0.000118
func EWMAVariance(returns []float64, lambda, initialVariance float64) (float64, error) {
	if err := checkSeries("returns", returns, 0); err != nil {
		return 0, err
	}
	if !reduce.AllFinite(lambda) || lambda <= 0 || lambda >= 1 {
		return 0, fmt.Errorf("%w: lambda = %v, must be in (0, 1)", ErrDomain, lambda)
	}
	if err := checkFinite("initialVariance", initialVariance); err != nil {
		return 0, err
	}
	if initialVariance < 0 {
		return 0, fmt.Errorf("%w: initialVariance = %v", ErrNegative, initialVariance)
	}
	v := initialVariance
	w := 1 - lambda
	for i, r := range returns {
		sq := r * r
		if sq > maxFloat {
			return 0, fmt.Errorf("%w: returns[%d]² is beyond float64", ErrOverflow, i)
		}
		v = lambda*v + w*sq
	}
	return v, nil
}
