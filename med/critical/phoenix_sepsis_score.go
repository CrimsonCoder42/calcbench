package critical

import "fmt"

// PhoenixInput holds the values for the Phoenix Sepsis Score of a child (under 18 years). A measurement that was
// not made is marked by its Measured flag and then scores 0 for its criterion.
type PhoenixInput struct {
	AgeMonths float64 // age in months, 0 ≤ age < 216 (under 18 years)

	// Respiratory
	PaO2FiO2Measured      bool
	PaO2FiO2              float64 // PaO2:FiO2, mmHg
	SpO2FiO2Measured      bool    // the caller measures it only when SpO2 ≤ 97% (table footnote)
	SpO2FiO2              float64 // SpO2:FiO2
	AnyRespiratorySupport bool    // any respiratory support (needed for 1 point)
	IMV                   bool    // invasive mechanical ventilation (needed for 2 and 3 points)

	// Cardiovascular
	VasoactiveMedications int // number of vasoactive medications, ≥ 0
	LactateMeasured       bool
	LactateMmolL          float64 // lactate, mmol/L
	MAPMeasured           bool
	MAPMmHg               float64 // mean arterial pressure, mmHg

	// Coagulation
	PlateletsMeasured  bool
	PlateletsKPerUL    float64 // platelets, ×10³/µL
	INRMeasured        bool
	INR                float64 // international normalised ratio
	DDimerMeasured     bool
	DDimerMgLFEU       float64 // D-dimer, mg/L FEU
	FibrinogenMeasured bool
	FibrinogenMgDL     float64 // fibrinogen, mg/dL

	// Neurological
	GCSMeasured          bool
	GCS                  int  // Glasgow Coma Scale total, 3–15
	FixedPupilsBilateral bool // "Fixed pupils bilaterally" (2 points)
}

// PhoenixScores holds each Phoenix system's points, the total (0–13) and the consensus criteria computed from it.
type PhoenixScores struct {
	Respiratory    int  // 0–3
	Cardiovascular int  // 0–6
	Coagulation    int  // 0–2
	Neurological   int  // 0–2
	Total          int  // 0–13
	ScoreAtLeast2  bool // "Phoenix Sepsis Score ≥2 points" (sepsis when infection is suspected, which the caller judges)
	SepticShockCV  bool // ScoreAtLeast2 and "≥1 cardiovascular point(s)" (septic shock when infection is suspected)
}

// phoenixMAPEdges are the upper ages (months, exclusive) and the MAP limits (mmHg) of the Phoenix table: a MAP
// below lo2 scores 2, a MAP from lo2 up to and including hi1 scores 1. Rows: < 1 mo 17–30; 1–11 mo 25–38;
// 1 to < 2 y 31–43; 2 to < 5 y 32–44; 5 to < 12 y 36–48; 12–17 y 38–51.
var phoenixMAPEdges = [...]struct{ ageBelowMonths, lo2, hi1 float64 }{
	{1, 17, 30}, {12, 25, 38}, {24, 31, 43}, {60, 32, 44}, {144, 36, 48}, {216, 38, 51},
}

