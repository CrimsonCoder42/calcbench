package invest

import "fmt"

// ProfitabilityIndex returns PV(inflows)/PV(outflows) of cashFlowsCents at ratePerPeriod
// (cashFlowsCents[0] at t = 0, not discounted; outflows negative); above 1 the project's NPV is
// positive. Both present values use Horner's rule in v = 1/(1+i). cashFlowsCents is read only.
// Errors: ErrRate if ratePerPeriod is NaN, ±Inf or ≤ −1; ErrDomain if cashFlowsCents is empty;
// ErrOutOfRange if a flow is outside ±2⁵³ cents; ErrZeroDenominator if there is no outflow;
// ErrOverflow if a present value exceeds float64.
//
// ID:      FIN-INVEST-021
// Time:    O(n)   n = len(cashFlowsCents)
// Space:   O(1), no allocation
// Source:  Dahlquist, Knight et al., Principles of Finance, 1st ed., OpenStax (Rice University), 2022, §16.4 "Alternative Methods" (PI = PV(Cash Inflows)/PV(Cash Outflows)), https://openstax.org/books/principles-finance/pages/16-4-alternative-methods
// Tags:    invest, capital-budgeting, profitability-index, benefit-cost-ratio, present-value
// Units:   cashFlowsCents in cents; ratePerPeriod decimal per period; result dimensionless
// Example: OpenStax §16.4 embroidery machine at 9 %: $18,836/$16,000 = 1.18 (hand 1.177226)
func ProfitabilityIndex(ratePerPeriod float64, cashFlowsCents []int64) (float64, error) {
	if err := checkRate("ratePerPeriod", ratePerPeriod); err != nil {
		return 0, err
	}
	if err := checkFlows(cashFlowsCents, 1); err != nil {
		return 0, err
	}
	v := 1 / (1 + ratePerPeriod)
	var in, out float64
	for t := len(cashFlowsCents) - 1; t >= 0; t-- {
		in *= v
		out *= v
		if c := cashFlowsCents[t]; c > 0 {
			in += float64(c)
		} else {
			out -= float64(c)
		}
	}
	if _, err := finiteResult(in); err != nil {
		return 0, err
	}
	if _, err := finiteResult(out); err != nil {
		return 0, err
	}
	if out == 0 {
		return 0, fmt.Errorf("%w: cashFlowsCents has no outflow", ErrZeroDenominator)
	}
	return finiteResult(in / out)
}
