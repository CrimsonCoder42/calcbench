package critical

import "fmt"

// SOFA2UrineOutput is the kidney row's urine-output criterion of SOFA-2 Table 2, as observed by the caller.
type SOFA2UrineOutput int

// SOFA-2 urine-output criteria (Ranzani 2025 Table 2, kidney row).
const (
	SOFA2UrineNone                     SOFA2UrineOutput = iota // no urine-output criterion met
	SOFA2UrineBelow05For6To12h                                 // "urine output <0.5 mL/kg/h for 6-12 h" (1 point)
	SOFA2UrineBelow05For12h                                    // "urine output <0.5 mL/kg/h for ≥12 h" (2 points)
	SOFA2UrineBelow03For24hOrAnuria12h                         // "urine output <0.3 mL/kg/h for ≥24 h or anuria for ≥12 h" (3 points)
)

// SOFA2Input holds one day's worst values for the SOFA-2 score.
type SOFA2Input struct {
	// Brain
	GCS           int  // Glasgow Coma Scale total, 3–15 (last before sedation when sedated)
	DeliriumDrugs bool // "need for drugs to treat delirium" (at least 1 point)

	// Respiratory
	OxygenRatioMmHg      float64 // PaO2:FiO2 in mmHg, or SpO2:FiO2 when RatioIsSpO2FiO2
	RatioIsSpO2FiO2      bool    // the ratio is SpO2:FiO2 (footnote f: only when PaO2:FiO2 is unavailable and SpO2 < 98%)
	AdvancedVentSupport  bool    // advanced ventilatory support (HFNC, CPAP, BiPAP, NIV, IMV, home ventilation)
	RespiratoryECMO      bool    // ECMO for respiratory failure (4 points)
	VentSupportPrecluded bool    // footnote h: support not available or precluded by a ceiling of treatment

	// Cardiovascular
	MAPMmHg                    float64 // mean arterial pressure, mmHg
	NorepiPlusEpiUgKgMin       float64 // sum of norepinephrine (as base) and epinephrine, µg/kg/min, continuous ≥ 1 h
	DopamineUgKgMin            float64 // dopamine, µg/kg/min
	OtherVasopressorOrInotrope bool    // any vasopressor or inotrope other than norepinephrine, epinephrine and dopamine
	MechanicalSupport          bool    // VA-ECMO, IABP, LVAD or microaxial pump
	VasoactivesPrecluded       bool    // footnote m: vasoactive drugs unavailable or precluded by a ceiling of treatment

	// Liver
	BilirubinMgDL float64 // total bilirubin, mg/dL

	// Kidney
	CreatinineMgDL float64          // serum creatinine, mg/dL
	UrineOutput    SOFA2UrineOutput // urine-output criterion met
	RRT            bool             // receiving RRT, or meets the RRT criteria of footnote p (4 points)

	// Hemostasis
	PlateletsKPerUL float64 // platelets, ×10³/µL
}

// SOFA2Scores holds each SOFA-2 system's points (0–4) and the total (0–24).
type SOFA2Scores struct {
	Brain          int
	Respiratory    int
	Cardiovascular int
	Liver          int
	Kidney         int
	Hemostasis     int
	Total          int
}

