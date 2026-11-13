package interest

import (
	"fmt"
	"math"
)

// DoublingTime returns the exact number of periods for an amount to double when compounded at
// ratePerPeriod: N = ln 2/ln(1+i), from FV = PV(1+i)ᴺ with FV/PV = 2. The result is fractional and
// not rounded up to a whole period. ln(1+i) is computed with math.Log1p.
// Errors: ErrRate if ratePerPeriod is NaN, ±Inf or ≤ −1; ErrDomain if ratePerPeriod ≤ 0; ErrOverflow
// if the result exceeds float64 (a rate so small that ln(1+i) underflows).
//
// ID:      FIN-INT-034
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  R. A. DeFusco, D. W. McLeavey, J. E. Pinto and D. E. Runkle, Quantitative Investment Analysis, 3rd ed., Wiley (CFA Institute Investment Series), 2015, ch. 1 "The Time Value of Money", §7.2 "Solving for the Number of Periods", N = ln(FV/PV)/ln(1 + r) with FV/PV = 2
// Tags:    interest, doubling-time, number-of-periods, compound, logarithm
// Units:   ratePerPeriod decimal per period; result periods
// Example: 7 % → 10.24 periods (QIA ch. 1, Example 19)
func DoublingTime(ratePerPeriod float64) (float64, error) {
	if err := checkRate("ratePerPeriod", ratePerPeriod); err != nil {
		return 0, err
	}
	if ratePerPeriod <= 0 {
		return 0, fmt.Errorf("%w: ratePerPeriod = %v, an amount doubles only at a rate > 0", ErrDomain, ratePerPeriod)
	}
	return factorResult(math.Ln2 / math.Log1p(ratePerPeriod))
}
