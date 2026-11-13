package respiratory

import "fmt"

// OttawaCOPDFindings are the ten items of the Ottawa COPD Risk Scale (Stiell 2014, Table 7), worded as in the
// table; the measured items are scored against the table's thresholds.
type OttawaCOPDFindings struct {
	CoronaryBypassGraft            bool    // History: "Coronary bypass graft" 1
	PeripheralVascularIntervention bool    // History: "Peripheral vascular disease intervention" 1
	IntubationForRespiratory       bool    // History: "Intubation for respiratory distress" 2
	HeartRateOnArrivalPerMin       float64 // Examination: "Heart rate on arrival in ED ≥ 110/min" 2
	TooIllForWalkTest              bool    // Examination: "Too ill to do walk test after treatment in ED (Sao2 < 90% or heart rate ≥ 120/min)" 2
	AcuteIschemicECG               bool    // Investigations: "Acute ischemic changes on ECG" 2
	PulmonaryCongestionCXR         bool    // Investigations: "Pulmonary congestion evident on chest radiography" 1
	HemoglobinGL                   float64 // Investigations: "Hemoglobin < 100 g/L" 3
	UreaMmolL                      float64 // Investigations: "Urea ≥ 12 mmol/L" 1
	SerumCO2MmolL                  float64 // Investigations: "Serum CO2 ≥ 35 mmol/L" 1
}

// OttawaCOPDRiskScale returns the Ottawa COPD Risk Scale total (0–16): coronary bypass graft 1, peripheral
// vascular disease intervention 1, intubation for respiratory distress 2, heart rate on arrival ≥ 110/min 2,
// too ill to do the walk test 2, acute ischemic ECG changes 2, pulmonary congestion on chest radiography 1,
// hemoglobin < 100 g/L 3, urea ≥ 12 mmol/L 1, serum CO₂ ≥ 35 mmol/L 1.
// Errors: ErrNotFinite for NaN or ±Inf; ErrOutOfRange outside Range.
//
// ID:      MED-RESP-030
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Stiell IG, Clement CM, Aaron SD, et al. Clinical characteristics associated with adverse events in patients with exacerbation of chronic obstructive pulmonary disease: a prospective cohort study. CMAJ. 2014;186(6):E193-E204. doi:10.1503/cmaj.130968. PMID 24549125, PMC3971051, Table 7 (the scale, "Total score (possible range 0–16)")
// Tags:    respiratory, copd, exacerbation, ottawa, emergency, score
// Units:   HeartRateOnArrivalPerMin beats/min; HemoglobinGL g/L (g/dL × 10); UreaMmolL mmol/L; SerumCO2MmolL mmol/L; the other items bool; result points 0–16
// Range:   HeartRateOnArrivalPerMin ≥ 0; HemoglobinGL > 0; UreaMmolL ≥ 0; SerumCO2MmolL ≥ 0
// Example: intubation history, heart rate 115, hemoglobin 95 g/L → 2 + 2 + 3 = 7 (hand-derived from Table 7)
func OttawaCOPDRiskScale(f OttawaCOPDFindings) (int, error) {
	if err := checkFiniteNonNegative("HeartRateOnArrivalPerMin", f.HeartRateOnArrivalPerMin); err != nil {
		return 0, err
	}
	if err := checkPositiveFinite("HemoglobinGL", f.HemoglobinGL); err != nil {
		return 0, err
	}
	if err := checkFiniteNonNegative("UreaMmolL", f.UreaMmolL); err != nil {
		return 0, err
	}
	if err := checkFiniteNonNegative("SerumCO2MmolL", f.SerumCO2MmolL); err != nil {
		return 0, err
	}
	s := 0
	add := func(present bool, pts int) {
		if present {
			s += pts
		}
	}
	add(f.CoronaryBypassGraft, 1)
	add(f.PeripheralVascularIntervention, 1)
	add(f.IntubationForRespiratory, 2)
	add(f.HeartRateOnArrivalPerMin >= 110, 2)
	add(f.TooIllForWalkTest, 2)
	add(f.AcuteIschemicECG, 2)
	add(f.PulmonaryCongestionCXR, 1)
	add(f.HemoglobinGL < 100, 3)
	add(f.UreaMmolL >= 12, 1)
	add(f.SerumCO2MmolL >= 35, 1)
	return s, nil
}

// OttawaCOPDRisk is a risk category of the Ottawa COPD Risk Scale; its value is the source's label, verbatim.
type OttawaCOPDRisk string

// The categories of Stiell 2014, Table 8 ("Risk category" column), verbatim.
const (
	OttawaCOPDLow      OttawaCOPDRisk = "Low"       // total 0
	OttawaCOPDMedium   OttawaCOPDRisk = "Medium"    // totals 1–2
	OttawaCOPDHigh     OttawaCOPDRisk = "High"      // totals 3–4
	OttawaCOPDVeryHigh OttawaCOPDRisk = "Very high" // totals 5 to "> 10"
)

// OttawaCOPDRiskCategory returns the Table 8 risk category for a total: 0 Low; 1–2 Medium; 3–4 High; 5 and
// above (the table's rows 5–10 and "> 10") Very high.
// Errors: ErrOutOfRange for a total outside 0–16.
//
// ID:      MED-RESP-031
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Stiell IG, Clement CM, Aaron SD, et al. Clinical characteristics associated with adverse events in patients with exacerbation of chronic obstructive pulmonary disease: a prospective cohort study. CMAJ. 2014;186(6):E193-E204. doi:10.1503/cmaj.130968. PMID 24549125, PMC3971051, Table 8 ("Total score | Risk of adverse event, % | Risk category")
// Tags:    respiratory, copd, ottawa, category
// Units:   total points 0–16; result the source's category label
// Range:   total 0–16
// Example: 0 → "Low"; 2 → "Medium"; 3 → "High"; 5 → "Very high" (Table 8)
func OttawaCOPDRiskCategory(total int) (OttawaCOPDRisk, error) {
	if total < 0 || total > 16 {
		return "", fmt.Errorf("%w: total = %d, want 0 ≤ total ≤ 16", ErrOutOfRange, total)
	}
	switch {
	case total == 0:
		return OttawaCOPDLow, nil
	case total <= 2:
		return OttawaCOPDMedium, nil
	case total <= 4:
		return OttawaCOPDHigh, nil
	}
	return OttawaCOPDVeryHigh, nil
}
