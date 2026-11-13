package invest

import (
	"fmt"
	"math"

	"github.com/CrimsonCoder42/calcbench/algo/reduce"
)

// checkPeriodsPerYear returns ErrPeriods unless periodsPerYear is finite and > 0.
func checkPeriodsPerYear(periodsPerYear float64) error {
	if !reduce.AllFinite(periodsPerYear) || periodsPerYear <= 0 {
		return fmt.Errorf("%w: periodsPerYear = %v, must be finite and > 0", ErrPeriods, periodsPerYear)
	}
	return nil
}

// AnnualizedReturn returns the geometric-average annual return of n simple per-period returns,
// [Π(1 + Rₜ)]^(periodsPerYear/n) − 1, computed as expm1((periodsPerYear/n)·Σ ln(1 + Rₜ)). The
// annualisation factor is the caller's (for example 252 trading days or 12 months per year). A return of
// −1 (total loss) gives −1. returns is read, never modified.
// Errors: ErrTooFew when returns is empty; ErrNotFinite for NaN or ±Inf; ErrRate for a return below −1;
// ErrPeriods when periodsPerYear is not finite and > 0; ErrOverflow when the result exceeds float64.
//
// ID:      FIN-INVEST-035
// Time:    O(n)   n = len(returns)
// Space:   O(1), no allocation
// Source:  Campbell, Lo, MacKinlay, The Econometrics of Financial Markets, Princeton University Press, 1997, §1.4.1 "Definitions and Notation", eq. (1.4.3) (annualised multiperiod return, the geometric average)
// Tags:    invest, returns, annualized-return, geometric-mean, cagr, time-series, no-alloc
// Units:   returns decimal per period; periodsPerYear periods per year; result decimal per year
// Example: hand-checked: quarterly returns 0.10, 0.10 with periodsPerYear = 4 → 1.1⁴ − 1 = 0.4641
func AnnualizedReturn(returns []float64, periodsPerYear float64) (float64, error) {
	if err := checkReturnSeries("returns", returns, 1); err != nil {
		return 0, err
	}
	if err := checkPeriodsPerYear(periodsPerYear); err != nil {
		return 0, err
	}
	s, total := sumLog1p(returns)
	if total {
		return -1, nil
	}
	return finiteResult(math.Expm1(s * (periodsPerYear / float64(len(returns)))))
}
