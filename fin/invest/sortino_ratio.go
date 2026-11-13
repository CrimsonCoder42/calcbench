package invest

import (
	"fmt"
	"math"

	"github.com/CrimsonCoder42/calcbench/algo/reduce"
)

// SortinoRatio returns (mean return − MAR)/downside deviation, with the downside deviation
// √( (1/n)·Σ min(0, Rₜ − MAR)² ) taken over all n periods. The minimum acceptable return MAR is the
// caller's. Per period: annualising is the caller's. returns is read, never modified.
// Errors: ErrTooFew when returns is empty; ErrNotFinite for NaN or ±Inf; ErrRate for a return below −1 or
// a MAR that is not finite and > −1; ErrZeroDenominator when no return is below MAR; ErrOverflow when an
// intermediate or the result exceeds float64.
//
// ID:      FIN-INVEST-046
// Time:    O(n)   n = len(returns); one pass
// Space:   O(1), no allocation
// Source:  Sortino, van der Meer, "Downside Risk", Journal of Portfolio Management, 1991, 17(4):27–31, doi:10.3905/jpm.1991.409343 (downside deviation below the minimum acceptable return, and the reward-to-downside-risk ratio)
// Source:  Sortino, Price, "Performance Measurement in a Downside Risk Framework", Journal of Investing, 1994, 3(3):59–64, doi:10.3905/joi.3.3.59 (the ratio with the n divisor)
// Tags:    invest, performance, sortino-ratio, downside-deviation, downside-risk, mar, risk-adjusted-return, no-alloc
// Units:   returns and MAR decimal per period; result dimensionless
// Example: hand-checked: returns 0.10, −0.05, 0.02, −0.01, MAR 0 → 0.015/√(0.0026/4) = 0.58834…
func SortinoRatio(returns []float64, minimumAcceptableReturn float64) (float64, error) {
	if err := checkReturnSeries("returns", returns, 1); err != nil {
		return 0, err
	}
	if err := checkRate("minimumAcceptableReturn", minimumAcceptableReturn); err != nil {
		return 0, err
	}
	var s, cs, d, cd float64
	for _, r := range returns {
		s, cs = reduce.NeumaierAdd(s, cs, r)
		if e := r - minimumAcceptableReturn; e < 0 {
			d, cd = reduce.NeumaierAdd(d, cd, e*e)
		}
	}
	n := float64(len(returns))
	dd := math.Sqrt((d + cd) / n)
	if math.IsInf(dd, 0) || math.IsNaN(dd) {
		return 0, fmt.Errorf("%w: downside deviation is beyond float64", ErrOverflow)
	}
	if dd == 0 {
		return 0, fmt.Errorf("%w: no return below minimumAcceptableReturn", ErrZeroDenominator)
	}
	return finiteResult(((s+cs)/n - minimumAcceptableReturn) / dd)
}
