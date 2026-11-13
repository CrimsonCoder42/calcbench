package float

import "math"

// ReplaceNonFinite returns values with every NaN replaced by nanValue, +Inf by posInfValue and −Inf by
// negInfValue (numpy nan_to_num with explicit replacements). The result is appended to dst[:0], so
// passing dst = values works in place and a dst with enough capacity avoids allocation.
//
// ID:      NUM-FLT-006
// Time:    O(n)   n = len(values)
// Space:   O(1) extra when cap(dst) ≥ n; otherwise O(n) for the result
// Source:  IEEE Std 754-2019, §6.1–6.2 (infinity and NaN encodings; x ≠ x identifies NaN)
// Tags:    floating-point, nan, infinity, cleaning, nan-to-num, slice
// Units:   replacements in the units of values
// Example: ReplaceNonFinite(nil, [1, NaN, +Inf, −Inf], 0, 9, −9) = [1, 0, 9, −9]
func ReplaceNonFinite(dst, values []float64, nanValue, posInfValue, negInfValue float64) []float64 {
	dst = dst[:0]
	for _, v := range values {
		switch {
		case v != v:
			v = nanValue
		case math.IsInf(v, 1):
			v = posInfValue
		case math.IsInf(v, -1):
			v = negInfValue
		}
		dst = append(dst, v)
	}
	return dst
}
