package invest

import (
	"fmt"
	"math"
)

// AnnualizedHoldingPeriodYield returns the effective annual rate of a holding-period return earned
// over days days, EAR = (1 + HPR)^m − 1 with m = 365/days holding periods per year (a 365-day year),
// computed as expm1(m·log1p(HPR)) so small returns keep their digits.
// Errors: ErrRate if holdingPeriodYield is NaN, ±Inf or ≤ −1; ErrPeriods if days < 1; ErrOverflow if
// the result exceeds float64.
//
// ID:      FIN-INVEST-011
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Dahlquist, Knight et al., Principles of Finance, 1st ed., OpenStax (Rice University), 2022, §15.1 "Risk and Return to an Individual Asset" (EAR = [1 + Holding Period Percentage Return]^m − 1, m = number of holding periods in a year), https://openstax.org/books/principles-finance/pages/15-1-risk-and-return-to-an-individual-asset
// Tags:    invest, return, annualized-return, effective-annual-rate, holding-period-return
// Units:   holdingPeriodYield decimal for the period; days days; result decimal per year
// Example: hand-checked: HPR 0.1462 over 73 days (m = 5) → 1.1462⁵ − 1 = 0.978345; OpenStax §15.1 gives m = 4 → 0.7260 (365/4 is not a whole number of days)
func AnnualizedHoldingPeriodYield(holdingPeriodYield float64, days int) (float64, error) {
	if err := checkRate("holdingPeriodYield", holdingPeriodYield); err != nil {
		return 0, err
	}
	if days < 1 {
		return 0, fmt.Errorf("%w: days = %d, want ≥ 1", ErrPeriods, days)
	}
	m := 365 / float64(days)
	return finiteResult(math.Expm1(m * math.Log1p(holdingPeriodYield)))
}
