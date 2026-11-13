package invest

import (
	"fmt"
	"math"
)

// checkFlows returns ErrDomain when cashFlowsCents has fewer than minLen flows and ErrOutOfRange for
// a flow outside ±2⁵³ cents.
func checkFlows(cashFlowsCents []int64, minLen int) error {
	if len(cashFlowsCents) < minLen {
		return fmt.Errorf("%w: cashFlowsCents has %d flows, want ≥ %d", ErrDomain, len(cashFlowsCents), minLen)
	}
	for t, c := range cashFlowsCents {
		if c > maxExactCents || c < -maxExactCents {
			return fmt.Errorf("%w: cashFlowsCents[%d] = %d is outside ±2^53 cents", ErrOutOfRange, t, c)
		}
	}
	return nil
}

// ModifiedInternalRateOfReturn returns the MIRR per period of cashFlowsCents (cashFlowsCents[0] at
// t = 0, outflows negative): (TV/PV)^(1/n) − 1, where PV is the present value at t = 0 of the outflows
// discounted at financeRatePerPeriod, TV the value at t = n of the inflows compounded at
// reinvestRatePerPeriod, and n = len(cashFlowsCents) − 1. With both rates equal to the cost of
// capital this is the textbook three-step MIRR. Both sums use Horner's rule; the root is
// expm1(ln(TV/PV)/n). cashFlowsCents is read only.
// Errors: ErrRate if a rate is NaN, ±Inf or ≤ −1; ErrDomain if there are fewer than 2 flows;
// ErrOutOfRange if a flow is outside ±2⁵³ cents; ErrNoSolution if there is no outflow or no inflow;
// ErrOverflow if a sum or the result exceeds float64.
//
// ID:      FIN-INVEST-020
// Time:    O(n)   n = len(cashFlowsCents)
// Space:   O(1), no allocation
// Source:  Dahlquist, Knight et al., Principles of Finance, 1st ed., OpenStax (Rice University), 2022, §16.4 "Alternative Methods" (modified internal rate of return), https://openstax.org/books/principles-finance/pages/16-4-alternative-methods
// Tags:    invest, capital-budgeting, mirr, modified-internal-rate-of-return, irr, reinvestment
// Units:   cashFlowsCents in cents; rates decimal per period; result decimal per period
// Example: OpenStax §16.4 embroidery machine at 9 %: −$16,000, $2,000, $4,000, $5,000 × 4 → 12 % (hand 0.1200476)
func ModifiedInternalRateOfReturn(cashFlowsCents []int64, financeRatePerPeriod, reinvestRatePerPeriod float64) (float64, error) {
	if err := checkRate("financeRatePerPeriod", financeRatePerPeriod); err != nil {
		return 0, err
	}
	if err := checkRate("reinvestRatePerPeriod", reinvestRatePerPeriod); err != nil {
		return 0, err
	}
	if err := checkFlows(cashFlowsCents, 2); err != nil {
		return 0, err
	}
	n := len(cashFlowsCents) - 1
	vf := 1 / (1 + financeRatePerPeriod)
	gr := 1 + reinvestRatePerPeriod
	var pv, tv float64
	for t := 0; t <= n; t++ { // TV = Σ c⁺ₜ·(1+r)^(n−t): Horner forward in (1+r)
		tv *= gr
		if c := cashFlowsCents[t]; c > 0 {
			tv += float64(c)
		}
	}
	for t := n; t >= 0; t-- { // PV = Σ |c⁻ₜ|·(1+f)^(−t): Horner backward in v
		pv *= vf
		if c := cashFlowsCents[t]; c < 0 {
			pv -= float64(c)
		}
	}
	if pv == 0 || tv == 0 {
		return 0, fmt.Errorf("%w: cashFlowsCents needs at least one outflow and one inflow", ErrNoSolution)
	}
	if _, err := finiteResult(tv); err != nil {
		return 0, err
	}
	if _, err := finiteResult(pv); err != nil {
		return 0, err
	}
	return finiteResult(math.Expm1(math.Log(tv/pv) / float64(n)))
}
