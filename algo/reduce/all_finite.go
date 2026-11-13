package reduce

import "math"

// AllFinite reports whether no value is NaN, +Inf or -Inf. With no arguments it returns true. It is
// variadic so one call checks scalars (AllFinite(a, b)) or a slice (AllFinite(xs...)) without allocating.
//
// ID:      ALG-RED-010
// Time:    O(n)   n = number of values; stops at the first non-finite value
// Space:   O(1) (no allocation)
// Source:  IEEE Std 754-2019, IEEE Standard for Floating-Point Arithmetic, §5.7.2 (isFinite, isNaN)
// Tags:    finite, nan, inf, validation, check, no-alloc
func AllFinite(values ...float64) bool {
	for _, x := range values {
		if math.IsNaN(x) || math.IsInf(x, 0) {
			return false
		}
	}
	return true
}
