package interest

import "time"

// civilDay returns the calendar date of t, read in t's own Location, as a day number: days from
// 1970-01-01 to that date at midnight UTC. UTC has no daylight-saving shifts, so the difference of two
// day numbers is the exact count of calendar days.
func civilDay(t time.Time) int64 {
	y, m, d := t.Date()
	return time.Date(y, m, d, 0, 0, 0, 0, time.UTC).Unix() / 86400
}

// YearFractionActual365 returns the fraction of a year from start to end on the Actual/365 basis: the
// actual number of calendar days between the dates (start counted, end not) divided by 365, in every
// year, leap or not. Each date is read in its own time.Location and the time of day is ignored, so a
// daylight-saving change never shifts the count.
// Errors: ErrDomain if end is on an earlier calendar date than start.
//
// ID:      FIN-INT-037
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  D. Kirk, Contemporary Mathematics, OpenStax, 2023, §6.3 "Simple Interest", https://openstax.org/books/contemporary-mathematics/pages/6-3-simple-interest: interest on a 365-day year with actual days ("Actual/365")
// Source:  12 CFR 1030 (Regulation DD) Appendix A Part I.A, read at https://www.consumerfinance.gov/rules-policy/regulations/1030/a/ (CFPB eRegulations; ecfr.gov blocked automated reading) 2026-11-12 (actual days in term over 365)
// Tags:    interest, day-count, actual-365, year-fraction, accrued-interest
// Units:   start and end calendar dates; result years (365-day basis)
// Example: 45 days → 45/365 (Contemporary Mathematics Example 6.30, $5,600 at 11.75 % → interest $81.12 half-even; the book rounds up to $81.13)
func YearFractionActual365(start, end time.Time) (float64, error) {
	if err := checkDateOrder(start, end); err != nil {
		return 0, err
	}
	days := civilDay(end) - civilDay(start)
	return float64(days) / 365, nil
}
