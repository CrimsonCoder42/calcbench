package renal

// KtV returns the dialysis dose index Kt/V (dimensionless): dialyzer urea clearance K times treatment time t,
// divided by the body urea distribution volume V. V is the caller's value; no estimate of V from body weight
// is applied here.
// Errors: ErrNotFinite for NaN or ±Inf; ErrOutOfRange for a negative clearance or time or a non-positive volume;
// ErrDomain when the result is beyond float64.
//
// ID:      MED-RENAL-024
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Gotch FA, Sargent JA. A mechanistic analysis of the National Cooperative Dialysis Study (NCDS). Kidney Int. 1985;28(3):526-534. doi:10.1038/ki.1985.160. PMID 3934452. Abstract: Kt/V, the "product of dialyzer urea clearance and treatment time divided by body urea volume"
// Tags:    renal, dialysis, kt/v, urea, adequacy
// Units:   clearanceMLPerMin mL/min; durationMin min; ureaVolumeML mL; result dimensionless
// Range:   clearanceMLPerMin ≥ 0; durationMin ≥ 0; ureaVolumeML > 0 (physical limits)
// Example: K 250 mL/min, t 240 min, V 42 000 mL → 250 × 240 / 42 000 ≈ 1.4286 (hand-derived)
func KtV(clearanceMLPerMin, durationMin, ureaVolumeML float64) (float64, error) {
	for _, e := range [...]error{
		checkFinite("clearanceMLPerMin", clearanceMLPerMin), checkFinite("durationMin", durationMin),
		checkFinite("ureaVolumeML", ureaVolumeML),
		checkNonNegative("clearanceMLPerMin", clearanceMLPerMin), checkNonNegative("durationMin", durationMin),
		checkPositive("ureaVolumeML", ureaVolumeML),
	} {
		if e != nil {
			return 0, e
		}
	}
	r := clearanceMLPerMin * durationMin / ureaVolumeML
	if err := checkResult(r); err != nil {
		return 0, err
	}
	return r, nil
}
