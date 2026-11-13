package stats

import (
	"fmt"

	"github.com/CrimsonCoder42/calcbench/algo/reduce"
)

// ZScore returns the standard score of x, z = (x − mean)/sd: how many standard deviations x lies above
// (z > 0) or below (z < 0) the reference mean. The caller supplies the reference mean and standard
// deviation. It states the number only; labels for |z| ranges are interpretation and are not given here.
// Errors: ErrNotFinite when x, mean or sd is NaN or ±Inf; ErrNegative when sd < 0; ErrZeroVariance when
// sd = 0; ErrOverflow when z exceeds the float64 range.
//
// ID:      NUM-STATS-028
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  NIST, Dataplot Reference Manual, "ISO 13528 ZSCORE", https://www.itl.nist.gov/div898/software/dataplot/refman2/auxillar/zscore.htm (read 2026-11-12)
// Source:  ISO 13528:2015, Statistical methods for use in proficiency testing by interlaboratory comparison, 2nd ed., ISO, §9.4 "z scores"
// Tags:    z-score, standard-score, standardize, normalization, descriptive
// Units:   x, mean and sd in one unit; result dimensionless
// Example: x = 12, mean = 10, sd = 2 → 1 (hand-checked)
func ZScore(x, mean, sd float64) (float64, error) {
	if err := checkScalar("x", x); err != nil {
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
	if sd == 0 {
		return 0, fmt.Errorf("%w: sd = 0", ErrZeroVariance)
	}
	return standardizedDifference(x, mean, sd, "z")
}

// standardizedDifference returns (a − b)/sd for finite a, b and sd > 0, dividing each term first when
// a − b overflows, or ErrOverflow naming result when the quotient itself exceeds the float64 range.
func standardizedDifference(a, b, sd float64, result string) (float64, error) {
	r := (a - b) / sd
	if !reduce.AllFinite(a - b) {
		r = a/sd - b/sd
	}
	if !reduce.AllFinite(r) {
		return 0, fmt.Errorf("%w: %s", ErrOverflow, result)
	}
	return r, nil
}