// PhoenixSepsisScore returns the Phoenix Sepsis Score of the 2024 international consensus criteria for paediatric
// sepsis. Respiratory: PaO2:FiO2 < 400 or SpO2:FiO2 < 292 on any respiratory support → 1; PaO2:FiO2 100–200 or
// SpO2:FiO2 148–220 with IMV → 2; PaO2:FiO2 < 100 or SpO2:FiO2 < 148 with IMV → 3 (the worse of the two ratios
// counts). Cardiovascular (sum, 0–6): 1 vasoactive medication → 1, ≥ 2 → 2; lactate 5–10.9 mmol/L → 1, ≥ 11 → 2;
// MAP by age (table rows above) → 1 or 2. Coagulation (1 each, at most 2): platelets < 100 ×10³/µL; INR > 1.3;
// D-dimer > 2 mg/L FEU; fibrinogen < 100 mg/dL. Neurological: GCS ≤ 10 → 1; fixed pupils bilaterally → 2.
// Errors: ErrNotFinite for NaN or ±Inf; ErrOutOfRange for an age outside 0 to < 216 months, a negative value,
// a negative medication count or a GCS outside 3–15.
//
// ID:      MED-CRIT-017
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Schlapbach LJ, Watson RS, Sorce LR, et al. International Consensus Criteria for Pediatric Sepsis and Septic Shock. JAMA. 2024;331(8):665-674. doi:10.1001/jama.2024.0179. PMID 38245889. PMC10900966, the Phoenix Sepsis Score table and "Sepsis: Suspected infection and Phoenix Sepsis Score ≥2 points"; "Septic shock: Sepsis with ≥1 cardiovascular point(s)"
// Tags:    critical, sepsis, pediatrics, phoenix, septic shock, score
// Units:   AgeMonths months; PaO2FiO2 mmHg; LactateMmolL mmol/L; MAPMmHg mmHg; PlateletsKPerUL ×10³/µL; DDimerMgLFEU mg/L FEU; FibrinogenMgDL mg/dL; GCS points; result points
// Range:   0 ≤ AgeMonths < 216; every other numeric input ≥ 0; 3 ≤ GCS ≤ 15 (children under 18 years, the source's population)
// Example: 30 months, SpO2:FiO2 200 on IMV, one vasoactive, lactate 6, MAP 40, platelets 80, GCS 12 → 2 + (1+1+1) + 1 + 0 = 6, septic-shock criterion met (HAND-DERIVED from the table)
func PhoenixSepsisScore(in PhoenixInput) (PhoenixScores, error) {
	for _, e := range [...]error{
		finite("AgeMonths", in.AgeMonths), finite("PaO2FiO2", in.PaO2FiO2), finite("SpO2FiO2", in.SpO2FiO2),
		finite("LactateMmolL", in.LactateMmolL), finite("MAPMmHg", in.MAPMmHg),
		finite("PlateletsKPerUL", in.PlateletsKPerUL), finite("INR", in.INR),
		finite("DDimerMgLFEU", in.DDimerMgLFEU), finite("FibrinogenMgDL", in.FibrinogenMgDL),
		atLeastZero("PaO2FiO2", in.PaO2FiO2), atLeastZero("SpO2FiO2", in.SpO2FiO2),
		atLeastZero("LactateMmolL", in.LactateMmolL), atLeastZero("MAPMmHg", in.MAPMmHg),
		atLeastZero("PlateletsKPerUL", in.PlateletsKPerUL), atLeastZero("INR", in.INR),
		atLeastZero("DDimerMgLFEU", in.DDimerMgLFEU), atLeastZero("FibrinogenMgDL", in.FibrinogenMgDL),
	} {
		if e != nil {
			return PhoenixScores{}, e
		}
	}
	if !(in.AgeMonths >= 0 && in.AgeMonths < 216) {
		return PhoenixScores{}, fmt.Errorf("%w: AgeMonths = %v, want 0 ≤ age < 216", ErrOutOfRange, in.AgeMonths)
	}
	if in.VasoactiveMedications < 0 {
		return PhoenixScores{}, fmt.Errorf("%w: VasoactiveMedications = %d, want ≥ 0", ErrOutOfRange, in.VasoactiveMedications)
	}
	if in.GCSMeasured {
		if err := checkGCS(in.GCS); err != nil {
			return PhoenixScores{}, err
		}
	}
	var s PhoenixScores

	pf, sf := in.PaO2FiO2Measured, in.SpO2FiO2Measured
	switch {
	case in.IMV && ((pf && in.PaO2FiO2 < 100) || (sf && in.SpO2FiO2 < 148)):
		s.Respiratory = 3
	case in.IMV && ((pf && in.PaO2FiO2 <= 200) || (sf && in.SpO2FiO2 <= 220)):
		s.Respiratory = 2
	case (in.AnyRespiratorySupport || in.IMV) && ((pf && in.PaO2FiO2 < 400) || (sf && in.SpO2FiO2 < 292)):
		s.Respiratory = 1
	}

	switch {
	case in.VasoactiveMedications >= 2:
		s.Cardiovascular += 2
	case in.VasoactiveMedications == 1:
		s.Cardiovascular++
	}
	if in.LactateMeasured {
		switch {
		case in.LactateMmolL >= 11:
			s.Cardiovascular += 2
		case in.LactateMmolL >= 5:
			s.Cardiovascular++
		}
	}
	if in.MAPMeasured {
		for _, row := range phoenixMAPEdges {
			if in.AgeMonths < row.ageBelowMonths {
				switch {
				case in.MAPMmHg < row.lo2:
					s.Cardiovascular += 2
				case in.MAPMmHg <= row.hi1:
					s.Cardiovascular++
				}
				break
			}
		}
	}

	for _, met := range [...]bool{
		in.PlateletsMeasured && in.PlateletsKPerUL < 100,
		in.INRMeasured && in.INR > 1.3,
		in.DDimerMeasured && in.DDimerMgLFEU > 2,
		in.FibrinogenMeasured && in.FibrinogenMgDL < 100,
	} {
		if met && s.Coagulation < 2 {
			s.Coagulation++
		}
	}

	switch {
	case in.FixedPupilsBilateral:
		s.Neurological = 2
	case in.GCSMeasured && in.GCS <= 10:
		s.Neurological = 1
	}

	s.Total = s.Respiratory + s.Cardiovascular + s.Coagulation + s.Neurological
	s.ScoreAtLeast2 = s.Total >= 2
	s.SepticShockCV = s.ScoreAtLeast2 && s.Cardiovascular >= 1
	return s, nil
}
