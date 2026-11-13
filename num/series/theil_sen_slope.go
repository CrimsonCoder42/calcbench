package series

import (
	"cmp"
	"errors"
	"fmt"
	"math"
	"math/rand/v2"
	"slices"
	"time"

	"github.com/CrimsonCoder42/calcbench/algo/reduce"
)

// Errors returned by this package. They are wrapped with the offending parameter; test with errors.Is.
var (
	// ErrTooFew: the input has fewer elements than the computation needs.
	ErrTooFew = errors.New("series: too few values")
	// ErrNotFinite: a value is NaN or ±Inf.
	ErrNotFinite = errors.New("series: value is NaN or infinite")
	// ErrLengthMismatch: paired inputs have different lengths.
	ErrLengthMismatch = errors.New("series: inputs have different lengths")
	// ErrUnsorted: points are not in non-decreasing time order.
	ErrUnsorted = errors.New("series: points are not sorted by time")
	// ErrOutOfRange: a parameter is outside its allowed range (a window, a stride, a day count, a weekday).
	ErrOutOfRange = errors.New("series: parameter out of range")
	// ErrOverflow: an intermediate or the result is beyond the float64 range.
	ErrOverflow = errors.New("series: result exceeds the float64 range")
	// ErrNilLocation: a *time.Location is nil.
	ErrNilLocation = errors.New("series: location is nil")
	// ErrNilRandom: a *rand.Rand is nil.
	ErrNilRandom = errors.New("series: random source is nil")
)

// Point is one observation of a time series: an instant and a value in the caller's unit.
type Point struct {
	At    time.Time
	Value float64
}

// TheilSenSlope returns the Theil–Sen slope of the points (x[i], y[i]): the median of the slopes
// (y[j] − y[i])/(x[j] − x[i]) over all pairs with x[i] ≠ x[j] (pairs with equal x are left out, as Sen
// does); with an even number of slopes, the mean of the two middle ones. It selects the median by
// randomized slope selection instead of listing the n(n−1)/2 slopes: the slopes below a value t are
// counted as the inversions of y − t·x in x order (merge sort), a random sample of the slopes inside the
// current interval narrows it around the median rank, and when at most max(4n, 4096) slopes remain they
// are listed and selected. rng drives the sampling only; the result does not depend on it. Slopes that
// coincide to within the rounding of y − t·x may be ordered among themselves differently than by their
// divided values; the result then differs from the listed median by at most that rounding.
// x and y are read, never modified.
// Errors: ErrNilRandom when rng is nil; ErrLengthMismatch when len(y) ≠ len(x); ErrTooFew when there are
// fewer than 2 points or every x is equal; ErrNotFinite when a value is NaN or ±Inf; ErrOutOfRange when
// len(x) > 2³¹ − 1; ErrOverflow when 2·max|y|/(least gap between distinct x), or that bound times max|x|,
// is beyond the float64 range (a slope or y − t·x could overflow).
//
// ID:      NUM-SERIES-001
// Time:    O(n log n) expected   n = len(x); a constant expected number of rounds, each O(n log n); O(n²) when n ≤ 90 (listing every slope is cheaper there)
// Space:   O(max(n, 2²²)) extra (the sorted points, two merge buffers, the sample and the final candidates; when rounding stalls the contraction on tied slopes, up to max(4n, 2²²) candidates); result is a scalar
// Source:  Sen PK, "Estimates of the regression coefficient based on Kendall's tau", Journal of the American Statistical Association, 1968, 63(324):1379–1389, doi:10.1080/01621459.1968.10480934 (the median of the pairwise slopes, pairs with equal x omitted)
// Source:  Dillencourt MB, Mount DM, Netanyahu NS, "A randomized algorithm for slope selection", International Journal of Computational Geometry & Applications, 1992, 2(1):1–27, doi:10.1142/S0218195992000020 (k-th smallest slope in expected O(n log n): inversion counting and random sampling of the slopes in an interval)
// Source:  Matoušek J, "Randomized optimal algorithm for slope selection", Information Processing Letters, 1991, 39(4):183–187, doi:10.1016/0020-0190(91)90177-J (the interval-contraction scheme)
// Units:   x in unit X; y in unit Y; result in unit Y per unit X
// Example: x = 0, 1, 2, 3, y = 1, 3, 2, 7 → slopes 2, 0.5, 2, −1, 2, 5 → sorted −1, 0.5, 2, 2, 2, 5 → 2 (hand-checked)
// Tags:    theil-sen, sen-slope, slope-selection, robust-regression, trend, median, time-series
func TheilSenSlope(x, y []float64, rng *rand.Rand) (float64, error) {
	if rng == nil {
		return 0, fmt.Errorf("%w: rng", ErrNilRandom)
	}
	if len(y) != len(x) {
		return 0, fmt.Errorf("%w: len(y) = %d, len(x) = %d", ErrLengthMismatch, len(y), len(x))
	}
	n := len(x)
	if n < 2 {
		return 0, fmt.Errorf("%w: len(x) = %d, want ≥ 2", ErrTooFew, n)
	}
	if n > math.MaxInt32 {
		return 0, fmt.Errorf("%w: len(x) = %d, want ≤ 2³¹ − 1", ErrOutOfRange, n)
	}
	for i := range x {
		if !reduce.AllFinite(x[i]) {
			return 0, fmt.Errorf("%w: x[%d] = %v", ErrNotFinite, i, x[i])
		}
		if !reduce.AllFinite(y[i]) {
			return 0, fmt.Errorf("%w: y[%d] = %v", ErrNotFinite, i, y[i])
		}
	}
	w := newSlopeWork(x, y)
	total := w.distinctXPairs()
	if total == 0 {
		return 0, fmt.Errorf("%w: every x equals %v, no slope is defined", ErrTooFew, x[0])
	}
	if err := w.checkRange(); err != nil {
		return 0, err
	}
	return w.median(total, rng), nil
}

