package scales

import "fmt"

// EpworthSleepinessCategory is an Epworth Sleepiness Scale band; its value is the official site's own label,
// verbatim. (The type is not named EpworthCategory because the function MED-SCALE-002 holds that name.)
type EpworthSleepinessCategory string

// The Epworth bands, labels verbatim from Johns MW, About the ESS, https://epworthsleepinessscale.com/about-the-ess/
// (read 2026-11-12).
const (
	EpworthLowerNormal       EpworthSleepinessCategory = "Lower Normal Daytime Sleepiness"       // 0–5
	EpworthHigherNormal      EpworthSleepinessCategory = "Higher Normal Daytime Sleepiness"      // 6–10
	EpworthMildExcessive     EpworthSleepinessCategory = "Mild Excessive Daytime Sleepiness"     // 11–12
	EpworthModerateExcessive EpworthSleepinessCategory = "Moderate Excessive Daytime Sleepiness" // 13–15
	EpworthSevereExcessive   EpworthSleepinessCategory = "Severe Excessive Daytime Sleepiness"   // 16–24
)

// EpworthCategory returns the Epworth Sleepiness Scale band for a total, with the official site's own label:
// 0–5 "Lower Normal Daytime Sleepiness"; 6–10 "Higher Normal Daytime Sleepiness"; 11–12 "Mild Excessive
// Daytime Sleepiness"; 13–15 "Moderate Excessive Daytime Sleepiness"; 16–24 "Severe Excessive Daytime
// Sleepiness".
// Errors: ErrOutOfRange when score is outside 0–24.
//
// ID:      MED-SCALE-002
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Johns MW. About the ESS. The Epworth Sleepiness Scale official website. https://epworthsleepinessscale.com/about-the-ess/ (read 2026-11-12)
// Tags:    scales, epworth, ess, sleep, sleepiness, category
// Units:   score points 0–24; result the official site's band label
// Range:   0 ≤ score ≤ 24
// Example: 5 → "Lower Normal Daytime Sleepiness"; 6 → "Higher Normal Daytime Sleepiness"; 16 → "Severe Excessive Daytime Sleepiness"
func EpworthCategory(score int) (EpworthSleepinessCategory, error) {
	if score < 0 || score > 24 {
		return "", fmt.Errorf("%w: score = %d, want 0–24", ErrOutOfRange, score)
	}
	switch {
	case score <= 5:
		return EpworthLowerNormal, nil
	case score <= 10:
		return EpworthHigherNormal, nil
	case score <= 12:
		return EpworthMildExcessive, nil
	case score <= 15:
		return EpworthModerateExcessive, nil
	}
	return EpworthSevereExcessive, nil
}
