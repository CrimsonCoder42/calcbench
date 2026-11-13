package invest

import (
	"fmt"
	"math"
)

// DiscountedPaybackPeriod returns the payback period of the flows discounted at ratePerPeriod:
// with dₜ = cashFlowsCents[t]/(1+i)ᵗ, the first t at which Σ_{k≤t} dₖ ≥ 0, counted fractionally as
// (t − 1) + (−cumulative(t−1))/dₜ. A series whose first flow is ≥ 0 returns 0. The discount factor is
// carried as a running product (relative error at most linear in n). cashFlowsCents is read only.
// Errors: ErrRate if ratePerPeriod is NaN, ±Inf or ≤ −1; ErrDomain if cashFlowsCents is empty;
// ErrOutOfRange if a flow is outside ±2⁵³ cents; ErrNoSolution if the discounted cumulative sum
// never reaches 0.
//
// ID:      FIN-INVEST-023
// Time:    O(n)   n = len(cashFlowsCents)
// Space:   O(1), no allocation
// Source:  Dahlquist, Knight et al., Principles of Finance, 1st ed., OpenStax (Rice University), 2022, §16.4 "Alternative Methods" (discounted payback period; fractional year "5 + 145.72/2,981.34 = 5.05 years"), https://openstax.org/books/principles-finance/pages/16-4-alternative-methods
// Tags:    invest, capital-budgeting, discounted-payback-period, payback-period, present-value
// Units:   cashFlowsCents in cents; ratePerPeriod decimal per period; result periods
// Example: OpenStax §16.4 embroidery machine at 9 %: 5 + 145.72/2,981.34 = 5.05 years (hand 5.048876)
func DiscountedPaybackPeriod(ratePerPeriod float64, cashFlowsCents []int64) (float64, error) {
	if err := checkRate("ratePerPeriod", ratePerPeriod); err != nil {
		return 0, err
	}
	if err := checkFlows(cashFlowsCents, 1); err != nil {
		return 0, err
	}
	v := 1 / (1 + ratePerPeriod)
	disc := 1.0
	cum := 0.0
	for t, c := range cashFlowsCents {
		d := float64(c) * disc
		if math.IsInf(d, 0) || math.IsNaN(d) {
			break // a later flow is beyond float64; the sum can no longer be judged
		}
		if cum+d >= 0 {
			if t == 0 {
				return 0, nil
			}
			return float64(t-1) + (-cum)/d, nil
		}
		cum += d
		disc *= v
	}
	return 0, fmt.Errorf("%w: the discounted cumulative cash flow never reaches 0", ErrNoSolution)
}
