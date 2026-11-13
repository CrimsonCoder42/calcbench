package series

import (
	"fmt"
	"time"
)

// DayBucket is one local calendar day that has points: its date in the location, and the half-open index
// range [First, End) of its points in the input slice.
type DayBucket struct {
	Year  int
	Month time.Month
	Day   int
	First int
	End   int
}

// LocalDayBuckets groups time-sorted points by their calendar date in loc and returns one bucket per date
// that has a point, in time order; a date with no point has no bucket. A point at 23:30 local stays on
// its local date whatever its UTC date, and a DST day keeps its 23 or 25 hours. Buckets index into
// points instead of copying them. points is read, never modified.
// Errors: ErrNilLocation when loc is nil; ErrUnsorted when a point's time is before the previous one's.
// Implementation: Go's time package — func (Time) In and func (Time) Date (the year, month and day of the instant in a location; proleptic Gregorian calendar, IANA time-zone rules).
//
// ID:      NUM-SERIES-005
// Time:    O(n)   n = len(points); two passes (count the dates, then fill)
// Space:   O(1) extra; O(d) for the result (one allocation), d = number of dates with a point
// Source:  Reingold EM, Dershowitz N, Calendrical Calculations, 3rd ed., Cambridge University Press, 2008, ch. 1 (fixed day numbers, floor and mod) and ch. 2 (the proleptic Gregorian calendar)
// Source:  ISO 8601-1:2019, Date and time — Representations for information interchange — Part 1: Basic rules (calendar day; local time and its offset from UTC)
// Units:   points' times instants; values are not read
// Example: 9 instants from 2026-11-01 02:00 UTC to 2026-11-02 15:00 UTC in America/New_York → 2026-10-31: 1, 2026-11-01: 7, 2026-11-02: 1
// Tags:    calendar, local-day, group-by-day, buckets, dst, time-zone, time-series
func LocalDayBuckets(points []Point, loc *time.Location) ([]DayBucket, error) {
	if loc == nil {
		return nil, fmt.Errorf("%w: loc", ErrNilLocation)
	}
	for i := 1; i < len(points); i++ {
		if points[i].At.Before(points[i-1].At) {
			return nil, fmt.Errorf("%w: points[%d].At = %v is before points[%d].At = %v", ErrUnsorted, i, points[i].At, i-1, points[i-1].At)
		}
	}
	days := 0
	var py, pd int
	var pm time.Month
	for i, p := range points {
		y, m, d := p.At.In(loc).Date()
		if i == 0 || y != py || m != pm || d != pd {
			days++
		}
		py, pm, pd = y, m, d
	}
	out := make([]DayBucket, 0, days)
	for i, p := range points {
		y, m, d := p.At.In(loc).Date()
		if n := len(out); n > 0 && out[n-1].Year == y && out[n-1].Month == m && out[n-1].Day == d {
			out[n-1].End = i + 1
			continue
		}
		out = append(out, DayBucket{Year: y, Month: m, Day: d, First: i, End: i + 1})
	}
	return out, nil
}
