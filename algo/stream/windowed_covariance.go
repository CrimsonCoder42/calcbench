package stream

import (
	"fmt"
	"math"

	"github.com/CrimsonCoder42/calcbench/algo/queue"
	"github.com/CrimsonCoder42/calcbench/algo/reduce"
)

// pair is one (x, y) observation held in a window.
type pair struct{ x, y float64 }

// WindowedCovariance is the covariance, correlation and beta of the last windowLen pairs (x, y) added
// (all of them while fewer have arrived), kept one pair at a time by the sliding update: each Add puts
// the new pair into the running means, the two M2 sums and the co-moment C = Σ(x − x̄)(y − ȳ) and, once
// the window is full, takes the oldest pair out, in O(1), never re-summing the window. The pairs live in
// a queue.CircularBuffer (ALG-QUEUE-002). The zero value is not usable: its Add returns
// ErrUninitialized. Not safe for concurrent use. Windows have no meaningful merge, so there is no Merge.
//
// Update. While filling, the k-th pair: δx = x − x̄; x̄′ = x̄ + δx/k; ȳ′ likewise; C ← C + δx·(y − ȳ′).
// Once full (w pairs), (x, y) replaces the oldest (xo, yo) (West 1979's weighted update, +1 for the new
// pair and −1 for the old, combined; for the co-moment the same algebra as Chan, Golub and LeVeque's
// updating formula): dx = x − xo, dy = y − yo; x̄′ = x̄ + dx/w; ȳ′ = ȳ + dy/w;
// C ← C + dx·(yo − ȳ) + dy·(x − x̄′); M2x ← M2x + dx·((x − x̄′) + (xo − x̄)), M2y likewise.
// Covariance C/(n − 1); correlation C/√(M2x·M2y); beta C/M2y, the least-squares slope of x on y (with x
// an asset's returns and y a benchmark's, the market-model beta). n = Len().
//
// Drift: as WindowedVariance (ALG-STREAM-005): the rounding of each Add accumulates over the whole
// stream; there is no recompute inside Add, so Add stays O(1) worst case; Resync, O(w), recomputes the
// moments from the held pairs when the caller chooses. M2x and M2y are clamped at 0.
//
// Methods: Add(x, y) adds a pair. MeanX, MeanY, Covariance (sample), Correlation and Beta read the
// window; each returns (0, false) when undefined (no pairs; fewer than two for the others; a zero
// spread in x or y for Correlation, in y for Beta). Len, WindowLen, Resync and Reset as in
// WindowedVariance.
//
// ID:      ALG-STREAM-006
// Time:    O(1) worst case per Add and per reader; Resync and Reset O(w)   w = windowLen
// Space:   O(w), allocated only in NewWindowedCovariance; Add never allocates
// Source:  West, "Updating mean and variance estimates: an improved method", Communications of the ACM, 1979, 22(9):532–535, doi:10.1145/359146.359153 (the weighted update; removal is the update with weight −1)
// Source:  Chan, Golub, LeVeque, "Algorithms for computing the sample variance: analysis and recommendations", The American Statistician, 1983, 37(3):242–247, doi:10.1080/00031305.1983.10483115 (the pairwise updating formula, applied to the co-moment)
// Tags:    stream, sliding-window, rolling-covariance, rolling-correlation, rolling-beta, covariance, correlation, beta, welford, west, ring-buffer, time-series
// Units:   x in any unit U, y in any unit V; MeanX in U, MeanY in V, Covariance in U·V, Correlation unitless, Beta in U/V
// Example: windowLen 3, Add (1, 2), (2, 4), (3, 6) → Covariance = 2, Correlation = 1, Beta = 0.5
type WindowedCovariance struct {
	window       *queue.CircularBuffer[pair]
	meanX, meanY float64
	m2X, m2Y, c  float64
}

// NewWindowedCovariance returns an empty WindowedCovariance over the last windowLen pairs.
// Errors: ErrWindowLen when windowLen < 1 or windowLen > 1<<24.
func NewWindowedCovariance(windowLen int) (*WindowedCovariance, error) {
	buf, err := newWindow[pair](windowLen)
	if err != nil {
		return nil, err
	}
	return &WindowedCovariance{window: buf}, nil
}

