package series

import (
	"fmt"
	"math"
	"time"

	"github.com/CrimsonCoder42/calcbench/algo/reduce"
)

// StrideBucket is the least, mean and greatest value and the count of the points in one fixed time stride
// that starts at Start (UTC).
type StrideBucket struct {
	Start time.Time
	Min   float64
	Mean  float64
	Max   float64
	Count int
}

// FixedStrideBuckets downsamples time-sorted points into fixed strides: each point's time is floored to a
// multiple of stride since the zero time (time.Truncate on the UTC instant, so a 5-minute stride starts on
// the same instants in every zone), and each non-empty stride gives one bucket with the least, mean and
// greatest value and the count, in time order. Strides that do not divide a day do not align with
// midnight; group by local day with NUM-SERIES-005 instead. Means use a compensated sum. points is read,
// never modified.
// Errors: ErrOutOfRange when stride ≤ 0; ErrUnsorted when a point's time is before the previous one's;
// ErrNotFinite when a value is NaN or ±Inf.
// Implementation: Go's time package — func (Time) Truncate ("the result of rounding t down to a multiple of d (since the zero time)").
//
// ID:      NUM-SERIES-006
// Time:    O(n)   n = len(points); two passes (count the strides, then fill)
// Space:   O(1) extra; O(b) for the result (one allocation), b = number of non-empty strides
// Source:  Reingold EM, Dershowitz N, Calendrical Calculations, 3rd ed., Cambridge University Press, 2008, ch. 1 (fixed day numbers, floor and mod) and ch. 2 (the proleptic Gregorian calendar)
// Source:  ISO 8601-1:2019, Date and time — Representations for information interchange — Part 1: Basic rules (calendar day; local time and its offset from UTC)
// Source:  Neumaier A, "Rundungsfehleranalyse einiger Verfahren zur Summation endlicher Summen", ZAMM, 1974, 54(1):39–51, doi:10.1002/zamm.19740540106 (compensated summation for the mean)
// Units:   values in any one unit; stride a time.Duration; Min, Mean, Max in the unit of the values
// Example: 60 at 10:00:10, 70 at 10:04:59, 80 at 10:05:00 UTC, stride 5 min → {10:00, 60, 65, 70, 2}, {10:05, 80, 80, 80, 1}
// Tags:    downsample, buckets, fixed-stride, resample, min-mean-max, time-series
func FixedStrideBuckets(points []Point, stride time.Duration) ([]StrideBucket, error) {
	if stride <= 0 {
		return nil, fmt.Errorf("%w: stride = %v, want > 0", ErrOutOfRange, stride)
	}
	if err := checkPoints(points); err != nil {
		return nil, err
	}
	count := 0
	var cur time.Time
	for i, p := range points {
		b := p.At.UTC().Truncate(stride)
		if i == 0 || !b.Equal(cur) {
			count++
		}
		cur = b
	}
	out := make([]StrideBucket, 0, count)
	first := 0
	for i := 1; i <= len(points); i++ {
		if i < len(points) && points[i].At.UTC().Truncate(stride).Equal(points[first].At.UTC().Truncate(stride)) {
			continue
		}
		out = append(out, strideBucket(points[first:i], points[first].At.UTC().Truncate(stride)))
		first = i
	}
	return out, nil
}

// strideBucket summarizes a non-empty run of finite points.
func strideBucket(run []Point, start time.Time) StrideBucket {
	b := StrideBucket{Start: start, Min: run[0].Value, Max: run[0].Value, Count: len(run)}
	b.Mean = meanOf(run)
	for _, p := range run[1:] {
		b.Min = math.Min(b.Min, p.Value)
		b.Max = math.Max(b.Max, p.Value)
	}
	return b
}

// meanOf returns the mean of a non-empty run of finite values, compensated; a sum beyond float64 is
// redone on the values scaled by 2⁻¹⁰, so the mean of finite values is always finite.
func meanOf(run []Point) float64 {
	var sum, comp float64
	for _, p := range run {
		sum, comp = reduce.NeumaierAdd(sum, comp, p.Value)
	}
	if reduce.AllFinite(sum) {
		return (sum + comp) / float64(len(run))
	}
	sum, comp = 0, 0
	for _, p := range run {
		sum, comp = reduce.NeumaierAdd(sum, comp, p.Value/1024)
	}
	return (sum + comp) / float64(len(run)) * 1024
}
