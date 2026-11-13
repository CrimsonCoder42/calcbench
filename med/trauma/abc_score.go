package trauma

import (
	"errors"
	"fmt"

	"github.com/CrimsonCoder42/calcbench/algo/reduce"
)

var (
	// ErrNotFinite is returned when a float input is NaN or ±Inf; it is checked before any range.
	ErrNotFinite = errors.New("trauma: value is NaN or infinite")
	// ErrOutOfRange is returned when an input is outside its physical range or the range the source validated.
	ErrOutOfRange = errors.New("trauma: input outside the range the source validated")
)

// finite returns ErrNotFinite, wrapped with name, when v is NaN or ±Inf.
func finite(name string, v float64) error {
	if !reduce.AllFinite(v) {
		return fmt.Errorf("%w: %s = %v", ErrNotFinite, name, v)
	}
	return nil
}

// within returns ErrOutOfRange, wrapped with name, unless lo ≤ v ≤ hi. NaN passes; check finite first.
func within(name string, v, lo, hi float64) error {
	if v < lo || v > hi {
		return fmt.Errorf("%w: %s = %v, want %v–%v", ErrOutOfRange, name, v, lo, hi)
	}
	return nil
}

// atLeastZero returns ErrOutOfRange, wrapped with name, when v < 0. NaN passes; check finite first.
func atLeastZero(name string, v float64) error {
	if v < 0 {
		return fmt.Errorf("%w: %s = %v, want ≥ 0", ErrOutOfRange, name, v)
	}
	return nil
}

// checkAge returns ErrNotFinite or ErrOutOfRange unless 0 ≤ ageYears ≤ 130 (physical limits).
func checkAge(ageYears float64) error {
	if err := finite("ageYears", ageYears); err != nil {
		return err
	}
	return within("ageYears", ageYears, 0, 130)
}

// checkGCS returns ErrOutOfRange unless 3 ≤ gcs ≤ 15 (the Glasgow Coma Scale total).
func checkGCS(gcs int) error {
	if gcs < 3 || gcs > 15 {
		return fmt.Errorf("%w: gcs = %d, want 3–15", ErrOutOfRange, gcs)
	}
	return nil
}

// count returns the number of true values.
func count(flags ...bool) int {
	n := 0
	for _, f := range flags {
		if f {
			n++
		}
	}
	return n
}

// ABCScore returns the Assessment of Blood Consumption (ABC) score, 0–4: one point each, "nonweighted", for a
// penetrating mechanism, a positive focused assessment sonography for trauma (FAST), an arrival systolic blood
// pressure of 90 mmHg or less and an arrival heart rate of 120 bpm or more. The abstract reports "ABC score of 2
// or greater was 75% sensitive and 86% specific for predicting MT" (massive transfusion).
// Errors: ErrNotFinite for NaN or ±Inf; ErrOutOfRange for a negative pressure or heart rate.
//
// ID:      MED-TRAUMA-001
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Nunez TC, Voskresensky IV, Dossett LA, Shinall R, Dutton WD, Cotton BA. Early prediction of massive transfusion in trauma: simple as ABC (assessment of blood consumption)? J Trauma. 2009;66(2):346-352. doi:10.1097/TA.0b013e3181961c35. PMID 19204506. Abstract: "ABC score was assigned based on four nonweighted parameters: penetrating mechanism, positive focused assessment sonography for trauma, arrival systolic blood pressure of 90 mm Hg or less, and arrival heart rate ≥120 bpm"
// Tags:    trauma, massive transfusion, abc, score
// Units:   systolicMmHg mmHg; heartRatePerMin beats/min; result points 0–4
// Range:   systolicMmHg ≥ 0; heartRatePerMin ≥ 0 (physical limits; the abstract states no validated range)
// Example: penetrating, FAST negative, SBP 85, HR 125 → 3 (HAND-DERIVED from the abstract's four parameters)
func ABCScore(penetratingMechanism, positiveFAST bool, systolicMmHg, heartRatePerMin float64) (int, error) {
	for _, e := range [...]error{
		finite("systolicMmHg", systolicMmHg), finite("heartRatePerMin", heartRatePerMin),
		atLeastZero("systolicMmHg", systolicMmHg), atLeastZero("heartRatePerMin", heartRatePerMin),
	} {
		if e != nil {
			return 0, e
		}
	}
	return count(penetratingMechanism, positiveFAST, systolicMmHg <= 90, heartRatePerMin >= 120), nil
}
