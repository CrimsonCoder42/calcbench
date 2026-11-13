package invest

import "fmt"

// InformationRatio returns the information ratio, the mean per-period active return over its sample
// standard deviation (TrackingError, FIN-INVEST-044): (mean of Rₚ − R_B)/ψ. Per period: multiply by
// √(periods per year) to annualise, as the source quotes it. The inputs are read, never modified.
// Errors: those of TrackingError, and ErrZeroDenominator when the tracking error is 0.
//
// ID:      FIN-INVEST-045
// Time:    O(n)   n = number of periods; two passes
// Space:   O(1), no allocation
// Source:  Grinold, Kahn, Active Portfolio Management, 2nd ed., McGraw-Hill, 2000, ch. 5 "Residual Risk and Return: The Information Ratio" (IR = expected residual return / residual risk, ex post: mean active return / tracking error)
// Tags:    invest, performance, information-ratio, active-return, tracking-error, benchmark, risk-adjusted-return, no-alloc
// Units:   returns decimal per period; result dimensionless, per period
// Example: hand-checked: portfolio 0.02, 0.03, 0.01, benchmark 0.01 each → mean 0.01 / ψ 0.01 = 1
func InformationRatio(portfolioReturns, benchmarkReturns []float64) (float64, error) {
	mean, sd, err := activeMoments(portfolioReturns, benchmarkReturns)
	if err != nil {
		return 0, err
	}
	if sd == 0 {
		return 0, fmt.Errorf("%w: tracking error is 0", ErrZeroDenominator)
	}
	return finiteResult(mean / sd)
}
