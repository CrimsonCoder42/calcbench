package interest

import (
	"fmt"

	"github.com/CrimsonCoder42/calcbench/algo/reduce"
)

// factorResult returns v, or ErrOverflow when v is ±Inf or NaN.
func factorResult(v float64) (float64, error) {
	if !reduce.AllFinite(v) {
		return 0, fmt.Errorf("%w: factor %v is beyond float64", ErrOverflow, v)
	}
	return v, nil
}

// SinglePaymentCompoundAmount returns the factor (F/P,i,n) = (1+i)ⁿ, the future value of one unit
// paid now. It is computed as exp(n·log1p(i)); periods = 0 gives exactly 1.
// Errors: ErrRate if ratePerPeriod is NaN, ±Inf or ≤ −1; ErrPeriods if periods < 0; ErrOverflow if
// the factor exceeds float64.
//
// ID:      FIN-INT-005
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  NCEES, FE Reference Handbook 10.0.1, NCEES, 2020, Engineering Economics, factor table p. 230, Single Payment Compound Amount (F/P, i%, n) = (1 + i)^n
// Source:  R. A. DeFusco, D. W. McLeavey, J. E. Pinto and D. E. Runkle, Quantitative Investment Analysis, 3rd ed., Wiley (CFA Institute Investment Series), 2015, ch. 1 "The Time Value of Money", §3, Eq. 2
// Tags:    interest, engineering-economy, factor, compound-amount, single-payment
// Units:   ratePerPeriod decimal per period; periods count; result dimensionless (per unit of P)
// Example: i = 1 %, n = 10 → 1.1046 (FE Reference Handbook 10.0.1, factor table p. 233)
func SinglePaymentCompoundAmount(ratePerPeriod float64, periods int) (float64, error) {
	if err := checkRate("ratePerPeriod", ratePerPeriod); err != nil {
		return 0, err
	}
	if err := checkPeriods("periods", periods, 0); err != nil {
		return 0, err
	}
	return factorResult(compound(ratePerPeriod, periods))
}
