package market

import (
	"fmt"

	"github.com/CrimsonCoder42/calcbench/algo/reduce"
)

// OrderFlowImbalance returns Cont, Kukanov and Stoikov's order flow imbalance over a run of N best-quote
// snapshots, OFI = Σ_{n=1}^{N−1} e_n, where each change of the best quotes contributes
//
//	e_n = 1{P^B_n ≥ P^B_{n−1}}·q^B_n − 1{P^B_n ≤ P^B_{n−1}}·q^B_{n−1}
//	    − 1{P^A_n ≤ P^A_{n−1}}·q^A_n + 1{P^A_n ≥ P^A_{n−1}}·q^A_{n−1}.
//
// Positive OFI is net buying pressure (bid queue added or ask queue removed). One snapshot gives 0.
// Snapshot n is (bidTicks[n], askTicks[n], bidSizes[n], askSizes[n]); the slices are read, never modified.
// Errors: ErrEmpty if bidTicks is empty; ErrLengthMismatch if the four lengths differ; ErrNegative or
// ErrOutOfRange for a price outside [0, 2⁵³] ticks; ErrDomain for a crossed snapshot; ErrNotFinite or
// ErrNegative for a bad size; ErrOverflow if the sum exceeds float64.
//
// ID:      FIN-MKT-005
// Time:    O(N)   N = number of snapshots
// Space:   O(1), no allocation
// Source:  Cont, Kukanov, Stoikov, "The price impact of order book events", Journal of Financial Econometrics 12(1):47–88, 2014, doi:10.1093/jjfinec/nbt003, §2.1 (definition of e_n and OFI_k = Σ e_n over the interval)
// Tags:    market, microstructure, order-flow-imbalance, ofi, order-book, price-impact
// Units:   prices in ticks; sizes in one size unit (shares, lots); result in that size unit
// Example: hand-checked: (bid, ask, q^B, q^A) = (100, 101, 5, 7) → (100, 101, 8, 6): e = 8 − 5 − 6 + 7 = 4
func OrderFlowImbalance(bidTicks, askTicks []int64, bidSizes, askSizes []float64) (float64, error) {
	n := len(bidTicks)
	if n == 0 {
		return 0, fmt.Errorf("%w: bidTicks", ErrEmpty)
	}
	if len(askTicks) != n || len(bidSizes) != n || len(askSizes) != n {
		return 0, fmt.Errorf("%w: len(bidTicks) = %d, len(askTicks) = %d, len(bidSizes) = %d, len(askSizes) = %d",
			ErrLengthMismatch, n, len(askTicks), len(bidSizes), len(askSizes))
	}
	var sum, comp float64
	for i := range n {
		if err := checkQuote(bidTicks[i], askTicks[i]); err != nil {
			return 0, fmt.Errorf("snapshot %d: %w", i, err)
		}
		if err := checkSize("bidSizes", bidSizes[i]); err != nil {
			return 0, fmt.Errorf("snapshot %d: %w", i, err)
		}
		if err := checkSize("askSizes", askSizes[i]); err != nil {
			return 0, fmt.Errorf("snapshot %d: %w", i, err)
		}
		if i == 0 {
			continue
		}
		var e float64
		if bidTicks[i] >= bidTicks[i-1] {
			e += bidSizes[i]
		}
		if bidTicks[i] <= bidTicks[i-1] {
			e -= bidSizes[i-1]
		}
		if askTicks[i] <= askTicks[i-1] {
			e -= askSizes[i]
		}
		if askTicks[i] >= askTicks[i-1] {
			e += askSizes[i-1]
		}
		sum, comp = reduce.NeumaierAdd(sum, comp, e)
	}
	ofi := sum + comp
	if !reduce.AllFinite(ofi) {
		return 0, fmt.Errorf("%w: OFI", ErrOverflow)
	}
	return ofi, nil
}
