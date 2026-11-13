package interest

import (
	"fmt"
	"math"
)

// AnnuityPeriods returns the number of level end-of-period payments of paymentCents that repay
// presentValueCents at ratePerPeriod: the N that solves PV = PMT·(P/A,i,N), which is
// N = −ln(1 − i·PV/PMT)/ln(1+i), and N = PV/PMT at i = 0. The result is fractional; the last
// payment of a real loan is smaller, so the payment count is ⌈N⌉. presentValueCents = 0 gives 0.
// ln(1 − x) and ln(1+i) are computed with math.Log1p, which keeps small rates exact.
// Errors: ErrOutOfRange if |presentValueCents| or |paymentCents| > 2⁵³; ErrRate if ratePerPeriod is
// NaN, ±Inf or ≤ −1; ErrNegative if presentValueCents < 0; ErrDomain if paymentCents ≤ 0;
// ErrNoSolution if paymentCents ≤ i·PV (the payment never covers the interest); ErrOverflow if the
// result is not finite.
//
// ID:      FIN-INT-043
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  NCEES, FE Reference Handbook 10.0.1, NCEES, 2020, Engineering Economics, factor table p. 230, Uniform Series Present Worth (P/A,i,n) = ((1+i)ⁿ − 1)/(i(1+i)ⁿ); setting (P/A,i,N) = PV/PMT and taking logarithms gives N
// Tags:    interest, time-value, annuity, number-of-periods, nper, loan-term, payoff
// Units:   presentValueCents and paymentCents in cents; ratePerPeriod decimal per period; result periods
// Example: $1,000 at 10 % per period repaid by $200 per period → ln 2/ln 1.1 = 7.2725 periods (hand computation)
func AnnuityPeriods(presentValueCents, paymentCents int64, ratePerPeriod float64) (float64, error) {
	if err := checkCents("presentValueCents", presentValueCents); err != nil {
		return 0, err
	}
	if err := checkCents("paymentCents", paymentCents); err != nil {
		return 0, err
	}
	if err := checkRate("ratePerPeriod", ratePerPeriod); err != nil {
		return 0, err
	}
	if presentValueCents < 0 {
		return 0, fmt.Errorf("%w: presentValueCents = %d", ErrNegative, presentValueCents)
	}
	if paymentCents <= 0 {
		return 0, fmt.Errorf("%w: paymentCents = %d, want > 0", ErrDomain, paymentCents)
	}
	if presentValueCents == 0 {
		return 0, nil
	}
	pv, pmt := float64(presentValueCents), float64(paymentCents)
	if ratePerPeriod == 0 {
		return pv / pmt, nil
	}
	x := ratePerPeriod * pv / pmt
	if x >= 1 {
		return 0, fmt.Errorf("%w: paymentCents = %d does not exceed the first period's interest %.2f cents", ErrNoSolution, paymentCents, ratePerPeriod*pv)
	}
	n := -math.Log1p(-x) / math.Log1p(ratePerPeriod)
	if math.IsNaN(n) || math.IsInf(n, 0) {
		return 0, fmt.Errorf("%w: periods %v", ErrOverflow, n)
	}
	return n, nil
}
