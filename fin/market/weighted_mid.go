package market

import (
	"fmt"

	"github.com/CrimsonCoder42/calcbench/algo/reduce"
)

// WeightedMid returns Stoikov's weighted mid-price M_w = I·P^a + (1 − I)·P^b, with I = Q^b/(Q^b + Q^a)
// the level-1 imbalance (FIN-MKT-004), in price units. A heavy bid queue pulls the estimate toward the
// ask. It is the imbalance-weighted mid only: the paper's full micro-price, the limit of the
// Markov-chain expected mid-price adjustments G¹ + G² + …, is NOT computed here.
// Errors: ErrNegative if a price or size is negative; ErrOutOfRange if a price is above 2⁵³ ticks or
// tickSize ≤ 0; ErrNotFinite if a size or tickSize is NaN or ±Inf; ErrDomain if askTicks < bidTicks
// or both sizes are 0; ErrOverflow if the result exceeds float64.
//
// ID:      FIN-MKT-002
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Stoikov, "The micro-price: a high-frequency estimator of future prices", Quantitative Finance 18(12):1959–1966, 2018, doi:10.1080/14697688.2018.1489139, §2 (imbalance I = Q^b/(Q^b + Q^a); weighted mid-price M_w = I·P^a + (1 − I)·P^b)
// Tags:    market, microstructure, weighted-mid, microprice, imbalance, order-book, fair-value
// Units:   bidTicks, askTicks in ticks; bidSize, askSize in shares (any one size unit); tickSize in price units per tick; result in price units
// Example: hand-checked: bid 100, ask 101 ticks, sizes 300 and 100, tick 0.01 → I = 0.75, M_w = (0.75·101 + 0.25·100)·0.01 = 1.0075
func WeightedMid(bidTicks, askTicks int64, bidSize, askSize, tickSize float64) (float64, error) {
	if err := checkQuote(bidTicks, askTicks); err != nil {
		return 0, err
	}
	if err := checkTickSize(tickSize); err != nil {
		return 0, err
	}
	imb, err := Level1Imbalance(bidSize, askSize)
	if err != nil {
		return 0, err
	}
	m := (float64(bidTicks) + imb*float64(askTicks-bidTicks)) * tickSize
	if !reduce.AllFinite(m) {
		return 0, fmt.Errorf("%w: weighted mid", ErrOverflow)
	}
	return m, nil
}
