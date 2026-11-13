package economics

import (
	"fmt"

	"github.com/CrimsonCoder42/calcbench/algo/reduce"
)

// InflationRate returns the rate of inflation between two price-index levels, currentIndex/previousIndex
// − 1, as a decimal (0.098 = 9.8 %). Both levels must be on the same base. A fall in the index gives a
// negative rate (deflation).
// Errors: ErrNotFinite if either index is NaN or ±Inf; ErrNegative if either index is negative;
// ErrZeroDenominator if previousIndex is 0; ErrOverflow if the rate is beyond float64 (e.g. 1e-300 → 1e308).
//
// ID:      FIN-ECON-002
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Greenlaw, Shapiro & MacDonald, Principles of Economics, 3rd ed., OpenStax (Rice University), 2022, §22.1 "Tracking Inflation" ((level in new year − level in previous year)/level in previous year), https://openstax.org/books/principles-economics-3e/pages/22-1-tracking-inflation (read 2026-11-12)
// Tags:    economics, inflation, price-index, cpi, percent-change
// Units:   both indexes in index points on a common base; result decimal per the interval between them
// Example: OpenStax §22.1: index 100 → 109.8 gives 0.098 (9.8 %)
func InflationRate(previousIndex, currentIndex float64) (float64, error) {
	if err := checkFinite("previousIndex", previousIndex); err != nil {
		return 0, err
	}
	if err := checkFinite("currentIndex", currentIndex); err != nil {
		return 0, err
	}
	if previousIndex < 0 {
		return 0, fmt.Errorf("%w: previousIndex = %v", ErrNegative, previousIndex)
	}
	if currentIndex < 0 {
		return 0, fmt.Errorf("%w: currentIndex = %v", ErrNegative, currentIndex)
	}
	if previousIndex == 0 {
		return 0, fmt.Errorf("%w: previousIndex = 0", ErrZeroDenominator)
	}
	r := currentIndex/previousIndex - 1
	if err := checkResult("inflation rate", r); err != nil {
		return 0, err
	}
	return r, nil
}

// checkResult returns ErrOverflow when a float64 result is NaN or ±Inf.
func checkResult(name string, v float64) error {
	if !reduce.AllFinite(v) {
		return fmt.Errorf("%w: %s = %v", ErrOverflow, name, v)
	}
	return nil
}
