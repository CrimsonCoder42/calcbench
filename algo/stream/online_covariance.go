package stream

import (
	"fmt"
	"math"

	"github.com/CrimsonCoder42/calcbench/algo/reduce"
)

// OnlineCovariance is the count, the two means, the covariance, the correlation and the beta of the
// pairs (x, y) added so far, kept in one pass without holding them (Welford's update extended to the
// co-moment), with Merge to combine two summaries of separate parts of a stream (Chan, Golub and
// LeVeque's pairwise formula). The zero value is an empty summary, ready to use. Not safe for
// concurrent use. It is the two-variable counterpart of RunningStats (ALG-STREAM-004); num/stats
// computes the covariance of whole slices in passes (NUM-STATS-038/039).
//
// Update, for the k-th pair: δx = x − x̄; x̄ ← x̄ + δx/k; ȳ ← ȳ + (y − ȳ)/k; C ← C + δx·(y − ȳ);
// M2x ← M2x + δx·(x − x̄), M2y likewise, with the new means. Merge of A and B: n = nA + nB;
// δx = x̄B − x̄A, δy = ȳB − ȳA; x̄ = x̄A + δx·nB/n; C = CA + CB + δx·δy·nA·nB/n; M2x, M2y as in
// RunningStats. Covariance C/(n − 1); correlation C/√(M2x·M2y); beta C/M2y, the least-squares slope of
// x on y (with x an asset's returns and y a benchmark's, the market-model beta).
//
// Methods: Add(x, y) adds a pair. Merge(other) adds every pair other summarises. Count, MeanX, MeanY,
// Covariance (sample), PopulationCovariance, Correlation and Beta read the summary; each returns
// (0, false) when it is undefined (no pairs; fewer than two for Covariance; a zero spread in x or y for
// Correlation, in y for Beta). Reset empties the summary.
//
// ID:      ALG-STREAM-008
// Time:    O(1) per Add, per Merge and per reader
// Space:   O(1), no allocation
// Source:  Welford, "Note on a method for calculating corrected sums of squares and products", Technometrics, 1962, 4(3):419–420, doi:10.1080/00401706.1962.10490022 (the corrected sum of products)
// Source:  Chan, Golub, LeVeque, "Algorithms for computing the sample variance: analysis and recommendations", The American Statistician, 1983, 37(3):242–247, doi:10.1080/00031305.1983.10483115 (the updating formula for combining two subsamples, applied to the co-moment)
// Tags:    stream, online-covariance, covariance, correlation, beta, welford, co-moment, merge, parallel, regression-slope
// Units:   x in any unit U, y in any unit V; MeanX in U, MeanY in V, Covariance and PopulationCovariance in U·V, Correlation unitless, Beta in U/V; Count a number of pairs
// Example: Add (1, 2), (2, 4), (3, 6) → Covariance = 2, PopulationCovariance = 4/3, Correlation = 1, Beta = 0.5
type OnlineCovariance struct {
	n            int64
	meanX, meanY float64
	m2X, m2Y, c  float64
}

// Add adds the pair (x, y). O(1).
// Errors: ErrNotFinite when x or y is NaN or ±Inf; ErrOverflow when the count would exceed
// math.MaxInt64 or a deviation or a moment would exceed the float64 range. On error nothing changes.
func (s *OnlineCovariance) Add(x, y float64) error {
	if math.IsNaN(x) || math.IsInf(x, 0) || math.IsNaN(y) || math.IsInf(y, 0) {
		return fmt.Errorf("%w: (x, y) = (%v, %v)", ErrNotFinite, x, y)
	}
	if s.n == math.MaxInt64 {
		return fmt.Errorf("%w: count", ErrOverflow)
	}
	k := float64(s.n + 1)
	dx, dy := x-s.meanX, y-s.meanY
	mx, my := s.meanX+dx/k, s.meanY+dy/k
	m2X, m2Y, c := s.m2X+dx*(x-mx), s.m2Y+dy*(y-my), s.c+dx*(y-my)
	if !reduce.AllFinite(dx, dy, mx, my, m2X, m2Y, c) {
		return fmt.Errorf("%w: (x, y) = (%v, %v)", ErrOverflow, x, y)
	}
	s.n++
	s.meanX, s.meanY, s.m2X, s.m2Y, s.c = mx, my, m2X, m2Y, c
	return nil
}

