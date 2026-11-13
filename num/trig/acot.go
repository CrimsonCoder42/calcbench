package trig

import (
	"fmt"
	"math"

	"github.com/CrimsonCoder42/calcbench/algo/reduce"
)

// Acot returns the principal inverse cotangent of x, arctan(1/x), in radians in (−π/2, π/2]: the DLMF
// branch, odd in x, which is NOT the (0, π) branch some texts use (π/2 − arctan x).
// At x = 0, where DLMF Table 4.23.1 gives ∓π/2, both +0 and −0 return +π/2, so the range is (−π/2, π/2].
// For tiny negative x the result rounds to −math.Pi/2, which is above the true −π/2.
// Errors: ErrNotFinite when x is NaN or ±Inf.
//
// ID:      NUM-TRIG-011
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  NIST Digital Library of Mathematical Functions, Release 1.2.8 (2026-09-15), eq. 4.23.9, 4.23.15 and Table 4.23.1, https://dlmf.nist.gov/4.23 (read 2026-11-12)
// Source:  Abramowitz and Stegun (eds.), Handbook of Mathematical Functions, NBS Applied Mathematics Series 55, 10th printing, 1972, eq. 4.4.8, p. 79
// Tags:    trig, inverse, arccotangent
// Units:   x dimensionless; result in radians, in (−π/2, π/2]
// Example: Acot(1) = π/4; Acot(−1) = −π/4 (DLMF Table 4.23.1)
func Acot(x float64) (float64, error) {
	if !reduce.AllFinite(x) {
		return 0, fmt.Errorf("%w: x", ErrNotFinite)
	}
	if x == 0 {
		return math.Pi / 2, nil
	}
	return math.Atan(1 / x), nil
}
