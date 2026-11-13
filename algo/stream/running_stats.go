package stream

import (
	"fmt"
	"math"
)

// RunningStats is the count, mean, variance (population and sample), minimum and maximum of the values
// added so far, kept in one pass without holding the values (Welford's update), with Merge to combine
// two summaries of separate parts of a stream (Chan, Golub and LeVeque's pairwise formula), for example
// per-shard summaries. The zero value is an empty summary, ready to use. Not safe for concurrent use.
// num/stats computes the variance of a whole slice in two passes; this type is its streaming,
// mergeable counterpart.
//
// Update, for the k-th value x (Welford 1962): δ = x − mean; mean ← mean + δ/k; M2 ← M2 + δ·(x − mean),
// where M2 is the sum of squared deviations from the mean. Merge of (nA, meanA, M2A) and (nB, meanB, M2B)
// (Chan, Golub, LeVeque): n = nA + nB; δ = meanB − meanA; mean = meanA + δ·nB/n;
// M2 = M2A + M2B + δ²·nA·nB/n. Population variance M2/n, sample variance M2/(n − 1).
//
// Methods: Add(x) adds a value. Merge(other) adds every value other summarises. Count, Mean,
// PopulationVariance, SampleVariance, Min, Max read the summary; each returns (0, false) when it is
// undefined (no values; fewer than two for SampleVariance).
//
// ID:      ALG-STREAM-004
// Time:    O(1) per Add, per Merge and per accessor
// Space:   O(1), no allocation
// Source:  Welford, "Note on a method for calculating corrected sums of squares and products", Technometrics, 1962, 4(3):419–420, doi:10.1080/00401706.1962.10490022
// Source:  Chan, Golub, LeVeque, "Algorithms for computing the sample variance: analysis and recommendations", The American Statistician, 1983, 37(3):242–247, doi:10.1080/00031305.1983.10483115 (the updating formula for combining two subsamples)
// Source:  Chan, Golub, LeVeque, "Updating formulae and a pairwise algorithm for computing sample variances", Stanford University, Computer Science Department, Technical Report STAN-CS-79-773, November 1979
// Tags:    stream, running-statistics, online-variance, welford, mean, variance, standard-deviation, min, max, merge, parallel
type RunningStats struct {
	n        int64
	mean, m2 float64
	min, max float64
}

// Add adds x to the summary. O(1).
// Errors: ErrNotFinite when x is NaN or ±Inf; ErrOverflow when the count would exceed math.MaxInt64 or a
// deviation or the sum of squared deviations would exceed the float64 range. On error nothing changes.
func (s *RunningStats) Add(x float64) error {
	if math.IsNaN(x) || math.IsInf(x, 0) {
		return fmt.Errorf("%w: x = %v", ErrNotFinite, x)
	}
	if s.n == math.MaxInt64 {
		return fmt.Errorf("%w: count", ErrOverflow)
	}
	n := s.n + 1
	delta := x - s.mean
	mean := s.mean + delta/float64(n)
	m2 := s.m2 + delta*(x-mean)
	if math.IsInf(delta, 0) || math.IsInf(m2, 0) || math.IsNaN(m2) {
		return fmt.Errorf("%w: x = %v", ErrOverflow, x)
	}
	if s.n == 0 {
		s.min, s.max = x, x
	} else {
		s.min, s.max = min(s.min, x), max(s.max, x)
	}
	s.n, s.mean, s.m2 = n, mean, m2
	return nil
}

// Merge adds to s every value other summarises, as if they had been added to s one by one (up to
// rounding). other is not changed. O(1).
// Errors: ErrOverflow when the count would exceed math.MaxInt64 or the difference of the means or the
// combined sum of squared deviations would exceed the float64 range. On error nothing changes.
func (s *RunningStats) Merge(other RunningStats) error {
	if other.n == 0 {
		return nil
	}
	if s.n == 0 {
		*s = other
		return nil
	}
	if other.n > math.MaxInt64-s.n {
		return fmt.Errorf("%w: count %d + %d", ErrOverflow, s.n, other.n)
	}
	n := s.n + other.n
	delta := other.mean - s.mean
	fn := float64(n)
	mean := s.mean + delta*(float64(other.n)/fn)
	m2 := s.m2 + other.m2 + delta*delta*(float64(s.n)/fn*float64(other.n))
	if math.IsInf(delta, 0) || math.IsInf(mean, 0) || math.IsInf(m2, 0) || math.IsNaN(m2) {
		return fmt.Errorf("%w: merged variance", ErrOverflow)
	}
	s.n, s.mean, s.m2 = n, mean, m2
	s.min, s.max = min(s.min, other.min), max(s.max, other.max)
	return nil
}

// Count returns the number of values added. O(1).
func (s *RunningStats) Count() int64 { return s.n }

// Mean returns the arithmetic mean, or (0, false) when no value has been added. O(1).
func (s *RunningStats) Mean() (float64, bool) {
	if s.n == 0 {
		return 0, false
	}
	return s.mean, true
}

// PopulationVariance returns M2/n, the variance with divisor n, or (0, false) when no value has been
// added. O(1).
func (s *RunningStats) PopulationVariance() (float64, bool) {
	if s.n == 0 {
		return 0, false
	}
	return s.m2 / float64(s.n), true
}

// SampleVariance returns M2/(n − 1), the variance with divisor n − 1, or (0, false) when fewer than two
// values have been added. O(1).
func (s *RunningStats) SampleVariance() (float64, bool) {
	if s.n < 2 {
		return 0, false
	}
	return s.m2 / float64(s.n-1), true
}

// Min returns the least value added, or (0, false) when none has been. -0 counts as less than +0 (the
// order of Go's min and max built-ins). O(1).
func (s *RunningStats) Min() (float64, bool) {
	if s.n == 0 {
		return 0, false
	}
	return s.min, true
}

// Max returns the greatest value added, or (0, false) when none has been. O(1).
func (s *RunningStats) Max() (float64, bool) {
	if s.n == 0 {
		return 0, false
	}
	return s.max, true
}