// SOFA2Score returns the SOFA-2 score of Ranzani 2025 Table 2: six systems scored 0–4, total 0–24.
// Brain (GCS): 15 → 0; 13–14 or drugs for delirium → 1; 9–12 → 2; 6–8 → 3; 3–5 → 4.
// Respiratory, PaO2:FiO2 (mmHg): > 300 → 0; ≤ 300 → 1; ≤ 225 → 2; ≤ 150 with advanced ventilatory support → 3;
// ≤ 75 with advanced ventilatory support → 4; ECMO → 4. SpO2:FiO2 (footnote f): > 300 → 0; ≤ 300 → 1; ≤ 250 → 2;
// ≤ 200 with support → 3; ≤ 120 with support → 4. Without advanced support at most 2 points, unless support is
// unavailable or precluded (footnote h), when the ratio alone scores.
// Cardiovascular: MAP ≥ 70 without vasopressor or inotrope → 0; MAP < 70 → 1; norepinephrine + epinephrine
// ≤ 0.2 µg/kg/min, or any dose of another vasopressor or inotrope → 2; sum > 0.2 to ≤ 0.4, or low dose with
// another agent → 3; sum > 0.4, or medium dose with another agent, or mechanical support → 4. Dopamine as the
// single vasopressor (footnote l): ≤ 20 → 2; > 20 to ≤ 40 → 3; > 40 → 4; dopamine with norepinephrine or
// epinephrine counts as another agent. Vasoactives unavailable or precluded (footnote m): MAP ≥ 70 → 0; 60–69 → 1;
// 50–59 → 2; 40–49 → 3; < 40 → 4.
// Liver (bilirubin mg/dL): ≤ 1.20 → 0; ≤ 3.0 → 1; ≤ 6.0 → 2; ≤ 12.0 → 3; > 12 → 4.
// Kidney: creatinine ≤ 1.20 → 0; ≤ 2.0 → 1; ≤ 3.50 → 2; > 3.50 → 3; the urine-output criterion when higher;
// RRT or RRT criteria → 4. Hemostasis (platelets ×10³/µL): > 150 → 0; ≤ 150 → 1; ≤ 100 → 2; ≤ 80 → 3; ≤ 50 → 4.
// Errors: ErrNotFinite for NaN or ±Inf; ErrOutOfRange for a negative value, a GCS outside 3–15 or an unknown
// UrineOutput.
//
// ID:      MED-CRIT-016
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Ranzani OT, Singer M, Salluh JIF, et al. Development and Validation of the Sequential Organ Failure Assessment (SOFA)-2 Score. JAMA. 2025;334(23):2090-2103. doi:10.1001/jama.2025.20516. PMID 41159833. PMC12573119, Table 2 and its footnotes f, h, j, k, l, m, o, p, q
// Tags:    critical, sepsis, sofa-2, sofa, organ failure, score
// Units:   GCS points; OxygenRatioMmHg mmHg; MAPMmHg mmHg; vasopressors µg/kg/min; BilirubinMgDL mg/dL; CreatinineMgDL mg/dL; PlateletsKPerUL ×10³/µL; result points
// Range:   every numeric input ≥ 0; 3 ≤ GCS ≤ 15 (physical limits; adult ICU population of the source)
// Example: GCS 14, PaO2:FiO2 200 without advanced support, MAP 65, no vasoactives, bilirubin 2.0, creatinine 1.5, platelets 90 → 1+2+1+1+1+2 = 8 (HAND-DERIVED from Table 2)
func SOFA2Score(in SOFA2Input) (SOFA2Scores, error) {
	for _, e := range [...]error{
		finite("OxygenRatioMmHg", in.OxygenRatioMmHg), finite("MAPMmHg", in.MAPMmHg),
		finite("NorepiPlusEpiUgKgMin", in.NorepiPlusEpiUgKgMin), finite("DopamineUgKgMin", in.DopamineUgKgMin),
		finite("BilirubinMgDL", in.BilirubinMgDL), finite("CreatinineMgDL", in.CreatinineMgDL),
		finite("PlateletsKPerUL", in.PlateletsKPerUL),
		atLeastZero("OxygenRatioMmHg", in.OxygenRatioMmHg), atLeastZero("MAPMmHg", in.MAPMmHg),
		atLeastZero("NorepiPlusEpiUgKgMin", in.NorepiPlusEpiUgKgMin), atLeastZero("DopamineUgKgMin", in.DopamineUgKgMin),
		atLeastZero("BilirubinMgDL", in.BilirubinMgDL), atLeastZero("CreatinineMgDL", in.CreatinineMgDL),
		atLeastZero("PlateletsKPerUL", in.PlateletsKPerUL), checkGCS(in.GCS),
	} {
		if e != nil {
			return SOFA2Scores{}, e
		}
	}
	if in.UrineOutput < SOFA2UrineNone || in.UrineOutput > SOFA2UrineBelow03For24hOrAnuria12h {
		return SOFA2Scores{}, fmt.Errorf("%w: UrineOutput = %d, want 0–3", ErrOutOfRange, in.UrineOutput)
	}
	var s SOFA2Scores

	switch g := in.GCS; {
	case g <= 5:
		s.Brain = 4
	case g <= 8:
		s.Brain = 3
	case g <= 12:
		s.Brain = 2
	case g <= 14:
		s.Brain = 1
	}
	if in.DeliriumDrugs && s.Brain < 1 {
		s.Brain = 1
	}

	s.Respiratory = sofa2Respiratory(in)
	s.Cardiovascular = sofa2Cardiovascular(in)

	switch b := in.BilirubinMgDL; {
	case b > 12:
		s.Liver = 4
	case b > 6:
		s.Liver = 3
	case b > 3:
		s.Liver = 2
	case b > 1.2:
		s.Liver = 1
	}

	switch c := in.CreatinineMgDL; {
	case c > 3.5:
		s.Kidney = 3
	case c > 2:
		s.Kidney = 2
	case c > 1.2:
		s.Kidney = 1
	}
	if u := int(in.UrineOutput); u > s.Kidney {
		s.Kidney = u
	}
	if in.RRT {
		s.Kidney = 4
	}

	switch p := in.PlateletsKPerUL; {
	case p <= 50:
		s.Hemostasis = 4
	case p <= 80:
		s.Hemostasis = 3
	case p <= 100:
		s.Hemostasis = 2
	case p <= 150:
		s.Hemostasis = 1
	}

	s.Total = s.Brain + s.Respiratory + s.Cardiovascular + s.Liver + s.Kidney + s.Hemostasis
	return s, nil
}

