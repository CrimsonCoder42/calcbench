package invest

import (
	"fmt"

	"github.com/CrimsonCoder42/calcbench/num/prob"
)

// ParametricVaR returns the normal (variance-covariance) value at risk at the given confidence, as a
// positive loss fraction: −(μ + σ·Φ⁻¹(1 − confidence)) = σ·z_c − μ, with Φ⁻¹ the standard normal
// quantile (NUM-PROB-061 NormalQuantile). The confidence level is the caller's. σ = 0 gives −μ. Scale
// μ and σ to the horizon before the call (σ·√h under i.i.d. returns).
// Errors: ErrNotFinite when meanReturn or stdDevReturn is NaN or ±Inf; ErrNegative when stdDevReturn < 0;
// ErrDomain when confidence is not in (0, 1); ErrOverflow when the result exceeds float64.
//
// ID:      FIN-INVEST-050
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Jorion, Value at Risk: The New Benchmark for Managing Financial Risk, 3rd ed., McGraw-Hill, 2007, ch. 5 "Computing VAR", §5.2 (VAR from the normal distribution: VAR = −(μ − ασ), α the standard normal deviate for confidence c)
// Tags:    invest, risk, value-at-risk, var, parametric, normal, variance-covariance, delta-normal, no-alloc
// Units:   meanReturn and stdDevReturn decimal per horizon; confidence a probability in (0, 1); result decimal loss per horizon
// Example: hand-checked: μ = 0, σ = 0.02, confidence 0.95 → 0.02 × 1.6448536… = 0.0328970…
func ParametricVaR(meanReturn, stdDevReturn, confidence float64) (float64, error) {
	if err := checkFinite("meanReturn", meanReturn); err != nil {
		return 0, err
	}
	if err := checkFinite("stdDevReturn", stdDevReturn); err != nil {
		return 0, err
	}
	if stdDevReturn < 0 {
		return 0, fmt.Errorf("%w: stdDevReturn = %v", ErrNegative, stdDevReturn)
	}
	if err := checkConfidence(confidence); err != nil {
		return 0, err
	}
	if stdDevReturn == 0 {
		return -meanReturn, nil
	}
	q, err := prob.NormalQuantile(1-confidence, meanReturn, stdDevReturn)
	if err != nil {
		return 0, fmt.Errorf("%w: %w", ErrOverflow, err)
	}
	return -q, nil
}
