package critical

// MSOFAInput holds the values for the Modified Sequential Organ Failure Assessment (MSOFA) score.
type MSOFAInput struct {
	SpO2FiO2              float64 // SpO2 (%) divided by FiO2 (fraction), as Grissom 2010 states
	Jaundice              bool    // "Scleral icterus or jaundice"
	MAPMmHg               float64 // mean arterial pressure, mmHg
	DopamineUgKgMin       float64 // dopamine, µg/kg/min (0 when not given)
	Dobutamine            bool    // dobutamine at any dose
	EpinephrineUgKgMin    float64 // epinephrine, µg/kg/min (0 when not given)
	NorepinephrineUgKgMin float64 // norepinephrine, µg/kg/min (0 when not given)
	GCS                   int     // Glasgow Coma Scale total, 3–15
	CreatinineMgDL        float64 // serum creatinine, mg/dL
}

// MSOFAScores holds each MSOFA component and their sum (0–19).
type MSOFAScores struct {
	Respiratory    int // 0–4
	Liver          int // 0 or 3
	Cardiovascular int // 0–4
	CNS            int // 0–4
	Renal          int // 0–4
	Total          int // 0–19
}

// MSOFAScore returns the Modified Sequential Organ Failure Assessment score of Grissom 2010 Table 2: five
// systems, the hematologic system removed, "with a maximum score of 19". Respiratory (SpO2/FiO2): > 400 → 0;
// ≤ 400 → 1; ≤ 315 → 2; ≤ 235 → 3; ≤ 150 → 4. Liver: no scleral icterus or jaundice → 0; scleral icterus or
// jaundice → 3 (the only non-zero cell; four systems at 4 plus this one give the stated maximum of 19).
// Cardiovascular: no hypotension → 0; MAP < 70 → 1; dopamine ≤ 5 or dobutamine any dose → 2; dopamine > 5,
// epinephrine ≤ 0.1 or norepinephrine ≤ 0.1 → 3; dopamine > 15, epinephrine > 0.1 or norepinephrine > 0.1 → 4.
// CNS (GCS): 15 → 0; 13–14 → 1; 10–12 → 2; 6–9 → 3; < 6 → 4. Renal (creatinine mg/dL): < 1.2 → 0; 1.2–1.9 → 1;
// 2.0–3.4 → 2; 3.5–4.9 → 3; ≥ 5.0 → 4 (printed "> 5.0"; a value between printed ranges scores with the lower
// range, so 5.0 itself scores 4).
// Errors: ErrNotFinite for NaN or ±Inf; ErrOutOfRange for a negative value or a GCS outside 3–15.
//
// ID:      MED-CRIT-015
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Grissom CK, Brown SM, Kuttler KG, et al. A modified sequential organ failure assessment score for critical care triage. Disaster Med Public Health Prep. 2010;4(4):277-284. doi:10.1001/dmp.2010.40. PMID 21149228. PMC3811929, Table 2 "Modified Sequential Organ Failure Assessment (MSOFA) Score" and Methods ("5 organ systems are scored on a 0 to 4 scale … with a maximum score of 19"; "the SpO2 is divided by the fraction of inspired oxygen (FIO2)")
// Tags:    critical, sepsis, msofa, sofa, triage, organ failure, score
// Units:   SpO2FiO2 % per fraction; MAPMmHg mmHg; catecholamines µg/kg/min; GCS points; CreatinineMgDL mg/dL; result points
// Range:   every numeric input ≥ 0; 3 ≤ GCS ≤ 15 (physical limits; adult ICU population of the source)
// Example: SpO2 94% on FiO2 0.40 (ratio 235), jaundice, MAP 65, no vasoactives, GCS 12, creatinine 1.5 → 3+3+1+2+1 = 10 (HAND-DERIVED from Table 2)
func MSOFAScore(in MSOFAInput) (MSOFAScores, error) {
	for _, e := range [...]error{
		finite("SpO2FiO2", in.SpO2FiO2), finite("MAPMmHg", in.MAPMmHg),
		finite("DopamineUgKgMin", in.DopamineUgKgMin), finite("EpinephrineUgKgMin", in.EpinephrineUgKgMin),
		finite("NorepinephrineUgKgMin", in.NorepinephrineUgKgMin), finite("CreatinineMgDL", in.CreatinineMgDL),
		atLeastZero("SpO2FiO2", in.SpO2FiO2), atLeastZero("MAPMmHg", in.MAPMmHg),
		atLeastZero("DopamineUgKgMin", in.DopamineUgKgMin), atLeastZero("EpinephrineUgKgMin", in.EpinephrineUgKgMin),
		atLeastZero("NorepinephrineUgKgMin", in.NorepinephrineUgKgMin), atLeastZero("CreatinineMgDL", in.CreatinineMgDL),
		checkGCS(in.GCS),
	} {
		if e != nil {
			return MSOFAScores{}, e
		}
	}
	var s MSOFAScores
	switch r := in.SpO2FiO2; {
	case r <= 150:
		s.Respiratory = 4
	case r <= 235:
		s.Respiratory = 3
	case r <= 315:
		s.Respiratory = 2
	case r <= 400:
		s.Respiratory = 1
	}
	if in.Jaundice {
		s.Liver = 3
	}
	s.Cardiovascular = sofaCardiovascular(in.MAPMmHg, in.DopamineUgKgMin, in.Dobutamine, in.EpinephrineUgKgMin, in.NorepinephrineUgKgMin)
	s.CNS = sofaCNS(in.GCS)
	s.Renal = sofaCreatinine(in.CreatinineMgDL)
	s.Total = s.Respiratory + s.Liver + s.Cardiovascular + s.CNS + s.Renal
	return s, nil
}
