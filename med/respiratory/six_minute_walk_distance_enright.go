package respiratory

import (
	"fmt"
	"math"

	"github.com/CrimsonCoder42/calcbench/algo/reduce"
)

// checkClosedRange returns ErrNotFinite for NaN or ±Inf and ErrOutOfRange unless lo ≤ v ≤ hi.
// hi may be +Inf for "no upper limit".
func checkClosedRange(name string, v, lo, hi float64) error {
	if !reduce.AllFinite(v) {
		return fmt.Errorf("%w: %s = %v", ErrNotFinite, name, v)
	}
	if v < lo || v > hi {
		return fmt.Errorf("%w: %s = %v, want %v ≤ %s ≤ %v", ErrOutOfRange, name, v, lo, name, hi)
	}
	return nil
}

// checkPositiveAtMost returns ErrNotFinite for NaN or ±Inf and ErrOutOfRange unless 0 < v ≤ hi.
func checkPositiveAtMost(name string, v, hi float64) error {
	if err := checkPositiveFinite(name, v); err != nil {
		return err
	}
	if v > hi {
		return fmt.Errorf("%w: %s = %v, want ≤ %v", ErrOutOfRange, name, v, hi)
	}
	return nil
}

// Exact conversion factors (until unit/* exists).
const (
	// cmPerInch is exact by the 1959 international yard and pound; NIST SP 811 (2008) §B.8.
	cmPerInch = 2.54
	// kgPerPound is exact by the 1959 international yard and pound; NIST SP 811 (2008) §B.8.
	kgPerPound = 0.45359237
)

// SixMinuteWalkDistanceEnright returns the reference six-minute walk distance, in metres, for a healthy adult
// by the Enright–Sherrill sex-specific regression equations:
// men 7.57·heightCm − 5.02·ageYears − 1.76·weightKg − 309; women 2.11·heightCm − 2.29·weightKg − 5.78·ageYears + 667.
// The result is not rounded.
// Errors: ErrNotFinite for NaN or ±Inf; ErrOutOfRange outside Range.
//
// ID:      MED-RESP-002
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Enright PL, Sherrill DL. Reference equations for the six-minute walk in healthy adults. Am J Respir Crit Care Med. 1998;158(5 Pt 1):1384-1387. doi:10.1164/ajrccm.158.5.9710086. PMID 9817683, abstract: "for men, 6MWD = (7.57 x heightcm) - (5.02 x age) - (1.76 x weightkg) - 309 m"; "for women, 6MWD = (2.11 x heightcm) - (2.29 x weightkg) - (5.78 x age) + 667 m"
// Tags:    respiratory, six-minute-walk, 6mwd, exercise, reference-equation, enright
// Units:   heightCm cm; weightKg kg; ageYears years; female true for the women's equation; result metres
// Range:   ageYears 40–80 (the study's "healthy men and ... women, aged 40 to 80 yr"); heightCm > 0; weightKg > 0
// Example: man, 170.18 cm, 77.1107 kg, 60 y → 1288.2626 − 301.2 − 135.7148 − 309 = 542.35 m (hand-derived from the abstract's equation)
func SixMinuteWalkDistanceEnright(heightCm, weightKg, ageYears float64, female bool) (float64, error) {
	if err := checkPositiveFinite("heightCm", heightCm); err != nil {
		return 0, err
	}
	if err := checkPositiveFinite("weightKg", weightKg); err != nil {
		return 0, err
	}
	if err := checkClosedRange("ageYears", ageYears, 40, 80); err != nil {
		return 0, err
	}
	d := 2.11*heightCm - 2.29*weightKg - 5.78*ageYears + 667
	if !female {
		d = 7.57*heightCm - 5.02*ageYears - 1.76*weightKg - 309
	}
	if math.IsInf(d, 0) {
		return 0, fmt.Errorf("%w: the regression overflows for these inputs", ErrDomain)
	}
	return d, nil
}

// SixMinuteWalkDistanceEnrightImperial returns the reference six-minute walk distance, in metres, from height in
// inches and weight in pounds: it converts once at entry (exact NIST factors) and applies SixMinuteWalkDistanceEnright.
// Errors: as SixMinuteWalkDistanceEnright.
//
// ID:      MED-RESP-003
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Enright PL, Sherrill DL. Reference equations for the six-minute walk in healthy adults. Am J Respir Crit Care Med. 1998;158(5 Pt 1):1384-1387. doi:10.1164/ajrccm.158.5.9710086. PMID 9817683, abstract equations
// Source:  Thompson A, Taylor BN. Guide for the Use of the International System of Units (SI). NIST Special Publication 811, 2008 ed., §B.8 (inch = 2.54 cm, pound = 0.45359237 kg, exact)
// Tags:    respiratory, six-minute-walk, 6mwd, exercise, reference-equation, enright, imperial
// Units:   heightIn inches; weightLb pounds; ageYears years; result metres
// Range:   ageYears 40–80; heightIn > 0; weightLb > 0
// Example: man, 67 in, 170 lb, 60 y → 542.35 m (hand-derived; 67 in = 170.18 cm, 170 lb = 77.1107 kg)
func SixMinuteWalkDistanceEnrightImperial(heightIn, weightLb, ageYears float64, female bool) (float64, error) {
	if err := checkPositiveFinite("heightIn", heightIn); err != nil {
		return 0, err
	}
	if err := checkPositiveFinite("weightLb", weightLb); err != nil {
		return 0, err
	}
	return SixMinuteWalkDistanceEnright(heightIn*cmPerInch, weightLb*kgPerPound, ageYears, female)
}
