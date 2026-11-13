package cron

import (
	"fmt"
	"math/bits"
	"time"
)

// searchYears is the number of calendar years Next scans, counting the starting year. A February 29
// specification has its longest wait across a skipped century leap year (2096 → 2104, eight years), so
// nine calendar years always contain the next firing of any specification ParseCron accepts.
const searchYears = 9

// Next returns the earliest instant strictly after after at which c fires: the first instant whose wall
// clock in after.Location() has a matching minute, hour, month and day; its seconds and nanoseconds are zero
// and it is in after.Location().
// Daylight saving: a matching wall time that does not exist (spring-forward gap) fires at the first instant
// after the gap (the transition instant, e.g. 03:00 EDT); several matching times inside one gap fire once.
// A matching wall time that occurs twice (fall-back overlap) fires once, at its first (earlier-offset)
// occurrence; the second occurrence never fires.
// The search walks civil fields year → month → day → hour → minute with carries, each next allowed value
// taken from its bitset with one trailing-zero count; transitions are found with Time.ZoneBounds rather than
// trusting time.Date's unspecified choice inside a gap or an overlap.
// Day rule (common cron, doc/ASSUMPTIONS.md A-068): the month always restricts; day of month and day of
// week are ORed only when both are not literally `*`. This agrees with POSIX's example "0 0 1,15 * 1" but
// not with POSIX's literal wording when the month is also restricted.
// Errors: ErrNeverFires when c is the zero CronSpec or nothing fires within the nine calendar years starting
// with the year of after's wall clock (which never happens for a CronSpec returned by ParseCron).
//
// ID:      ALG-CRON-002
// Time:    O(1)   bounded: at most 9 years × 12 months × 31 days candidate days, each with one trailing-zero count per field and at most one time.Date and two ZoneBounds per candidate wall time
// Space:   O(1), no allocation on success
// Source:  IEEE Std 1003.1-2024 (POSIX.1-2024), Shell & Utilities, "crontab", INPUT FILES (field matching; day of month and day of week both restricted match when either matches), https://pubs.opengroup.org/onlinepubs/9799919799/utilities/crontab.html (read 2026-11-13)
// Source:  The Go Authors, package time, https://pkg.go.dev/time (func Date: in a gap or overlap the offset chosen is not guaranteed; func (Time) ZoneBounds: the instants the zone in effect begins and ends) (read 2026-11-13)
// Tags:    cron, crontab, schedule, next-run, time-zone, daylight-saving
// Units:   after is an instant in any location; the result is an instant in after.Location() with zero seconds and nanoseconds
// Example: "0 10 * * 1-5" after 2026-11-13T10:00Z (a Friday) → 2026-11-16T10:00Z; "30 2 * * *" in America/New_York after 2026-03-07T03:00-05:00 → 2026-03-08T03:00-04:00 (02:30 is in the gap)
func (c CronSpec) Next(after time.Time) (time.Time, error) {
	if c.minute == 0 || c.hour == 0 || c.month == 0 || (c.dom == 0 && c.dow == 0) {
		return time.Time{}, fmt.Errorf("%w: empty specification", ErrNeverFires)
	}
	loc := after.Location()
	y, mo, d := after.Date()
	h, mi, _ := after.Clock()
	mi++ // the first candidate is the next whole wall minute; carries are handled by the search

	// If after is inside the second pass of a fall-back overlap, every wall time up to the end of the overlap
	// already fired (first occurrence) before after: start at the end of the overlap.
	if start, _ := after.ZoneBounds(); !start.IsZero() {
		_, cur := after.Zone()
		_, prev := start.Add(-time.Nanosecond).Zone()
		if prev > cur {
			end := start.Add(time.Duration(prev-cur) * time.Second)
			if after.Before(end) {
				y, mo, d = end.Date()
				var s int
				h, mi, s = end.Clock()
				if s != 0 || end.Nanosecond() != 0 {
					mi++
				}
			}
		}
	}

	lastYear := y + searchYears - 1
	for y <= lastYear {
		m := nextBit(uint64(c.month), int(mo), 12)
		if m < 0 {
			y, mo, d, h, mi = y+1, 1, 1, 0, 0
			continue
		}
		if time.Month(m) != mo {
			mo, d, h, mi = time.Month(m), 1, 0, 0
		}
		dd := nextBit(uint64(c.dayMask(y, mo)), d, 31)
		if dd < 0 {
			mo, d, h, mi = mo+1, 1, 0, 0
			continue
		}
		if dd != d {
			d, h, mi = dd, 0, 0
		}
		hh := nextBit(uint64(c.hour), h, 23)
		if hh < 0 {
			d, h, mi = d+1, 0, 0
			continue
		}
		if hh != h {
			h, mi = hh, 0
		}
		mm := nextBit(c.minute, mi, 59)
		if mm < 0 {
			h, mi = h+1, 0
			continue
		}
		mi = mm
		if f := fireInstant(y, mo, d, h, mi, loc); f.After(after) {
			return f, nil
		}
		mi++ // first occurrence not after after (only possible around an overlap): keep searching
	}
	return time.Time{}, fmt.Errorf("%w: no match in the %d calendar years from %d", ErrNeverFires, searchYears, lastYear-searchYears+1)
}

