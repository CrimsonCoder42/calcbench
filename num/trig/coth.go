package trig

import (
	"fmt"
	"math"

	"github.com/CrimsonCoder42/calcbench/algo/reduce"
)

// Coth returns the hyperbolic cotangent of x, 1/tanh x, with |result| > 1 (it rounds to ±1 for |x|
// above about 19). It is odd.
// Errors: ErrNotFinite when x is NaN or ±Inf; ErrDomain when x is ±0 (the pole); ErrOverflow when
// 0 < |x| < about 5.6e-309, where the result is beyond float64.
//
// ID:      NUM-TRIG-014
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  NIST Digital Library of Mathematical Functions, Release 1.2.8 (2026-09-15), eq. 4.28.7, https://dlmf.nist.gov/4.28 (read 2026-11-12)
// Source:  Abramowitz and Stegun (eds.), Handbook of Mathematical Functions, NBS Applied Mathematics Series 55, 10th printing, 1972, eq. 4.5.6, p. 83
// Tags:    trig, hyperbolic, cotangent
// Units:   x dimensionless; result dimensionless
// Example: Coth(1) = 1.3130352854993312 (hand-derived from 4.28.7)
func Coth(x float64) (float64, error) {
	if !reduce.AllFinite(x) {
		return 0, fmt.Errorf("%w: x", ErrNotFinite)
	}
	if x == 0 {
		return 0, fmt.Errorf("%w: x is 0", ErrDomain)
	}
	r := 1 / math.Tanh(x)
	if math.IsInf(r, 0) {
		return 0, fmt.Errorf("%w: x", ErrOverflow)
	}
	return r, nil
}
