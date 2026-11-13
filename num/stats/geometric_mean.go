package stats

import (
	"fmt"
	"math"
)

// GeometricMean returns the geometric mean (x₁·x₂·…·xₙ)^(1/n) of non-negative values. The product is kept
// scaled: each value is split by math.Frexp into a mantissa in [1/2, 1) and a power of two, the running
// mantissa product is renormalised after every factor and the exponents are summed exactly as an integer
// E, so no product overflows or underflows. The result is (Π m)^(1/n)·2^(E/n), with E/n split into an
// integer part q (applied exactly by math.Ldexp) and a remainder r/n in [0, 1) (math.Exp2). Each factor
// adds one rounding to the mantissa product, and the n-th root divides that relative error by n, so the
// result is within a few ulps whatever the magnitude of the values (unlike exp(mean(ln xᵢ)), whose error
// grows with |ln GM|). A zero value makes the product, and so the result, 0.
// Errors: ErrEmpty for no values; ErrNotFinite for NaN or ±Inf; ErrNegative for a negative value (the
// n-th root of a product with negative factors is not a real mean).
//
// ID:      NUM-STATS-049
// Time:    O(n)   n = len(values); one pass, two math.Frexp calls per value
// Space:   O(1) extra, no allocation
// Source:  NIST, Dataplot Reference Manual, "GEOMETRIC MEAN" (XGM = (∏ Xᵢ)^(1/n)), https://itl.nist.gov/div898/software/dataplot/refman2/auxillar/geommean.htm (read 2026-11-13)
// Source:  IEEE Std 754-2019, §5.3.3 (logB and scaleB: exact power-of-two scaling, here math.Frexp and math.Ldexp)
// Tags:    mean, geometric-mean, average, location, central-tendency, descriptive, scaled-product
// Units:   values and the result in one unit
// Example: 2, 8 → 4; 1, 3, 9 → 3; 10, 2, 38, 23, 38, 23, 21 → 16.412764443111… (hand-checked)
func GeometricMean(values []float64) (float64, error) {
	if len(values) == 0 {
		return 0, ErrEmpty
	}
	if err := checkFinite("values", values); err != nil {
		return 0, err
	}
	zero := false
	prod := 1.0  // running mantissa product, kept in [1/2, 1)
	var exp2 int // exact sum of the binary exponents (|E| ≤ 1100·n, no overflow for any slice length)
	for i, v := range values {
		switch {
		case v < 0:
			return 0, fmt.Errorf("%w: values[%d] = %v", ErrNegative, i, v)
		case v == 0:
			zero = true
		default:
			m, e := math.Frexp(v)
			pm, pe := math.Frexp(prod * m)
			prod, exp2 = pm, exp2+e+pe
		}
	}
	if zero {
		return 0, nil
	}
	n := len(values)
	q, r := exp2/n, exp2%n
	if r < 0 { // floor division, so that 0 ≤ r < n
		q, r = q-1, r+n
	}
	return math.Ldexp(math.Pow(prod, 1/float64(n))*math.Exp2(float64(r)/float64(n)), q), nil
}
