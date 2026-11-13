package invest

import (
	"math"

	"github.com/CrimsonCoder42/calcbench/algo/reduce"
)

// RealizedVariance returns the realized variance of a period, Σ rⱼ², the sum of the squared intraperiod
// returns (no mean is removed and no divisor applied), by compensated summation. returns is read, never
// modified.
// Errors: ErrTooFew when returns is empty; ErrNotFinite for NaN or ±Inf; ErrOverflow when the sum
// exceeds float64.
//
// ID:      FIN-INVEST-037
// Time:    O(n)   n = len(returns)
// Space:   O(1), no allocation
// Source:  Andersen, Bollerslev, Diebold, Labys, "Modeling and Forecasting Realized Volatility", Econometrica, 2003, 71(2):579–625, doi:10.1111/1468-0262.00418, §2 (realized variance, the sum of squared intraday returns)
// Tags:    invest, risk, volatility, realized-variance, high-frequency, intraday, no-alloc
// Units:   returns decimal per intraperiod interval (log returns in the source); result decimal² per period
// Example: hand-checked: 0.01, −0.02, 0.02 → 0.0001 + 0.0004 + 0.0004 = 0.0009
func RealizedVariance(returns []float64) (float64, error) {
	if err := checkSeries("returns", returns, 1); err != nil {
		return 0, err
	}
	var s, c float64
	for _, r := range returns {
		s, c = reduce.NeumaierAdd(s, c, r*r)
	}
	return finiteResult(s + c)
}

// RealizedVolatility returns the realized volatility of a period, √(Σ rⱼ²), the square root of
// RealizedVariance (FIN-INVEST-037). returns is read, never modified.
// Errors: those of RealizedVariance.
//
// ID:      FIN-INVEST-038
// Time:    O(n)   n = len(returns)
// Space:   O(1), no allocation
// Source:  Andersen, Bollerslev, Diebold, Labys, "Modeling and Forecasting Realized Volatility", Econometrica, 2003, 71(2):579–625, doi:10.1111/1468-0262.00418, §2 (realized volatility, the square root of realized variance)
// Tags:    invest, risk, volatility, realized-volatility, high-frequency, intraday, no-alloc
// Units:   returns decimal per intraperiod interval; result decimal per period
// Example: hand-checked: 0.01, −0.02, 0.02 → √0.0009 = 0.03
func RealizedVolatility(returns []float64) (float64, error) {
	v, err := RealizedVariance(returns)
	if err != nil {
		return 0, err
	}
	return math.Sqrt(v), nil
}
