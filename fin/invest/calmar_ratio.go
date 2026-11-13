package invest

import "fmt"

// CalmarRatio returns the Calmar ratio, the annualised return over the maximum drawdown,
// annualizedReturn/maxDrawdown. The source measures both over the latest 36 months; the window is the
// caller's (compute the inputs with AnnualizedReturn, FIN-INVEST-035, and MaxDrawdown, FIN-INVEST-047).
// Errors: ErrRate when annualizedReturn is not finite and > −1; ErrNotFinite when maxDrawdown is NaN or
// ±Inf; ErrDomain when maxDrawdown is outside [0, 1]; ErrZeroDenominator when it is 0.
//
// ID:      FIN-INVEST-048
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Young, "Calmar Ratio: A Smoother Tool", Futures, 1991, 20(15):40 (annual return over the maximum drawdown, both over 36 months)
// Tags:    invest, performance, calmar-ratio, drawdown, risk-adjusted-return, no-alloc
// Units:   annualizedReturn decimal per year; maxDrawdown decimal in [0, 1]; result dimensionless
// Example: hand-checked: 0.15 a year, maximum drawdown 0.30 → 0.5
func CalmarRatio(annualizedReturn, maxDrawdown float64) (float64, error) {
	if err := checkRate("annualizedReturn", annualizedReturn); err != nil {
		return 0, err
	}
	if err := checkFinite("maxDrawdown", maxDrawdown); err != nil {
		return 0, err
	}
	if maxDrawdown < 0 || maxDrawdown > 1 {
		return 0, fmt.Errorf("%w: maxDrawdown = %v, must be in [0, 1]", ErrDomain, maxDrawdown)
	}
	if maxDrawdown == 0 {
		return 0, fmt.Errorf("%w: maxDrawdown = 0", ErrZeroDenominator)
	}
	return finiteResult(annualizedReturn / maxDrawdown)
}
