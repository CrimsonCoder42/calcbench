package economics

import (
	"fmt"
	"math"

	"github.com/CrimsonCoder42/calcbench/algo/reduce"
)

// HerfindahlHirschmanIndex returns the Herfindahl-Hirschman Index of a market, the sum of the squares of
// the firms' market shares stated in percent, from 0 (atomistic) to 10 000 (one firm). The squares are
// summed with Neumaier compensation. The shares are not required to total exactly 100 (published shares
// are rounded); the DOJ/FTC thresholds (1,800; a change of 100) are policy and are not applied here.
// Errors: ErrDomain if marketSharesPercent is empty; ErrNotFinite if a share is NaN or ±Inf;
// ErrOutOfRange if a share is below 0 or above 100.
//
// ID:      FIN-ECON-004
// Time:    O(n)   n = len(marketSharesPercent)
// Space:   O(1), no allocation
// Source:  U.S. Department of Justice and Federal Trade Commission, Merger Guidelines, December 18, 2023, §2.1 "Guideline 1", p. 5 ("the sum of the squares of the market shares … reaching 10,000 in a market with a single firm") and n. 13, https://www.justice.gov/d9/2023-12/2023%20Merger%20Guidelines.pdf (read 2026-11-12)
// Tags:    economics, antitrust, concentration, hhi, market-share, merger
// Units:   each share in percent (0–100); result in HHI points (squared percent)
// Example: Merger Guidelines n. 13: five equal firms (20 % each) → 2,000; six equal firms (16.67 %) → 1,667
func HerfindahlHirschmanIndex(marketSharesPercent []float64) (float64, error) {
	if len(marketSharesPercent) == 0 {
		return 0, fmt.Errorf("%w: marketSharesPercent is empty", ErrDomain)
	}
	var sum, comp float64
	for k, s := range marketSharesPercent {
		if math.IsNaN(s) || math.IsInf(s, 0) {
			return 0, fmt.Errorf("%w: marketSharesPercent[%d] = %v", ErrNotFinite, k, s)
		}
		if s < 0 || s > 100 {
			return 0, fmt.Errorf("%w: marketSharesPercent[%d] = %v, want 0…100", ErrOutOfRange, k, s)
		}
		sum, comp = reduce.NeumaierAdd(sum, comp, float64(s*s))
	}
	return sum + comp, nil
}
