package interest

import "math"

// UniformGradientUniformSeries returns the factor (A/G,i,n) = 1/i − n/((1+i)ⁿ − 1), the level
// end-of-period series equivalent to the gradient series 0, G, 2G, …, (n−1)G, per unit of G. It is
// (F/G,i,n)·(A/F,i,n); when n·ln(1+i) > 1 it is evaluated directly as 1/i − n/((1+i)ⁿ − 1), which has
// no cancellation there. i = 0 gives exactly (n−1)/2, and periods = 1 gives 0.
// The factor is (F/G)/(F/A), a weighted mean of 0, 1, …, n−1 with positive weights (1+i)ⁿ⁻ᵏ, so it
// lies in [0, n−1] and never overflows.
// Errors: ErrRate if ratePerPeriod is NaN, ±Inf or ≤ −1; ErrPeriods if periods < 1.
//
// ID:      FIN-INT-013
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  NCEES, FE Reference Handbook 10.0.1, NCEES, 2020, Engineering Economics, factor table p. 230, Uniform Gradient Uniform Series (A/G, i%, n) = 1/i − n/[(1 + i)^n − 1]
// Tags:    interest, engineering-economy, factor, gradient, uniform-series, arithmetic-gradient
// Units:   ratePerPeriod decimal per period; periods count; result dimensionless (per unit of G)
// Example: i = 1 %, n = 10 → 4.4179 (FE Reference Handbook 10.0.1, factor table p. 233)
func UniformGradientUniformSeries(ratePerPeriod float64, periods int) (float64, error) {
	if err := checkRate("ratePerPeriod", ratePerPeriod); err != nil {
		return 0, err
	}
	if err := checkPeriods("periods", periods, 1); err != nil {
		return 0, err
	}
	if periods == 1 {
		return 0, nil
	}
	if ratePerPeriod == 0 {
		return float64(periods-1) / 2, nil
	}
	x := float64(periods) * math.Log1p(ratePerPeriod)
	if x > 1 {
		return 1/ratePerPeriod - float64(periods)/math.Expm1(x), nil
	}
	return uniformGradientFutureWorth(ratePerPeriod, periods) / uniformSeriesCompoundAmount(ratePerPeriod, periods), nil
}