// Add appends the pair (x, y) to the window, dropping the oldest pair when the window is full. O(1),
// no loop.
// Errors: ErrNotFinite when x or y is NaN or ±Inf; ErrOverflow when a difference, a mean or a moment
// would leave the float64 range; ErrUninitialized on the zero value. On error nothing changes.
func (v *WindowedCovariance) Add(x, y float64) error {
	if v.window == nil {
		return ErrUninitialized
	}
	if math.IsNaN(x) || math.IsInf(x, 0) || math.IsNaN(y) || math.IsInf(y, 0) {
		return fmt.Errorf("%w: (x, y) = (%v, %v)", ErrNotFinite, x, y)
	}
	n := v.window.Len()
	var dx, dy, mx, my, m2X, m2Y, c float64
	switch {
	case n < v.window.Cap():
		k := float64(n + 1)
		dx, dy = x-v.meanX, y-v.meanY
		mx, my = v.meanX+dx/k, v.meanY+dy/k
		m2X, m2Y, c = v.m2X+dx*(x-mx), v.m2Y+dy*(y-my), v.c+dx*(y-my)
	case n == 1:
		mx, my = x, y
	default:
		old, _ := v.window.Front()
		w := float64(n)
		dx, dy = x-old.x, y-old.y
		mx, my = v.meanX+dx/w, v.meanY+dy/w
		m2X = v.m2X + dx*((x-mx)+(old.x-v.meanX))
		m2Y = v.m2Y + dy*((y-my)+(old.y-v.meanY))
		c = v.c + dx*(old.y-v.meanY) + dy*(x-mx)
	}
	if !reduce.AllFinite(dx, dy, mx, my, m2X, m2Y, c) {
		return fmt.Errorf("%w: (x, y) = (%v, %v)", ErrOverflow, x, y)
	}
	v.window.Push(pair{x, y}) // Overwrite policy: never fails on a constructed buffer
	v.meanX, v.meanY, v.m2X, v.m2Y, v.c = mx, my, max(m2X, 0), max(m2Y, 0), c
	return nil
}

// Resync recomputes the means and moments from the pairs held, in two passes (compensated sums for the
// means, then the centred products), discarding the drift of earlier Adds. If a pass overflows, the
// state is kept. O(w). No-op on the zero value.
func (v *WindowedCovariance) Resync() {
	if v.window == nil || v.window.Len() == 0 {
		return
	}
	n := v.window.Len()
	var sx, cx, sy, cy float64
	for i := range n {
		p, _ := v.window.At(i)
		sx, cx = reduce.NeumaierAdd(sx, cx, p.x)
		sy, cy = reduce.NeumaierAdd(sy, cy, p.y)
	}
	mx, my := (sx+cx)/float64(n), (sy+cy)/float64(n)
	var m2X, m2Y, c float64
	for i := range n {
		p, _ := v.window.At(i)
		dx, dy := p.x-mx, p.y-my
		m2X += dx * dx
		m2Y += dy * dy
		c += dx * dy
	}
	if reduce.AllFinite(sx, sy, mx, my, m2X, m2Y, c) {
		v.meanX, v.meanY, v.m2X, v.m2Y, v.c = mx, my, m2X, m2Y, c
	}
}

// Reset empties the window and keeps its length. O(w).
func (v *WindowedCovariance) Reset() {
	if v.window == nil {
		return
	}
	clearWindow(v.window)
	v.meanX, v.meanY, v.m2X, v.m2Y, v.c = 0, 0, 0, 0, 0
}

// MeanX returns the mean of the x values in the window, or (0, false) when it is empty. O(1).
func (v *WindowedCovariance) MeanX() (float64, bool) {
	if v.Len() == 0 {
		return 0, false
	}
	return v.meanX, true
}

// MeanY returns the mean of the y values in the window, or (0, false) when it is empty. O(1).
func (v *WindowedCovariance) MeanY() (float64, bool) {
	if v.Len() == 0 {
		return 0, false
	}
	return v.meanY, true
}

// Covariance returns the sample covariance C/(n − 1) of the window, or (0, false) when it holds fewer
// than two pairs. O(1).
func (v *WindowedCovariance) Covariance() (float64, bool) {
	n := v.Len()
	if n < 2 {
		return 0, false
	}
	return v.c / float64(n-1), true
}

// Correlation returns Pearson's r = C/√(M2x·M2y) of the window, clamped to [−1, 1] against rounding,
// or (0, false) when it holds fewer than two pairs or x or y has no spread. O(1).
func (v *WindowedCovariance) Correlation() (float64, bool) {
	if v.Len() < 2 {
		return 0, false
	}
	return correlation(v.c, v.m2X, v.m2Y)
}

// Beta returns C/M2y, the covariance of the window over the variance of its y values (the least-squares
// slope of x on y), or (0, false) when it holds fewer than two pairs, y has no spread, or the slope
// exceeds the float64 range. O(1).
func (v *WindowedCovariance) Beta() (float64, bool) {
	if v.Len() < 2 {
		return 0, false
	}
	return beta(v.c, v.m2Y)
}

// Len returns the number of pairs in the window: min(pairs added since Reset, WindowLen()). O(1).
func (v *WindowedCovariance) Len() int {
	if v.window == nil {
		return 0
	}
	return v.window.Len()
}

// WindowLen returns the window length given to NewWindowedCovariance (0 for the zero value). O(1).
func (v *WindowedCovariance) WindowLen() int {
	if v.window == nil {
		return 0
	}
	return v.window.Cap()
}
