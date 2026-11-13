package interest

import (
	"fmt"
	"time"
)

// checkDateOrder returns ErrDomain when end is on an earlier calendar date than start, each date read
// in its own time.Location.
func checkDateOrder(start, end time.Time) error {
	y1, m1, d1 := start.Date()
	y2, m2, d2 := end.Date()
	if y2 < y1 || (y2 == y1 && (m2 < m1 || (m2 == m1 && d2 < d1))) {
		return fmt.Errorf("%w: end %s is before start %s", ErrDomain, end.Format(time.DateOnly), start.Format(time.DateOnly))
	}
	return nil
}

// YearFraction30360 returns the fraction of a year from start to end on the 30/360 basis of MSRB Rule
// G-33(e): Number of Days = (Y2 − Y1)·360 + (M2 − M1)·30 + (D2 − D1), where D2 = 31 becomes 30 when
// D1 is 30 or 31, and then D1 = 31 becomes 30; the result is Number of Days/360. The start day is
// counted and the end day is not. Each date is read in its own time.Location and the time of day is
// ignored. The rule has no end-of-February adjustment, and none is applied.
// Errors: ErrDomain if end is on an earlier calendar date than start.
//
// ID:      FIN-INT-036
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  MSRB Rule G-33(e) "Day Counting", https://www.msrb.org/Rules-and-Interpretations/MSRB-Rules/General/Rule-G-33, read 2026-11-12: Number of Days = (Y2 − Y1)360 + (M2 − M1)30 + (D2 − D1), with the D1/D2 = 31 rules
// Tags:    interest, day-count, 30-360, year-fraction, accrued-interest, municipal-securities
// Units:   start and end calendar dates; result years (360-day basis)
// Example: hand-checked: 2026-01-31 to 2026-03-31 → D1 31→30, D2 31→30, 60 days → 0.16667
func YearFraction30360(start, end time.Time) (float64, error) {
	if err := checkDateOrder(start, end); err != nil {
		return 0, err
	}
	y1, m1, d1 := start.Date()
	y2, m2, d2 := end.Date()
	if d2 == 31 && (d1 == 30 || d1 == 31) {
		d2 = 30
	}
	if d1 == 31 {
		d1 = 30
	}
	days := (y2-y1)*360 + (int(m2)-int(m1))*30 + (d2 - d1)
	return float64(days) / 360, nil
}
