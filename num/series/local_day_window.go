package series

import (
	"fmt"
	"time"
)

// LocalDayWindow returns the half-open window [from, to) covering the days local calendar days in loc
// that end with asOf's local day: from is local midnight of (asOf's date − days + 1), to is local midnight
// of the day after asOf's date. Calendar days are counted in loc, so a day with a DST change is 23 or 25
// hours long, and adjacent windows share an edge without double-counting a point.
// Errors: ErrNilLocation when loc is nil; ErrOutOfRange when days < 1.
//
// ID:      NUM-SERIES-004
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  The Go standard library, package time, func Date ("normalizes" out-of-range day values; proleptic Gregorian calendar; the zone offset is the one in effect in loc at that civil time), https://pkg.go.dev/time#Date (read 2026-11-12)
// Units:   asOf an instant; days calendar days; from and to instants in loc
// Example: asOf 2026-11-01 12:00 America/New_York, days 1 → [2026-11-01 00:00 EDT, 2026-11-02 00:00 EST), 25 h (the US DST end)
// Tags:    calendar, local-day, window, half-open, dst, time-zone, time-series
func LocalDayWindow(asOf time.Time, loc *time.Location, days int) (from, to time.Time, err error) {
	if loc == nil {
		return time.Time{}, time.Time{}, fmt.Errorf("%w: loc", ErrNilLocation)
	}
	if days < 1 {
		return time.Time{}, time.Time{}, fmt.Errorf("%w: days = %d, want ≥ 1", ErrOutOfRange, days)
	}
	y, m, d := asOf.In(loc).Date()
	return time.Date(y, m, d-days+1, 0, 0, 0, 0, loc), time.Date(y, m, d+1, 0, 0, 0, 0, loc), nil
}
