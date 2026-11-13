package series

import (
	"fmt"
	"time"

	"github.com/CrimsonCoder42/calcbench/algo/reduce"
)

// WeekSum is one local week: local midnight of its first day, and the sum and count of the values of the
// points in it.
type WeekSum struct {
	Start time.Time
	Sum   float64
	Count int
}

// WeeklySums sums time-sorted points by local week (weeks begin on weekStart, in loc) over the half-open
// window [from, to): one WeekSum for every week the window touches, in order, a week with no point
// included with Count 0. Only points inside [from, to) are summed. Week boundaries are civil midnights
// (NUM-SERIES-007), so a week with a DST change keeps its 167 or 169 hours. Sums are compensated.
// points is read, never modified.
// Errors: ErrNilLocation when loc is nil; ErrOutOfRange when weekStart is not Sunday … Saturday or when to
// is not after from; ErrUnsorted when a point's time is before the previous one's; ErrNotFinite when a
// value is NaN or ±Inf; ErrOverflow when a week's sum is beyond the float64 range.
// Implementation: Go's time package — func Date (normalizes day d + 7k into the right month; zone offset in effect at that civil time).
//
// ID:      NUM-SERIES-008
// Time:    O(n + w)   n = len(points), w = weeks touched by [from, to); one merge-like pass
// Space:   O(1) extra; O(w) for the result (one allocation)
// Source:  Reingold EM, Dershowitz N, Calendrical Calculations, 3rd ed., Cambridge University Press, 2008, ch. 1, "Days of the week" (k-day on or before)
// Source:  ISO 8601-1:2019, Date and time — Representations for information interchange — Part 1: Basic rules (calendar day; local time and its offset from UTC)
// Units:   values in any one unit; Sum in that unit; Count points
// Example: 2 on Mon 2026-10-05, 3 on Sat 10-10 23:30, 1 on Sun 10-11 00:30 New York, window [10-04, 10-25), Sunday weeks → 10-04 (5, 2), 10-11 (1, 1), 10-18 (0, 0)
// Tags:    calendar, week, weekly, buckets, sum, dst, time-zone, time-series
func WeeklySums(points []Point, from, to time.Time, loc *time.Location, weekStart time.Weekday) ([]WeekSum, error) {
	first, err := WeekStart(from, loc, weekStart)
	if err != nil {
		return nil, err
	}
	if !to.After(from) {
		return nil, fmt.Errorf("%w: to = %v, want after from = %v", ErrOutOfRange, to, from)
	}
	if err := checkPoints(points); err != nil {
		return nil, err
	}
	y, m, d := first.Date()
	weeks := 0
	for s := first; s.Before(to); s = time.Date(y, m, d+7*(weeks), 0, 0, 0, 0, loc) {
		weeks++
	}
	out := make([]WeekSum, weeks)
	for k := range out {
		out[k].Start = time.Date(y, m, d+7*k, 0, 0, 0, 0, loc)
	}
	k := 0
	var sum, comp float64
	flush := func() error {
		total := sum + comp
		if !reduce.AllFinite(sum) || !reduce.AllFinite(total) {
			return fmt.Errorf("%w: week starting %v", ErrOverflow, out[k].Start)
		}
		out[k].Sum = total
		sum, comp = 0, 0
		return nil
	}
	for _, p := range points {
		if p.At.Before(from) {
			continue
		}
		if !p.At.Before(to) {
			break
		}
		for k+1 < weeks && !p.At.Before(out[k+1].Start) {
			if err := flush(); err != nil {
				return nil, err
			}
			k++
		}
		sum, comp = reduce.NeumaierAdd(sum, comp, p.Value)
		out[k].Count++
	}
	if err := flush(); err != nil {
		return nil, err
	}
	return out, nil
}
