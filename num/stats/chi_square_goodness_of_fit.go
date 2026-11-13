package stats

import (
	"fmt"
	"math"

	"github.com/CrimsonCoder42/calcbench/algo/reduce"
)

// ChiSquareGoodnessOfFit returns Pearson's chi-square statistic χ² = Σ (Oᵢ − Eᵢ)²/Eᵢ for observed bin counts
// Oᵢ and expected counts Eᵢ (the e-Handbook's Eᵢ = N·(F(upper) − F(lower)) for a hypothesised
// distribution F; any positive expected counts are accepted). Compare it with the chi-square quantile at
// 1 − α for k − c degrees of freedom (k non-empty bins, c estimated parameters), or take its p-value from
// the chi-square survival function; neither is computed here. The terms are summed with compensation.
// observed and expected are read, never modified.
// Errors: ErrEmpty when observed is empty; ErrLengthMismatch when len(expected) ≠ len(observed);
// ErrNegative when an observed count is negative; ErrNotFinite when an expected count is NaN or ±Inf;
// ErrOutOfRange when an expected count is ≤ 0; ErrOverflow when the statistic is beyond the float64 range.
//
// ID:      NUM-STATS-033
// Time:    O(k)   k = len(observed)
// Space:   O(1) extra, no allocation
// Source:  NIST/SEMATECH, e-Handbook of Statistical Methods, §1.3.5.15 "Chi-Square Goodness-of-Fit Test", https://www.itl.nist.gov/div898/handbook/eda/section3/eda35f.htm (read 2026-11-12) (χ² = Σᵢ₌₁ᵏ (Oᵢ − Eᵢ)²/Eᵢ; reject if χ² > χ²(1−α, k−c))
// Tags:    chi-square, goodness-of-fit, pearson, hypothesis-test, counts, frequency, inference
// Units:   observed counts; expected counts (may be fractional); result dimensionless
// Example: observed 10, 20, 30, 40 against expected 25 each → (225 + 25 + 25 + 225)/25 = 20 (hand-checked)
func ChiSquareGoodnessOfFit(observed []int, expected []float64) (float64, error) {
	if len(observed) == 0 {
		return 0, fmt.Errorf("%w: observed", ErrEmpty)
	}
	if len(expected) != len(observed) {
		return 0, fmt.Errorf("%w: len(expected) = %d, len(observed) = %d", ErrLengthMismatch, len(expected), len(observed))
	}
	for i, e := range expected {
		if !reduce.AllFinite(e) {
			return 0, fmt.Errorf("%w: expected[%d] = %v", ErrNotFinite, i, e)
		}
		if e <= 0 {
			return 0, fmt.Errorf("%w: expected[%d] = %v, want > 0", ErrOutOfRange, i, e)
		}
	}
	var s compensatedSum
	for i, o := range observed {
		if o < 0 {
			return 0, fmt.Errorf("%w: observed[%d] = %d", ErrNegative, i, o)
		}
		d := float64(o) - expected[i]
		// (d/√e)² instead of d²/e: d² overflows first when e is large.
		r := d / math.Sqrt(expected[i])
		s.add(r * r)
	}
	chi := s.total()
	if !reduce.AllFinite(chi) {
		return 0, fmt.Errorf("%w: chi-square statistic", ErrOverflow)
	}
	return chi, nil
}
