package cardio

// CARTScore returns the Cardiac Arrest Risk Triage (CART) score, 0–57, from four ward vital signs: respiratory
// rate < 21 → 0, 21–23 → 8, 24–25 → 12, 26–29 → 15, > 29 → 22; heart rate < 110 → 0, 110–139 → 4, > 139 → 13;
// diastolic pressure > 49 → 0, 40–49 → 4, 35–39 → 6, < 35 → 13; age < 55 → 0, 55–69 → 4, > 69 → 9. The paper reports
// the cut-off > 17 (specificity 89.9%, sensitivity 53.4%); no named categories. Band edges between the printed
// integers are closed upward (e.g. a respiratory rate of 20.5 scores 0, 21 scores 8).
// Errors: ErrNotFinite for NaN or ±Inf; ErrOutOfRange outside Range.
//
// ID:      MED-CARDIO-047
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Churpek MM, Yuen TC, Park SY, Meltzer DO, Hall JB, Edelson DP. Derivation of a cardiac arrest prediction model using ward vital signs. Crit Care Med. 2012;40(7):2102-2108. doi:10.1097/CCM.0b013e318250aa5a. PMID 22584764. PMC3378796, Table 3 (CART score points)
// Tags:    cardio, cardiac-arrest, ward, early-warning, cart, risk-score
// Units:   respiratoryRatePerMin breaths/min; heartRateBpm beats/min; diastolicMmHg mmHg; ageYears years; result points 0–57
// Range:   respiratoryRatePerMin 0–80; heartRateBpm 0–300; diastolicMmHg 0–200; ageYears 18–110 (hospitalised adults with ward vital signs, November 2008 to January 2011)
// Example: RR 24, HR 115, DBP 45, age 72 → 12 + 4 + 4 + 9 = 29 (HAND-DERIVED from Table 3)
func CARTScore(respiratoryRatePerMin, heartRateBpm, diastolicMmHg, ageYears float64) (int, error) {
	if err := checkBounded([]boundedInput{
		{"respiratoryRatePerMin", respiratoryRatePerMin, 0, 80},
		{"heartRateBpm", heartRateBpm, 0, 300},
		{"diastolicMmHg", diastolicMmHg, 0, 200},
		{"ageYears", ageYears, 18, 110},
	}); err != nil {
		return 0, err
	}
	s := 0
	switch {
	case respiratoryRatePerMin > 29:
		s += 22
	case respiratoryRatePerMin >= 26:
		s += 15
	case respiratoryRatePerMin >= 24:
		s += 12
	case respiratoryRatePerMin >= 21:
		s += 8
	}
	switch {
	case heartRateBpm > 139:
		s += 13
	case heartRateBpm >= 110:
		s += 4
	}
	switch {
	case diastolicMmHg < 35:
		s += 13
	case diastolicMmHg < 40:
		s += 6
	case diastolicMmHg <= 49:
		s += 4
	}
	switch {
	case ageYears > 69:
		s += 9
	case ageYears >= 55:
		s += 4
	}
	return s, nil
}

// boundedInput is one named float input with its validated range, for checkBounded.
type boundedInput struct {
	name      string
	v, lo, hi float64
}

// checkBounded returns ErrNotFinite for the first NaN or ±Inf input (all inputs are checked for finiteness
// before any range), then ErrOutOfRange for the first input outside lo–hi.
func checkBounded(ins []boundedInput) error {
	for _, in := range ins {
		if err := checkFinite(in.name, in.v); err != nil {
			return err
		}
	}
	for _, in := range ins {
		if err := checkRange(in.name, in.v, in.lo, in.hi); err != nil {
			return err
		}
	}
	return nil
}
