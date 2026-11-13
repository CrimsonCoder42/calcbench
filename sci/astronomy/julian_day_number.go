package astronomy

import (
	"errors"
	"fmt"
)

var (
	// ErrInvalidDate reports a month outside 1..12 or a day outside the month's length.
	ErrInvalidDate = errors.New("astronomy: invalid date")
	// ErrOutOfRange reports a year or Julian Day Number outside the supported range.
	ErrOutOfRange = errors.New("astronomy: out of range")
)

const (
	minYear = -4712
	maxYear = 1_000_000
	// minJDN is the JDN of −4712-01-01 Gregorian (proleptic), the first supported day.
	minJDN = 38
	// maxJDN is the JDN of 1000000-12-31 Gregorian.
	maxJDN = 366_963_925
)

// floorDiv returns ⌊a/b⌋ for b > 0.
func floorDiv(a, b int64) int64 {
	q := a / b
	if (a%b != 0) && (a < 0) {
		q--
	}
	return q
}

// isGregorianLeap reports the Gregorian leap-year rule: divisible by 4, except centuries not divisible by 400.
func isGregorianLeap(year int64) bool {
	return year%4 == 0 && (year%100 != 0 || year%400 == 0)
}

// daysInMonth returns the length of a Gregorian month (month 1..12).
func daysInMonth(year int64, month int) int {
	switch month {
	case 2:
		if isGregorianLeap(year) {
			return 29
		}
		return 28
	case 4, 6, 9, 11:
		return 30
	default:
		return 31
	}
}

func checkDate(year int64, month, day int) error {
	if year < minYear || year > maxYear {
		return fmt.Errorf("%w: year = %d, want %d..%d", ErrOutOfRange, year, minYear, maxYear)
	}
	if month < 1 || month > 12 {
		return fmt.Errorf("%w: month = %d, want 1..12", ErrInvalidDate, month)
	}
	if day < 1 || day > daysInMonth(year, month) {
		return fmt.Errorf("%w: day = %d for %d-%02d", ErrInvalidDate, day, year, month)
	}
	return nil
}

// GregorianToJulianDayNumber returns the Julian Day Number (the integer JD at noon) of a proleptic Gregorian
// date. It is Meeus's eq. 7.1, JD = INT(365.25(Y+4716)) + INT(30.6001(M+1)) + D + B − 1524.5 with
// A = INT(Y/100), B = 2 − A + INT(A/4) (January and February counted as months 13 and 14 of the year before),
// plus 0.5, done in exact integer arithmetic: INT(365.25k) = ⌊1461k/4⌋ and INT(30.6001m) = ⌊306001m/10000⌋.
// The Gregorian branch is used for every date (proleptic), so dates before 1582-10-15 are not the Julian
// calendar dates Meeus's text switches to.
// Errors: ErrOutOfRange for a year outside −4712..1,000,000; ErrInvalidDate for an impossible month or day.
//
// ID:      SCI-ASTRO-001
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Meeus J. Astronomical Algorithms, 2nd ed., Willmann-Bell, 1998, ch. 7 "Julian Day", eq. (7.1) and Example 7.a
// Tags:    julian-day, calendar, gregorian, date, day-number, astronomy
// Units:   year (astronomical numbering), month 1..12, day 1..31; result days (JDN)
// Range:   year −4712..1,000,000
// Example: 2000-01-01 → 2451545 (J2000.0 is JD 2451545.0, noon of that day); 1957-10-04 → 2436116 (Meeus Example 7.a, JD 2436116.31 at 4.81 d)
func GregorianToJulianDayNumber(year int64, month, day int) (int64, error) {
	if err := checkDate(year, month, day); err != nil {
		return 0, err
	}
	return jdn(year, month, day), nil
}

// jdn is eq. 7.1 on a validated date.
func jdn(year int64, month, day int) int64 {
	y, m := year, int64(month)
	if m <= 2 {
		y--
		m += 12
	}
	a := floorDiv(y, 100)
	b := 2 - a + floorDiv(a, 4)
	return floorDiv(1461*(y+4716), 4) + floorDiv(306001*(m+1), 10000) + int64(day) + b - 1524
}

// JulianDayNumberToGregorian returns the proleptic Gregorian date (year, month 1..12, day) whose noon is the
// Julian Date jdn.0. It is Meeus's ch. 7 inverse method with the Gregorian correction applied to every day:
// α = INT((Z − 1867216.25)/36524.25), A = Z + 1 + α − INT(α/4), B = A + 1524, C = INT((B − 122.1)/365.25),
// D = INT(365.25C), E = INT((B − D)/30.6001), day = B − D − INT(30.6001E), month = E − 1 (E < 14) or E − 13,
// year = C − 4716 (month > 2) or C − 4715; in exact integer arithmetic.
// Errors: ErrOutOfRange for jdn outside 38..366,963,925 (−4712-01-01 .. 1000000-12-31).
//
// ID:      SCI-ASTRO-002
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Meeus J. Astronomical Algorithms, 2nd ed., Willmann-Bell, 1998, ch. 7 "Julian Day", "Calculation of the Calendar Date from the JD" and Example 7.c
// Tags:    julian-day, calendar, gregorian, date, day-number, astronomy
// Units:   jdn days; result year (astronomical numbering), month 1..12, day 1..31
// Range:   jdn 38..366,963,925
// Example: 2436116 → 1957-10-04 (Meeus Example 7.c, JD 2436116.31 → 1957 October 4.81); 2451545 → 2000-01-01
func JulianDayNumberToGregorian(jdn int64) (year int64, month, day int, err error) {
	if jdn < minJDN || jdn > maxJDN {
		return 0, 0, 0, fmt.Errorf("%w: jdn = %d, want %d..%d", ErrOutOfRange, jdn, int64(minJDN), int64(maxJDN))
	}
	y, m, d := gregorianOf(jdn)
	return y, m, d, nil
}

// gregorianOf is the inverse method on a validated JDN.
func gregorianOf(z int64) (int64, int, int) {
	alpha := floorDiv(4*z-7468865, 146097) // INT((Z − 1867216.25)/36524.25)
	a := z + 1 + alpha - floorDiv(alpha, 4)
	b := a + 1524
	c := floorDiv(20*b-2442, 7305) // INT((B − 122.1)/365.25)
	d := floorDiv(1461*c, 4)       // INT(365.25C)
	e := floorDiv(10000*(b-d), 306001)
	day := b - d - floorDiv(306001*e, 10000)
	month := e - 1
	if e >= 14 {
		month = e - 13
	}
	year := c - 4715
	if month > 2 {
		year = c - 4716
	}
	return year, int(month), int(day)
}
