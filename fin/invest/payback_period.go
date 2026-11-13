package invest

import "fmt"

// PaybackPeriod returns the number of periods until the cumulative cash flow first reaches zero,
// with the last period counted fractionally: if the cumulative sum is negative after period t−1 and
// non-negative after period t, the result is (t − 1) + (−cumulative(t−1))/cashFlowsCents[t].
// cashFlowsCents[0] is at t = 0 (usually the negative investment); a series that starts at or above 0
// returns 0. Flows after recovery are ignored. The running sum is exact in int64. Read only.
// Errors: ErrDomain if cashFlowsCents is empty; ErrOutOfRange if a flow is outside ±2⁵³ cents;
// ErrNoSolution if the cumulative sum never reaches 0; ErrOverflow if the running sum leaves int64.
//
// ID:      FIN-INVEST-022
// Time:    O(n)   n = len(cashFlowsCents)
// Space:   O(1), no allocation
// Source:  Dahlquist, Knight et al., Principles of Finance, 1st ed., OpenStax (Rice University), 2022, §16.1 "Payback Period Method", https://openstax.org/books/principles-finance/pages/16-1-payback-period-method
// Source:  Dahlquist, Knight et al., Principles of Finance, 1st ed., OpenStax (Rice University), 2022, §16.4 "Alternative Methods" (fractional year: "5 + 145.72/2,981.34 = 5.05 years"), https://openstax.org/books/principles-finance/pages/16-4-alternative-methods
// Tags:    invest, capital-budgeting, payback-period, cash-flow, break-even
// Units:   cashFlowsCents in cents; result periods
// Example: OpenStax §16.1 Sam's Sporting Goods: −$16,000, $2,000, $4,000, $5,000, $5,000, $5,000 → 4 years
func PaybackPeriod(cashFlowsCents []int64) (float64, error) {
	if err := checkFlows(cashFlowsCents, 1); err != nil {
		return 0, err
	}
	var cum int64
	for t, c := range cashFlowsCents {
		next := cum + c
		if (c > 0 && next < cum) || (c < 0 && next > cum) {
			return 0, fmt.Errorf("%w: cumulative cash flow leaves int64 at t = %d", ErrOverflow, t)
		}
		if next >= 0 {
			if t == 0 {
				return 0, nil
			}
			return float64(t-1) + float64(-cum)/float64(c), nil
		}
		cum = next
	}
	return 0, fmt.Errorf("%w: the cumulative cash flow never reaches 0", ErrNoSolution)
}
