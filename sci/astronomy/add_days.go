package astronomy

import "fmt"

// GregorianAddDays returns the proleptic Gregorian date that is days after the given date (before it when days
// is negative): the date of JDN(date) + days, by eq. 7.1 and its inverse.
// Errors: as GregorianToJulianDayNumber for the start date; ErrOutOfRange when the result leaves the supported
// range (also when days is so large the sum would overflow).
//
// ID:      SCI-ASTRO-006
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Meeus J. Astronomical Algorithms, 2nd ed., Willmann-Bell, 1998, ch. 7 "Julian Day", eq. (7.1) and the inverse method (Examples 7.a, 7.c)
// Tags:    add-days, date-arithmetic, calendar, gregorian, date, astronomy
// Units:   year, month 1..12, day 1..31; days (signed); result year, month, day
// Range:   start and result within years −4712..1,000,000
// Example: 1910-04-20 + 27689 days → 1986-02-09 (Meeus Example 7.d read backwards)
func GregorianAddDays(year int64, month, day int, days int64) (int64, int, int, error) {
	n, err := GregorianToJulianDayNumber(year, month, day)
	if err != nil {
		return 0, 0, 0, err
	}
	if days > maxJDN || days < -maxJDN {
		return 0, 0, 0, fmt.Errorf("%w: days = %d", ErrOutOfRange, days)
	}
	return JulianDayNumberToGregorian(n + days)
}
