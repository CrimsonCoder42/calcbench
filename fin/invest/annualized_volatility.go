package invest

import (
	"fmt"
	"math"

	"github.com/CrimsonCoder42/calcbench/num/stats"
)

// AnnualizedVolatility returns the annualised volatility of per-period returns, s·√periodsPerYear, with
// s the sample standard deviation (n − 1 divisor, NUM-STATS-021's corrected two-pass form). The
// annualisation factor is the caller's (for example 252 trading days per year). Hull computes s from log
// returns (LogReturns); the function accepts any return series. returns is read, never modified.
// Errors: ErrTooFew when there are fewer than 2 returns; ErrNotFinite for NaN or ±Inf; ErrPeriods when
// periodsPerYear is not finite and > 0; ErrOverflow when the result exceeds float64.
//
// ID:      FIN-INVEST-036
// Time:    O(n)   n = len(returns)
// Space:   O(1), no allocation
// Source:  Hull, Options, Futures, and Other Derivatives, 10th ed., Pearson, 2018, §15.4 "Estimating Volatility from Historical Data" (σ̂ = s/√τ, τ the interval length in years = 1/periodsPerYear)
// Tags:    invest, risk, volatility, annualized-volatility, standard-deviation, historical-volatility, no-alloc
// Units:   returns decimal per period; periodsPerYear periods per year; result decimal per √year
// Example: hand-checked: returns 0.01, −0.01, periodsPerYear = 252 → √(0.0002 × 252) = 0.2244994432…
func AnnualizedVolatility(returns []float64, periodsPerYear float64) (float64, error) {
	if err := checkSeries("returns", returns, 2); err != nil {
		return 0, err
	}
	if err := checkPeriodsPerYear(periodsPerYear); err != nil {
		return 0, err
	}
	sd, err := stats.SampleStandardDeviation(returns)
	if err != nil {
		return 0, fmt.Errorf("%w: %w", ErrOverflow, err)
	}
	return finiteResult(sd * math.Sqrt(periodsPerYear))
}
