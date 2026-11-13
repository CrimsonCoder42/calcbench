package critical

import (
	"fmt"
	"math"
)

// SAPS3Location is the intra-hospital location before ICU admission (Moreno 2005, Table 1, Box I).
type SAPS3Location int

const (
	SAPS3LocationNotListed     SAPS3Location = iota // a location Table 1 does not list (no points)
	SAPS3LocationEmergencyRoom                      // "Emergency room", 5 points
	SAPS3LocationOtherICU                           // "Other ICU", 7 points
	SAPS3LocationOther                              // "Other" (ward, other), 8 points
)

// SAPS3SurgicalStatus is the surgical status at ICU admission (Table 1, Box II).
type SAPS3SurgicalStatus int

const (
	SAPS3ScheduledSurgery SAPS3SurgicalStatus = iota // "Scheduled surgery", 0 points
	SAPS3NoSurgery                                   // "No surgery" (not submitted to surgery), 5 points
	SAPS3EmergencySurgery                            // "Emergency surgery", 6 points
)

// SAPS3SurgicalSite is the anatomical site of surgery (Table 2).
type SAPS3SurgicalSite int

const (
	SAPS3SiteAllOthers        SAPS3SurgicalSite = iota // "All others", 0 points
	SAPS3SiteTransplantation                           // transplantation surgery (liver, kidney, pancreas, kidney and pancreas, other), −11
	SAPS3SiteTrauma                                    // trauma – other, isolated (thorax, abdomen, limb); trauma – multiple, −8
	SAPS3SiteCABGWithoutValve                          // cardiac surgery: CABG without valvular repair, −6
	SAPS3SiteNeurosurgeryCVA                           // neurosurgery: cerebrovascular accident, 5
)

// SAPS3Input holds the SAPS 3 admission items other than bilirubin and creatinine, which are passed in the
// unit system of the function called. Physiology is the worst value within 1 h before or after ICU admission.
type SAPS3Input struct {
	// Box I.
	AgeYears                       float64
	CancerTherapy                  bool // 3 (chemotherapy, other immunosuppression, radiotherapy, steroid treatment)
	ChronicHeartFailureNYHAIV      bool // 6
	HaematologicalCancer           bool // 6
	Cirrhosis                      bool // 8
	AIDS                           bool // 8
	Cancer                         bool // 11 (metastatic cancer)
	HospitalDaysBeforeICU          float64
	Location                       SAPS3Location
	VasoactiveDrugsBeforeAdmission bool // 3

	// Box II.
	UnplannedAdmission bool // 3
	// Reasons for ICU admission (Table 2), all that apply.
	RhythmDisturbances              bool // −5 (−4 only when Seizures is also present)
	Seizures                        bool // −4
	HypovolemicShock                bool // 3 (hemorrhagic or non hemorrhagic)
	DigestiveAcuteAbdomenOrOther    bool // 3
	ComaStuporConfusionAgitation    bool // 4 (coma, stupor, obtunded, vigilance disturbances, confusion, agitation, delirium)
	SepticShock                     bool // 5
	AnaphylacticMixedUndefinedShock bool // 5
	LiverFailure                    bool // 6
	FocalNeurologicDeficit          bool // 7
	SeverePancreatitis              bool // 9
	IntracranialMassEffect          bool // 10
	SurgicalStatus                  SAPS3SurgicalStatus
	SurgicalSite                    SAPS3SurgicalSite
	NosocomialInfection             bool // 4
	RespiratoryInfection            bool // 5

	// Box III.
	GCS                   int     // lowest estimated Glasgow Coma Scale, 3–15
	TemperatureC          float64 // body temperature, °C (Table 1 labels it "highest")
	HeartRatePerMin       float64 // highest, beats/min
	LeukocytesGL          float64 // highest, G/L (10⁹/L)
	PH                    float64 // lowest arterial pH
	PlateletsGL           float64 // lowest, G/L (10⁹/L)
	SystolicMmHg          float64 // lowest systolic blood pressure, mmHg
	MechanicalVentilation bool
	PaO2MmHg              float64 // lowest arterial PaO2, mmHg
	FiO2                  float64 // inspired O2 fraction, 0.21–1; used only with MechanicalVentilation
}

