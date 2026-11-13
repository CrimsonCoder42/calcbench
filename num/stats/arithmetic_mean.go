package stats

import (
	"errors"
	"fmt"
)

// Errors returned by this package. They are wrapped with the offending parameter; test with errors.Is.
var (
	// ErrEmpty: the input has no elements.
	ErrEmpty = errors.New("stats: empty input")
	// ErrNotFinite: a value is NaN or ±Inf.
	ErrNotFinite = errors.New("stats: value is NaN or infinite")
	// ErrOverflow: the result, or an intermediate the formula cannot avoid, is beyond the float64 range
	// (for a total of counts, beyond the int range).
	ErrOverflow = errors.New("stats: result exceeds the float64 range")
	// ErrTooFew: fewer elements than the formula needs (for example one value for a sample variance).
	ErrTooFew = errors.New("stats: too few values")
	// ErrLengthMismatch: two inputs that pair element by element have different lengths.
	ErrLengthMismatch = errors.New("stats: inputs have different lengths")
	// ErrNegative: a count, weight or standard deviation is negative.
	ErrNegative = errors.New("stats: negative value")
	// ErrZeroTotal: counts or weights that must have a positive total sum to zero.
	ErrZeroTotal = errors.New("stats: total is zero")
	// ErrZeroMean: the mean is zero, so a ratio to it is undefined.
	ErrZeroMean = errors.New("stats: mean is zero")
	// ErrZeroVariance: a standard deviation or spread that divides is zero.
	ErrZeroVariance = errors.New("stats: zero variance")
	// ErrOutOfRange: a parameter is outside the values the function accepts.
	ErrOutOfRange = errors.New("stats: out of range")
	// ErrNotAscending: class edges are not strictly ascending.
	ErrNotAscending = errors.New("stats: edges not strictly ascending")
)

// ArithmeticMean returns the arithmetic mean of values: their sum divided by their count.
// The sum is compensated (Neumaier), and a sum beyond the float64 range is redone on values scaled by a
// power of two, so the mean of finite values is always returned. values is read, never modified.
// Errors: ErrEmpty when values is empty; ErrNotFinite when a value is NaN or ±Inf.
//
// ID:      NUM-STATS-001
// Time:    O(n)   n = len(values); one pass (two only when the plain sum overflows)
// Space:   O(1) extra, no allocation
// Source:  NIST/SEMATECH, e-Handbook of Statistical Methods, §1.3.5.1 "Measures of Location", https://www.itl.nist.gov/div898/handbook/eda/section3/eda351.htm (read 2026-11-12)
// Source:  Neumaier, "Rundungsfehleranalyse einiger Verfahren zur Summation endlicher Summen", ZAMM, 1974, 54(1):39–51, doi:10.1002/zamm.19740540106
// Tags:    mean, average, location, central-tendency, descriptive, compensated-summation
// Units:   values in any one unit; result in that unit
// Example: NIST StRD Univariate NumAcc1 (10000001, 10000003, 10000002) → 10000002, the certified mean
func ArithmeticMean(values []float64) (float64, error) {
	if len(values) == 0 {
		return 0, fmt.Errorf("%w: values", ErrEmpty)
	}
	if err := checkFinite("values", values); err != nil {
		return 0, err
	}
	return meanOf(values), nil
}
