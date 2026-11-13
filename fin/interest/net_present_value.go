package interest

import (
	"fmt"
	"math"

	"github.com/CrimsonCoder42/calcbench/num/algebra"
)

// npvChunk is the number of coefficients npvFloat hands to NUM-ALGEB-001 per call, from a stack array.
const npvChunk = 64

// npvFloat returns Σ cₜ·vᵗ with v = 1/(1+rate), evaluated by Horner's rule from the last flow back
// to t = 0 with NUM-ALGEB-001 algebra.PolynomialEvaluate (one fused multiply-add per step). The flows
// are converted to float64 npvChunk at a time into a stack array (no allocation); the running value
// of the higher-order flows rides as the top coefficient of the next chunk, so the steps are exactly
// one Horner pass over all n flows. It is unrounded; it is NaN when the value exceeds float64. It
// assumes checked inputs.
func npvFloat(rate float64, cashFlowsCents []int64) float64 {
	v := 1 / (1 + rate)
	var buf [npvChunk + 1]float64
	hi := len(cashFlowsCents) // flows [lo, hi) are the chunk; acc carries flows [hi, n)
	acc, carry := 0.0, false
	for hi > 0 {
		lo := max(0, hi-npvChunk)
		k := 0
		for t := lo; t < hi; t++ {
			buf[k] = float64(cashFlowsCents[t])
			k++
		}
		if carry {
			buf[k] = acc
			k++
		}
		p, err := algebra.PolynomialEvaluate(buf[:k], v)
		if err != nil {
			return math.NaN() // the flows and v are finite, so the only error is ErrOverflow
		}
		acc, carry, hi = p, true, lo
	}
	return acc
}

// checkCashFlows returns ErrDomain for an empty slice and ErrOutOfRange for a flow outside ±2⁵³.
func checkCashFlows(cashFlowsCents []int64) error {
	if len(cashFlowsCents) == 0 {
		return fmt.Errorf("%w: cashFlowsCents is empty", ErrDomain)
	}
	for t, c := range cashFlowsCents {
		if c > maxExactCents || c < -maxExactCents {
			return fmt.Errorf("%w: cashFlowsCents[%d] = %d is outside ±2^53 cents", ErrOutOfRange, t, c)
		}
	}
	return nil
}

// NetPresentValue returns Σₜ cashFlowsCents[t]/(1+ratePerPeriod)ᵗ for t = 0 … n−1, rounded half-even
// to the cent. cashFlowsCents[0] is at t = 0 (now, not discounted), cashFlowsCents[t] at the end of
// period t; outflows are negative. cashFlowsCents is read only. The sum is evaluated by Horner's rule
// in the discount factor v = 1/(1+i), whose relative error grows at most linearly in n.
// The result is the float64 value rounded half-even; it equals the correctly rounded exact value
// except when that value lies within about n·ε·|x| of a half cent, and is within 1 cent of it for
// |result| ≤ 2⁴³ cents.
// Errors: ErrRate if ratePerPeriod is NaN, ±Inf or ≤ −1; ErrDomain if cashFlowsCents is empty;
// ErrOutOfRange if a flow is outside ±2⁵³ cents; ErrOverflow if the result exceeds ±2⁵³ cents.
//
// ID:      FIN-INT-038
// Time:    O(n)   n = len(cashFlowsCents)
// Space:   O(1), no allocation
// Source:  J. Dahlquist and R. Knight, Principles of Finance, OpenStax, 2022, §16.2 "Net Present Value (NPV) Method", Eq. 16.1, https://openstax.org/books/principles-finance/pages/16-2-net-present-value-npv-method
// Source:  R. A. DeFusco, D. W. McLeavey, J. E. Pinto and D. E. Runkle, Quantitative Investment Analysis, 3rd ed., Wiley (CFA Institute Investment Series), 2015, ch. 1 "The Time Value of Money", §6.4 "The Present Value of a Series of Unequal Cash Flows"
// Tags:    interest, time-value, npv, net-present-value, discounted-cash-flow, capital-budgeting
// Units:   cashFlowsCents and the result in cents; ratePerPeriod decimal per period
// Example: −16,000; 2,000; 4,000; 5,000 ×4 at 9 % → $2,835.62 (book: $2,835.63, from rounded row values; Principles of Finance §16.2 Table 16.3); Table 16.5 → $26,946.90
func NetPresentValue(ratePerPeriod float64, cashFlowsCents []int64) (int64, error) {
	if err := checkRate("ratePerPeriod", ratePerPeriod); err != nil {
		return 0, err
	}
	if err := checkCashFlows(cashFlowsCents); err != nil {
		return 0, err
	}
	v := npvFloat(ratePerPeriod, cashFlowsCents)
	if math.IsInf(v, 0) || math.IsNaN(v) {
		return 0, fmt.Errorf("%w: net present value is beyond float64", ErrOverflow)
	}
	return roundCents(v)
}
