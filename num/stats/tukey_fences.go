package stats

import (
	"fmt"
	"math"
)

// Fences are the quartiles of a sample, its interquartile range IQ = Q3 − Q1, and the inner fences
// (Q1 − 1.5·IQ, Q3 + 1.5·IQ) and outer fences (Q1 − 3·IQ, Q3 + 3·IQ). In the source's words, a point
// beyond an inner fence is a "mild outlier" and one beyond an outer fence an "extreme outlier".
type Fences struct {
	Q1, Q3, IQ          float64
	InnerLow, InnerHigh float64
	OuterLow, OuterHigh float64
}

// TukeyFences returns the lower and upper quartiles of values (the 25th and 75th percentiles by the
// chosen Hyndman–Fan method, NUM-STATS-015), the interquartile range and the inner and outer fences
// Q1 − 1.5·IQ, Q3 + 1.5·IQ, Q1 − 3·IQ, Q3 + 3·IQ. The e-Handbook defines the quartiles as the 25th and
// 75th percentiles and its own percentile is type 6; type 7 is DefaultQuantileMethod. values is read,
// never modified.
// Errors: ErrEmpty when values is empty; ErrNotFinite when a value is NaN or ±Inf; ErrOutOfRange when
// method is not HyndmanFan1 … HyndmanFan9; ErrOverflow when IQ or a fence is beyond the float64 range.
//
// ID:      NUM-STATS-032
// Time:    O(n)   n = len(values); two linear-time selections
// Space:   O(n) extra (two allocations, one copy per selection) plus O(log n) stack
// Source:  NIST/SEMATECH, e-Handbook of Statistical Methods, §7.1.6 "What are outliers in the data?", https://www.itl.nist.gov/div898/handbook/prc/section1/prc16.htm (read 2026-11-12) (lower inner fence Q1 − 1.5·IQ, upper inner fence Q3 + 1.5·IQ, lower outer fence Q1 − 3·IQ, upper outer fence Q3 + 3·IQ; Q1, Q3 the 25th and 75th percentiles)
// Source:  Hyndman, Fan, "Sample quantiles in statistical packages", The American Statistician, 1996, 50(4):361–365, doi:10.1080/00031305.1996.10473566 (the percentile definitions, via NUM-STATS-015)
// Tags:    outliers, fences, tukey, quartiles, interquartile-range, iqr, box-plot, robust, descriptive
// Units:   values in any one unit; every field of the result in that unit
// Example: e-Handbook §7.1.6: quartiles 429.75 and 742.25 (IQ 312.5) → inner fences −39.0 and 1211.0, outer −507.75 and 1679.75; 62, 64, 63, 61, 65, 118 with type 7 → Q1 62.25, Q3 64.75, inner high 68.5, outer high 72.25
func TukeyFences(values []float64, method QuantileMethod) (Fences, error) {
	q1, err := Percentile(values, 0.25, method)
	if err != nil {
		return Fences{}, err
	}
	q3, err := Percentile(values, 0.75, method)
	if err != nil {
		return Fences{}, err
	}
	return fencesFromQuartiles(q1, q3)
}

// fencesFromQuartiles computes the fences of §7.1.6 from finite quartiles q1 ≤ q3.
func fencesFromQuartiles(q1, q3 float64) (Fences, error) {
	iq := q3 - q1
	f := Fences{
		Q1: q1, Q3: q3, IQ: iq,
		InnerLow: q1 - 1.5*iq, InnerHigh: q3 + 1.5*iq,
		OuterLow: q1 - 3*iq, OuterHigh: q3 + 3*iq,
	}
	for _, v := range [...]float64{f.IQ, f.InnerLow, f.InnerHigh, f.OuterLow, f.OuterHigh} {
		if math.IsInf(v, 0) {
			return Fences{}, fmt.Errorf("%w: fences of quartiles %v and %v", ErrOverflow, q1, q3)
		}
	}
	return f, nil
}