// Merge adds to s every pair other summarises, as if they had been added to s one by one (up to
// rounding). other is not changed. O(1).
// Errors: ErrOverflow when the count would exceed math.MaxInt64 or a difference of means or a combined
// moment would exceed the float64 range. On error nothing changes.
func (s *OnlineCovariance) Merge(other OnlineCovariance) error {
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
	fn := float64(n)
	wB := float64(other.n) / fn                 // nB/n
	wAB := float64(s.n) / fn * float64(other.n) // nA·nB/n
	dx, dy := other.meanX-s.meanX, other.meanY-s.meanY
	mx, my := s.meanX+dx*wB, s.meanY+dy*wB
	m2X := s.m2X + other.m2X + dx*dx*wAB
	m2Y := s.m2Y + other.m2Y + dy*dy*wAB
	c := s.c + other.c + dx*dy*wAB
	if !reduce.AllFinite(dx, dy, mx, my, m2X, m2Y, c) {
		return fmt.Errorf("%w: merged moments", ErrOverflow)
	}
	s.n, s.meanX, s.meanY, s.m2X, s.m2Y, s.c = n, mx, my, m2X, m2Y, c
	return nil
}

// Reset empties the summary. O(1).
func (s *OnlineCovariance) Reset() { *s = OnlineCovariance{} }

// Count returns the number of pairs added. O(1).
func (s *OnlineCovariance) Count() int64 { return s.n }

// MeanX returns the mean of the x values, or (0, false) when no pair has been added. O(1).
func (s *OnlineCovariance) MeanX() (float64, bool) {
	if s.n == 0 {
		return 0, false
	}
	return s.meanX, true
}

// MeanY returns the mean of the y values, or (0, false) when no pair has been added. O(1).
func (s *OnlineCovariance) MeanY() (float64, bool) {
	if s.n == 0 {
		return 0, false
	}
	return s.meanY, true
}

// Covariance returns the sample covariance C/(n − 1), or (0, false) when fewer than two pairs have been
// added. O(1).
func (s *OnlineCovariance) Covariance() (float64, bool) {
	if s.n < 2 {
		return 0, false
	}
	return s.c / float64(s.n-1), true
}

// PopulationCovariance returns C/n, or (0, false) when no pair has been added. O(1).
func (s *OnlineCovariance) PopulationCovariance() (float64, bool) {
	if s.n == 0 {
		return 0, false
	}
	return s.c / float64(s.n), true
}

// Correlation returns Pearson's r = C/√(M2x·M2y), clamped to [−1, 1] against rounding, or (0, false)
// when fewer than two pairs have been added or x or y has no spread. O(1).
func (s *OnlineCovariance) Correlation() (float64, bool) {
	if s.n < 2 {
		return 0, false
	}
	return correlation(s.c, s.m2X, s.m2Y)
}

// Beta returns C/M2y, the slope of the least-squares line of x on y, or (0, false) when fewer than two
// pairs have been added, y has no spread, or the slope exceeds the float64 range. O(1).
func (s *OnlineCovariance) Beta() (float64, bool) {
	if s.n < 2 {
		return 0, false
	}
	return beta(s.c, s.m2Y)
}

// correlation returns c/√(m2x·m2y) clamped to [−1, 1], or false when either moment is not positive.
// The square roots are taken apart so their product cannot overflow.
func correlation(c, m2x, m2y float64) (float64, bool) {
	if !(m2x > 0 && m2y > 0) {
		return 0, false
	}
	r := c / math.Sqrt(m2x) / math.Sqrt(m2y)
	if math.IsNaN(r) {
		return 0, false
	}
	return max(-1, min(1, r)), true
}

// beta returns c/m2y, or false when m2y is not positive or the quotient is not finite.
func beta(c, m2y float64) (float64, bool) {
	if !(m2y > 0) {
		return 0, false
	}
	b := c / m2y
	if math.IsNaN(b) || math.IsInf(b, 0) {
		return 0, false
	}
	return b, true
}
