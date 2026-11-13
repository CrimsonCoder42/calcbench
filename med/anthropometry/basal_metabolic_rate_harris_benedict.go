package anthropometry

import "fmt"

// BasalMetabolicRateHarrisBenedict returns the predicted basal heat production, in kcal per 24 h, by the
// Harris–Benedict multiple regression equations of 1918, with the coefficients as printed:
// men h = 66.4730 + 13.7516 w + 5.0033 s − 6.7550 a; women h = 655.0955 + 9.5634 w + 1.8496 s − 4.6756 a
// (w weight in kg, s stature in cm, a age in years). Not rounded. No activity factor is applied (the source
// gives none).
// Errors: ErrNotFinite for NaN or ±Inf; ErrOutOfRange outside Range.
//
// ID:      MED-ANTHRO-010
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Harris JA, Benedict FG. A biometric study of human basal metabolism. Proc Natl Acad Sci USA. 1918;4(12):370-373. doi:10.1073/pnas.4.12.370. PMID 16576330. PMC1091498, p. 373: the two equations; "tabulated for values of weight from 25.0 to 124.9 kgm., for stature from 151 to 200 cm., and for age from 21 to 70 years"
// Tags:    anthropometry, energy, basal-metabolic-rate, bmr, bee, harris-benedict, metric
// Units:   weightKg kg; heightCm cm; ageYears years; result kcal/24 h (the source's "calories")
// Range:   weightKg 25.0–124.9; heightCm 151–200; ageYears 21–70 (the ranges the source tabulated)
// Example: man, 70 kg, 175 cm, 40 y → 66.4730 + 962.612 + 875.5775 − 270.2 = 1634.4625; woman, 60 kg, 165 cm, 30 y → 1393.8155 (HAND-DERIVED from p. 373)
func BasalMetabolicRateHarrisBenedict(weightKg, heightCm, ageYears float64, female bool) (float64, error) {
	if err := checkFinite("weightKg", weightKg); err != nil {
		return 0, err
	}
	if err := checkFinite("heightCm", heightCm); err != nil {
		return 0, err
	}
	if err := checkFinite("ageYears", ageYears); err != nil {
		return 0, err
	}
	if !(weightKg >= 25.0 && weightKg <= 124.9) {
		return 0, fmt.Errorf("%w: weightKg = %v, want 25.0–124.9", ErrOutOfRange, weightKg)
	}
	if !(heightCm >= 151 && heightCm <= 200) {
		return 0, fmt.Errorf("%w: heightCm = %v, want 151–200", ErrOutOfRange, heightCm)
	}
	if !(ageYears >= 21 && ageYears <= 70) {
		return 0, fmt.Errorf("%w: ageYears = %v, want 21–70", ErrOutOfRange, ageYears)
	}
	if female {
		return 655.0955 + 9.5634*weightKg + 1.8496*heightCm - 4.6756*ageYears, nil
	}
	return 66.4730 + 13.7516*weightKg + 5.0033*heightCm - 6.7550*ageYears, nil
}