// SAPS3Result holds the SAPS 3 admission score and the hospital mortality the general equation gives for it.
type SAPS3Result struct {
	Score              int     // offset 16 plus Boxes I–III; the paper's text says 0–217, but Table 1 as printed (points for all applicable combinations) sums to at most 251, which this code follows (decision with Devin)
	Logit              float64 // −32.6659 + ln(Score + 20.5958) × 7.3068
	ProbabilityOfDeath float64 // e^logit / (1 + e^logit), fraction 0–1
}

// SAPS3Score returns the SAPS 3 admission score of Moreno 2005 (Tables 1–2), bilirubin and creatinine in mg/dL,
// and the hospital mortality from the general equation. Points (Table 1 Box I): age < 40 0, 40–< 60 5, 60–< 70 9,
// 70–< 75 13, 75–< 80 15, ≥ 80 18; comorbidities cancer therapy 3, chronic HF NYHA IV 6, haematological cancer 6,
// cirrhosis 8, AIDS 8, cancer 11 (all that apply); days before ICU < 14 0, 14–< 28 6, ≥ 28 7; location emergency
// room 5, other ICU 7, other 8; vasoactive drugs 3. Box II: 16 for being admitted; unplanned 3; reasons (Table 2)
// rhythm disturbances −5, seizures −4 (both: −4), hypovolemic shock 3, digestive acute abdomen/other 3, coma etc. 4,
// septic shock 5, anaphylactic/mixed/undefined shock 5, liver failure 6, focal neurologic deficit 7, severe
// pancreatitis 9, intracranial mass effect 10 (all that apply); surgery scheduled 0, none 5, emergency 6; site
// transplantation −11, trauma −8, CABG without valve −6, neurosurgery CVA 5, others 0; infection nosocomial 4,
// respiratory 5. Box III: GCS 3–4 15, 5 10, 6 7, 7–12 2, ≥ 13 0; bilirubin < 2 0, 2–< 6 4, ≥ 6 5; temperature
// < 35 7; creatinine < 1.2 0, 1.2–< 2 2, 2–< 3.5 7, ≥ 3.5 8; heart rate < 120 0, 120–< 160 5, ≥ 160 7; leukocytes
// ≥ 15 2; pH ≤ 7.25 3; platelets < 20 13, 20–< 50 8, 50–< 100 5; systolic < 40 11, 40–< 70 8, 70–< 120 3;
// oxygenation with ventilation PaO2/FiO2 < 100 11, ≥ 100 7, without PaO2 < 60 5.
// Errors: ErrNotFinite for NaN or ±Inf; ErrOutOfRange outside Range or for an unknown enum value.
//
// ID:      MED-CRIT-020
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Moreno RP, Metnitz PG, Almeida E, et al. SAPS 3—From evaluation of the patient to evaluation of the intensive care unit. Part 2: Development of a prognostic model for hospital mortality at ICU admission. Intensive Care Med. 2005;31(10):1345-1355. doi:10.1007/s00134-005-2763-5. PMID 16132892, PMC1315315, Table 1 (Boxes I–III and footnotes), Table 2 (offset 16, reasons for admission, anatomical site, footnotes 12–13), text ("Logit = −32.6659 +ln(SAPS 3 score +20.5958) ×7.3068"; "Probability of death = e^logit/(1+e^logit)")
// Tags:    critical care, icu, saps 3, severity, mortality, score
// Units:   AgeYears years; HospitalDaysBeforeICU days; TemperatureC °C; HeartRatePerMin beats/min; LeukocytesGL and PlateletsGL 10⁹/L; PH pH units; SystolicMmHg and PaO2MmHg mmHg; FiO2 fraction; bilirubinMgDL and creatinineMgDL mg/dL; Score points; ProbabilityOfDeath fraction
// Range:   0 ≤ AgeYears ≤ 130; HospitalDaysBeforeICU ≥ 0; 3 ≤ GCS ≤ 15; TemperatureC 10–50; HeartRatePerMin, LeukocytesGL, PlateletsGL, SystolicMmHg, PaO2MmHg, bilirubinMgDL, creatinineMgDL ≥ 0; PH 6–8; FiO2 0.21–1 when MechanicalVentilation (physical limits; the paper states none)
// Example: age 65, unplanned, no surgery, septic shock, respiratory infection, GCS 15, HR 125, SBP 85, PaO2 55 without ventilation, others in their 0 bands → 9 + 16 + 3 + 5 + 5 + 5 + 5 + 3 + 5 = 56, logit −32.6659 + ln(76.5958) × 7.3068 = −0.96504, probability 0.27587 (HAND-DERIVED)
func SAPS3Score(in SAPS3Input, bilirubinMgDL, creatinineMgDL float64) (SAPS3Result, error) {
	if err := finite("bilirubinMgDL", bilirubinMgDL); err != nil {
		return SAPS3Result{}, err
	}
	if err := finite("creatinineMgDL", creatinineMgDL); err != nil {
		return SAPS3Result{}, err
	}
	if err := atLeastZero("bilirubinMgDL", bilirubinMgDL); err != nil {
		return SAPS3Result{}, err
	}
	if err := atLeastZero("creatinineMgDL", creatinineMgDL); err != nil {
		return SAPS3Result{}, err
	}
	bili := 0
	switch {
	case bilirubinMgDL >= 6:
		bili = 5
	case bilirubinMgDL >= 2:
		bili = 4
	}
	creat := 0
	switch {
	case creatinineMgDL >= 3.5:
		creat = 8
	case creatinineMgDL >= 2:
		creat = 7
	case creatinineMgDL >= 1.2:
		creat = 2
	}
	return saps3(in, bili+creat)
}

