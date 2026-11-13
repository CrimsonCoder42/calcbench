package stats

import (
	"fmt"
	"math"

	"github.com/CrimsonCoder42/calcbench/algo/reduce"
)

// WeightedMean returns the weighted mean of values, Σ w·x / Σ w, with weights[i] the weight of
// values[i]. Weights need not be normalised; none may be negative and at least one must be positive.
// Both sums are compensated, and sums beyond the float64 range, or products below its normal range
// (weights and values both near 1e-300, say), are redone on rescaled terms, so the result is returned,
// accurate, for valid input. values and weights are read, never modified.
// Errors: ErrEmpty when values is empty; ErrLengthMismatch when len(weights) ≠ len(values); ErrNotFinite
// when a value or a weight is NaN or ±Inf; ErrNegative when a weight is negative; ErrZeroTotal when
// every weight is 0.
//
// ID:      NUM-STATS-026
// Time:    O(n)   n = len(values); two passes (four only when a plain sum overflows or a product underflows)
// Space:   O(1) extra, no allocation
// Source:  NIST, Dataplot Reference Manual, "WEIGHTED MEAN", https://www.itl.nist.gov/div898/software/dataplot/refman2/ch2/weigmean.pdf (read 2026-11-12)
// Source:  NIST, Dataplot Reference Manual, "WEIGHTED TRIMMED MEAN", https://www.itl.nist.gov/div898/software/dataplot/refman2/auxillar/weigtmea.htm (read 2026-11-12) (weights non-negative, at least one positive)
// Tags:    weighted-mean, weighted-average, mean, location, central-tendency, descriptive
// Units:   values in any one unit; weights dimensionless; result in the unit of values
// Example: values 80, 90 with weights 1, 3 → (80 + 270)/4 = 87.5 (hand-checked)
func WeightedMean(values, weights []float64) (float64, error) {
	if len(values) == 0 {
		return 0, fmt.Errorf("%w: values", ErrEmpty)
	}
	if len(weights) != len(values) {
		return 0, fmt.Errorf("%w: len(weights) = %d, len(values) = %d", ErrLengthMismatch, len(weights), len(values))
	}
	if err := checkFinite("values", values); err != nil {
		return 0, err
	}
	if err := checkFinite("weights", weights); err != nil {
		return 0, err
	}
	wMax := 0.0
	for i, w := range weights {
		if w < 0 {
			return 0, fmt.Errorf("%w: weights[%d] = %v", ErrNegative, i, w)
		}
		wMax = max(wMax, w)
	}
	if wMax == 0 {
		return 0, fmt.Errorf("%w: weights", ErrZeroTotal)
	}
	var num, den compensatedSum
	lost := false // a nonzero product fell below the normal range, so it lost bits or vanished
	for i, v := range values {
		p := weights[i] * v
		if math.Abs(p) < 0x1p-1022 && v != 0 && weights[i] != 0 {
			lost = true
		}
		num.add(p)
		den.add(weights[i])
	}
	if n, d := num.total(), den.total(); !lost && reduce.AllFinite(n) && reduce.AllFinite(d) && d > 0 {
		if r := n / d; reduce.AllFinite(r) {
			return r, nil
		}
	}
	// A product or a sum overflowed, or a product underflowed: normalise the weights by the largest, so
	// each is in [0, 1] and their sum in [1, n], and scale the values by 2^-k so that the largest |value|
	// is in [1/2, 1) (exact: a power of two). A product still below 2^-1022 then errs by at most 2^-1075
	// against a denominator of at least 1, negligible beside values of magnitude up to 1. The weighted
	// mean lies between the least and greatest value, so it is finite.
	vMax := 0.0
	for _, v := range values {
		vMax = max(vMax, math.Abs(v))
	}
	_, k := math.Frexp(vMax)
	k = max(k, minScaleExp)
	scale := math.Ldexp(1, -k)
	num, den = compensatedSum{}, compensatedSum{}
	for i, v := range values {
		w := weights[i] / wMax
		num.add(w * (v * scale))
		den.add(w)
	}
	r := math.Ldexp(num.total()/den.total(), k)
	if math.IsInf(r, 0) {
		return math.Copysign(math.MaxFloat64, r), nil
	}
	return r, nil
}