// nextBit returns the smallest set bit of set at index ≥ from and ≤ top, or −1 when there is none.
func nextBit(set uint64, from, top int) int {
	if from > top {
		return -1
	}
	set &^= uint64(1)<<from - 1
	set &= uint64(1)<<(top+1) - 1
	if set == 0 {
		return -1
	}
	return bits.TrailingZeros64(set)
}

// dayMask returns the days of month mo of year y that c allows (bit d set for day d), month aside.
func (c CronSpec) dayMask(y int, mo time.Month) uint32 {
	n := daysIn(y, mo)
	valid := uint32(1)<<(n+1) - 2 // bits 1..n
	if c.domStar && c.dowStar {
		return valid
	}
	w1 := uint(time.Date(y, mo, 1, 0, 0, 0, 0, time.UTC).Weekday())
	dw := uint32(c.dow)
	week := (dw>>w1 | dw<<(7-w1)) & 0x7f // bit k: weekday of day k+1 allowed
	byDow := (week | week<<7 | week<<14 | week<<21 | week<<28) << 1
	switch {
	case c.domStar:
		return byDow & valid
	case c.dowStar:
		return c.dom & valid
	default:
		return (c.dom | byDow) & valid
	}
}

// daysIn returns the number of days in month mo of proleptic Gregorian year y.
func daysIn(y int, mo time.Month) int {
	if mo == time.February {
		if y%4 == 0 && (y%100 != 0 || y%400 == 0) {
			return 29
		}
		return 28
	}
	return maxDays[mo]
}

// fireInstant returns the instant the wall time y-mo-d h:mi:00 fires in loc: the instant itself when it
// exists once, its first (earlier-offset) occurrence in an overlap, the transition instant ending a gap.
func fireInstant(y int, mo time.Month, d, h, mi int, loc *time.Location) time.Time {
	t := time.Date(y, mo, d, h, mi, 0, 0, loc)
	if wallIs(t, y, mo, d, h, mi) {
		// Exists. An earlier occurrence lies in the previous zone period when that offset was larger.
		start, _ := t.ZoneBounds()
		if start.IsZero() {
			return t
		}
		_, cur := t.Zone()
		_, prev := start.Add(-time.Nanosecond).Zone()
		if prev > cur {
			if u := t.Add(time.Duration(cur-prev) * time.Second); wallIs(u, y, mo, d, h, mi) {
				return u
			}
		}
		return t
	}
	// In a gap: time.Date applied one of the two offsets. If its wall clock came out later than asked, t is
	// past the transition (the old offset was used) and the transition is t's zone start; otherwise t is
	// before it and the transition is t's zone end.
	start, end := t.ZoneBounds()
	_, off := t.Zone()
	asked := time.Date(y, mo, d, h, mi, 0, 0, time.UTC).Unix()
	if t.Unix()+int64(off) > asked {
		if !start.IsZero() {
			return start
		}
	} else if !end.IsZero() {
		return end
	}
	return t
}

// wallIs reports whether t's wall clock in its own location reads exactly y-mo-d h:mi:00.000000000.
func wallIs(t time.Time, y int, mo time.Month, d, h, mi int) bool {
	ty, tmo, td := t.Date()
	th, tmi, ts := t.Clock()
	return ty == y && tmo == mo && td == d && th == h && tmi == mi && ts == 0 && t.Nanosecond() == 0
}
