package invest

import (
	"fmt"
	"math"

	"github.com/CrimsonCoder42/calcbench/algo/reduce"
)

// maxFloat is the largest finite float64.
const maxFloat = math.MaxFloat64

// checkHighLow checks paired high and low price series: finite, > 0, low ≤ high, at least one bar.
func checkHighLow(high, low []float64) error {
	if err := checkPairedSeries("high", high, "low", low, 1); err != nil {
		return err
	}
	for i := range high {
		if low[i] <= 0 || high[i] <= 0 {
			return fmt.Errorf("%w: bar %d: prices must be > 0", ErrDomain, i)
		}
		if low[i] > high[i] {
			return fmt.Errorf("%w: bar %d: low %v > high %v", ErrDomain, i, low[i], high[i])
		}
	}
	return nil
}

// ParkinsonVolatility returns the per-period volatility estimated from each period's high and low,
// σ = √( Σ (ln Hₜ/Lₜ)² / (4·n·ln 2) ). The constant 4 ln 2 is the source's: E[(ln H/L)²] = 4 ln 2 · σ² for
// a driftless geometric Brownian motion. Multiply by √(periods per year) to annualise. high and low are
// read, never modified.
// Errors: ErrTooFew when there is no bar; ErrLengthMismatch when len(low) ≠ len(high); ErrNotFinite for
// NaN or ±Inf; ErrDomain for a price ≤ 0 or a low above its high.
//
// ID:      FIN-INVEST-040
// Time:    O(n)   n = number of bars
// Space:   O(1), no allocation
// Source:  Parkinson, "The Extreme Value Method for Estimating the Variance of the Rate of Return", Journal of Business, 1980, 53(1):61–65, doi:10.1086/296071 (the estimator σ² = (1/(4 ln 2))·mean (ln H/L)²)
// Tags:    invest, risk, volatility, range-based, high-low, parkinson, ohlc, no-alloc
// Units:   high and low in any one price unit; result decimal per period (bar)
// Example: hand-checked: one bar H = 110, L = 100 → √((ln 1.1)²/(4 ln 2)) = ln 1.1/√(4 ln 2) = 0.0572395…
func ParkinsonVolatility(high, low []float64) (float64, error) {
	if err := checkHighLow(high, low); err != nil {
		return 0, err
	}
	var s, c float64
	for i := range high {
		hl := logRatio(high[i], low[i])
		s, c = reduce.NeumaierAdd(s, c, hl*hl)
	}
	return finiteResult(math.Sqrt((s + c) / (4 * float64(len(high)) * math.Ln2)))
}
