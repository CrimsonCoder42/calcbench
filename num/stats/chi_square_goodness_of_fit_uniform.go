package stats

import (
	"fmt"
	"math"

	"github.com/CrimsonCoder42/calcbench/algo/reduce"
)

// ChiSquareGoodnessOfFitUniform returns Pearson's chi-square statistic for observed bin counts against
// the uniform hypothesis, every bin expecting E = N/k (N the total count, k the number of bins):
// χ² = Σ (Oᵢ − E)²/E. It is NUM-STATS-033 with equal expected counts, without the slice; the
// deviations Oᵢ − E are formed from the exact integer total, then summed with compensation. The test has
// k − 1 degrees of freedom (no parameter is estimated). observed is read, never modified.
// Errors: ErrEmpty when observed is empty; ErrNegative when a count is negative; ErrZeroTotal when every
// count is 0; ErrOverflow when the total exceeds the int64 range or the statistic the float64 range.
//
// ID:      NUM-STATS-034
// Time:    O(k)   k = len(observed); two passes (total, then the sum)
// Space:   O(1) extra, no allocation
// Source:  NIST/SEMATECH, e-Handbook of Statistical Methods, §1.3.5.15 "Chi-Square Goodness-of-Fit Test", https://www.itl.nist.gov/div898/handbook/eda/section3/eda35f.htm (read 2026-11-12) (χ² = Σ (Oᵢ − Eᵢ)²/Eᵢ with Eᵢ = N·(F(yᵢ₊₁) − F(yᵢ)); for the discrete uniform F over k bins, Eᵢ = N/k)
// Tags:    chi-square, goodness-of-fit, uniform, pearson, hypothesis-test, counts, frequency, inference
// Units:   observed counts; result dimensionless
// Example: 10, 20, 30, 40 → E = 25, χ² = 20; 25, 25, 25, 25 → 0 (hand-checked)
func ChiSquareGoodnessOfFitUniform(observed []int) (float64, error) {
	k := len(observed)
	if k == 0 {
		return 0, fmt.Errorf("%w: observed", ErrEmpty)
	}
	var total int64
	for i, o := range observed {
		if o < 0 {
			return 0, fmt.Errorf("%w: observed[%d] = %d", ErrNegative, i, o)
		}
		if int64(o) > math.MaxInt64-total {
			return 0, fmt.Errorf("%w: total of observed", ErrOverflow)
		}
		total += int64(o)
	}
	if total == 0 {
		return 0, fmt.Errorf("%w: observed", ErrZeroTotal)
	}
	// χ² = (k/N)·Σ (Oᵢ − N/k)² = Σ (k·Oᵢ − N)²/(k·N); k·Oᵢ − N is formed in float64 from exact integers.
	kf, nf := float64(k), float64(total)
	scale := 1 / math.Sqrt(kf*nf)
	var s compensatedSum
	for _, o := range observed {
		r := (kf*float64(o) - nf) * scale
		s.add(r * r)
	}
	chi := s.total()
	if !reduce.AllFinite(chi) {
		return 0, fmt.Errorf("%w: chi-square statistic", ErrOverflow)
	}
	return chi, nil
}
