package stats

import (
	"fmt"
	"math"

	"github.com/CrimsonCoder42/calcbench/algo/reduce"
)

// RawScoreFromZScore returns the value x whose standard score is z, x = mean + z·sd: ZScore solved for x,
// computed with one rounding (fused multiply-add).
// Errors: ErrNotFinite when z, mean or sd is NaN or ±Inf; ErrNegative when sd < 0; ErrOverflow when x
// exceeds the float64 range.
//
// ID:      NUM-STATS-029
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  NIST, Dataplot Reference Manual, "ISO 13528 ZSCORE", https://www.itl.nist.gov/div898/software/dataplot/refman2/auxillar/zscore.htm (read 2026-11-12) (Z(i) = (X(i) − Xref)/sigma, solved for X(i))
// Source:  ISO 13528:2015, Statistical methods for use in proficiency testing by interlaboratory comparison, 2nd ed., ISO, §9.4 "z scores"
// Tags:    z-score, raw-score, standard-score, inverse, descriptive
// Units:   z dimensionless; mean and sd in one unit; result in that unit
// Example: z = 1, mean = 10, sd = 2 → 12, the inverse of the ZScore example (hand-checked)
func RawScoreFromZScore(z, mean, sd float64) (float64, error) {
	if err := checkScalar("z", z); err != nil {
		return 0, err
	}
	if err := checkScalar("mean", mean); err != nil {
		return 0, err
	}
	if err := checkScalar("sd", sd); err != nil {
		return 0, err
	}
	if sd < 0 {
		return 0, fmt.Errorf("%w: sd = %v", ErrNegative, sd)
	}
	x := math.FMA(z, sd, mean)
	if !reduce.AllFinite(x) {
		return 0, fmt.Errorf("%w: x", ErrOverflow)
	}
	return x, nil
}
