package prob

import (
	"errors"
	"fmt"

	"github.com/CrimsonCoder42/calcbench/algo/reduce"
)

// Errors returned by this package. They are wrapped with the offending parameter; test with errors.Is.
var (
	// ErrNotFinite: an input is NaN or ±Inf.
	ErrNotFinite = errors.New("prob: value is NaN or infinite")
	// ErrProbability: a probability is outside [0, 1].
	ErrProbability = errors.New("prob: probability outside [0, 1]")
	// ErrOutOfRange: a count or parameter is outside its domain (k > n, σ ≤ 0, n < 0, unequal lengths, …).
	ErrOutOfRange = errors.New("prob: argument out of range")
	// ErrInconsistent: the inputs cannot come from one probability space, e.g. P(A∩B) > P(A).
	ErrInconsistent = errors.New("prob: inconsistent probabilities")
	// ErrNotDistribution: weights do not sum to 1 within the package tolerance (see the package comment).
	ErrNotDistribution = errors.New("prob: probabilities do not sum to 1")
	// ErrOverflow: finite inputs whose true result is beyond the float64 range.
	ErrOverflow = errors.New("prob: result exceeds the float64 range")
)

// checkFinite returns ErrNotFinite, naming the parameter, when x is NaN or ±Inf.
func checkFinite(name string, x float64) error {
	if !reduce.AllFinite(x) {
		return fmt.Errorf("%w: %s = %v", ErrNotFinite, name, x)
	}
	return nil
}

// checkProbability returns ErrNotFinite for NaN or ±Inf and ErrProbability outside [0, 1].
func checkProbability(name string, p float64) error {
	if !reduce.AllFinite(p) {
		return fmt.Errorf("%w: %s = %v", ErrNotFinite, name, p)
	}
	if p < 0 || p > 1 {
		return fmt.Errorf("%w: %s = %v", ErrProbability, name, p)
	}
	return nil
}

// checkPositive returns ErrNotFinite for NaN or ±Inf and ErrOutOfRange when x ≤ 0.
func checkPositive(name string, x float64) error {
	if !reduce.AllFinite(x) {
		return fmt.Errorf("%w: %s = %v", ErrNotFinite, name, x)
	}
	if x <= 0 {
		return fmt.Errorf("%w: %s = %v, want > 0", ErrOutOfRange, name, x)
	}
	return nil
}

// ProbabilityFromCounts returns the classical probability favorable/total of an event that holds in
// favorable of total equally likely outcomes. That the outcomes are equally likely is the caller's
// assumption; it cannot be checked from two counts.
// Errors: ErrOutOfRange when total < 1, favorable < 0 or favorable > total.
//
// ID:      NUM-PROB-001
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Ross, A First Course in Probability, 10th ed., Pearson, 2019, §2.5 "Sample Spaces Having Equally Likely Outcomes"
// Tags:    probability, classical, counting, equally-likely
// Units:   favorable, total: counts of outcomes; result: a probability in [0, 1]
// Example: a fair die, the event "even": 3 of 6 → 0.5
func ProbabilityFromCounts(favorable, total int) (float64, error) {
	if total < 1 {
		return 0, fmt.Errorf("%w: total = %d, want >= 1", ErrOutOfRange, total)
	}
	if favorable < 0 || favorable > total {
		return 0, fmt.Errorf("%w: favorable = %d, want 0..%d", ErrOutOfRange, favorable, total)
	}
	// Each conversion rounds at most once (exactly below 2⁵³) and the quotient once more.
	return float64(favorable) / float64(total), nil
}
