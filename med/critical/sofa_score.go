package critical

import "fmt"

// SOFAInput holds one day's worst values for the Sequential (Sepsis-related) Organ Failure Assessment score.
type SOFAInput struct {
	PaO2FiO2MmHg          float64 // PaO2/FiO2 ratio, mmHg
	RespiratorySupport    bool    // "with respiratory support" (scores 3 and 4 need it)
	PlateletsKPerUL       float64 // platelets, ×10³/µL
	BilirubinMgDL         float64 // total bilirubin, mg/dL
	MAPMmHg               float64 // mean arterial pressure, mmHg
	DopamineUgKgMin       float64 // dopamine, µg/kg/min (0 when not given)
	Dobutamine            bool    // dobutamine at any dose
	EpinephrineUgKgMin    float64 // epinephrine, µg/kg/min (0 when not given)
	NorepinephrineUgKgMin float64 // norepinephrine, µg/kg/min (0 when not given)
	GCS                   int     // Glasgow Coma Scale total, 3–15
	CreatinineMgDL        float64 // serum creatinine, mg/dL
	UrineOutputMeasured   bool    // whether UrineOutputMLPerDay is known
	UrineOutputMLPerDay   float64 // urine output, mL/day (read only when UrineOutputMeasured)
}

// SOFAScores holds each organ system's SOFA score (0–4) and their sum (0–24).
type SOFAScores struct {
	Respiration    int
	Coagulation    int
	Liver          int
	Cardiovascular int
	CNS            int
	Renal          int
	Total          int
}

// SOFAScore returns the Sequential (Sepsis-related) Organ Failure Assessment score: each organ system's
// points (0–4) and the total (0–24), from the table of Vincent 1996 as reproduced in Sepsis-3 Table 1.
// Respiration: ≥ 400 → 0; < 400 → 1; < 300 → 2; < 200 with respiratory support → 3; < 100 with respiratory
// support → 4 (a ratio < 200 without support scores 2). Coagulation (platelets ×10³/µL): ≥ 150 → 0; < 150 → 1;
// < 100 → 2; < 50 → 3; < 20 → 4. Liver (bilirubin mg/dL): < 1.2 → 0; 1.2–1.9 → 1; 2.0–5.9 → 2; 6.0–11.9 → 3;
// ≥ 12.0 → 4. Cardiovascular: MAP ≥ 70 → 0; MAP < 70 → 1; dopamine ≤ 5 or dobutamine (any) → 2; dopamine > 5
// or epinephrine ≤ 0.1 or norepinephrine ≤ 0.1 → 3; dopamine > 15 or epinephrine > 0.1 or norepinephrine > 0.1
// → 4. CNS (GCS): 15 → 0; 13–14 → 1; 10–12 → 2; 6–9 → 3; < 6 → 4. Renal (creatinine mg/dL): < 1.2 → 0; 1.2–1.9
// → 1; 2.0–3.4 → 2; 3.5–4.9 or urine output < 500 mL/day → 3; ≥ 5.0 or urine output < 200 mL/day → 4. The
// table prints its ranges at 0.1 resolution; a value between two printed ranges scores with the lower range
// (e.g. creatinine 1.95 → 1), the reading where every printed lower edge is inclusive.
// Errors: ErrNotFinite for NaN or ±Inf; ErrOutOfRange for a negative value or a GCS outside 3–15.
//
// ID:      MED-CRIT-014
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Vincent JL, Moreno R, Takala J, et al. The SOFA (Sepsis-related Organ Failure Assessment) score to describe organ dysfunction/failure. Intensive Care Med. 1996;22(7):707-710. doi:10.1007/BF01709751. PMID 8844239 (primary; not opened)
// Source:  Singer M, Deutschman CS, Seymour CW, et al. The Third International Consensus Definitions for Sepsis and Septic Shock (Sepsis-3). JAMA. 2016;315(8):801-810. doi:10.1001/jama.2016.0287. PMID 26903338. PMC4968574, Table 1 "Sequential [Sepsis-Related] Organ Failure Assessment Score", adapted from Vincent et al
// Source:  Lambden S, Laterre PF, Levy MM, Francois B. The SOFA score-development, utility and challenges of accurate assessment in clinical trials. Crit Care. 2019;23(1):374. doi:10.1186/s13054-019-2663-7. PMID 31775846. PMC6880479, Table 1, "Modified from Vincent et al"
// Verified: secondary (Singer 2016 Table 1; Lambden 2019 Table 1)
// Tags:    critical, sepsis, sofa, organ failure, score
// Units:   PaO2FiO2MmHg mmHg; PlateletsKPerUL ×10³/µL; BilirubinMgDL mg/dL; MAPMmHg mmHg; catecholamines µg/kg/min; GCS points; CreatinineMgDL mg/dL; UrineOutputMLPerDay mL/day; result points
// Range:   every numeric input ≥ 0; 3 ≤ GCS ≤ 15 (physical limits; adult ICU population of the source)
// Example: PaO2/FiO2 250 without support, platelets 120, bilirubin 1.5, MAP 65 without vasoactives, GCS 14, creatinine 2.1 → 2+1+1+1+1+2 = 8 (HAND-DERIVED from Table 1)
func SOFAScore(in SOFAInput) (SOFAScores, error) {
	for _, e := range [...]error{
		finite("PaO2FiO2MmHg", in.PaO2FiO2MmHg), finite("PlateletsKPerUL", in.PlateletsKPerUL),
		finite("BilirubinMgDL", in.BilirubinMgDL), finite("MAPMmHg", in.MAPMmHg),
		finite("DopamineUgKgMin", in.DopamineUgKgMin), finite("EpinephrineUgKgMin", in.EpinephrineUgKgMin),
		finite("NorepinephrineUgKgMin", in.NorepinephrineUgKgMin), finite("CreatinineMgDL", in.CreatinineMgDL),
		finite("UrineOutputMLPerDay", in.UrineOutputMLPerDay),
		atLeastZero("PaO2FiO2MmHg", in.PaO2FiO2MmHg), atLeastZero("PlateletsKPerUL", in.PlateletsKPerUL),
		atLeastZero("BilirubinMgDL", in.BilirubinMgDL), atLeastZero("MAPMmHg", in.MAPMmHg),
		atLeastZero("DopamineUgKgMin", in.DopamineUgKgMin), atLeastZero("EpinephrineUgKgMin", in.EpinephrineUgKgMin),
		atLeastZero("NorepinephrineUgKgMin", in.NorepinephrineUgKgMin), atLeastZero("CreatinineMgDL", in.CreatinineMgDL),
		atLeastZero("UrineOutputMLPerDay", in.UrineOutputMLPerDay), checkGCS(in.GCS),
	} {
		if e != nil {
			return SOFAScores{}, e
		}
	}
	var s SOFAScores

	switch r := in.PaO2FiO2MmHg; {
	case r < 100 && in.RespiratorySupport:
		s.Respiration = 4
	case r < 200 && in.RespiratorySupport:
		s.Respiration = 3
	case r < 300:
		s.Respiration = 2
	case r < 400:
		s.Respiration = 1
	}

	switch p := in.PlateletsKPerUL; {
	case p < 20:
		s.Coagulation = 4
	case p < 50:
		s.Coagulation = 3
	case p < 100:
		s.Coagulation = 2
	case p < 150:
		s.Coagulation = 1
	}

	switch b := in.BilirubinMgDL; {
	case b >= 12:
		s.Liver = 4
	case b >= 6:
		s.Liver = 3
	case b >= 2:
		s.Liver = 2
	case b >= 1.2:
		s.Liver = 1
	}

	s.Cardiovascular = sofaCardiovascular(in.MAPMmHg, in.DopamineUgKgMin, in.Dobutamine, in.EpinephrineUgKgMin, in.NorepinephrineUgKgMin)
	s.CNS = sofaCNS(in.GCS)
	s.Renal = sofaCreatinine(in.CreatinineMgDL)
	if in.UrineOutputMeasured {
		switch u := in.UrineOutputMLPerDay; {
		case u < 200:
			s.Renal = 4
		case u < 500 && s.Renal < 3:
			s.Renal = 3
		}
	}
	s.Total = s.Respiration + s.Coagulation + s.Liver + s.Cardiovascular + s.CNS + s.Renal
	return s, nil
}

