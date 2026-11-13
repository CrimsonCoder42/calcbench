package stats

import "fmt"

// RelativeFrequency returns each class count divided by the total count: the fractions of the
// observations in each class, which sum to 1 (multiply by 100 for percent). counts is typically the
// result of FrequencyDistribution. counts is read, never modified.
// Errors: ErrEmpty when counts is empty; ErrNegative when a count is negative; ErrZeroTotal when every
// count is 0; ErrOverflow when the total exceeds the int range.
//
// ID:      NUM-STATS-002
// Time:    O(b)   b = len(counts); two passes
// Space:   O(b) for the result (one allocation); O(1) extra
// Source:  NIST/SEMATECH, e-Handbook of Statistical Methods, §1.3.3.14 "Histogram", https://www.itl.nist.gov/div898/handbook/eda/section3/histogra.htm (read 2026-11-12)
// Tags:    relative-frequency, frequency, proportion, histogram, descriptive
// Units:   counts dimensionless; result dimensionless fractions
// Example: counts 1, 3, 1 → 0.2, 0.6, 0.2 (hand-checked)
func RelativeFrequency(counts []int) ([]float64, error) {
	total, err := countTotal(counts)
	if err != nil {
		return nil, err
	}
	out := make([]float64, len(counts))
	t := float64(total)
	for i, c := range counts {
		out[i] = float64(c) / t
	}
	return out, nil
}

// countTotal returns the sum of counts, with ErrEmpty, ErrNegative, ErrOverflow or ErrZeroTotal.
func countTotal(counts []int) (int, error) {
	if len(counts) == 0 {
		return 0, fmt.Errorf("%w: counts", ErrEmpty)
	}
	total := 0
	for i, c := range counts {
		if c < 0 {
			return 0, fmt.Errorf("%w: counts[%d] = %d", ErrNegative, i, c)
		}
		if total > maxInt-c {
			return 0, fmt.Errorf("%w: total of counts", ErrOverflow)
		}
		total += c
	}
	if total == 0 {
		return 0, fmt.Errorf("%w: counts", ErrZeroTotal)
	}
	return total, nil
}

// maxInt is the largest int.
const maxInt = int(^uint(0) >> 1)
