package invest

import (
	"fmt"

	"github.com/CrimsonCoder42/calcbench/algo/reduce"
)

// KellyFraction returns the growth-optimal fraction of capital to stake on a repeated binary bet that
// pays netOdds per unit staked with probability winProbability and loses the stake otherwise:
// f* = p − q/b, q = 1 − p. A result ≤ 0 means the bet has no edge and nothing should be staked; the value
// is returned as computed so the caller sees the size of the disadvantage. With b = 1 (even money) it is
// Kelly's p − q.
// Errors: ErrNotFinite when an input is NaN or ±Inf; ErrDomain when winProbability is outside [0, 1] or
// netOdds ≤ 0; ErrOverflow when the result exceeds float64.
//
// ID:      FIN-INVEST-052
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Kelly, "A New Interpretation of Information Rate", Bell System Technical Journal, 1956, 35(4):917–926, doi:10.1002/j.1538-7305.1956.tb03809.x (growth rate G = p log(1 + ℓ) + q log(1 − ℓ), maximised at ℓ = p − q for even money)
// Source:  Thorp, "The Kelly Criterion in Blackjack, Sports Betting, and the Stock Market", in Zenios, Ziemba (eds.), Handbook of Asset and Liability Management, vol. 1, Elsevier, 2006, ch. 9, §3 (f* = p − q/b for odds b to 1)
// Tags:    invest, betting, kelly-criterion, position-sizing, growth-optimal, no-alloc
// Units:   winProbability a probability in [0, 1]; netOdds currency won per currency staked; result a fraction of capital
// Example: hand-checked: p = 0.6, b = 1 → 0.2; p = 0.5, b = 2 → 0.25
func KellyFraction(winProbability, netOdds float64) (float64, error) {
	if !reduce.AllFinite(winProbability, netOdds) {
		return 0, fmt.Errorf("%w: winProbability = %v, netOdds = %v", ErrNotFinite, winProbability, netOdds)
	}
	if winProbability < 0 || winProbability > 1 {
		return 0, fmt.Errorf("%w: winProbability = %v, must be in [0, 1]", ErrDomain, winProbability)
	}
	if netOdds <= 0 {
		return 0, fmt.Errorf("%w: netOdds = %v, must be > 0", ErrDomain, netOdds)
	}
	return finiteResult(winProbability - (1-winProbability)/netOdds)
}
