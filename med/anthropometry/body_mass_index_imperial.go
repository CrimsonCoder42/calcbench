package anthropometry

// The CDC's own imperial factor (BMI Frequently Asked Questions, https://www.cdc.gov/bmi/faq/index.html, read
// 2026-11-12: "weight (pounds)/[height (in)]2 x 703"). It is the CDC's rounding of the exact
// 0.45359237 / 0.0254² = 703.0695796…, so this form differs from MED-ANTHRO-001 on converted inputs by
// about 0.0099%.
const cdcImperialFactor = 703

// BodyMassIndexImperial returns the body mass index, in kg/m², from weight in pounds and height in inches by
// the CDC's imperial form: weightLb / heightIn² × 703. The factor 703 is the CDC's; the result agrees with
// BodyMassIndex (MED-ANTHRO-001) on the same person to within 0.01%, not exactly.
// Errors: ErrNotFinite for NaN or ±Inf; ErrOutOfRange outside Range; ErrDomain when the result is beyond
// float64.
//
// ID:      MED-ANTHRO-002
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Centers for Disease Control and Prevention. BMI Frequently Asked Questions. https://www.cdc.gov/bmi/faq/index.html (last reviewed 2024-06-28; read 2026-11-12): "weight (pounds)/[height (in)]2 x 703"
// Tags:    anthropometry, bmi, body-mass-index, imperial
// Units:   weightLb lb; heightIn in; result kg/m²
// Range:   weightLb > 0; heightIn > 0 (the source states no limits)
// Example: 150 lb, 65 in → 150 / 4225 × 703 = 24.958579881656806 (HAND-DERIVED)
func BodyMassIndexImperial(weightLb, heightIn float64) (float64, error) {
	if err := checkFinite("weightLb", weightLb); err != nil {
		return 0, err
	}
	if err := checkFinite("heightIn", heightIn); err != nil {
		return 0, err
	}
	if err := checkPositive("weightLb", weightLb); err != nil {
		return 0, err
	}
	if err := checkPositive("heightIn", heightIn); err != nil {
		return 0, err
	}
	r := weightLb / (heightIn * heightIn) * cdcImperialFactor
	if err := checkResult(r); err != nil {
		return 0, err
	}
	return r, nil
}
