package series

import (
	"fmt"
	"time"
)

// WeekStart returns local midnight in loc of the first day of t's week, for weeks that begin on
// weekStart: the latest date on or before t's local date whose weekday is weekStart. Calendar arithmetic is
// done on the civil date in loc, so the result is midnight even across a DST change.
// Errors: ErrNilLocation when loc is nil; ErrOutOfRange when weekStart is not time.Sunday … time.Saturday.
// Implementation: Go's time package — func Date and func (Time) Weekday (proleptic Gregorian calendar; IANA time-zone rules).
//
// ID:      NUM-SERIES-007
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Reingold EM, Dershowitz N, Calendrical Calculations, 3rd ed., Cambridge University Press, 2008, ch. 1, "Days of the week" (k-day on or before: date − ((date − k) mod 7))
// Source:  ISO 8601-1:2019, Date and time — Representations for information interchange — Part 1: Basic rules (calendar day; local time and its offset from UTC)
// Units:   t an instant; result an instant (local midnight in loc)
// Example: 2026-10-04 (a Sunday) with weekStart Monday → 2026-09-28 00:00 local; with weekStart Sunday → 2026-10-04 00:00 local
// Tags:    calendar, week, week-start, weekday, dst, time-zone, time-series
func WeekStart(t time.Time, loc *time.Location, weekStart time.Weekday) (time.Time, error) {
	if loc == nil {
		return time.Time{}, fmt.Errorf("%w: loc", ErrNilLocation)
	}
	if weekStart < time.Sunday || weekStart > time.Saturday {
		return time.Time{}, fmt.Errorf("%w: weekStart = %d, want 0 (Sunday) … 6 (Saturday)", ErrOutOfRange, int(weekStart))
	}
	return weekStartIn(t, loc, weekStart), nil
}

// weekStartIn is WeekStart for a valid loc and weekStart.
func weekStartIn(t time.Time, loc *time.Location, weekStart time.Weekday) time.Time {
	l := t.In(loc)
	back := (int(l.Weekday()) - int(weekStart) + 7) % 7
	y, m, d := l.Date()
	return time.Date(y, m, d-back, 0, 0, 0, 0, loc)
}
