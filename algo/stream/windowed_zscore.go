package stream

import "math"

// WindowedZScore is the standard score of the newest value against the last windowLen values,
// z = (x − x̄)/s, where x̄ and s are the mean and the sample standard deviation of the window that
// includes x itself (the score of a member of the sample, as NUM-STATS-028 ZScore computes it for a
// whole slice). It wraps a WindowedVariance (ALG-STREAM-005), so every Add is O(1) with no loop over the
// window and its drift behaviour and Resync are that type's. The zero value is not usable: its Add
// returns ErrUninitialized. Not safe for concurrent use. No Merge (windows have no meaningful merge).
//
// Methods: Add(x) adds a value. ZScore returns the score of the newest value, Mean and Std the window's;
// each returns (0, false) when undefined (fewer than two values for ZScore and Std, or a window with no
// spread for ZScore; no value for Mean). Len, WindowLen, Resync and Reset as in WindowedVariance.
//
// ID:      ALG-STREAM-007
// Time:    O(1) worst case per Add and per reader; Resync and Reset O(w)   w = windowLen
// Space:   O(w), allocated only in NewWindowedZScore; Add never allocates
// Source:  NIST, Dataplot Reference Manual, "ISO 13528 ZSCORE", https://www.itl.nist.gov/div898/software/dataplot/refman2/auxillar/zscore.htm (read 2026-11-12) (z = (x − mean)/sd)
// Source:  West, "Updating mean and variance estimates: an improved method", Communications of the ACM, 1979, 22(9):532–535, doi:10.1145/359146.359153 (the sliding mean and variance, via ALG-STREAM-005)
// Tags:    stream, sliding-window, rolling-zscore, z-score, standard-score, standardize, anomaly, mean-reversion, time-series
// Units:   values x in any unit U; Mean and Std in U, ZScore unitless (standard deviations)
// Example: windowLen 3, Add 1, 2, 3 → Mean = 2, Std = 1, ZScore = (3 − 2)/1 = 1
type WindowedZScore struct {
	stats  WindowedVariance
	newest float64
}

// NewWindowedZScore returns an empty WindowedZScore over the last windowLen values.
// Errors: ErrWindowLen when windowLen < 1 or windowLen > 1<<24.
func NewWindowedZScore(windowLen int) (*WindowedZScore, error) {
	buf, err := newWindow[float64](windowLen)
	if err != nil {
		return nil, err
	}
	return &WindowedZScore{stats: WindowedVariance{window: buf}}, nil
}

// Add appends x to the window, dropping the oldest value when the window is full. O(1), no loop.
// Errors: those of WindowedVariance.Add (ErrNotFinite, ErrOverflow, ErrUninitialized). On error
// nothing changes.
func (z *WindowedZScore) Add(x float64) error {
	if err := z.stats.Add(x); err != nil {
		return err
	}
	z.newest = x
	return nil
}

// ZScore returns (x − x̄)/s for the newest value x, or (0, false) when the window holds fewer than two
// values, has no spread (s = 0), or the score exceeds the float64 range. O(1).
func (z *WindowedZScore) ZScore() (float64, bool) {
	s, ok := z.stats.Std()
	if !ok || !(s > 0) {
		return 0, false
	}
	score := (z.newest - z.stats.mean) / s
	if math.IsInf(score, 0) || math.IsNaN(score) {
		return 0, false
	}
	return score, true
}

// Mean returns the mean of the window, or (0, false) when it is empty. O(1).
func (z *WindowedZScore) Mean() (float64, bool) { return z.stats.Mean() }

// Std returns the sample standard deviation of the window, or (0, false) when it holds fewer than two
// values. O(1).
func (z *WindowedZScore) Std() (float64, bool) { return z.stats.Std() }

// Resync recomputes the window's mean and variance from the held values (WindowedVariance.Resync). O(w).
func (z *WindowedZScore) Resync() { z.stats.Resync() }

// Reset empties the window and keeps its length. O(w).
func (z *WindowedZScore) Reset() {
	z.stats.Reset()
	z.newest = 0
}

// Len returns the number of values in the window. O(1).
func (z *WindowedZScore) Len() int { return z.stats.Len() }

// WindowLen returns the window length given to NewWindowedZScore (0 for the zero value). O(1).
func (z *WindowedZScore) WindowLen() int { return z.stats.WindowLen() }
