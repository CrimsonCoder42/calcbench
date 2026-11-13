package interest

import (
	"fmt"
	"math"
)

// AnnualPercentageYield returns the annual percentage yield of a deposit by the general formula of
// Regulation DD, 12 CFR 1030 Appendix A Part I.A: (1 + Interest/Principal)^(365/Days in term) − 1,
// as a decimal (0.0617 = 6.17 %) and unrounded. The regulation's APY is 100 times this, rounded to
// the nearest 0.01 percentage point (12 CFR 1030.3(f)(1)); that rounding is left to the caller.
// The other Appendix A cases (Part I.E and Part II) are not covered.
// Errors: ErrOutOfRange if |principalCents| or |interestCents| > 2⁵³; ErrDomain if principalCents ≤ 0
// or interestCents ≤ −principalCents; ErrPeriods if daysInTerm < 1; ErrOverflow if the result
// exceeds float64.
//
// ID:      FIN-INT-031
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  12 CFR 1030 (Regulation DD) Appendix A Part I.A, read at https://www.consumerfinance.gov/rules-policy/regulations/1030/a/ (CFPB eRegulations; ecfr.gov blocked automated reading) 2026-11-12: APY = 100[(1 + Interest/Principal)^(365/Days in term) − 1]; rounding 12 CFR 1030.3(f)(1)
// Tags:    interest, deposit, apy, annual-percentage-yield, regulation-dd, truth-in-savings, yield
// Units:   principalCents and interestCents in cents; daysInTerm days; result decimal per year
// Example: $61.68 on $1,000 for 365 days → 6.17 %; $30.37 on $1,000 for 182 days → 6.18 % (Appendix A Part I.A)
func AnnualPercentageYield(principalCents, interestCents int64, daysInTerm int) (float64, error) {
	if err := checkCents("principalCents", principalCents); err != nil {
		return 0, err
	}
	if err := checkCents("interestCents", interestCents); err != nil {
		return 0, err
	}
	if principalCents <= 0 {
		return 0, fmt.Errorf("%w: principalCents = %d, want > 0", ErrDomain, principalCents)
	}
	if interestCents <= -principalCents {
		return 0, fmt.Errorf("%w: interestCents = %d loses the whole principal", ErrDomain, interestCents)
	}
	if err := checkPeriods("daysInTerm", daysInTerm, 1); err != nil {
		return 0, err
	}
	ratio := float64(interestCents) / float64(principalCents)
	return factorResult(math.Expm1(365 / float64(daysInTerm) * math.Log1p(ratio)))
}
