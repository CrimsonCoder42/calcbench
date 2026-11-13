package series

import (
	"fmt"
	"math"

	"github.com/CrimsonCoder42/calcbench/algo/reduce"
)

// AverageRealVariability returns the average real variability of a sequence of readings in time order:
// the mean absolute change between consecutive readings, ARV = (1/(N − 1))·Σₖ₌₁ᴺ⁻¹ |xₖ₊₁ − xₖ|. Order
// matters (it is not a statistic of the set); the caller passes the readings in time order. This is the
// unweighted form; the time-weighted form with each change weighted by the interval between the two
// readings is not computed here. The sum is compensated. values is read, never modified.
// Errors: ErrTooFew when there are fewer than 2 values; ErrNotFinite when a value is NaN or ±Inf;
// ErrOverflow when the result is beyond the float64 range (consecutive values near ∓MaxFloat64).
//
// ID:      NUM-SERIES-009
// Time:    O(n)   n = len(values); one pass (two only when a difference overflows)
// Space:   O(1) extra, no allocation
// Source:  Mena LJ, Felix VG, Melgarejo JD, Maestre GE, "24-Hour Blood Pressure Variability Assessed by Average Real Variability: A Systematic Review and Meta-Analysis", Journal of the American Heart Association, 2017, 6(10):e006895, doi:10.1161/JAHA.117.006895, PMID 29051214, PMC5721878, Introduction, equation 1 ("ARV = 1/(N−1) Σ_{k=1}^{N−1} |BP_{k+1} − BP_k|, where N denotes the number of valid BP measurements, and k is the order of measurements")
// Source:  Kulkarni S, et al., "Blood pressure variability: a review", Journal of Hypertension, 2025, PMID 40084481, PMC12052075, Table 1 (ARV = Σ|xᵢ − xᵢ₋₁|/(N − 1), consecutive readings)
// Tags:    variability, average-real-variability, arv, successive-differences, time-series, dispersion
// Units:   values in any one unit; result in that unit
// Example: 120, 130, 125, 135 → (10 + 5 + 10)/3 = 8.333… (hand-checked)
func AverageRealVariability(values []float64) (float64, error) {
	n := len(values)
	if n < 2 {
		return 0, fmt.Errorf("%w: len(values) = %d, want ≥ 2", ErrTooFew, n)
	}
	for i, v := range values {
		if !reduce.AllFinite(v) {
			return 0, fmt.Errorf("%w: values[%d] = %v", ErrNotFinite, i, v)
		}
	}
	var sum, comp float64
	for i := 1; i < n; i++ {
		sum, comp = reduce.NeumaierAdd(sum, comp, math.Abs(values[i]-values[i-1]))
	}
	if t := sum + comp; reduce.AllFinite(t) {
		return t / float64(n-1), nil
	}
	// A difference or the sum overflowed: redo on values scaled by 2^-k with 2^k ≥ 2n (exact for normal
	// numbers), so each |difference| ≤ MaxFloat64/n and the sum cannot overflow.
	k := 1
	for 1<<k < 2*n {
		k++
	}
	sum, comp = 0, 0
	for i := 1; i < n; i++ {
		sum, comp = reduce.NeumaierAdd(sum, comp, math.Abs(math.Ldexp(values[i], -k)-math.Ldexp(values[i-1], -k)))
	}
	r := math.Ldexp((sum+comp)/float64(n-1), k)
	if !reduce.AllFinite(r) {
		return 0, fmt.Errorf("%w: average real variability", ErrOverflow)
	}
	return r, nil
}
