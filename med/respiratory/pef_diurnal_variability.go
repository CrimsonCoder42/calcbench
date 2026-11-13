package respiratory

import (
	"fmt"
	"math"

	"github.com/CrimsonCoder42/calcbench/algo/reduce"
)

// PEFAmplitudePercentMean returns one day's diurnal peak expiratory flow variability, in percent, as the
// daily amplitude percent mean of GINA: (day's highest − day's lowest) / mean of day's highest and lowest ×
// 100. GINA computes it from twice-daily readings, each the highest of 3 blows; which readings make up the
// day is the caller's choice (group a local calendar day with NUM-SERIES-005 LocalDayBuckets). Nothing is
// rounded.
// Errors: ErrNotFinite for NaN or ±Inf; ErrOutOfRange for dayLowestLPerMin < 0, dayHighestLPerMin ≤ 0 or
// dayHighestLPerMin < dayLowestLPerMin.
//
// ID:      MED-RESP-050
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Global Initiative for Asthma. Global Strategy for Asthma Management and Prevention, 2023 update. https://ginasthma.org/wp-content/uploads/2023/07/GINA-2023-Full-report-23_07_06-WMS.pdf (read 2026-11-12), Box 1-2 footnote * and section "How much variation in expiratory airflow is consistent with asthma?"
// Tags:    respiratory, asthma, peak-flow, pef, diurnal-variability, amplitude-percent-mean, gina
// Units:   dayHighestLPerMin, dayLowestLPerMin L/min (any one flow unit for both: only their ratio is used); result percent (0–200)
// Range:   0 ≤ dayLowestLPerMin ≤ dayHighestLPerMin; dayHighestLPerMin > 0
// Example: highest 440, lowest 400 L/min → 40 / 420 × 100 = 9.5238 % (HAND-DERIVED from the GINA formula)
func PEFAmplitudePercentMean(dayHighestLPerMin, dayLowestLPerMin float64) (float64, error) {
	if err := checkDayPEF(-1, dayHighestLPerMin, dayLowestLPerMin); err != nil {
		return 0, err
	}
	return amplitudePercentMean(dayHighestLPerMin, dayLowestLPerMin), nil
}

// PEFDiurnalVariabilityPercent returns the average daily diurnal peak expiratory flow variability, in
// percent: the mean over the days given of each day's amplitude percent mean (MED-RESP-050), as GINA
// averages it ("averaged over one week" in Box 1-2; "over 1–2 weeks" in the text). dayHighestLPerMin[i] and
// dayLowestLPerMin[i] are day i's highest and lowest readings; the function averages every day it is given
// and does not check how many days that is. Both slices are read, never modified. Nothing is rounded.
// Errors: ErrNotFinite for a NaN or ±Inf reading; ErrOutOfRange for no day, slices of different lengths, or
// a day outside the Range of MED-RESP-050.
//
// ID:      MED-RESP-051
// Time:    O(d)   d = number of days = len(dayHighestLPerMin)
// Space:   O(1), no allocation
// Source:  Global Initiative for Asthma. Global Strategy for Asthma Management and Prevention, 2023 update. https://ginasthma.org/wp-content/uploads/2023/07/GINA-2023-Full-report-23_07_06-WMS.pdf (read 2026-11-12), Box 1-2 footnote * and section "How much variation in expiratory airflow is consistent with asthma?"
// Tags:    respiratory, asthma, peak-flow, pef, diurnal-variability, amplitude-percent-mean, gina
// Units:   dayHighestLPerMin, dayLowestLPerMin L/min (any one flow unit for all: only ratios are used); result percent (0–200)
// Range:   at least one day; equal lengths; each day 0 ≤ lowest ≤ highest, highest > 0
// Example: days (440, 400) and (460, 380) L/min → (9.5238 + 19.0476) / 2 = 14.2857 % (HAND-DERIVED from the GINA formula)
func PEFDiurnalVariabilityPercent(dayHighestLPerMin, dayLowestLPerMin []float64) (float64, error) {
	d := len(dayHighestLPerMin)
	if d == 0 {
		return 0, fmt.Errorf("%w: dayHighestLPerMin is empty, want ≥ 1 day", ErrOutOfRange)
	}
	if len(dayLowestLPerMin) != d {
		return 0, fmt.Errorf("%w: len(dayLowestLPerMin) = %d, want len(dayHighestLPerMin) = %d", ErrOutOfRange, len(dayLowestLPerMin), d)
	}
	// Neumaier compensated sum of the daily values via ALG-RED-002 NeumaierAdd (rules/numeric.md).
	var sum, comp float64
	for i := 0; i < d; i++ {
		hi, lo := dayHighestLPerMin[i], dayLowestLPerMin[i]
		if err := checkDayPEF(i, hi, lo); err != nil {
			return 0, err
		}
		v := amplitudePercentMean(hi, lo)
		sum, comp = reduce.NeumaierAdd(sum, comp, v)
	}
	return (sum + comp) / float64(d), nil
}

