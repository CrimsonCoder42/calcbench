package stream

import (
	"fmt"
	"math"

	"github.com/CrimsonCoder42/calcbench/algo/queue"
	"github.com/CrimsonCoder42/calcbench/algo/reduce"
)

// WindowedVariance is the mean and variance of the last windowLen values added (all of them while fewer
// have arrived), kept one value at a time by the sliding Welford update: each Add puts the new value
// into the running mean and M2 (the sum of squared deviations from the mean) and, once the window is
// full, takes the oldest out, in O(1), never re-summing the window. The values live in a
// queue.CircularBuffer (ALG-QUEUE-002). The zero value is not usable: its Add returns ErrUninitialized.
// Not safe for concurrent use. Windows have no meaningful merge, so there is no Merge.
//
// Update. While the window is filling, the k-th value x: δ = x − mean; mean ← mean + δ/k;
// M2 ← M2 + δ·(x − mean). Once full (w values), x replaces the oldest value o (West 1979, the weighted
// update with weight +1 for x and −1 for o, combined): d = x − o; mean′ = mean + d/w;
// M2 ← M2 + d·((x − mean′) + (o − mean)); mean ← mean′. A window of length 1 holds x with M2 = 0.
// Sample variance M2/(n − 1), population variance M2/n, n = Len().
//
// Drift. Each Add rounds; the errors of mean and M2 add up as a random walk over the whole stream, not
// only over the window. Removing a value of magnitude L can leave an absolute error of order u·L² in M2
// (u = 2⁻⁵³), so after a large outlier leaves, a window of small spread can read a variance that is
// mostly rounding. The error is of the size of the data that has passed through, not of the window's
// own variance: for a very short window (w ≤ 3) whose values momentarily agree, the relative error of
// the variance can be large (a simulation over 1,000 values 50 ± 10 measured up to 3e-4 relative at
// w = 2, under 1e-10 at w = 3 and under 5e-13 at w ≥ 7). M2 is clamped at 0. There is deliberately no periodic recompute inside Add:
// Add stays O(1) worst case, which a latency-bound caller needs (an amortised recompute would make one
// Add in w cost O(w)). A caller that wants the drift bounded calls Resync, O(w), when it chooses (for
// example after an outlier leaves, or every w Adds for O(1) amortised).
//
// Methods: Add(x) adds a value. Mean, Variance (sample), PopulationVariance and Std (sample standard
// deviation) read the window; each returns (0, false) when undefined (no values; fewer than two for
// Variance and Std). Len is the number of values held, WindowLen the window length. Resync recomputes
// mean and M2 from the held values in two passes. Reset empties the window.
//
// ID:      ALG-STREAM-005
// Time:    O(1) worst case per Add and per reader; Resync and Reset O(w)   w = windowLen
// Space:   O(w), allocated only in NewWindowedVariance; Add never allocates
// Source:  West, "Updating mean and variance estimates: an improved method", Communications of the ACM, 1979, 22(9):532–535, doi:10.1145/359146.359153 (the weighted update; removal is the update with weight −1)
// Source:  Welford, "Note on a method for calculating corrected sums of squares and products", Technometrics, 1962, 4(3):419–420, doi:10.1080/00401706.1962.10490022 (the filling phase)
// Tags:    stream, sliding-window, rolling-variance, moving-variance, rolling-std, welford, west, online-variance, ring-buffer, time-series
// Units:   values x in any unit U; Mean and Std in U, Variance and PopulationVariance in U²
// Example: windowLen 3, Add 1, 2, 3, 4 → window {2, 3, 4}: Mean = 3, Variance = 1, Std = 1, PopulationVariance = 2/3
type WindowedVariance struct {
	window   *queue.CircularBuffer[float64]
	mean, m2 float64
}

// newWindow returns an overwriting ring buffer of windowLen slots, or ErrWindowLen.
func newWindow[T any](windowLen int) (*queue.CircularBuffer[T], error) {
	if windowLen < 1 || windowLen > maxWindowLen {
		return nil, fmt.Errorf("%w: windowLen = %d", ErrWindowLen, windowLen)
	}
	buf, err := queue.NewCircularBuffer[T](windowLen, queue.Overwrite)
	if err != nil {
		return nil, fmt.Errorf("%w: windowLen = %d: %w", ErrWindowLen, windowLen, err)
	}
	return buf, nil
}

// clearWindow pops every element. O(Len).
func clearWindow[T any](b *queue.CircularBuffer[T]) {
	for b.Len() > 0 {
		b.Pop()
	}
}

// NewWindowedVariance returns an empty WindowedVariance over the last windowLen values.
// Errors: ErrWindowLen when windowLen < 1 or windowLen > 1<<24.
func NewWindowedVariance(windowLen int) (*WindowedVariance, error) {
	buf, err := newWindow[float64](windowLen)
	if err != nil {
		return nil, err
	}
	return &WindowedVariance{window: buf}, nil
}