// slopeBound is a slope value or ±∞ (inf = −1 or +1).
type slopeBound struct {
	t   float64
	inf int
}

// keyed is a point index with its comparison key at one slope bound, compared lexicographically.
type keyed struct {
	k1, k2 float64
	p      int32
}

func keyLess(a, b keyed) bool { return a.k1 < b.k1 || (a.k1 == b.k1 && a.k2 < b.k2) }

// slopeWork holds the points sorted by (x, y) and the buffers of one TheilSenSlope call.
type slopeWork struct {
	xs, ys   []float64
	u        []float64
	perm     []int32
	seq, buf []keyed
}

func newSlopeWork(x, y []float64) *slopeWork {
	n := len(x)
	perm := make([]int32, n)
	for i := range perm {
		perm[i] = int32(i)
	}
	slices.SortFunc(perm, func(a, b int32) int {
		if c := cmp.Compare(x[a], x[b]); c != 0 {
			return c
		}
		return cmp.Compare(y[a], y[b])
	})
	w := &slopeWork{xs: make([]float64, n), ys: make([]float64, n), u: make([]float64, n), perm: perm,
		seq: make([]keyed, n), buf: make([]keyed, n)}
	for i, p := range perm {
		w.xs[i], w.ys[i] = x[p], y[p]
	}
	return w
}

// distinctXPairs returns the number of pairs with different x: n(n−1)/2 minus the pairs inside each run
// of equal x.
func (w *slopeWork) distinctXPairs() int64 {
	n := int64(len(w.xs))
	total := n * (n - 1) / 2
	run := int64(1)
	for i := 1; i <= len(w.xs); i++ {
		if i < len(w.xs) && w.xs[i] == w.xs[i-1] {
			run++
			continue
		}
		total -= run * (run - 1) / 2
		run = 1
	}
	return total
}

// checkRange returns ErrOverflow when a slope, or y − t·x for t a slope, could leave the float64 range.
func (w *slopeWork) checkRange() error {
	var maxX, maxY float64
	minGap := math.Inf(1)
	for i := range w.xs {
		maxX = math.Max(maxX, math.Abs(w.xs[i]))
		maxY = math.Max(maxY, math.Abs(w.ys[i]))
		if i > 0 && w.xs[i] != w.xs[i-1] {
			minGap = math.Min(minGap, w.xs[i]-w.xs[i-1])
		}
	}
	bound := 2 * maxY / minGap
	if math.IsInf(minGap, 0) || !reduce.AllFinite(bound) || !reduce.AllFinite(bound*maxX+maxY) || !reduce.AllFinite(maxX*2) {
		return fmt.Errorf("%w: slope bound 2·max|y|/min Δx = %v with max|x| = %v", ErrOverflow, bound, maxX)
	}
	return nil
}

func (w *slopeWork) slope(p, q int32) float64 {
	return (w.ys[q] - w.ys[p]) / (w.xs[q] - w.xs[p])
}

// count returns the number of slopes strictly below t.
func (w *slopeWork) count(t float64) int64 {
	c, _ := w.pass(slopeBound{inf: -1}, slopeBound{t: t}, nil, false, nil)
	return c
}

