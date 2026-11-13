package algebra

import (
	"fmt"
	"math"

	"github.com/CrimsonCoder42/calcbench/algo/reduce"
)

// RoundDecimals returns x rounded to the given number of decimal places, ties to even, as
// RoundToEven(x·10ᵈ)/10ᵈ (d < 0 rounds to tens, hundreds, …). As with any binary floating-point rounding
// to decimals, a value whose decimal expansion is not exact (2.675 is stored as 2.67499…) rounds by its
// stored value. When x·10ᵈ is already an integer beyond 2⁵³ or overflows, x is returned unchanged.
// Errors: ErrNotFinite for NaN or ±Inf x; ErrDegenerate when |decimals| > 308.
//
// ID:      NUM-ALGEB-022
// Time:    O(1)
// Space:   O(1)
// Source:  IEEE Std 754-2019, IEEE Standard for Floating-Point Arithmetic, §5.3.1 (roundToIntegralTiesToEven) and §4.3.1 (roundTiesToEven)
// Source:  Muller et al., Handbook of Floating-Point Arithmetic, 2nd ed., Birkhäuser, 2018, §2.2 "Rounding"
// Tags:    round, decimals, half-even, bankers-rounding, scalar, algebra
// Units:   x and the result in one unit
// Example: 2.5 at 0 places → 2; 3.5 → 4; 1234.5678 at 2 → 1234.57; at −2 → 1200 (hand-checked)
func RoundDecimals(x float64, decimals int) (float64, error) {
	if !reduce.AllFinite(x) {
		return 0, fmt.Errorf("%w: x = %v", ErrNotFinite, x)
	}
	if decimals < -308 || decimals > 308 {
		return 0, fmt.Errorf("%w: decimals = %d", ErrDegenerate, decimals)
	}
	if decimals >= 0 {
		p := math.Pow10(decimals)
		y := x * p
		if !reduce.AllFinite(y) || math.Abs(y) >= 0x1p53 {
			return x, nil
		}
		return math.RoundToEven(y) / p, nil
	}
	p := math.Pow10(-decimals)
	return math.RoundToEven(x/p) * p, nil
}
