package invest

import (
	"errors"
	"fmt"
	"math"

	"github.com/CrimsonCoder42/calcbench/algo/reduce"
)

// Sentinel errors for the return-series functions (FIN-INVEST-032 onward).
var (
	// ErrTooFew: a series has fewer values than the function needs.
	ErrTooFew = errors.New("invest: too few values")
	// ErrLengthMismatch: paired series have different lengths.
	ErrLengthMismatch = errors.New("invest: series have different lengths")
)

// checkSeries returns ErrTooFew when len(values) < minLen and ErrNotFinite (with the index) for NaN or ±Inf.
func checkSeries(name string, values []float64, minLen int) error {
	if len(values) < minLen {
		return fmt.Errorf("%w: len(%s) = %d, need at least %d", ErrTooFew, name, len(values), minLen)
	}
	for i, v := range values {
		if math.IsNaN(v) || math.IsInf(v, 0) {
			return fmt.Errorf("%w: %s[%d] = %v", ErrNotFinite, name, i, v)
		}
	}
	return nil
}

// checkPairedSeries checks both series with checkSeries and returns ErrLengthMismatch when their lengths differ.
func checkPairedSeries(xName string, x []float64, yName string, y []float64, minLen int) error {
	if len(x) != len(y) {
		return fmt.Errorf("%w: len(%s) = %d, len(%s) = %d", ErrLengthMismatch, xName, len(x), yName, len(y))
	}
	if err := checkSeries(xName, x, minLen); err != nil {
		return err
	}
	return checkSeries(yName, y, minLen)
}

// checkReturnSeries is checkSeries plus ErrRate for a return below −1 (a loss of more than everything).
// A return of exactly −1 (total loss) is allowed.
func checkReturnSeries(name string, returns []float64, minLen int) error {
	if err := checkSeries(name, returns, minLen); err != nil {
		return err
	}
	for i, r := range returns {
		if r < -1 {
			return fmt.Errorf("%w: %s[%d] = %v is below -1", ErrRate, name, i, r)
		}
	}
	return nil
}

// checkPriceSeries is checkSeries plus ErrDomain for a price ≤ 0.
func checkPriceSeries(name string, prices []float64, minLen int) error {
	if err := checkSeries(name, prices, minLen); err != nil {
		return err
	}
	for i, p := range prices {
		if p <= 0 {
			return fmt.Errorf("%w: %s[%d] = %v, prices must be > 0", ErrDomain, name, i, p)
		}
	}
	return nil
}

// sizedBuffer returns dst resized to n, reusing its backing array when cap(dst) ≥ n.
func sizedBuffer(dst []float64, n int) []float64 {
	if cap(dst) >= n {
		return dst[:n]
	}
	return make([]float64, n)
}

// SimpleReturns writes the simple net returns Rₜ = Pₜ/Pₜ₋₁ − 1 of a price series into dst and returns
// it, len(prices) − 1 values. dst is reused when cap(dst) ≥ len(prices) − 1, else one new slice is made;
// dst must not overlap prices. For the return of a single pair of cent prices use CapitalGainsYield
// (FIN-INVEST-012) or HoldingPeriodYield (FIN-INVEST-010).
// Errors: ErrTooFew when there are fewer than 2 prices; ErrNotFinite for NaN or ±Inf; ErrDomain for a
// price ≤ 0; ErrOverflow when a ratio exceeds the float64 range.
//
// ID:      FIN-INVEST-032
// Time:    O(n)   n = len(prices)
// Space:   O(1) extra; no allocation when cap(dst) ≥ n − 1, else one allocation for the result
// Source:  Campbell, Lo, MacKinlay, The Econometrics of Financial Markets, Princeton University Press, 1997, §1.4.1 "Definitions and Notation", eq. (1.4.1) (simple net return)
// Tags:    invest, returns, simple-return, price-series, time-series, dst
// Units:   prices in any one currency unit per share (or index points); results decimal per period (0.01 = 1 %)
// Example: hand-checked: prices 100, 110, 99 → 0.10, −0.10
func SimpleReturns(dst, prices []float64) ([]float64, error) {
	if err := checkPriceSeries("prices", prices, 2); err != nil {
		return dst, err
	}
	out := sizedBuffer(dst, len(prices)-1)
	for i := 1; i < len(prices); i++ {
		r := prices[i]/prices[i-1] - 1
		if math.IsInf(r, 0) {
			return dst, fmt.Errorf("%w: prices[%d]/prices[%d] is beyond float64", ErrOverflow, i, i-1)
		}
		out[i-1] = r
	}
	return out, nil
}

// logRatio returns ln(a/b) for a, b > 0 finite, falling back to ln a − ln b when a/b overflows or underflows.
func logRatio(a, b float64) float64 {
	r := a / b
	if r == 0 || math.IsInf(r, 0) {
		return math.Log(a) - math.Log(b)
	}
	return math.Log(r)
}

// sumLog1p returns Σ ln(1 + rᵢ) by compensated summation, and total = true when some rᵢ = −1.
func sumLog1p(returns []float64) (sum float64, total bool) {
	var s, c float64
	for _, r := range returns {
		if r == -1 {
			return 0, true
		}
		s, c = reduce.NeumaierAdd(s, c, math.Log1p(r))
	}
	return s + c, false
}
