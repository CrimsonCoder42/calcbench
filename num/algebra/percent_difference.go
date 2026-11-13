package algebra

import (
	"fmt"
	"math"

	"github.com/CrimsonCoder42/calcbench/algo/reduce"
)

// PercentDifference returns the symmetric percent difference of two positive values, 100·|x − y| / ((x + y)/2):
// the change measured against their arithmetic mean, so swapping x and y gives the same result.
// Errors: ErrNotFinite for NaN or ±Inf input; ErrDegenerate when x ≤ 0 or y ≤ 0 (the source defines the
// measure for positive quantities only).
//
// ID:      NUM-ALGEB-034
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Törnqvist, Vartia and Vartia, "How Should Relative Changes Be Measured?", The American Statistician, 1985, 39(1):43–46, doi:10.1080/00031305.1985.10479385 (relative change with the arithmetic mean of x and y as the reference value)
// Tags:    percent, percent-difference, relative-difference, symmetric, scalar, algebra
// Units:   x and y in one unit; the result in percent
// Example: 40 and 50 → 100·10/45 = 22.2222… % (hand-checked)
func PercentDifference(x, y float64) (float64, error) {
	if !reduce.AllFinite(x, y) {
		return 0, fmt.Errorf("%w: x = %v, y = %v", ErrNotFinite, x, y)
	}
	if x <= 0 || y <= 0 {
		return 0, fmt.Errorf("%w: x = %v, y = %v, both must be > 0", ErrDegenerate, x, y)
	}
	// x/2 + y/2 cannot overflow; the quotient is at most 200.
	return 100 * (math.Abs(x-y) / (x/2 + y/2)), nil
}