// Add appends x to the window, dropping the oldest value when the window is full. O(1), no loop.
// Errors: ErrNotFinite when x is NaN or ±Inf; ErrOverflow when a difference, the mean or M2 would leave
// the float64 range; ErrUninitialized on the zero value. On error nothing changes.
func (v *WindowedVariance) Add(x float64) error {
	if v.window == nil {
		return ErrUninitialized
	}
	if math.IsNaN(x) || math.IsInf(x, 0) {
		return fmt.Errorf("%w: x = %v", ErrNotFinite, x)
	}
	mean, m2, d := slideMoment(v.window, v.mean, v.m2, x)
	if math.IsInf(d, 0) || math.IsInf(mean, 0) || math.IsInf(m2, 0) || math.IsNaN(m2) {
		return fmt.Errorf("%w: x = %v", ErrOverflow, x)
	}
	v.window.Push(x) // Overwrite policy: never fails on a constructed buffer
	v.mean, v.m2 = mean, max(m2, 0)
	return nil
}

// slideMoment returns the mean and M2 after x enters window (and its oldest value leaves when full),
// and the difference that was formed, for the caller's overflow check. It does not push x. O(1).
func slideMoment(window *queue.CircularBuffer[float64], mean, m2, x float64) (newMean, newM2, d float64) {
	n := window.Len()
	if n < window.Cap() {
		d = x - mean
		newMean = mean + d/float64(n+1)
		return newMean, m2 + d*(x-newMean), d
	}
	if n == 1 {
		return x, 0, 0
	}
	old, _ := window.Front()
	d = x - old
	newMean = mean + d/float64(n)
	return newMean, m2 + d*((x-newMean)+(old-mean)), d
}

// Resync recomputes the mean and M2 from the values held, in two passes (a compensated sum for the mean,
// then the squared deviations), discarding the drift of earlier Adds. If a pass overflows, the state is
// kept. O(w). No-op on the zero value.
func (v *WindowedVariance) Resync() {
	if v.window == nil || v.window.Len() == 0 {
		return
	}
	mean, m2, ok := twoPassWindow(v.window)
	if ok {
		v.mean, v.m2 = mean, m2
	}
}

// twoPassWindow returns the mean and M2 of the values in b, and false when a pass overflows. O(Len).
func twoPassWindow(b *queue.CircularBuffer[float64]) (mean, m2 float64, ok bool) {
	n := b.Len()
	var sum, comp float64
	for i := range n {
		x, _ := b.At(i)
		sum, comp = reduce.NeumaierAdd(sum, comp, x)
	}
	mean = (sum + comp) / float64(n)
	for i := range n {
		x, _ := b.At(i)
		m2 += (x - mean) * (x - mean)
	}
	if math.IsInf(sum, 0) || math.IsInf(mean, 0) || math.IsNaN(mean) || math.IsInf(m2, 0) {
		return 0, 0, false
	}
	return mean, m2, true
}

// Reset empties the window and keeps its length. O(w).
func (v *WindowedVariance) Reset() {
	if v.window == nil {
		return
	}
	clearWindow(v.window)
	v.mean, v.m2 = 0, 0
}

// Mean returns the mean of the values in the window, or (0, false) when there are none. O(1).
func (v *WindowedVariance) Mean() (float64, bool) {
	if v.Len() == 0 {
		return 0, false
	}
	return v.mean, true
}

// Variance returns the sample variance M2/(n − 1) of the window, or (0, false) when it holds fewer
// than two values. O(1).
func (v *WindowedVariance) Variance() (float64, bool) {
	n := v.Len()
	if n < 2 {
		return 0, false
	}
	return v.m2 / float64(n-1), true
}

// PopulationVariance returns M2/n, the variance of the window with divisor n, or (0, false) when it is
// empty. O(1).
func (v *WindowedVariance) PopulationVariance() (float64, bool) {
	n := v.Len()
	if n == 0 {
		return 0, false
	}
	return v.m2 / float64(n), true
}

// Std returns the sample standard deviation √(M2/(n − 1)) of the window, or (0, false) when it holds
// fewer than two values. O(1).
func (v *WindowedVariance) Std() (float64, bool) {
	s2, ok := v.Variance()
	if !ok {
		return 0, false
	}
	return math.Sqrt(s2), true
}

// Len returns the number of values in the window: min(values added since Reset, WindowLen()). O(1).
func (v *WindowedVariance) Len() int {
	if v.window == nil {
		return 0
	}
	return v.window.Len()
}

// WindowLen returns the window length given to NewWindowedVariance (0 for the zero value). O(1).
func (v *WindowedVariance) WindowLen() int {
	if v.window == nil {
		return 0
	}
	return v.window.Cap()
}
