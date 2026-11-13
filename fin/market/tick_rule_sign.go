package market

import "fmt"

// Trade signs written by TickRuleSign and LeeReadySign.
const (
	// Buy marks a buyer-initiated trade.
	Buy int8 = 1
	// Sell marks a seller-initiated trade.
	Sell int8 = -1
	// Unsigned marks a trade the rule cannot classify (no earlier price change).
	Unsigned int8 = 0
)

// TickRuleSign classifies each trade by the tick test: a trade above the previous trade price is a buy
// (uptick), below it a sell (downtick); a trade at the same price (zero tick) takes the sign of the last
// price change. Trades before the first price change are Unsigned. The signs are written into dst,
// reused when cap(dst) ≥ len(priceTicks), else one new slice is allocated; priceTicks is never modified.
// Errors: ErrEmpty if priceTicks is empty; ErrNegative or ErrOutOfRange for a price outside [0, 2⁵³]
// ticks.
//
// ID:      FIN-MKT-006
// Time:    O(n)   n = len(priceTicks)
// Space:   O(1) extra; O(n) for the result when dst is too small
// Source:  Lee, Ready, "Inferring trade direction from intraday data", Journal of Finance 46(2):733–746, 1991, doi:10.1111/j.1540-6261.1991.tb02683.x, §I (the tick test: uptick and zero-uptick buys, downtick and zero-downtick sells)
// Tags:    market, microstructure, trade-classification, tick-rule, tick-test, trade-sign
// Units:   priceTicks in ticks; result signs +1 (Buy), −1 (Sell), 0 (Unsigned)
// Example: hand-checked: 100, 101, 101, 100, 100, 102 → 0, +1, +1, −1, −1, +1
func TickRuleSign(dst []int8, priceTicks []int64) ([]int8, error) {
	if len(priceTicks) == 0 {
		return dst[:0], fmt.Errorf("%w: priceTicks", ErrEmpty)
	}
	for i, p := range priceTicks {
		if err := checkTicks("priceTicks", p); err != nil {
			return dst[:0], fmt.Errorf("trade %d: %w", i, err)
		}
	}
	dst = resizeSigns(dst, len(priceTicks))
	last := Unsigned
	dst[0] = Unsigned
	for i := 1; i < len(priceTicks); i++ {
		dst[i] = tickStep(priceTicks[i], priceTicks[i-1], &last)
	}
	return dst, nil
}

// tickStep returns the tick-test sign of price p after prev, updating last (the last nonzero change).
func tickStep(p, prev int64, last *int8) int8 {
	switch {
	case p > prev:
		*last = Buy
	case p < prev:
		*last = Sell
	}
	return *last
}

// resizeSigns returns dst with length n, reusing its storage when it is large enough.
func resizeSigns(dst []int8, n int) []int8 {
	if cap(dst) >= n {
		return dst[:n]
	}
	return make([]int8, n)
}
