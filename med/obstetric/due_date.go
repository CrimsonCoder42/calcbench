package obstetric

import (
	"fmt"
	"time"
)

// lmpToEDDDays is from ACOG Committee Opinion No. 700 (2017), Background: "By convention, the EDD is 280 days after the
// first day of the LMP." Read in the committee opinion's text, 2026-11-13.
const lmpToEDDDays = 280

// civilDay returns the number of days from 1970-01-01 to t's calendar date (year, month, day as read in t's own
// location), so two dates in different locations compare by the date each shows.
func civilDay(t time.Time) int64 {
	y, m, d := t.Date()
	return time.Date(y, m, d, 0, 0, 0, 0, time.UTC).Unix() / 86400
}

// DueDateFromLMP returns the estimated due date (EDD) from the first day of the last menstrual period (LMP):
// the calendar date 280 days after it, by the ACOG convention. The result is midnight of that date in lmp's
// location (proleptic Gregorian calendar, Go's time package). The time of day of lmp is ignored. The LMP
// estimate does not account for cycle length or ovulation timing (the source says so); the caller chooses the
// method.
//
// ID:      MED-OB-003
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  American College of Obstetricians and Gynecologists, AIUM, SMFM. Committee Opinion No 700: Methods for Estimating the Due Date. Obstet Gynecol. 2017;129(5):e150-e154. doi:10.1097/AOG.0000000000002046. PMID 28426621. Background: "By convention, the EDD is 280 days after the first day of the LMP."
// Tags:    obstetric, pregnancy, due-date, edd, lmp, naegele, acog
// Units:   lmp a calendar date (time.Time; its date in its own location); result a calendar date (midnight, same location)
// Range:   any date
// Example: LMP 2026-01-01 → EDD 2026-10-08; LMP 2024-02-29 → EDD 2024-12-05 (HAND-DERIVED: 280 days on the Gregorian calendar)
func DueDateFromLMP(lmp time.Time) time.Time {
	y, m, d := lmp.Date()
	return time.Date(y, m, d+lmpToEDDDays, 0, 0, 0, 0, lmp.Location())
}

// GestationalAgeFromLMP returns the gestational age on a given date, counted from the first day of the last
// menstrual period, as completed weeks and days past the completed weeks (the w d/7 form of MED-OB-001): the
// whole calendar days from lmp's date to on's date, split into weeks and days. Day 280 is 40 0/7 weeks, the
// due date of DueDateFromLMP (MED-OB-003). Each date is read in its own location; times of day are ignored.
// Errors: ErrOutOfRange when on's date is before lmp's date.
//
// ID:      MED-OB-004
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  American College of Obstetricians and Gynecologists, AIUM, SMFM. Committee Opinion No 700: Methods for Estimating the Due Date. Obstet Gynecol. 2017;129(5):e150-e154. doi:10.1097/AOG.0000000000002046. PMID 28426621. Background: "By convention, the EDD is 280 days after the first day of the LMP" (gestational age counted from the LMP)
// Tags:    obstetric, pregnancy, gestational-age, lmp, acog
// Units:   lmp, on calendar dates (time.Time); weeks completed weeks; days 0–6 past the completed weeks
// Range:   on's date ≥ lmp's date (no upper limit is applied)
// Example: LMP 2026-01-01, on 2026-10-08 → 40 weeks 0 days; on 2026-03-15 → 73 days = 10 weeks 3 days (HAND-DERIVED)
func GestationalAgeFromLMP(lmp, on time.Time) (weeks, days int, err error) {
	n := civilDay(on) - civilDay(lmp)
	if n < 0 {
		return 0, 0, fmt.Errorf("%w: on (%s) is before lmp (%s)", ErrOutOfRange, on.Format(time.DateOnly), lmp.Format(time.DateOnly))
	}
	return int(n / 7), int(n % 7), nil
}
