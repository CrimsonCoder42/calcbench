package stats

import (
	"fmt"
	"slices"
)

// Percentiles returns the p-quantile of values for every p in ps, by the Hyndman–Fan definition method,
// from one sort of a copy: result[i] equals Percentile(values, ps[i], method). Use it instead of calling
// Percentile repeatedly when there are more than about log n probabilities. values and ps are read,
// never modified.
// Errors: ErrEmpty when values is empty; ErrNotFinite when a value is NaN or ±Inf; ErrOutOfRange when some
// p is not in [0, 1] (NaN included) or method is not HyndmanFan1 … HyndmanFan9.
//
// ID:      NUM-STATS-042
// Time:    O(n log n + m)   n = len(values), m = len(ps); one sort, then O(1) per p
// Space:   O(n) extra (the sorted copy) and O(m) for the result; two allocations
// Source:  Hyndman, Fan, "Sample quantiles in statistical packages", The American Statistician, 1996, 50(4):361–365, doi:10.1080/00031305.1996.10473566 (eq. 1; definitions 1–3 §2, 4–9 §3 and Table 3)
// Source:  NIST/SEMATECH, e-Handbook of Statistical Methods, §7.2.6.2 "Percentiles", https://www.itl.nist.gov/div898/handbook/prc/section2/prc262.htm (read 2026-11-13)
// Tags:    percentile, quantile, order-statistic, hyndman-fan, descriptive, batch
// Units:   values in any one unit; each p dimensionless in [0, 1]; results in the unit of values
// Example: e-Handbook §7.2.6.2, the 12 resistivity values, ps = 0.90 by type 6 → 95.1981 (as printed)
func Percentiles(values, ps []float64, method QuantileMethod) ([]float64, error) {
	if err := checkSample(values, 1); err != nil {
		return nil, err
	}
	for i, p := range ps {
		if !(p >= 0 && p <= 1) {
			return nil, fmt.Errorf("%w: ps[%d] = %v, want 0..1", ErrOutOfRange, i, p)
		}
	}
	if method < HyndmanFan1 || method > HyndmanFan9 {
		return nil, fmt.Errorf("%w: method = %d, want HyndmanFan1..HyndmanFan9", ErrOutOfRange, int(method))
	}
	sorted := slices.Clone(values)
	slices.Sort(sorted)
	n := len(sorted)
	out := make([]float64, len(ps))
	for i, p := range ps {
		j, gamma := quantilePosition(n, p, method)
		lo := min(max(j, 1), n)
		hi := min(max(j+1, 1), n)
		switch {
		case gamma == 0 || lo == hi:
			if gamma == 1 {
				lo = hi
			}
			out[i] = sorted[lo-1]
		default:
			out[i] = lerp(sorted[lo-1], sorted[hi-1], gamma)
		}
	}
	return out, nil
}