// PEFDiurnalVariabilityExcessiveGINA reports whether an average daily diurnal PEF variability meets the
// GINA Box 1-2 criterion "Excessive variability in twice-daily PEF over 2 weeks": "Adults: average daily
// diurnal PEF variability >10%"; "Children: average daily diurnal PEF variability >13%". The box covers
// "adults, adolescents, and children 6–11 years"; ages 6 to < 12 take the children's > 13 %, ages ≥ 12 the
// adults' > 10 %, reading the box's "Children" as its title's "children 6–11 years".
// Both thresholds are strict: exactly 10 % (adults) or 13 % (children) is false.
// Errors: ErrNotFinite for NaN or ±Inf; ErrOutOfRange for variabilityPercent outside 0–200 or ageYears < 6.
//
// ID:      MED-RESP-052
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Global Initiative for Asthma. Global Strategy for Asthma Management and Prevention, 2023 update. https://ginasthma.org/wp-content/uploads/2023/07/GINA-2023-Full-report-23_07_06-WMS.pdf (read 2026-11-12), Box 1-2 "Diagnostic criteria for asthma in adults, adolescents, and children 6–11 years", row "Excessive variability in twice-daily PEF over 2 weeks"
// Tags:    respiratory, asthma, peak-flow, pef, diurnal-variability, gina, criterion
// Units:   variabilityPercent percent (MED-RESP-051's result); ageYears years (decimal); result true when the box's threshold is exceeded
// Range:   0 ≤ variabilityPercent ≤ 200 (the amplitude percent mean's arithmetic bounds); ageYears ≥ 6 (the box's youngest group)
// Example: 10.0 % at 40 years → false; 10.1 % at 40 years → true; 12.9 % at 8 years → false; 13.1 % at 8 years → true
func PEFDiurnalVariabilityExcessiveGINA(variabilityPercent, ageYears float64) (bool, error) {
	if err := checkClosedRange("variabilityPercent", variabilityPercent, 0, 200); err != nil {
		return false, err
	}
	if err := checkClosedRange("ageYears", ageYears, 6, math.Inf(1)); err != nil {
		return false, err
	}
	// GINA Box 1-2: children (6–11 years) > 13 %, adults and adolescents > 10 %.
	if ageYears < 12 {
		return variabilityPercent > 13, nil
	}
	return variabilityPercent > 10, nil
}

// amplitudePercentMean is (hi − lo) / ((hi + lo) / 2) × 100. The mean is formed as hi/2 + lo/2 so no finite
// input overflows; with 0 ≤ lo ≤ hi the quotient is at most 2.
func amplitudePercentMean(hi, lo float64) float64 {
	return (hi - lo) / (hi/2 + lo/2) * 100
}

// dayName names a parameter: the scalar name for day < 0, else name[day]. Called only on an error path.
func dayName(name string, day int) string {
	if day < 0 {
		return name
	}
	return fmt.Sprintf("%s[%d]", name, day)
}

// checkDayPEF checks one day's highest and lowest readings; day < 0 names the scalar parameters.
// It allocates only when it returns an error.
func checkDayPEF(day int, hi, lo float64) error {
	const hiName, loName = "dayHighestLPerMin", "dayLowestLPerMin"
	switch {
	case math.IsNaN(hi) || math.IsInf(hi, 0):
		return fmt.Errorf("%w: %s = %v", ErrNotFinite, dayName(hiName, day), hi)
	case math.IsNaN(lo) || math.IsInf(lo, 0):
		return fmt.Errorf("%w: %s = %v", ErrNotFinite, dayName(loName, day), lo)
	case lo < 0:
		return fmt.Errorf("%w: %s = %v, want ≥ 0", ErrOutOfRange, dayName(loName, day), lo)
	case !(hi > 0):
		return fmt.Errorf("%w: %s = %v, want > 0", ErrOutOfRange, dayName(hiName, day), hi)
	case hi < lo:
		return fmt.Errorf("%w: %s = %v < %s = %v, want highest ≥ lowest", ErrOutOfRange, dayName(hiName, day), hi, dayName(loName, day), lo)
	}
	return nil
}
