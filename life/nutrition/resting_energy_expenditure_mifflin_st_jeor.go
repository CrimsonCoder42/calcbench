package nutrition

import (
	"errors"
	"fmt"

	"github.com/CrimsonCoder42/calcbench/algo/reduce"
)

var (
	// ErrNotFinite is returned when an input is NaN or ±Inf; it is checked before any range.
	ErrNotFinite = errors.New("nutrition: value is NaN or infinite")
	// ErrOutOfRange is returned when an input is outside the range the source defines or validated.
	ErrOutOfRange = errors.New("nutrition: input outside the range the source validated")
	// ErrDomain is returned when inputs valid one by one have no value in the source together (for example a
	// pregnancy row asked for a male) or give a result beyond float64.
	ErrDomain = errors.New("nutrition: inputs valid one by one but the source gives no value for them together")
)

// checkFinite returns ErrNotFinite, wrapped with name, when v is NaN or ±Inf.
func checkFinite(name string, v float64) error {
	if !reduce.AllFinite(v) {
		return fmt.Errorf("%w: %s = %v", ErrNotFinite, name, v)
	}
	return nil
}

// checkFinitePositive returns ErrNotFinite for NaN or ±Inf, then ErrOutOfRange unless v > 0.
func checkFinitePositive(name string, v float64) error {
	if err := checkFinite(name, v); err != nil {
		return err
	}
	if !(v > 0) {
		return fmt.Errorf("%w: %s = %v, want > 0", ErrOutOfRange, name, v)
	}
	return nil
}

// checkMifflinInputs applies the shared checks of the two Mifflin–St Jeor equations.
func checkMifflinInputs(weightKg, heightCm, ageYears float64) error {
	if err := checkFinitePositive("weightKg", weightKg); err != nil {
		return err
	}
	if err := checkFinitePositive("heightCm", heightCm); err != nil {
		return err
	}
	if err := checkFinite("ageYears", ageYears); err != nil {
		return err
	}
	// "aged 19-78 y": completed ages 19 through 78, so 19 ≤ age < 79 in decimal years.
	if !(ageYears >= 19 && ageYears < 79) {
		return fmt.Errorf("%w: ageYears = %v, want 19 ≤ age < 79", ErrOutOfRange, ageYears)
	}
	return nil
}

// RestingEnergyExpenditureMifflinStJeor returns the predicted resting energy expenditure (REE), in kcal/day,
// by the Mifflin–St Jeor 1990 equation derived for both sexes:
// REE = 9.99 × weight + 6.25 × height − 4.92 × age + 166 × sex − 161, sex = 1 for males and 0 for females.
// Not rounded. No activity factor is applied (the source gives none).
// Errors: ErrNotFinite for NaN or ±Inf; ErrOutOfRange outside Range; ErrDomain when the result is beyond float64.
//
// ID:      LIF-NUTR-001
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Mifflin MD, St Jeor ST, Hill LA, Scott BJ, Daugherty SA, Koh YO. A new predictive equation for resting energy expenditure in healthy individuals. Am J Clin Nutr. 1990;51(2):241-247. doi:10.1093/ajcn/51.2.241. PMID 2305711, abstract: "REE = 9.99 x weight + 6.25 x height - 4.92 x age + 166 x sex (males, 1; females, 0) - 161"; "aged 19-78 y"
// Tags:    nutrition, energy, resting-energy-expenditure, ree, bmr, mifflin-st-jeor, metric
// Units:   weightKg kg; heightCm cm; ageYears years (the abstract's simplified forms state kg, cm, y); female true for the source's sex 0; result kcal/day (the abstract prints no REE unit; kcal/day is the only unit the magnitude fits, e.g. ≈ 1650 for a 70 kg man, and Eklo RR et al., Am J Clin Nutr 2026;123(3):101209, PMC12975359, Table 1 "Predictive equations utilized to estimate REE and BMR (kcal/d)" reproduces this equation in kcal/d)
// Range:   ageYears 19 ≤ age < 79 ("aged 19-78 y"); weightKg > 0 and heightCm > 0 (the abstract states no weight or height limits)
// Example: man, 70 kg, 175 cm, 30 y → 699.3 + 1093.75 − 147.6 + 166 − 161 = 1650.45 (HAND-DERIVED from the abstract's equation)
func RestingEnergyExpenditureMifflinStJeor(weightKg, heightCm, ageYears float64, female bool) (float64, error) {
	if err := checkMifflinInputs(weightKg, heightCm, ageYears); err != nil {
		return 0, err
	}
	sex := 1.0
	if female {
		sex = 0
	}
	r := 9.99*weightKg + 6.25*heightCm - 4.92*ageYears + 166*sex - 161
	if !reduce.AllFinite(r) {
		return 0, fmt.Errorf("%w: result is beyond float64", ErrDomain)
	}
	return r, nil
}

// RestingEnergyExpenditureMifflinStJeorSimplified returns the predicted resting energy expenditure (REE), in
// kcal/day, by the simplified sex-specific Mifflin–St Jeor 1990 equations: males
// 10 × weight + 6.25 × height − 5 × age + 5; females 10 × weight + 6.25 × height − 5 × age − 161. Not rounded.
// The abstract states the simplification "did not affect its predictive value"; it differs from
// RestingEnergyExpenditureMifflinStJeor (LIF-NUTR-001) by a few kcal/day.
// Errors: ErrNotFinite for NaN or ±Inf; ErrOutOfRange outside Range; ErrDomain when the result is beyond float64.
//
// ID:      LIF-NUTR-002
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Mifflin MD, St Jeor ST, Hill LA, Scott BJ, Daugherty SA, Koh YO. A new predictive equation for resting energy expenditure in healthy individuals. Am J Clin Nutr. 1990;51(2):241-247. doi:10.1093/ajcn/51.2.241. PMID 2305711, abstract: "REE (males) = 10 x weight (kg) + 6.25 x height (cm) - 5 x age (y) + 5"; "REE (females) = 10 x weight (kg) + 6.25 x height (cm) - 5 x age (y) - 161"
// Tags:    nutrition, energy, resting-energy-expenditure, ree, bmr, mifflin-st-jeor, metric
// Units:   weightKg kg; heightCm cm; ageYears years; result kcal/day (the abstract prints no REE unit; see LIF-NUTR-001 for the corroborating source)
// Range:   ageYears 19 ≤ age < 79 ("aged 19-78 y"); weightKg > 0 and heightCm > 0 (no limits stated)
// Example: man, 70 kg, 175 cm, 30 y → 700 + 1093.75 − 150 + 5 = 1648.75; woman, 60 kg, 165 cm, 40 y → 600 + 1031.25 − 200 − 161 = 1270.25 (HAND-DERIVED)
func RestingEnergyExpenditureMifflinStJeorSimplified(weightKg, heightCm, ageYears float64, female bool) (float64, error) {
	if err := checkMifflinInputs(weightKg, heightCm, ageYears); err != nil {
		return 0, err
	}
	c := 5.0
	if female {
		c = -161
	}
	r := 10*weightKg + 6.25*heightCm - 5*ageYears + c
	if !reduce.AllFinite(r) {
		return 0, fmt.Errorf("%w: result is beyond float64", ErrDomain)
	}
	return r, nil
}