// SAPS3ScoreSI returns the SAPS 3 admission score of Moreno 2005 with bilirubin and creatinine in µmol/L, using
// the µmol/L rows Table 1 prints (bilirubin < 34.2 0, 34.2–< 102.6 4, ≥ 102.6 5; creatinine < 106.1 0,
// 106.1–< 176.8 2, 176.8–< 309.4 7, ≥ 309.4 8; no conversion is applied); every other item as SAPS3Score.
// Errors: ErrNotFinite for NaN or ±Inf; ErrOutOfRange outside Range or for an unknown enum value.
//
// ID:      MED-CRIT-021
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Moreno RP, Metnitz PG, Almeida E, et al. Intensive Care Med. 2005;31(10):1345-1355. doi:10.1007/s00134-005-2763-5. PMID 16132892, PMC1315315, Table 1 rows "Total bilirubine (highest), µmol/L" and "Creatinine (highest), µmol/L", Table 2, text (logit equation)
// Tags:    critical care, icu, saps 3, severity, mortality, score, si units
// Units:   bilirubinUmolL and creatinineUmolL µmol/L; every other input and the result as SAPS3Score
// Range:   bilirubinUmolL, creatinineUmolL ≥ 0; every other input as SAPS3Score
// Example: the SAPS3Score example with bilirubin 20 µmol/L and creatinine 80 µmol/L → 56 (HAND-DERIVED)
func SAPS3ScoreSI(in SAPS3Input, bilirubinUmolL, creatinineUmolL float64) (SAPS3Result, error) {
	if err := finite("bilirubinUmolL", bilirubinUmolL); err != nil {
		return SAPS3Result{}, err
	}
	if err := finite("creatinineUmolL", creatinineUmolL); err != nil {
		return SAPS3Result{}, err
	}
	if err := atLeastZero("bilirubinUmolL", bilirubinUmolL); err != nil {
		return SAPS3Result{}, err
	}
	if err := atLeastZero("creatinineUmolL", creatinineUmolL); err != nil {
		return SAPS3Result{}, err
	}
	bili := 0
	switch {
	case bilirubinUmolL >= 102.6:
		bili = 5
	case bilirubinUmolL >= 34.2:
		bili = 4
	}
	creat := 0
	switch {
	case creatinineUmolL >= 309.4:
		creat = 8
	case creatinineUmolL >= 176.8:
		creat = 7
	case creatinineUmolL >= 106.1:
		creat = 2
	}
	return saps3(in, bili+creat)
}

