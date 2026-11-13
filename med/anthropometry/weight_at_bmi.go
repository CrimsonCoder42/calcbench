package anthropometry

import "fmt"

// WeightAtBMI returns the body weight, in kilograms, at which a person of the given height has the given body
// mass index: bmiKgM2 × heightM², the CDC BMI definition solved for weight. With the CDC adult category edges
// (18.5 and 25, MED-ANTHRO-003) it gives the weights that bound a category at a height.
// Errors: ErrNotFinite for NaN or ±Inf; ErrOutOfRange outside Range; ErrDomain when the result is beyond
// float64.
//
// ID:      MED-ANTHRO-022
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Centers for Disease Control and Prevention. BMI Frequently Asked Questions. https://www.cdc.gov/bmi/faq/index.html (last reviewed 2024-06-28; read 2026-11-12): "weight (kg)/[height (m)]2", solved for weight
// Tags:    anthropometry, bmi, body-mass-index, weight, inverse, metric
// Units:   bmiKgM2 kg/m²; heightM m; result kg
// Range:   bmiKgM2 > 0; heightM > 0 (the source states no limits)
// Example: BMI 25, 1.778 m → 25 × 3.161284 = 79.0321 kg; BMI 18.5, 1.778 m → 58.483754 kg (HAND-DERIVED)
func WeightAtBMI(bmiKgM2, heightM float64) (float64, error) {
	if err := checkFinite("bmiKgM2", bmiKgM2); err != nil {
		return 0, err
	}
	if err := checkFinite("heightM", heightM); err != nil {
		return 0, err
	}
	if err := checkPositive("bmiKgM2", bmiKgM2); err != nil {
		return 0, err
	}
	if err := checkPositive("heightM", heightM); err != nil {
		return 0, err
	}
	w := bmiKgM2 * heightM * heightM
	if err := checkResult(w); err != nil {
		return 0, err
	}
	if !(w > 0) {
		return 0, fmt.Errorf("%w: result underflows to 0", ErrDomain)
	}
	return w, nil
}

// WeightAtBMIImperial returns the body weight, in pounds, at which a person of the given height in inches has
// the given body mass index, by the CDC's imperial form solved for weight: bmiKgM2 × heightIn² / 703. It is the
// exact inverse of BodyMassIndexImperial (MED-ANTHRO-002), and so agrees with WeightAtBMI to within 0.01%, not
// exactly.
// Errors: ErrNotFinite for NaN or ±Inf; ErrOutOfRange outside Range; ErrDomain when the result is beyond
// float64.
//
// ID:      MED-ANTHRO-023
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Centers for Disease Control and Prevention. BMI Frequently Asked Questions. https://www.cdc.gov/bmi/faq/index.html (last reviewed 2024-06-28; read 2026-11-12): "weight (pounds)/[height (in)]2 x 703", solved for weight
// Tags:    anthropometry, bmi, body-mass-index, weight, inverse, imperial
// Units:   bmiKgM2 kg/m²; heightIn in; result lb
// Range:   bmiKgM2 > 0; heightIn > 0 (the source states no limits)
// Example: BMI 25, 70 in → 25 × 4900 / 703 = 174.25320056899005 lb; BMI 18.5, 70 in → 128.94736842105263 lb (HAND-DERIVED)
func WeightAtBMIImperial(bmiKgM2, heightIn float64) (float64, error) {
	if err := checkFinite("bmiKgM2", bmiKgM2); err != nil {
		return 0, err
	}
	if err := checkFinite("heightIn", heightIn); err != nil {
		return 0, err
	}
	if err := checkPositive("bmiKgM2", bmiKgM2); err != nil {
		return 0, err
	}
	if err := checkPositive("heightIn", heightIn); err != nil {
		return 0, err
	}
	w := bmiKgM2 * heightIn * heightIn / 703
	if err := checkResult(w); err != nil {
		return 0, err
	}
	if !(w > 0) {
		return 0, fmt.Errorf("%w: result underflows to 0", ErrDomain)
	}
	return w, nil
}
