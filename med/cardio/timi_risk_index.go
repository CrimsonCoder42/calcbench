package cardio

// TIMIRiskIndex returns the TIMI risk index for ST-elevation MI, (heart rate × (age/10)²) / systolic blood
// pressure, unitless. The quintile cut-offs of the paper were not read and are not applied.
// Errors: ErrNotFinite for NaN or ±Inf; ErrOutOfRange outside Range.
//
// ID:      MED-CARDIO-006
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Morrow DA, Antman EM, Giugliano RP, et al. A simple risk index for rapid initial triage of patients with ST-elevation myocardial infarction: an InTIME II substudy. Lancet. 2001;358(9293):1571-1575. doi:10.1016/S0140-6736(01)06649-1. PMID 11716882. Abstract: "(heart rate x [age/10](2))/SBP"
// Tags:    cardio, acs, stemi, timi, risk-index
// Units:   heartRateBpm beats/min; ageYears years; systolicMmHg mmHg; result unitless
// Range:   ageYears 18–120; heartRateBpm > 0; systolicMmHg > 0 (13,253 persons with STEMI of InTIME II; the abstract states no input limits)
// Example: HR 80, age 60, SBP 120 → 80 × 36 / 120 = 24 (HAND-DERIVED from the abstract's formula)
func TIMIRiskIndex(heartRateBpm, ageYears, systolicMmHg float64) (float64, error) {
	for _, p := range [...]struct {
		n string
		v float64
	}{{"heartRateBpm", heartRateBpm}, {"ageYears", ageYears}, {"systolicMmHg", systolicMmHg}} {
		if err := checkFinite(p.n, p.v); err != nil {
			return 0, err
		}
	}
	if err := checkPositive("heartRateBpm", heartRateBpm); err != nil {
		return 0, err
	}
	if err := checkRange("ageYears", ageYears, 18, 120); err != nil {
		return 0, err
	}
	if err := checkPositive("systolicMmHg", systolicMmHg); err != nil {
		return 0, err
	}
	a := ageYears / 10
	return heartRateBpm * a * a / systolicMmHg, nil
}