// saps3 checks the shared inputs, adds labPoints (bilirubin + creatinine) and applies the logit equation.
func saps3(in SAPS3Input, labPoints int) (SAPS3Result, error) {
	for _, v := range [...]struct {
		name string
		v    float64
	}{
		{"AgeYears", in.AgeYears}, {"HospitalDaysBeforeICU", in.HospitalDaysBeforeICU}, {"TemperatureC", in.TemperatureC},
		{"HeartRatePerMin", in.HeartRatePerMin}, {"LeukocytesGL", in.LeukocytesGL}, {"PH", in.PH},
		{"PlateletsGL", in.PlateletsGL}, {"SystolicMmHg", in.SystolicMmHg}, {"PaO2MmHg", in.PaO2MmHg}, {"FiO2", in.FiO2},
	} {
		if err := finite(v.name, v.v); err != nil {
			return SAPS3Result{}, err
		}
	}
	if err := within("AgeYears", in.AgeYears, 0, 130); err != nil {
		return SAPS3Result{}, err
	}
	if err := within("TemperatureC", in.TemperatureC, 10, 50); err != nil {
		return SAPS3Result{}, err
	}
	if err := within("PH", in.PH, 6, 8); err != nil {
		return SAPS3Result{}, err
	}
	for _, v := range [...]struct {
		name string
		v    float64
	}{
		{"HospitalDaysBeforeICU", in.HospitalDaysBeforeICU}, {"HeartRatePerMin", in.HeartRatePerMin},
		{"LeukocytesGL", in.LeukocytesGL}, {"PlateletsGL", in.PlateletsGL}, {"SystolicMmHg", in.SystolicMmHg},
		{"PaO2MmHg", in.PaO2MmHg},
	} {
		if err := atLeastZero(v.name, v.v); err != nil {
			return SAPS3Result{}, err
		}
	}
	if in.GCS < 3 || in.GCS > 15 {
		return SAPS3Result{}, fmt.Errorf("%w: GCS = %d, want 3–15", ErrOutOfRange, in.GCS)
	}
	if in.MechanicalVentilation {
		if err := within("FiO2", in.FiO2, 0.21, 1); err != nil {
			return SAPS3Result{}, err
		}
	}

	s := 16 + labPoints // Table 2 footnote 12: an offset of 16 points for being admitted

	// Box I.
	switch a := in.AgeYears; {
	case a >= 80:
		s += 18
	case a >= 75:
		s += 15
	case a >= 70:
		s += 13
	case a >= 60: // printed ">60<70"; 60 itself falls in no printed band otherwise (ledger note)
		s += 9
	case a >= 40:
		s += 5
	}
	s += pts(in.CancerTherapy, 3) + pts(in.ChronicHeartFailureNYHAIV, 6) + pts(in.HaematologicalCancer, 6) +
		pts(in.Cirrhosis, 8) + pts(in.AIDS, 8) + pts(in.Cancer, 11)
	switch d := in.HospitalDaysBeforeICU; {
	case d >= 28:
		s += 7
	case d >= 14:
		s += 6
	}
	switch in.Location {
	case SAPS3LocationNotListed:
	case SAPS3LocationEmergencyRoom:
		s += 5
	case SAPS3LocationOtherICU:
		s += 7
	case SAPS3LocationOther:
		s += 8
	default:
		return SAPS3Result{}, fmt.Errorf("%w: Location = %d", ErrOutOfRange, in.Location)
	}
	s += pts(in.VasoactiveDrugsBeforeAdmission, 3)

	// Box II.
	s += pts(in.UnplannedAdmission, 3)
	switch { // Table 2 footnote 13: "If both reasons for admission are present, only the worse value (–4) is scored."
	case in.RhythmDisturbances && in.Seizures:
		s -= 4
	case in.RhythmDisturbances:
		s -= 5
	case in.Seizures:
		s -= 4
	}
	s += pts(in.HypovolemicShock, 3) + pts(in.DigestiveAcuteAbdomenOrOther, 3) + pts(in.ComaStuporConfusionAgitation, 4) +
		pts(in.SepticShock, 5) + pts(in.AnaphylacticMixedUndefinedShock, 5) + pts(in.LiverFailure, 6) +
		pts(in.FocalNeurologicDeficit, 7) + pts(in.SeverePancreatitis, 9) + pts(in.IntracranialMassEffect, 10)
	switch in.SurgicalStatus {
	case SAPS3ScheduledSurgery:
	case SAPS3NoSurgery:
		s += 5
	case SAPS3EmergencySurgery:
		s += 6
	default:
		return SAPS3Result{}, fmt.Errorf("%w: SurgicalStatus = %d", ErrOutOfRange, in.SurgicalStatus)
	}
	switch in.SurgicalSite {
	case SAPS3SiteAllOthers:
	case SAPS3SiteTransplantation:
		s -= 11
	case SAPS3SiteTrauma:
		s -= 8
	case SAPS3SiteCABGWithoutValve:
		s -= 6
	case SAPS3SiteNeurosurgeryCVA:
		s += 5
	default:
		return SAPS3Result{}, fmt.Errorf("%w: SurgicalSite = %d", ErrOutOfRange, in.SurgicalSite)
	}
	s += pts(in.NosocomialInfection, 4) + pts(in.RespiratoryInfection, 5)

	// Box III.
	switch g := in.GCS; {
	case g <= 4:
		s += 15
	case g == 5:
		s += 10
	case g == 6:
		s += 7
	case g <= 12:
		s += 2
	}
	if in.TemperatureC < 35 {
		s += 7
	}
	switch h := in.HeartRatePerMin; {
	case h >= 160:
		s += 7
	case h >= 120:
		s += 5
	}
	if in.LeukocytesGL >= 15 {
		s += 2
	}
	if in.PH <= 7.25 {
		s += 3
	}
	switch p := in.PlateletsGL; {
	case p < 20:
		s += 13
	case p < 50:
		s += 8
	case p < 100:
		s += 5
	}
	switch b := in.SystolicMmHg; {
	case b < 40:
		s += 11
	case b < 70:
		s += 8
	case b < 120:
		s += 3
	}
	if in.MechanicalVentilation {
		if in.PaO2MmHg/in.FiO2 < 100 {
			s += 11
		} else {
			s += 7
		}
	} else if in.PaO2MmHg < 60 {
		s += 5
	}

	if s < 0 { // cannot happen with the printed points (minimum 16 − 5 − 11 = 0); guards the logarithm
		return SAPS3Result{}, fmt.Errorf("%w: score %d below the paper's minimum 0", ErrOutOfRange, s)
	}
	logit := -32.6659 + math.Log(float64(s)+20.5958)*7.3068
	return SAPS3Result{Score: s, Logit: logit, ProbabilityOfDeath: 1 / (1 + math.Exp(-logit))}, nil
}

// pts returns p when present, else 0.
func pts(present bool, p int) int {
	if present {
		return p
	}
	return 0
}
