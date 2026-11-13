package market

import (
	"fmt"

	"github.com/CrimsonCoder42/calcbench/algo/reduce"
)

// VWAP returns the volume-weighted average price Σ pᵢ·vᵢ / Σ vᵢ of a run of trades, in price units
// (ticks times tickSize). Both sums are compensated. Zero-size trades are allowed and carry no weight.
// Errors: ErrEmpty if priceTicks is empty; ErrLengthMismatch if the lengths differ; ErrNegative or
// ErrOutOfRange for a price outside [0, 2⁵³] ticks or a bad tickSize; ErrNotFinite or ErrNegative for a
// bad size; ErrDomain if every size is 0; ErrOverflow if a sum exceeds float64.
//
// ID:      FIN-MKT-011
// Time:    O(n)   n = number of trades
// Space:   O(1), no allocation
// Source:  Berkowitz, Logue, Noser, "The total cost of transactions on the NYSE", Journal of Finance 43(1):97–112, 1988, doi:10.1111/j.1540-6261.1988.tb02591.x, §I (the volume-weighted average price benchmark)
// Source:  same mathematics as NUM-STATS-026 WeightedMean (Σw·x / Σw) with volume weights on int64 ticks; NIST, Dataplot Reference Manual, "WEIGHTED MEAN", https://www.itl.nist.gov/div898/software/dataplot/refman2/ch2/weigmean.pdf (read 2026-11-12)
// Tags:    market, execution, vwap, volume-weighted-average-price, benchmark, transaction-cost
// Units:   priceTicks in ticks; sizes in one size unit; tickSize in price units per tick; result in price units
// Example: hand-checked: prices 100, 102 ticks, sizes 1, 3, tick 0.5 → (100 + 306)/4·0.5 = 50.75
func VWAP(priceTicks []int64, sizes []float64, tickSize float64) (float64, error) {
	n := len(priceTicks)
	if n == 0 {
		return 0, fmt.Errorf("%w: priceTicks", ErrEmpty)
	}
	if len(sizes) != n {
		return 0, fmt.Errorf("%w: len(priceTicks) = %d, len(sizes) = %d", ErrLengthMismatch, n, len(sizes))
	}
	if err := checkTickSize(tickSize); err != nil {
		return 0, err
	}
	var pv, pvc, vs, vsc float64
	for i := range n {
		if err := checkTicks("priceTicks", priceTicks[i]); err != nil {
			return 0, fmt.Errorf("trade %d: %w", i, err)
		}
		if err := checkSize("sizes", sizes[i]); err != nil {
			return 0, fmt.Errorf("trade %d: %w", i, err)
		}
		pv, pvc = reduce.NeumaierAdd(pv, pvc, float64(priceTicks[i])*sizes[i])
		vs, vsc = reduce.NeumaierAdd(vs, vsc, sizes[i])
	}
	num, den := pv+pvc, vs+vsc
	if !reduce.AllFinite(num, den) {
		return 0, fmt.Errorf("%w: Σ p·v or Σ v", ErrOverflow)
	}
	if den == 0 {
		return 0, fmt.Errorf("%w: every size is 0", ErrDomain)
	}
	v := num / den * tickSize
	if !reduce.AllFinite(v) {
		return 0, fmt.Errorf("%w: VWAP", ErrOverflow)
	}
	return v, nil
}
