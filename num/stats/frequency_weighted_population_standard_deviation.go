package stats

import (
	"fmt"
	"math"
)

// checkFrequencies validates values with integer frequencies and returns N = Σ counts as an int64.
// Errors: ErrEmpty, ErrLengthMismatch, ErrNotFinite, ErrNegative, ErrZeroTotal, ErrOverflow (N beyond int64).
func checkFrequencies(values []float64, counts []int) (int64, error) {
	if len(values) == 0 {
		return 0, fmt.Errorf("%w: values", ErrEmpty)
	}
	if len(counts) != len(values) {
		return 0, fmt.Errorf("%w: len(counts) = %d, len(values) = %d", ErrLengthMismatch, len(counts), len(values))
	}
	if err := checkFinite("values", values); err != nil {
		return 0, err
	}
	var total int64
	for i, c := range counts {
		if c < 0 {
			return 0, fmt.Errorf("%w: counts[%d] = %d", ErrNegative, i, c)
		}
		if int64(c) > math.MaxInt64-total {
			return 0, fmt.Errorf("%w: total of counts", ErrOverflow)
		}
		total += int64(c)
	}
	if total == 0 {
		return 0, fmt.Errorf("%w: counts", ErrZeroTotal)
	}
	return total, nil
}

// frequencySquares returns ss and e with Σ fᵢ(xᵢ − x̄)² = ss·4^e, x̄ = Σ fᵢxᵢ/N, for checked input.
// Every value is first multiplied by 2^(−e), e chosen so that the largest |value| lands in [½, 1) (a power
// of two: exact, except that values below 2^(−1074+e) lose low bits, negligible beside the largest).
// Then fᵢ·xᵢ ≤ fᵢ and fᵢ·dᵢ² ≤ 4fᵢ, so no sum can overflow or underflow harmfully. The mean and the
// corrected two-pass sum Σfd² − (Σfd)²/N are compensated.
func frequencySquares(values []float64, counts []int, total int64) (ss float64, e int) {
	maxAbs := 0.0
	for _, v := range values {
		maxAbs = max(maxAbs, math.Abs(v))
	}
	if maxAbs == 0 {
		return 0, 0
	}
	_, e = math.Frexp(maxAbs) // maxAbs = frac·2^e, frac in [½, 1)
	n := float64(total)
	var sx compensatedSum
	for i, v := range values {
		if counts[i] != 0 {
			sx.add(float64(counts[i]) * math.Ldexp(v, -e))
		}
	}
	mean := sx.total() / n
	var sd, sd2 compensatedSum
	for i, v := range values {
		if counts[i] == 0 {
			continue
		}
		f := float64(counts[i])
		d := math.Ldexp(v, -e) - mean
		sd.add(f * d)
		sd2.add(f * d * d)
	}
	c := sd.total()
	ss = sd2.total() - c*c/n
	return max(ss, 0), e
}

// FrequencyWeightedPopulationStandardDeviation returns the population standard deviation of data given as
// distinct values with their frequencies (a frequency table): σ = √(Σ fᵢ(xᵢ − μ)²/N), μ = Σ fᵢxᵢ/N,
// N = Σ fᵢ. It equals NUM-STATS-022 applied to the data with each value repeated fᵢ times, without
// expanding it. Zero frequencies are allowed (the value is ignored). values and counts are read, never
// modified.
// Errors: ErrEmpty when values is empty; ErrLengthMismatch when len(counts) ≠ len(values); ErrNotFinite
// when a value is NaN or ±Inf; ErrNegative when a count is negative; ErrZeroTotal when every count is 0;
// ErrOverflow when the total count exceeds the int64 range or σ the float64 range (values near ±MaxFloat64).
//
// ID:      NUM-STATS-035
// Time:    O(m)   m = len(values); three passes (scale, mean, centred sums)
// Space:   O(1) extra, no allocation
// Source:  Illowsky B, Dean S, et al., Introductory Statistics 2e, OpenStax, 2023, §2.7 "Measures of the Spread of the Data", https://openstax.org/books/introductory-statistics-2e/pages/2-7-measures-of-the-spread-of-the-data (read 2026-11-12) (frequency-table formulas σ = √(Σf(x − μ)²/N), s = √(Σf(x − x̄)²/(n − 1)); Example 2.32)
// Source:  Chan, Golub, LeVeque, "Algorithms for computing the sample variance: analysis and recommendations", The American Statistician, 1983, 37(3):242–247, doi:10.1080/00031305.1983.10483115 (the corrected two-pass algorithm)
// Tags:    standard-deviation, population-standard-deviation, frequency-table, weighted, counts, dispersion, descriptive
// Units:   values in any one unit; counts dimensionless; result in the unit of values
// Example: OpenStax Example 2.32 ages 9, 9.5, 10, 10.5, 11, 11.5 with frequencies 1, 2, 4, 4, 6, 3 → Σf(x − x̄)² = 9.7375, σ = √(9.7375/20) ≈ 0.6977643
func FrequencyWeightedPopulationStandardDeviation(values []float64, counts []int) (float64, error) {
	total, err := checkFrequencies(values, counts)
	if err != nil {
		return 0, err
	}
	ss, e := frequencySquares(values, counts, total)
	sd := math.Ldexp(math.Sqrt(ss/float64(total)), e)
	if math.IsInf(sd, 0) {
		return 0, fmt.Errorf("%w: standard deviation of values", ErrOverflow)
	}
	return sd, nil
}
