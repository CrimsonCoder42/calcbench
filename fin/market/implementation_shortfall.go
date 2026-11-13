package market

import (
	"fmt"

	"github.com/CrimsonCoder42/calcbench/algo/reduce"
)

// ImplementationShortfall returns Perold's implementation shortfall of an order, the paper portfolio's
// gain minus the real portfolio's, split into its two parts:
//
//	executionCost   = side·Σ qᵢ·(pᵢ − p_d)·tickSize        (the price paid on the shares filled)
//	opportunityCost = side·(X − Σ qᵢ)·(p_end − p_d)·tickSize (the move missed on the shares not filled)
//
// with side = +1 (Buy) or −1 (Sell), p_d the decision price, (pᵢ, qᵢ) the fills, X the order size and
// p_end the price at the end of the horizon. Positive is a cost. Commissions and fees are not included:
// the caller adds them (Perold's fixed costs).
// Errors: ErrOutOfRange if side is not Buy or Sell, a price is above 2⁵³ ticks or tickSize is bad;
// ErrLengthMismatch if the fill lengths differ; ErrNegative or ErrNotFinite for a bad price, size or
// orderSize; ErrDomain if the fills exceed orderSize; ErrOverflow if a result exceeds float64.
//
// ID:      FIN-MKT-012
// Time:    O(n)   n = number of fills
// Space:   O(1), no allocation
// Source:  Perold, "The implementation shortfall: paper versus reality", Journal of Portfolio Management 14(3):4–9, 1988, doi:10.3905/jpm.1988.409150 (shortfall = paper return − real return = execution cost + opportunity cost + fixed costs)
// Tags:    market, execution, implementation-shortfall, transaction-cost, opportunity-cost, slippage, tca
// Units:   prices in ticks; fillSizes and orderSize in shares; tickSize in price units per tick; costs in price units × shares (currency)
// Example: hand-checked: buy 100, decision 1000 ticks, fills 60 @ 1002 and 20 @ 1005, end 1010, tick 0.01 → execution (120 + 100)·0.01 = 2.2, opportunity 20·10·0.01 = 2.0
func ImplementationShortfall(side int8, decisionPriceTicks int64, fillPriceTicks []int64, fillSizes []float64,
	orderSize float64, endPriceTicks int64, tickSize float64) (executionCost, opportunityCost float64, err error) {
	if side != Buy && side != Sell {
		return 0, 0, fmt.Errorf("%w: side = %d, want Buy (+1) or Sell (−1)", ErrOutOfRange, side)
	}
	if err := checkTicks("decisionPriceTicks", decisionPriceTicks); err != nil {
		return 0, 0, err
	}
	if err := checkTicks("endPriceTicks", endPriceTicks); err != nil {
		return 0, 0, err
	}
	if err := checkSize("orderSize", orderSize); err != nil {
		return 0, 0, err
	}
	if err := checkTickSize(tickSize); err != nil {
		return 0, 0, err
	}
	if len(fillSizes) != len(fillPriceTicks) {
		return 0, 0, fmt.Errorf("%w: len(fillPriceTicks) = %d, len(fillSizes) = %d", ErrLengthMismatch, len(fillPriceTicks), len(fillSizes))
	}
	var ex, exc, filled, fc float64
	for i, p := range fillPriceTicks {
		if err := checkTicks("fillPriceTicks", p); err != nil {
			return 0, 0, fmt.Errorf("fill %d: %w", i, err)
		}
		if err := checkSize("fillSizes", fillSizes[i]); err != nil {
			return 0, 0, fmt.Errorf("fill %d: %w", i, err)
		}
		ex, exc = reduce.NeumaierAdd(ex, exc, fillSizes[i]*float64(p-decisionPriceTicks))
		filled, fc = reduce.NeumaierAdd(filled, fc, fillSizes[i])
	}
	filled += fc
	if filled > orderSize {
		return 0, 0, fmt.Errorf("%w: filled %v > orderSize %v", ErrDomain, filled, orderSize)
	}
	s := float64(side)
	executionCost = s * (ex + exc) * tickSize
	opportunityCost = s * (orderSize - filled) * float64(endPriceTicks-decisionPriceTicks) * tickSize
	if !reduce.AllFinite(executionCost, opportunityCost) {
		return 0, 0, fmt.Errorf("%w: shortfall", ErrOverflow)
	}
	return executionCost, opportunityCost, nil
}
