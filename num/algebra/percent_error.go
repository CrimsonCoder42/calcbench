package algebra

import (
	"fmt"
	"math"

	"github.com/CrimsonCoder42/calcbench/algo/reduce"
)

// PercentError returns the relative error of an observed (approximate) value against the true value, as a
// percent: 100·|observed − trueValue| / |trueValue|.
// Errors: ErrNotFinite for NaN or ±Inf input; ErrDegenerate when trueValue = 0 (relative error is defined
// only for a nonzero true value); ErrOverflow when the result is outside the float64 range.
//
// ID:      NUM-ALGEB-035
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Burden, Faires and Burden, Numerical Analysis, 10th ed., Cengage, 2016, §1.2, Definition 1.15 (relative error |p − p*|/|p|, p ≠ 0)
// Tags:    percent-error, relative-error, error, accuracy, measurement, scalar, algebra
// Units:   observed and trueValue in one unit; the result in percent
// Example: p = 0.3000×10¹, p* = 0.3100×10¹ gives relative error 0.3333×10⁻¹, i.e. 3.333… % (Burden §1.2, Example 1)
func PercentError(observed, trueValue float64) (float64, error) {
	if !reduce.AllFinite(observed, trueValue) {
		return 0, fmt.Errorf("%w: observed = %v, trueValue = %v", ErrNotFinite, observed, trueValue)
	}
	if trueValue == 0 {
		return 0, fmt.Errorf("%w: trueValue = 0", ErrDegenerate)
	}
	return finiteResult(100*(math.Abs(observed-trueValue)/math.Abs(trueValue)), "percent error")
}
