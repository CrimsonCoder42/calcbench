package cardio

import "fmt"

// TRiPCastTrauma is the trauma item of TRiP(cast); one of the three must apply.
type TRiPCastTrauma int

// Trauma levels (Nemeth 2020, Table 1). The zero value is not valid.
const (
	TRiPCastTraumaLow          TRiPCastTrauma = iota + 1 // 1: single malleolar fracture, tarsal bones, non-severe sprain, muscle injury
	TRiPCastTraumaIntermediate                           // 2: bi/trimalleolar or patellar fracture, ankle dislocation, severe knee/ankle sprain
	TRiPCastTraumaHigh                                   // 3: fibula/tibia shaft or tibial plateau fracture, Achilles rupture
)

// TRiPCastImmobilization is the immobilisation item of TRiP(cast).
type TRiPCastImmobilization int

// Immobilisation levels (Table 1).
const (
	TRiPCastOtherWithPlantarSupport TRiPCastImmobilization = iota // 0
	TRiPCastFootCastOrSemiRigid                                   // "Foot cast (ankle free) or any semi-rigid without plantar support": 1
	TRiPCastLowerLegCast                                          // 2
	TRiPCastUpperLegCast                                          // 3
)

// TRiPCastFactors are the yes/no person factors of TRiP(cast) (Table 1; points in the comments).
type TRiPCastFactors struct {
	FamilyHistoryVTE            bool // 2
	PersonalVTEOrThrombophilia  bool // "Personal history of VTE or known major thrombophilia": 4
	OralContraceptiveOrEstrogen bool // "Current use of oral contraceptives or Estrogenic hormone therapy": 4
	Cancer5Years                bool // "Cancer diagnosis within the past 5 years": 3
	PregnancyOrPuerperium       bool // 3
	OtherImmobilization3Months  bool // "Immobilization (other) within the past 3 months": 2
	Surgery3Months              bool // 2
	Comorbidity                 bool // heart failure, rheumatoid arthritis, CKD, COPD or IBD: 1
	ChronicVenousInsufficiency  bool // "Chronic venous insufficiency (varicose veins)": 1
}

// TRiPCastScore returns the TRiP(cast) score for venous thromboembolism after lower-limb immobilisation for
// trauma: trauma 1–3; immobilisation 0–3; age < 35 0, 35 to < 55 1, 55 to < 75 2, ≥ 75 3; male 1; BMI ≥ 25 and
// < 35 1, ≥ 35 2; and the person factors. The paper's cut-off: < 7 low, ≥ 7 high.
// Errors: ErrNotFinite for NaN or ±Inf; ErrOutOfRange when ageYears is outside 18–120, bmiKgM2 is negative,
// or an enumerated item is not one of its constants.
//
// ID:      MED-CARDIO-045
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Nemeth B, Douillet D, le Cessie S, et al. Clinical risk assessment model to predict venous thromboembolism risk after immobilization for lower-limb trauma. EClinicalMedicine. 2020;20:100270. doi:10.1016/j.eclinm.2020.100270. PMID 32300734. PMC7152815, Table 1
// Tags:    vte, trauma, immobilization, cast, trip-cast, risk-score
// Units:   ageYears years; bmiKgM2 kg/m²; result points
// Range:   ageYears 18–120 (POT-CAST validation: adults ≥ 18 with a lower-leg cast); bmiKgM2 ≥ 0
// Example: tibial shaft fracture, lower-leg cast, man aged 40, BMI 27 → 3 + 2 + 1 + 1 + 1 = 8 (HAND-DERIVED from Table 1)
func TRiPCastScore(trauma TRiPCastTrauma, immob TRiPCastImmobilization, ageYears float64, male bool, bmiKgM2 float64, f TRiPCastFactors) (int, error) {
	if trauma < TRiPCastTraumaLow || trauma > TRiPCastTraumaHigh {
		return 0, fmt.Errorf("%w: trauma = %d, want 1, 2 or 3", ErrOutOfRange, trauma)
	}
	if immob < TRiPCastOtherWithPlantarSupport || immob > TRiPCastUpperLegCast {
		return 0, fmt.Errorf("%w: immob = %d, want 0–3", ErrOutOfRange, immob)
	}
	if err := checkFinite("ageYears", ageYears); err != nil {
		return 0, err
	}
	if err := checkFinite("bmiKgM2", bmiKgM2); err != nil {
		return 0, err
	}
	if err := checkRange("ageYears", ageYears, 18, 120); err != nil {
		return 0, err
	}
	if err := checkRange("bmiKgM2", bmiKgM2, 0, maxFinite); err != nil {
		return 0, err
	}
	s := int(trauma) + int(immob)
	switch {
	case ageYears >= 75:
		s += 3
	case ageYears >= 55:
		s += 2
	case ageYears >= 35:
		s++
	}
	switch {
	case bmiKgM2 >= 35:
		s += 2
	case bmiKgM2 >= 25:
		s++
	}
	add := func(present bool, pts int) {
		if present {
			s += pts
		}
	}
	add(male, 1)
	add(f.FamilyHistoryVTE, 2)
	add(f.PersonalVTEOrThrombophilia, 4)
	add(f.OralContraceptiveOrEstrogen, 4)
	add(f.Cancer5Years, 3)
	add(f.PregnancyOrPuerperium, 3)
	add(f.OtherImmobilization3Months, 2)
	add(f.Surgery3Months, 2)
	add(f.Comorbidity, 1)
	add(f.ChronicVenousInsufficiency, 1)
	return s, nil
}
