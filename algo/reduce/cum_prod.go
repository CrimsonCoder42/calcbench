package reduce

import (
	"fmt"
	"math"
)

// CumProd returns the running products of values: out[i] = values[0] · … · values[i]. The running product
// is held scaled (as in Prod), so out[i] is finite whenever the true prefix product is representable. The
// result is written into dst (see the package comment). An empty input gives an empty result.
// Errors: ErrNotFinite (with the index) when a value is NaN or ±Inf, checked before anything is written;
// ErrOverflow when a prefix product exceeds ±MaxFloat64, in which case dst may have been partly written.
// On error the returned slice is nil.
//
// ID:      ALG-RED-006
// Time:    O(n)   n = len(values)
// Space:   O(1) extra (no allocation when cap(dst) >= n, else one); O(n) for the result
// Source:  Higham, Accuracy and Stability of Numerical Algorithms, 2nd ed., SIAM, 2002, §3.1 (error of a product) and §27.8 (scaling to avoid overflow)
// Tags:    cumprod, prefix-product, running-product, scan, scaling, dst
func CumProd(dst, values []float64) ([]float64, error) {
	if err := checkFinite(values); err != nil {
		return nil, err
	}
	out := resize(dst, len(values))
	mant, exp := 1.0, 0
	for i, x := range values {
		mant, exp = scaledMul(mant, exp, x)
		p := math.Ldexp(mant, exp)
		if math.IsInf(p, 0) {
			return nil, fmt.Errorf("%w: product through values[%d]", ErrOverflow, i)
		}
		out[i] = p
	}
	return out, nil
}
