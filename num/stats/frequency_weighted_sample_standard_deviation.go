package stats

import (
	"fmt"
	"math"
)

// FrequencyWeightedSampleStandardDeviation returns the sample standard deviation of data given as distinct
// values with their frequencies (a frequency table): s = √(Σ fᵢ(xᵢ − x̄)²/(n − 1)), x̄ = Σ fᵢxᵢ/n,
// n = Σ fᵢ. It equals NUM-STATS-019 applied to the data with each value repeated fᵢ times, without
// expanding it. values and counts are read, never modified.
// Errors: as NUM-STATS-035, and ErrTooFew when the total count is 1.
//
// ID:      NUM-STATS-036
// Time:    O(m)   m = len(values); three passes (scale, mean, centred sums)
// Space:   O(1) extra, no allocation
// Source:  Illowsky B, Dean S, et al., Introductory Statistics 2e, OpenStax, 2023, §2.7 "Measures of the Spread of the Data", https://openstax.org/books/introductory-statistics-2e/pages/2-7-measures-of-the-spread-of-the-data (read 2026-11-12) (s = √(Σf(x − x̄)²/(n − 1)); Example 2.32: s² = 9.7375/(20 − 1) = 0.5125, s = 0.715891)
// Source:  Chan, Golub, LeVeque, "Algorithms for computing the sample variance: analysis and recommendations", The American Statistician, 1983, 37(3):242–247, doi:10.1080/00031305.1983.10483115 (the corrected two-pass algorithm)
// Tags:    standard-deviation, sample-standard-deviation, frequency-table, weighted, counts, dispersion, descriptive
// Units:   values in any one unit; counts dimensionless; result in the unit of values
// Example: OpenStax Example 2.32 ages 9, 9.5, 10, 10.5, 11, 11.5 with frequencies 1, 2, 4, 4, 6, 3 → s = 0.715891 (printed)
func FrequencyWeightedSampleStandardDeviation(values []float64, counts []int) (float64, error) {
	total, err := checkFrequencies(values, counts)
	if err != nil {
		return 0, err
	}
	if total < 2 {
		return 0, fmt.Errorf("%w: total count = %d, want ≥ 2", ErrTooFew, total)
	}
	ss, e := frequencySquares(values, counts, total)
	sd := math.Ldexp(math.Sqrt(ss/float64(total-1)), e)
	if math.IsInf(sd, 0) {
		return 0, fmt.Errorf("%w: standard deviation of values", ErrOverflow)
	}
	return sd, nil
}
