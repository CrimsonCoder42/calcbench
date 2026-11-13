package market

import "fmt"

// LeeReadySign classifies each trade by the Lee–Ready algorithm: a trade above the prevailing quote
// midpoint is a buy, below it a sell, and a trade exactly at the midpoint is signed by the tick test
// (FIN-MKT-006) on the trade prices. Quote i must be the quote prevailing for trade i: Lee and Ready
// lag quotes by five seconds against trade reports; aligning the quotes is the caller's step. The
// midpoint comparison is exact in integers (p − bid against ask − p). Signs are written into dst,
// reused when cap(dst) ≥ len(priceTicks), else one new slice is allocated; inputs are never modified.
// Errors: ErrEmpty if priceTicks is empty; ErrLengthMismatch if the three lengths differ; ErrNegative or
// ErrOutOfRange for a price outside [0, 2⁵³] ticks; ErrDomain for a crossed quote.
//
// ID:      FIN-MKT-007
// Time:    O(n)   n = len(priceTicks)
// Space:   O(1) extra; O(n) for the result when dst is too small
// Source:  Lee, Ready, "Inferring trade direction from intraday data", Journal of Finance 46(2):733–746, 1991, doi:10.1111/j.1540-6261.1991.tb02683.x, §I–II (quote test against the midpoint; tick test for midpoint trades; five-second quote lag)
// Tags:    market, microstructure, trade-classification, lee-ready, quote-rule, tick-test, trade-sign
// Units:   all prices in ticks; result signs +1 (Buy), −1 (Sell), 0 (Unsigned)
// Example: hand-checked: quote 100/102 throughout, trades 101, 102, 101, 100, 101 → 0, +1, −1, −1, +1
func LeeReadySign(dst []int8, priceTicks, bidTicks, askTicks []int64) ([]int8, error) {
	n := len(priceTicks)
	if n == 0 {
		return dst[:0], fmt.Errorf("%w: priceTicks", ErrEmpty)
	}
	if len(bidTicks) != n || len(askTicks) != n {
		return dst[:0], fmt.Errorf("%w: len(priceTicks) = %d, len(bidTicks) = %d, len(askTicks) = %d",
			ErrLengthMismatch, n, len(bidTicks), len(askTicks))
	}
	for i := range n {
		if err := checkTicks("priceTicks", priceTicks[i]); err != nil {
			return dst[:0], fmt.Errorf("trade %d: %w", i, err)
		}
		if err := checkQuote(bidTicks[i], askTicks[i]); err != nil {
			return dst[:0], fmt.Errorf("trade %d: %w", i, err)
		}
	}
	dst = resizeSigns(dst, n)
	last := Unsigned
	for i := range n {
		tick := Unsigned
		if i > 0 {
			tick = tickStep(priceTicks[i], priceTicks[i-1], &last)
		}
		up := priceTicks[i] - bidTicks[i]   // distance above the bid
		down := askTicks[i] - priceTicks[i] // distance below the ask
		switch {
		case up > down:
			dst[i] = Buy
		case up < down:
			dst[i] = Sell
		default:
			dst[i] = tick
		}
	}
	return dst, nil
}
