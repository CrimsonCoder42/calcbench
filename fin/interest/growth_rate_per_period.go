package interest

import (
	"fmt"
	"math"
)

// GrowthRatePerPeriod returns the constant compound rate per period that grows presentValueCents to
// futureValueCents in periods periods: g = (FV/PV)^(1/N) − 1, evaluated as expm1(ln(FV/PV)/N). This is
// the compound annual growth rate (CAGR) when the period is a year, the yield of a zero-coupon
// instrument bought at PV and redeemed at FV, and the implied return of a lump sum. A decline gives a
// negative rate. Both amounts must be non-zero and of the same sign.
// Errors: ErrOutOfRange if |presentValueCents| or |futureValueCents| > 2⁵³; ErrDomain if either is 0
// or their signs differ; ErrPeriods if periods < 1.
//
// ID:      FIN-INT-040
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  R. A. DeFusco, D. W. McLeavey, J. E. Pinto and D. E. Runkle, Quantitative Investment Analysis, 3rd ed., Wiley (CFA Institute Investment Series), 2015, ch. 1 "The Time Value of Money", §7.1 "Solving for Interest Rates and Growth Rates", Eq. 14 g = (FV_N/PV)^(1/N) − 1
// Tags:    interest, growth-rate, cagr, compound-annual-growth-rate, implied-return, zero-coupon-yield
// Units:   presentValueCents and futureValueCents in cents; periods count; result decimal per period
// Example: 10,503.0 → 14,146.4 in 4 years → 0.077291 (QIA ch. 1, Example 17); 8.52 → 7.35 in 5 → −0.029111 (Example 18)
func GrowthRatePerPeriod(presentValueCents, futureValueCents int64, periods int) (float64, error) {
	if err := checkCents("presentValueCents", presentValueCents); err != nil {
		return 0, err
	}
	if err := checkCents("futureValueCents", futureValueCents); err != nil {
		return 0, err
	}
	if presentValueCents == 0 || futureValueCents == 0 || (presentValueCents < 0) != (futureValueCents < 0) {
		return 0, fmt.Errorf("%w: presentValueCents = %d and futureValueCents = %d must be non-zero and of one sign",
			ErrDomain, presentValueCents, futureValueCents)
	}
	if err := checkPeriods("periods", periods, 1); err != nil {
		return 0, err
	}
	ratio := float64(futureValueCents) / float64(presentValueCents)
	return math.Expm1(math.Log(ratio) / float64(periods)), nil
}
