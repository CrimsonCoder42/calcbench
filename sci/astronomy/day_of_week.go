package astronomy

import "time"

// GregorianDayOfWeek returns the day of the week of a proleptic Gregorian date. Meeus: the weekday is the
// remainder of (JD + 1.5) / 7, 0 = Sunday, with JD at 0h; with the noon-based JDN that is (JDN + 1) mod 7.
// The result is a time.Weekday (time.Sunday = 0, the same numbering).
// Errors: as GregorianToJulianDayNumber.
//
// ID:      SCI-ASTRO-003
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Meeus J. Astronomical Algorithms, 2nd ed., Willmann-Bell, 1998, ch. 7 "Julian Day", "Day of the Week" and Example 7.e
// Tags:    day-of-week, weekday, calendar, gregorian, date, astronomy
// Units:   year, month 1..12, day 1..31; result weekday (0 = Sunday)
// Range:   year −4712..1,000,000
// Example: 1954-06-30 → Wednesday (Meeus Example 7.e, JD 2434923.5)
func GregorianDayOfWeek(year int64, month, day int) (time.Weekday, error) {
	n, err := GregorianToJulianDayNumber(year, month, day)
	if err != nil {
		return 0, err
	}
	return time.Weekday((n + 1) % 7), nil
}
