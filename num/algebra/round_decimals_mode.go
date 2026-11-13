package algebra

import (
	"fmt"
	"math"

	"github.com/CrimsonCoder42/calcbench/algo/reduce"
)

// RoundMode names an IEEE 754-2019 rounding-direction attribute, applied by RoundDecimalsMode at a
// decimal position.
type RoundMode int

const (
	// RoundTiesToEven: nearest; a tie goes to the even neighbour (§4.3.1).
	RoundTiesToEven RoundMode = iota
	// RoundTiesToAway: nearest; a tie goes away from zero (§4.3.1).
	RoundTiesToAway
	// RoundTiesTowardZero: nearest; a tie goes toward zero (§9.5, augmented operations).
	RoundTiesTowardZero
	// RoundTowardPositive: the least value ≥ x, the ceiling (§4.3.2).
	RoundTowardPositive
	// RoundTowardNegative: the greatest value ≤ x, the floor (§4.3.2).
	RoundTowardNegative
	// RoundTowardZero: the nearest value no larger in magnitude, truncation (§4.3.2).
	RoundTowardZero
)

// RoundDecimalsMode returns x rounded to the given number of decimal places under mode, as
// round(x·10ᵈ)/10ᵈ (d < 0 rounds to tens, hundreds, …). As with RoundDecimals, the rounding acts on
// the stored binary value of x·10ᵈ, so a decimal tie that is not exactly representable (2.675 is
// 2.67499… in binary) is not a tie. When x·10ᵈ is already an integer beyond 2⁵³ or overflows, x is
// returned unchanged.
// Errors: ErrNotFinite for NaN or ±Inf x; ErrDegenerate when |decimals| > 308 or mode is not one of the
// six RoundMode constants.
//
// ID:      NUM-ALGEB-037
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  IEEE Std 754-2019, IEEE Standard for Floating-Point Arithmetic, §4.3.1 (roundTiesToEven, roundTiesToAway), §4.3.2 (roundTowardPositive, roundTowardNegative, roundTowardZero), §9.5 (roundTiesToZero for augmented operations)
// Tags:    round, rounding-mode, decimals, half-even, half-away, ceiling, floor, truncate, scalar, algebra
// Units:   x and the result in one unit
// Example: −2.5 at 0 places → −2 (ties to even), −3 (ties away), −2 (ties toward zero), −2 (ceiling), −3 (floor), −2 (toward zero) (hand-checked)
func RoundDecimalsMode(x float64, decimals int, mode RoundMode) (float64, error) {
	if !reduce.AllFinite(x) {
		return 0, fmt.Errorf("%w: x = %v", ErrNotFinite, x)
	}
	if decimals < -308 || decimals > 308 {
		return 0, fmt.Errorf("%w: decimals = %d", ErrDegenerate, decimals)
	}
	if mode < RoundTiesToEven || mode > RoundTowardZero {
		return 0, fmt.Errorf("%w: mode = %d", ErrDegenerate, mode)
	}
	if decimals >= 0 {
		p := math.Pow10(decimals)
		y := x * p
		if !reduce.AllFinite(y) || math.Abs(y) >= 0x1p53 {
			return x, nil
		}
		return roundIntegral(y, mode) / p, nil
	}
	p := math.Pow10(-decimals)
	return roundIntegral(x/p, mode) * p, nil
}

// roundIntegral rounds y to an integer under mode.
func roundIntegral(y float64, mode RoundMode) float64 {
	switch mode {
	case RoundTiesToAway:
		return math.Round(y)
	case RoundTiesTowardZero:
		t := math.Trunc(y)
		if math.Abs(y-t) == 0.5 {
			return t
		}
		return math.Round(y)
	case RoundTowardPositive:
		return math.Ceil(y)
	case RoundTowardNegative:
		return math.Floor(y)
	case RoundTowardZero:
		return math.Trunc(y)
	default:
		return math.RoundToEven(y)
	}
}
