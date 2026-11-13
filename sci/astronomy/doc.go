// Package astronomy holds the calendar arithmetic of astronomy: the Julian Day Number of a proleptic Gregorian
// date and back, the day of the week, the day of the year, and day counts between dates, by the methods of
// Meeus, Astronomical Algorithms, 2nd ed., 1998, ch. 7 "Julian Day".
//
// ID short form: ASTRO (IDs are SCI-ASTRO-NNN).
//
// Units: dates are (year, month 1..12, day 1..31) in the proleptic Gregorian calendar, astronomical year
// numbering (year 0 = 1 BC). A Julian Day Number (JDN) is an integer count of days; JDN n is the day whose noon
// is Julian Date n.0. Dates are restricted to years −4712..1,000,000 (JDN ≥ 0); outside it, ErrOutOfRange.
// An impossible date (month 13, 30 February) returns ErrInvalidDate. Nothing reads the clock.
//
// Functions:
//
//	SCI-ASTRO-001  GregorianToJulianDayNumber  JDN of a Gregorian date (Meeus eq. 7.1, integer form)
//	SCI-ASTRO-002  JulianDayNumberToGregorian  Gregorian date of a JDN (Meeus ch. 7, inverse method)
//	SCI-ASTRO-003  GregorianDayOfWeek          weekday of a Gregorian date ((JD + 1.5) mod 7)
//	SCI-ASTRO-004  GregorianDayOfYear          ordinal day 1..366 of a Gregorian date
//	SCI-ASTRO-005  DaysBetweenDates            signed day count from one Gregorian date to another
//	SCI-ASTRO-006  GregorianAddDays            the Gregorian date a signed number of days after another
package astronomy
