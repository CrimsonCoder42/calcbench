package invest

import (
	"fmt"
	"math"

	"github.com/CrimsonCoder42/calcbench/algo/reduce"
)

// activeMoments returns the mean and the sample standard deviation (n − 1 divisor, two passes with
// compensated sums) of the active returns pᵢ − bᵢ, without materialising them.
func activeMoments(portfolioReturns, benchmarkReturns []float64) (mean, sd float64, err error) {
	if err := checkPairedSeries("portfolioReturns", portfolioReturns, "benchmarkReturns", benchmarkReturns, 2); err != nil {
		return 0, 0, err
	}
	n := float64(len(portfolioReturns))
	var s, c float64
	for i, p := range portfolioReturns {
		d := p - benchmarkReturns[i]
		if math.IsInf(d, 0) {
			return 0, 0, fmt.Errorf("%w: active return %d is beyond float64", ErrOverflow, i)
		}
		s, c = reduce.NeumaierAdd(s, c, d)
	}
	mean = (s + c) / n
	s, c = 0, 0
	for i, p := range portfolioReturns {
		e := (p - benchmarkReturns[i]) - mean
		s, c = reduce.NeumaierAdd(s, c, e*e)
	}
	v := (s + c) / (n - 1)
	if math.IsInf(v, 0) || math.IsNaN(v) {
		return 0, 0, fmt.Errorf("%w: active-return variance is beyond float64", ErrOverflow)
	}
	return mean, math.Sqrt(v), nil
}

// TrackingError returns the tracking error of a portfolio against its benchmark, the sample standard
// deviation (n − 1 divisor) of the per-period active returns Rₚ − R_B. Per period: multiply by
// √(periods per year) to annualise. The inputs are read, never modified.
// Errors: ErrLengthMismatch when the lengths differ; ErrTooFew when there are fewer than 2 periods;
// ErrNotFinite for NaN or ±Inf; ErrOverflow when an active return or the variance exceeds float64.
//
// ID:      FIN-INVEST-044
// Time:    O(n)   n = number of periods; two passes
// Space:   O(1), no allocation
// Source:  Grinold, Kahn, Active Portfolio Management, 2nd ed., McGraw-Hill, 2000, ch. 5 "Residual Risk and Return: The Information Ratio" (active risk ψ, the standard deviation of active return)
// Tags:    invest, risk, tracking-error, active-risk, benchmark, portfolio, no-alloc
// Units:   returns decimal per period; result decimal per period
// Example: hand-checked: portfolio 0.02, 0.03, 0.01, benchmark 0.01 each → active 0.01, 0.02, 0 → 0.01
func TrackingError(portfolioReturns, benchmarkReturns []float64) (float64, error) {
	_, sd, err := activeMoments(portfolioReturns, benchmarkReturns)
	return sd, err
}
