package astronomy

// DaysBetweenDates returns the signed number of days from the first proleptic Gregorian date to the second:
// JDN(to) − JDN(from). It is negative when to is earlier; the end date is not counted (add 1 to include it).
// Errors: as GregorianToJulianDayNumber, for either date.
//
// ID:      SCI-ASTRO-005
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Meeus J. Astronomical Algorithms, 2nd ed., Willmann-Bell, 1998, ch. 7 "Julian Day", "the interval in days between two dates is the difference of their JDs" and Example 7.d
// Tags:    days-between, day-count, date-difference, calendar, gregorian, astronomy
// Units:   two dates (year, month 1..12, day 1..31); result days
// Range:   years −4712..1,000,000
// Example: 1910-04-20 → 1986-02-09 = 27689 days (Meeus Example 7.d, Halley's Comet perihelia)
func DaysBetweenDates(fromYear int64, fromMonth, fromDay int, toYear int64, toMonth, toDay int) (int64, error) {
	a, err := GregorianToJulianDayNumber(fromYear, fromMonth, fromDay)
	if err != nil {
		return 0, err
	}
	b, err := GregorianToJulianDayNumber(toYear, toMonth, toDay)
	if err != nil {
		return 0, err
	}
	return b - a, nil
}
