package invest

import (
	"math"

	"github.com/CrimsonCoder42/calcbench/algo/reduce"
)

// RogersSatchellVolatility returns the per-period volatility estimated from open, high, low and close,
// σ = √( (1/n)·Σ [ ln(Hₜ/Cₜ)·ln(Hₜ/Oₜ) + ln(Lₜ/Cₜ)·ln(Lₜ/Oₜ) ] ), which, unlike Parkinson and
// Garman–Klass, is unbiased when the price has a drift. Each term is ≥ 0 for a valid bar. Multiply by
// √(periods per year) to annualise. The inputs are read, never modified.
// Errors: ErrTooFew when there is no bar; ErrLengthMismatch when the lengths differ; ErrNotFinite for NaN
// or ±Inf; ErrDomain for a price ≤ 0, a low above its high, or an open or close outside [low, high].
//
// ID:      FIN-INVEST-042
// Time:    O(n)   n = number of bars
// Space:   O(1), no allocation
// Source:  Rogers, Satchell, "Estimating Variance from High, Low and Closing Prices", Annals of Applied Probability, 1991, 1(4):504–512, doi:10.1214/aoap/1177005835 (the drift-independent estimator)
// Tags:    invest, risk, volatility, range-based, ohlc, rogers-satchell, drift-independent, no-alloc
// Units:   prices in any one unit; result decimal per period (bar)
// Example: hand-checked: one bar O = 100, H = 110, L = 95, C = 105 → √(ln(110/105)·ln 1.1 + ln(95/105)·ln 0.95) = 0.09781…
func RogersSatchellVolatility(open, high, low, close []float64) (float64, error) {
	if err := checkOHLC(open, high, low, close); err != nil {
		return 0, err
	}
	var s, comp float64
	for i := range open {
		term := logRatio(high[i], close[i])*logRatio(high[i], open[i]) +
			logRatio(low[i], close[i])*logRatio(low[i], open[i])
		s, comp = reduce.NeumaierAdd(s, comp, term)
	}
	v := (s + comp) / float64(len(open))
	if v < 0 { // rounding only: every term is ≥ 0 for a valid bar
		v = 0
	}
	return finiteResult(math.Sqrt(v))
}
