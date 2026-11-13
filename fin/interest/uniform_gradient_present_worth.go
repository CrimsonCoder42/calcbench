package interest

import "math"

// UniformGradientPresentWorth returns the factor (P/G,i,n) = ((1+i)ⁿ − 1)/(i²(1+i)ⁿ) − n/(i(1+i)ⁿ),
// the present worth of the gradient series 0, G, 2G, …, (n−1)G paid at the ends of periods 1…n, per
// unit of G. It is (F/G,i,n)·(1+i)⁻ⁿ; when n·ln(1+i) > 1 it is evaluated as ((P/A,i,n) − n(1+i)⁻ⁿ)/i,
// which has no cancellation there and keeps the limit 1/i² for large n. i = 0 gives exactly
// n(n−1)/2, and periods ≤ 1 gives 0.
// Errors: ErrRate if ratePerPeriod is NaN, ±Inf or ≤ −1; ErrPeriods if periods < 0; ErrOverflow if
// the factor exceeds float64.
//
// ID:      FIN-INT-011
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  NCEES, FE Reference Handbook 10.0.1, NCEES, 2020, Engineering Economics, factor table p. 230, Uniform Gradient Present Worth (P/G, i%, n) = [(1 + i)^n − 1]/[i²(1 + i)^n] − n/[i(1 + i)^n]
// Tags:    interest, engineering-economy, factor, gradient, present-worth, arithmetic-gradient
// Units:   ratePerPeriod decimal per period; periods count; result dimensionless (per unit of G)
// Example: i = 1 %, n = 10 → 41.8435 (FE Reference Handbook 10.0.1, factor table p. 233)
func UniformGradientPresentWorth(ratePerPeriod float64, periods int) (float64, error) {
	if err := checkRate("ratePerPeriod", ratePerPeriod); err != nil {
		return 0, err
	}
	if err := checkPeriods("periods", periods, 0); err != nil {
		return 0, err
	}
	if periods <= 1 {
		return 0, nil
	}
	x := float64(periods) * math.Log1p(ratePerPeriod)
	if x > 1 {
		pa := uniformSeriesPresentWorth(ratePerPeriod, periods)
		return factorResult((pa - float64(periods)*math.Exp(-x)) / ratePerPeriod)
	}
	return factorResult(uniformGradientFutureWorth(ratePerPeriod, periods) * math.Exp(-x))
}
