package float

import (
	"fmt"
	"math"
)

// Ulp returns the unit in the last place of x: the gap 2^(e−52) between consecutive doubles in the
// binade [2^e, 2^(e+1)) that holds |x| (Goldberg's definition), and 2⁻¹⁰⁷⁴ for zero and subnormals.
// It is always positive (numpy's spacing carries the sign of x; Ulp(x) = |spacing(x)| except at
// MaxFloat64, where spacing is +Inf and Ulp is the finite gap 2⁹⁷¹).
//
// ID:      NUM-FLT-005
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Muller et al., Handbook of Floating-Point Arithmetic, 2nd ed., Birkhäuser, 2018, §2.6.1 (Definition 2.6, Goldberg's ulp)
// Source:  IEEE Std 754-2019, §3.3 (binary64 format)
// Tags:    floating-point, ulp, spacing, ieee-754, precision
// Units:   same unit as x
// Range:   x finite. Errors: ErrNotFinite.
// Example: Ulp(1) = 2⁻⁵² = 2.220446049250313e-16; Ulp(0) = 5e-324
func Ulp(x float64) (float64, error) {
	if math.IsNaN(x) || math.IsInf(x, 0) {
		return 0, fmt.Errorf("%w: x = %v", ErrNotFinite, x)
	}
	a := math.Abs(x)
	if a == math.MaxFloat64 {
		return a - math.Nextafter(a, 0), nil
	}
	return math.Nextafter(a, math.Inf(1)) - a, nil
}