// sofa2Respiratory scores the SOFA-2 respiratory row with footnotes f and h.
func sofa2Respiratory(in SOFA2Input) int {
	if in.RespiratoryECMO {
		return 4
	}
	// Cut-offs for 1, 2, 3 and 4 points: PaO2:FiO2 300/225/150/75; SpO2:FiO2 300/250/200/120.
	c1, c2, c3, c4 := 300.0, 225.0, 150.0, 75.0
	if in.RatioIsSpO2FiO2 {
		c2, c3, c4 = 250, 200, 120
	}
	r := in.OxygenRatioMmHg
	pts := 0
	switch {
	case r <= c4:
		pts = 4
	case r <= c3:
		pts = 3
	case r <= c2:
		pts = 2
	case r <= c1:
		pts = 1
	}
	if pts > 2 && !in.AdvancedVentSupport && !in.VentSupportPrecluded {
		pts = 2
	}
	return pts
}

// sofa2Cardiovascular scores the SOFA-2 cardiovascular row with footnotes l and m.
func sofa2Cardiovascular(in SOFA2Input) int {
	if in.VasoactivesPrecluded {
		switch m := in.MAPMmHg; {
		case m < 40:
			return 4
		case m < 50:
			return 3
		case m < 60:
			return 2
		case m < 70:
			return 1
		}
		return 0
	}
	if in.MechanicalSupport {
		return 4
	}
	ne, dopa := in.NorepiPlusEpiUgKgMin, in.DopamineUgKgMin
	if ne == 0 && dopa > 0 && !in.OtherVasopressorOrInotrope {
		switch {
		case dopa > 40:
			return 4
		case dopa > 20:
			return 3
		}
		return 2
	}
	other := in.OtherVasopressorOrInotrope || (dopa > 0 && ne > 0)
	switch {
	case ne > 0.4:
		return 4
	case ne > 0.2 && other:
		return 4
	case ne > 0.2:
		return 3
	case ne > 0 && other:
		return 3
	case ne > 0 || other:
		return 2
	case in.MAPMmHg < 70:
		return 1
	}
	return 0
}
