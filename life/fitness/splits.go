package fitness

import (
	"fmt"
	"math"
	"time"
)

// DistanceSample is one sample of a cumulative distance stream: the total distance covered by At, in
// kilometres or miles (one unit for the whole stream).
type DistanceSample struct {
	At             time.Time
	DistanceKmOrMi float64
}

// Split is one completed whole unit of distance: its 1-based index k (the unit ending at k·unit), its own
// duration, and the elapsed time from the first sample to its end.
type Split struct {
	Index          int
	Seconds        float64
	ElapsedSeconds float64
}

// Splits returns the split times of a cumulative distance stream for whole units of unitKmOrMi (1 for
// kilometre splits of a stream in km, or mile splits of a stream in mi): for each boundary k·unitKmOrMi
// crossed after the first sample, the crossing instant is interpolated linearly between the two samples
// around it, and the split is the time since the previous crossing (the first split's time is measured
// from the first sample). A partial last unit is not a split. samples is read, never modified.
// Errors: ErrNotFinite for a NaN or ±Inf distance or unit; ErrOutOfRange when unitKmOrMi ≤ 0, when there
// are fewer than 2 samples, when a distance is negative, or when a sample's time or distance is less than
// the previous one's; ErrDomain when the stream would give more than 2³¹ − 1 splits.
//
// ID:      LIF-FIT-010
// Time:    O(n + s)   n = len(samples), s = splits; two passes (count, then fill)
// Space:   O(1) extra; O(s) for the result (one allocation)
// Source:  Burden RL, Faires JD, Numerical Analysis, 9th ed., Brooks/Cole, 2011, §3.1 "Interpolation and the Lagrange Polynomial" (the linear interpolating polynomial through two points)
// Tags:    fitness, running, cycling, splits, pace, interpolation
// Units:   samples' distances km or mi; unitKmOrMi in the same unit; Seconds and ElapsedSeconds s
// Range:   unitKmOrMi > 0; at least 2 samples, non-decreasing in time and in distance, distances ≥ 0
// Example: (0 s, 0 mi), (300 s, 0.5), (700 s, 1.25), (1000 s, 2.0), unit 1 mi → mile 1 at 300 + 400·(0.5/0.75) = 566.667 s; mile 2 at 1000 s, split 433.333 s (HAND-DERIVED)
func Splits(samples []DistanceSample, unitKmOrMi float64) ([]Split, error) {
	if err := checkFinite("unitKmOrMi", unitKmOrMi); err != nil {
		return nil, err
	}
	if unitKmOrMi <= 0 {
		return nil, fmt.Errorf("%w: unitKmOrMi = %v, want > 0", ErrOutOfRange, unitKmOrMi)
	}
	if len(samples) < 2 {
		return nil, fmt.Errorf("%w: len(samples) = %d, want ≥ 2", ErrOutOfRange, len(samples))
	}
	for i, s := range samples {
		if math.IsNaN(s.DistanceKmOrMi) || math.IsInf(s.DistanceKmOrMi, 0) {
			return nil, fmt.Errorf("%w: samples[%d].DistanceKmOrMi = %v", ErrNotFinite, i, s.DistanceKmOrMi)
		}
		if s.DistanceKmOrMi < 0 {
			return nil, fmt.Errorf("%w: samples[%d].DistanceKmOrMi = %v, want ≥ 0", ErrOutOfRange, i, s.DistanceKmOrMi)
		}
		if i > 0 && (s.At.Before(samples[i-1].At) || s.DistanceKmOrMi < samples[i-1].DistanceKmOrMi) {
			return nil, fmt.Errorf("%w: samples[%d] goes back in time or distance", ErrOutOfRange, i)
		}
	}
	firstK := math.Floor(samples[0].DistanceKmOrMi/unitKmOrMi) + 1
	if est := math.Floor(samples[len(samples)-1].DistanceKmOrMi/unitKmOrMi) - firstK + 1; est > math.MaxInt32 {
		return nil, fmt.Errorf("%w: about %v splits of %v, want ≤ 2³¹ − 1", ErrDomain, est, unitKmOrMi)
	}
	count := walkSplits(samples, unitKmOrMi, int(firstK), nil)
	out := make([]Split, 0, count)
	walkSplits(samples, unitKmOrMi, int(firstK), &out)
	return out, nil
}

// walkSplits walks the boundaries k·unit, k = firstK, firstK+1, …, that the stream crosses, appending
// each split to out when out is not nil, and returns how many there are.
func walkSplits(samples []DistanceSample, unit float64, firstK int, out *[]Split) int {
	t0 := samples[0].At
	k, prev, n := firstK, 0.0, 0
	for i := 1; i < len(samples); i++ {
		a, b := samples[i-1], samples[i]
		for b.DistanceKmOrMi >= float64(k)*unit {
			if out != nil {
				target := float64(k) * unit
				frac := 0.0
				if b.DistanceKmOrMi > a.DistanceKmOrMi {
					frac = (target - a.DistanceKmOrMi) / (b.DistanceKmOrMi - a.DistanceKmOrMi)
				}
				at := a.At.Sub(t0).Seconds() + frac*b.At.Sub(a.At).Seconds()
				*out = append(*out, Split{Index: k, Seconds: at - prev, ElapsedSeconds: at})
				prev = at
			}
			k++
			n++
		}
	}
	return n
}
