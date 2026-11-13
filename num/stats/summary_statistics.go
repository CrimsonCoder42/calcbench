package stats

import (
	"fmt"
	"math"

	"github.com/CrimsonCoder42/calcbench/algo/reduce"
)

// Summary holds the usual descriptive statistics of a sample, each as defined by its own function.
type Summary struct {
	Count             int     // n, the number of values
	Mean              float64 // ArithmeticMean, NUM-STATS-001
	StandardDeviation float64 // SampleStandardDeviation (divisor n − 1), NUM-STATS-019
	Min               float64 // the smallest value
	Max               float64 // the largest value
	Median            float64 // Median, NUM-STATS-012
	Range             float64 // Range (Max − Min), NUM-STATS-018
}

// SummaryStatistics returns the count, mean, sample standard deviation, minimum, maximum, median and
// range of values in one call. The set is this package's choice, not a named list from the source; each
// field is computed by the same internal code as its own function, so the numbers always agree with it.
// values is read, never modified: one copy is reordered for the median.
// Errors: ErrEmpty when values is empty; ErrTooFew when it has one value; ErrNotFinite when a value is
// NaN or ±Inf; ErrOverflow when the standard deviation or the range exceeds the float64 range.
//
// ID:      NUM-STATS-020
// Time:    O(n) worst case   n = len(values); four passes and one linear-time selection
// Space:   O(n) extra (one allocation, the median's copy) plus O(log n) stack
// Source:  NIST/SEMATECH, e-Handbook of Statistical Methods, §1.3.5.1 "Measures of Location", https://www.itl.nist.gov/div898/handbook/eda/section3/eda351.htm (read 2026-11-12)
// Source:  NIST/SEMATECH, e-Handbook of Statistical Methods, §1.3.5.6 "Measures of Scale", https://www.itl.nist.gov/div898/handbook/eda/section3/eda356.htm (read 2026-11-12)
// Tags:    summary-statistics, descriptive, mean, standard-deviation, median, range, min, max
// Units:   values in any one unit; Count dimensionless; every other field in the unit of values
// Example: NIST StRD Univariate NumAcc1 → Count 3, Mean 10000002, StandardDeviation 1 (certified); Min 10000001, Max 10000003, Median 10000002, Range 2 (hand-checked)
func SummaryStatistics(values []float64) (Summary, error) {
	if err := checkSample(values, 2); err != nil {
		return Summary{}, err
	}
	s := Summary{Count: len(values), Mean: meanOf(values)}
	ss, e := centredSquares(values, s.Mean)
	s.StandardDeviation = math.Ldexp(math.Sqrt(ss/float64(len(values)-1)), e)
	if !reduce.AllFinite(s.StandardDeviation) {
		return Summary{}, fmt.Errorf("%w: standard deviation of values", ErrOverflow)
	}
	s.Min, s.Max = minMax(values)
	s.Range = s.Max - s.Min
	if !reduce.AllFinite(s.Range) {
		return Summary{}, fmt.Errorf("%w: range of values", ErrOverflow)
	}
	s.Median = medianOf(values)
	return s, nil
}