// checkGCS returns ErrOutOfRange unless 3 ≤ gcs ≤ 15.
func checkGCS(gcs int) error {
	if gcs < 3 || gcs > 15 {
		return fmt.Errorf("%w: GCS = %d, want 3–15", ErrOutOfRange, gcs)
	}
	return nil
}

// sofaCardiovascular scores the SOFA cardiovascular row (shared with MSOFA, Grissom 2010 Table 2, which
// prints the same row): catecholamine doses in µg/kg/min.
func sofaCardiovascular(mapMmHg, dopamine float64, dobutamine bool, epinephrine, norepinephrine float64) int {
	switch {
	case dopamine > 15 || epinephrine > 0.1 || norepinephrine > 0.1:
		return 4
	case dopamine > 5 || epinephrine > 0 || norepinephrine > 0:
		return 3
	case dopamine > 0 || dobutamine:
		return 2
	case mapMmHg < 70:
		return 1
	}
	return 0
}

// sofaCNS scores the SOFA central-nervous-system row from a GCS total of 3–15.
func sofaCNS(gcs int) int {
	switch {
	case gcs < 6:
		return 4
	case gcs < 10:
		return 3
	case gcs < 13:
		return 2
	case gcs < 15:
		return 1
	}
	return 0
}

// sofaCreatinine scores the SOFA renal row's creatinine column (mg/dL), shared with MSOFA.
func sofaCreatinine(creatinineMgDL float64) int {
	switch c := creatinineMgDL; {
	case c >= 5:
		return 4
	case c >= 3.5:
		return 3
	case c >= 2:
		return 2
	case c >= 1.2:
		return 1
	}
	return 0
}
