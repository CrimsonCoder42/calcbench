package pediatrics

import (
	"fmt"
	"math"
)

// RecommendedSleepHoursAASM returns the American Academy of Sleep Medicine's recommended sleep per 24 hours for
// a child's age, as the band's minimum and maximum in hours: 4–12 months "12 to 16 hours per 24 hours
// (including naps)"; 1–2 years "11 to 14 hours" (including naps); 3–5 years "10 to 13 hours" (including
// naps); 6–12 years "9 to 12 hours"; 13–18 years "8 to 10 hours". The statement excludes infants younger than
// 4 months. Ages are read as completed months: 12 months is "1–2 years", 36 months "3–5 years", 72 months
// "6–12 years", 156 months "13–18 years", and 13–18 years runs until the 19th birthday (228 months).
// Errors: ErrNotFinite for NaN or ±Inf; ErrOutOfRange outside Range.
//
// ID:      MED-PEDS-002
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Paruthi S, Brooks LJ, D'Ambrosio C, et al. Recommended amount of sleep for pediatric populations: a consensus statement of the American Academy of Sleep Medicine. J Clin Sleep Med. 2016;12(6):785-786. doi:10.5664/jcsm.5866, PMID 27250809, recommendations
// Tags:    pediatrics, sleep, aasm, reference range
// Units:   ageMonths months (decimal allowed); results hours per 24 h
// Range:   4 ≤ ageMonths < 228 (4 months to the end of age 18, the statement's population)
// Example: 6 months → 12–16 h; 24 months → 11–14 h; 48 → 10–13; 100 → 9–12; 200 → 8–10 (Paruthi 2016)
func RecommendedSleepHoursAASM(ageMonths float64) (minHours, maxHours float64, err error) {
	if math.IsNaN(ageMonths) || math.IsInf(ageMonths, 0) {
		return 0, 0, fmt.Errorf("%w: ageMonths = %v", ErrNotFinite, ageMonths)
	}
	if ageMonths < 4 || ageMonths >= 228 {
		return 0, 0, fmt.Errorf("%w: ageMonths = %v, want 4 ≤ ageMonths < 228", ErrOutOfRange, ageMonths)
	}
	switch {
	case ageMonths < 12:
		return 12, 16, nil
	case ageMonths < 36:
		return 11, 14, nil
	case ageMonths < 72:
		return 10, 13, nil
	case ageMonths < 156:
		return 9, 12, nil
	}
	return 8, 10, nil
}