// pass lists the points in their order at lo (by y − lo·x, then x, then y), keys them at hi, and merge
// sorts the keys counting strict inversions: each inversion is one pair whose slope lies in [lo, hi). It
// returns that count and appends to out the slopes of the inversions whose discovery ranks are in ranks
// (ascending), or of every inversion when all is set.
func (w *slopeWork) pass(lo, hi slopeBound, ranks []int64, all bool, out []float64) (int64, []float64) {
	n := len(w.xs)
	perm := w.perm
	for i := range perm {
		perm[i] = int32(i)
	}
	if lo.inf == 0 {
		for i := range w.u {
			w.u[i] = math.FMA(-lo.t, w.xs[i], w.ys[i])
		}
		slices.SortFunc(perm, func(a, b int32) int {
			if c := cmp.Compare(w.u[a], w.u[b]); c != 0 {
				return c
			}
			if c := cmp.Compare(w.xs[a], w.xs[b]); c != 0 {
				return c
			}
			return cmp.Compare(w.ys[a], w.ys[b])
		})
	}
	for i, p := range perm {
		if hi.inf > 0 {
			w.seq[i] = keyed{k1: -w.xs[p], k2: w.ys[p], p: p}
		} else {
			w.seq[i] = keyed{k1: math.FMA(-hi.t, w.xs[p], w.ys[p]), p: p}
		}
	}
	src, dst := w.seq, w.buf
	var inv int64
	ri := 0
	for width := 1; width < n; width *= 2 {
		for start := 0; start < n; start += 2 * width {
			mid := min(start+width, n)
			end := min(start+2*width, n)
			i, j, k := start, mid, start
			for i < mid && j < end {
				if keyLess(src[j], src[i]) {
					r := int64(mid - i)
					if all {
						for off := i; off < mid; off++ {
							out = append(out, w.slope(src[off].p, src[j].p))
						}
					} else {
						for ri < len(ranks) && ranks[ri] < inv+r {
							out = append(out, w.slope(src[i+int(ranks[ri]-inv)].p, src[j].p))
							ri++
						}
					}
					inv += r
					dst[k] = src[j]
					j++
				} else {
					dst[k] = src[i]
					i++
				}
				k++
			}
			k += copy(dst[k:], src[i:mid])
			copy(dst[k:], src[j:end])
		}
		src, dst = dst, src
	}
	return inv, out
}

// median returns the median of the total slopes (total ≥ 1).
func (w *slopeWork) median(total int64, rng *rand.Rand) float64 {
	n := len(w.xs)
	r1, r2 := (total-1)/2, total/2
	enumCap := int64(max(4*n, 4096))
	if total <= enumCap {
		all := make([]float64, 0, total)
		for i := 0; i < n; i++ {
			for j := i + 1; j < n; j++ {
				if w.xs[j] != w.xs[i] {
					all = append(all, w.slope(int32(i), int32(j)))
				}
			}
		}
		slices.Sort(all)
		return mean2(all[r1], all[r2])
	}
	lo, hi := slopeBound{inf: -1}, slopeBound{inf: 1}
	cLo, cHi := int64(0), total
	m := max(n, 1024)
	ranks := make([]int64, m)
	var sample []float64
	stall := 0
	for round := 0; cHi-cLo > enumCap && round < 64 && stall < 4; round++ {
		k := cHi - cLo
		for i := range ranks {
			ranks[i] = rng.Int64N(k)
		}
		slices.Sort(ranks)
		_, sample = w.pass(lo, hi, ranks, false, sample[:0])
		slices.Sort(sample)
		fm, fk := float64(len(sample)), float64(k)
		d := 3 * math.Sqrt(fm)
		iLo := int(math.Floor(float64(r1-cLo)/fk*fm - d))
		iHi := int(math.Ceil(float64(r2-cLo)/fk*fm + d))
		if iLo >= 0 && iLo < len(sample) {
			t := sample[iLo]
			if (lo.inf < 0 || t > lo.t) && (hi.inf > 0 || t < hi.t) {
				if c := w.count(t); c <= r1 {
					lo, cLo = slopeBound{t: t}, c
				}
			}
		}
		if iHi >= 0 && iHi < len(sample) {
			t := math.Nextafter(sample[iHi], math.Inf(1))
			if (hi.inf > 0 || t < hi.t) && (lo.inf < 0 || t > lo.t) {
				if c := w.count(t); c > r2 {
					hi, cHi = slopeBound{t: t}, c
				}
			}
		}
		if cHi-cLo > k-k/10 {
			stall++
		} else {
			stall = 0
		}
	}
	if cHi-cLo <= max(enumCap, 1<<22) {
		_, list := w.pass(lo, hi, nil, true, make([]float64, 0, cHi-cLo))
		if len(list) > 0 {
			slices.Sort(list)
			last := int64(len(list) - 1)
			return mean2(list[clamp(r1-cLo, 0, last)], list[clamp(r2-cLo, 0, last)])
		}
	}
	// Rounding stalled the contraction on a block of slopes it cannot separate; they agree to within that
	// rounding, so the sample's order statistics at the median position stand for them.
	if len(sample) == 0 {
		if lo.inf == 0 {
			return lo.t
		}
		return hi.t
	}
	fm, fk, last := float64(len(sample)), float64(cHi-cLo), int64(len(sample)-1)
	a := sample[clamp(int64(float64(r1-cLo)/fk*fm), 0, last)]
	b := sample[clamp(int64(float64(r2-cLo)/fk*fm), 0, last)]
	return mean2(a, b)
}

func clamp(v, lo, hi int64) int64 { return max(lo, min(v, hi)) }

// mean2 returns (a + b)/2 without overflow.
func mean2(a, b float64) float64 {
	if (a >= 0) == (b >= 0) {
		return a + (b-a)/2
	}
	return (a + b) / 2
}
