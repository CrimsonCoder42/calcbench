package market

import (
	"fmt"
	"math"

	"github.com/CrimsonCoder42/calcbench/num/stats"
)

// RollSpread returns Roll's implicit effective bid-ask spread s = 2√(−cov(Δpₜ, Δpₜ₋₁)) from a run of
// trade prices, in price units. The serial covariance is the sample covariance (NUM-STATS-038) of the
// pairs (Δpₜ, Δpₜ₋₁), t = 2 … n−1. The estimator is undefined when that covariance is ≥ 0, which
// happens often in practice; that case is ErrDomain, never a silent 0 or NaN.
// Errors: ErrTooFew if fewer than 4 prices (2 pairs); ErrNegative or ErrOutOfRange for a price outside
// [0, 2⁵³] ticks or a bad tickSize; ErrDomain if the serial covariance is ≥ 0; ErrOverflow if the
// result exceeds float64.
//
// ID:      FIN-MKT-014
// Time:    O(n)   n = len(priceTicks)
// Space:   O(n) for the price changes
// Source:  Roll, "A simple implicit measure of the effective bid-ask spread in an efficient market", Journal of Finance 39(4):1127–1139, 1984, doi:10.1111/j.1540-6261.1984.tb03897.x, §I (cov(Δpₜ, Δpₜ₋₁) = −s²/4, so s = 2√(−cov))
// Tags:    market, microstructure, roll-spread, effective-spread, serial-covariance, liquidity
// Units:   priceTicks in ticks; tickSize in price units per tick; result in price units
// Example: hand-checked: 100, 101, 100, 101, 100 ticks, tick 1 → Δp = 1, −1, 1, −1; cov = −4/3; s = 4/√3 ≈ 2.3094
func RollSpread(priceTicks []int64, tickSize float64) (float64, error) {
	n := len(priceTicks)
	if n < 4 {
		return 0, fmt.Errorf("%w: len(priceTicks) = %d, want ≥ 4", ErrTooFew, n)
	}
	if err := checkTickSize(tickSize); err != nil {
		return 0, err
	}
	for i, p := range priceTicks {
		if err := checkTicks("priceTicks", p); err != nil {
			return 0, fmt.Errorf("price %d: %w", i, err)
		}
	}
	dp := make([]float64, n-1)
	for i := range dp {
		dp[i] = float64(priceTicks[i+1] - priceTicks[i])
	}
	cov, err := stats.SampleCovariance(dp[1:], dp[:n-2])
	if err != nil {
		return 0, fmt.Errorf("%w: %w", ErrOverflow, err)
	}
	if cov >= 0 {
		return 0, fmt.Errorf("%w: serial covariance = %v ≥ 0", ErrDomain, cov)
	}
	s := 2 * math.Sqrt(-cov) * tickSize
	if math.IsInf(s, 0) {
		return 0, fmt.Errorf("%w: spread", ErrOverflow)
	}
	return s, nil
}
