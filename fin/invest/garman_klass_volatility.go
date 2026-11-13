package invest

import (
	"fmt"
	"math"

	"github.com/CrimsonCoder42/calcbench/algo/reduce"
)

// checkOHLC checks four paired price series: finite, > 0, and low ≤ open, close ≤ high in every bar.
func checkOHLC(open, high, low, close []float64) error {
	if err := checkHighLow(high, low); err != nil {
		return err
	}
	if err := checkPairedSeries("open", open, "high", high, 1); err != nil {
		return err
	}
	if err := checkPairedSeries("close", close, "high", high, 1); err != nil {
		return err
	}
	for i := range open {
		if open[i] < low[i] || open[i] > high[i] || close[i] < low[i] || close[i] > high[i] {
			return fmt.Errorf("%w: bar %d: open and close must lie in [low, high]", ErrDomain, i)
		}
	}
	return nil
}

// GarmanKlassVolatility returns the per-period volatility estimated from open, high, low and close,
// σ = √( (1/n)·Σ [ ½(ln Hₜ/Lₜ)² − (2 ln 2 − 1)(ln Cₜ/Oₜ)² ] ). Each term is ≥ 0 for a valid bar because
// |ln C/O| ≤ ln H/L and 2 ln 2 − 1 < ½. The constants are the source's. Multiply by √(periods per year)
// to annualise. The inputs are read, never modified.
// Errors: ErrTooFew when there is no bar; ErrLengthMismatch when the lengths differ; ErrNotFinite for NaN
// or ±Inf; ErrDomain for a price ≤ 0, a low above its high, or an open or close outside [low, high].
//
// ID:      FIN-INVEST-041
// Time:    O(n)   n = number of bars
// Space:   O(1), no allocation
// Source:  Garman, Klass, "On the Estimation of Security Price Volatilities from Historical Data", Journal of Business, 1980, 53(1):67–78, doi:10.1086/296072 (the practical estimator ½(u − d)² − (2 ln 2 − 1)c², u − d = ln H/L, c = ln C/O)
// Tags:    invest, risk, volatility, range-based, ohlc, garman-klass, no-alloc
// Units:   prices in any one unit; result decimal per period (bar)
// Example: hand-checked: one bar O = 100, H = 110, L = 100, C = 105 → √(½(ln 1.1)² − (2 ln 2 − 1)(ln 1.05)²) = 0.06018…
func GarmanKlassVolatility(open, high, low, close []float64) (float64, error) {
	if err := checkOHLC(open, high, low, close); err != nil {
		return 0, err
	}
	k := 2*math.Ln2 - 1
	var s, comp float64
	for i := range open {
		hl := logRatio(high[i], low[i])
		co := logRatio(close[i], open[i])
		s, comp = reduce.NeumaierAdd(s, comp, 0.5*hl*hl-k*co*co)
	}
	v := (s + comp) / float64(len(open))
	if v < 0 { // rounding only: every term is ≥ 0 for a valid bar
		v = 0
	}
	return finiteResult(math.Sqrt(v))
}
