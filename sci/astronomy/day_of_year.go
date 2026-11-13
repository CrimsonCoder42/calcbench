package astronomy

// GregorianDayOfYear returns the ordinal day of a Gregorian date, 1 on 1 January to 365 or 366 on
// 31 December: N = INT(275M/9) − K·INT((M + 9)/12) + D − 30, K = 1 in a leap year and 2 otherwise.
// Errors: as GregorianToJulianDayNumber.
//
// ID:      SCI-ASTRO-004
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Meeus J. Astronomical Algorithms, 2nd ed., Willmann-Bell, 1998, ch. 7 "Julian Day", "Day of the Year" and Example 7.f
// Tags:    day-of-year, ordinal-date, calendar, gregorian, date, astronomy
// Units:   year, month 1..12, day 1..31; result day 1..366
// Range:   year −4712..1,000,000
// Example: 1978-11-14 → 318; 1988-04-22 → 113 (Meeus Example 7.f)
func GregorianDayOfYear(year int64, month, day int) (int, error) {
	if err := checkDate(year, month, day); err != nil {
		return 0, err
	}
	k := 2
	if isGregorianLeap(year) {
		k = 1
	}
	return 275*month/9 - k*((month+9)/12) + day - 30, nil
}
