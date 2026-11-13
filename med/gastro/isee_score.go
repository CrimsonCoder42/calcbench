package gastro

import "fmt"

// ISEESymptoms is the symptom-frequency cell of the I-SEE "Symptoms and complications" row (Dellon 2022, Table 1).
type ISEESymptoms int

const (
	ISEESymptomsNone     ISEESymptoms = iota // 0 points
	ISEESymptomsWeekly                       // 1 point: "Weekly"
	ISEESymptomsDaily                        // 2 points: "Daily"
	ISEESymptomsMultiple                     // 4 points: "Multiple times per day or disrupting social functioning"
)

// ISEEExtent is a three-level cell of the I-SEE inflammatory endoscopy row.
type ISEEExtent int

const (
	ISEEExtentNone      ISEEExtent = iota // 0 points
	ISEEExtentLocalized                   // 1 point: "Localized"
	ISEEExtentDiffuse                     // 2 points: "Diffuse"
)

// ISEEFibrostenosis is the I-SEE fibrostenotic endoscopy row (rings, strictures).
type ISEEFibrostenosis int

const (
	ISEEFibrostenosisNone           ISEEFibrostenosis = iota // 0 points
	ISEEFibrostenosisPassesEasily                            // 1 point: "Present, but endoscope passes easily"
	ISEEFibrostenosisDilationOrSnug                          // 2 points: "Present, but requires dilation or a snug fit when passing a standard endoscope"
	ISEEFibrostenosisCannotPass                              // 15 points: "Cannot pass standard upper endoscope; repeated dilations (in an adult ≥18 years); or any dilation (in a child < 18 years)"
)

// ISEEFindings holds the I-SEE items of Dellon 2022, Table 1. "For boxes with more than one element, each
// selected feature gets points."
type ISEEFindings struct {
	Symptoms                ISEESymptoms      // 0, 1, 2 or 4 points
	FoodImpactionERVisit    bool              // "Food impaction with ER visit or endoscopy": 2 points (person ≥ 18 y), 4 points (< 18 y)
	Under18                 bool              // the person is under 18 years (sets the food-impaction points)
	HospitalizationDueToEoE bool              // "Hospitalization due to EoE": 4 points
	EsophagealPerforation   bool              // "Esophageal perforation": 15 points
	Malnutrition            bool              // "Malnutrition with body mass <5th percentile or decreased growth trajectory": 15 points
	PersistentInflammation  bool              // "Persistent inflammation requiring elemental formula, or systemic corticosteroid, or immunomodulatory treatments": 15 points
	EndoscopicInflammation  ISEEExtent        // edema, furrows and/or exudates: 0, 1 or 2 points
	PeakEosinophilsPerHPF   float64           // peak eosinophil count, eos/hpf: 15–60: 1 point, > 60: 2 points
	EndoscopicFibrostenosis ISEEFibrostenosis // rings, strictures: 0, 1, 2 or 15 points
	FibrostenoticHistology  bool              // "BZH or LPF (or DEC/SEA if no LP)": 2 points
}

// ISEEScore returns the Index of Severity for Eosinophilic Esophagitis (I-SEE) total: the sum of the points of
// every selected feature in the symptoms and complications, inflammatory and fibrostenotic domains of Table 1.
// A peak count below 15 eos/hpf scores 0.
// Errors: ErrNotFinite for NaN or ±Inf; ErrOutOfRange outside Range or for an unknown level.
//
// ID:      MED-GASTRO-006
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Dellon ES, Khoury P, Muir AB, et al. A clinical severity index for eosinophilic esophagitis: development, consensus, and future directions. Gastroenterology. 2022;163(1):59-76. doi:10.1053/j.gastro.2022.03.025. PMID 35606197. PMC9233087, Table 1 and its notes
// Tags:    gastro, eosinophilic-esophagitis, eoe, i-see, severity, score
// Units:   PeakEosinophilsPerHPF eosinophils per high-power field; result points
// Range:   PeakEosinophilsPerHPF ≥ 0; each level within its constants
// Example: daily symptoms, adult food impaction, diffuse inflammation, 70 eos/hpf, rings passing easily → 2 + 2 + 2 + 2 + 1 = 9 (HAND-DERIVED from Table 1)
func ISEEScore(f ISEEFindings) (int, error) {
	if err := checkFinite("PeakEosinophilsPerHPF", f.PeakEosinophilsPerHPF); err != nil {
		return 0, err
	}
	if !(f.PeakEosinophilsPerHPF >= 0) {
		return 0, fmt.Errorf("%w: PeakEosinophilsPerHPF = %v, want ≥ 0", ErrOutOfRange, f.PeakEosinophilsPerHPF)
	}
	points := 0
	switch f.Symptoms {
	case ISEESymptomsNone:
	case ISEESymptomsWeekly:
		points++
	case ISEESymptomsDaily:
		points += 2
	case ISEESymptomsMultiple:
		points += 4
	default:
		return 0, fmt.Errorf("%w: Symptoms = %d, want 0–3", ErrOutOfRange, f.Symptoms)
	}
	switch f.EndoscopicInflammation {
	case ISEEExtentNone:
	case ISEEExtentLocalized:
		points++
	case ISEEExtentDiffuse:
		points += 2
	default:
		return 0, fmt.Errorf("%w: EndoscopicInflammation = %d, want 0–2", ErrOutOfRange, f.EndoscopicInflammation)
	}
	switch f.EndoscopicFibrostenosis {
	case ISEEFibrostenosisNone:
	case ISEEFibrostenosisPassesEasily:
		points++
	case ISEEFibrostenosisDilationOrSnug:
		points += 2
	case ISEEFibrostenosisCannotPass:
		points += 15
	default:
		return 0, fmt.Errorf("%w: EndoscopicFibrostenosis = %d, want 0–3", ErrOutOfRange, f.EndoscopicFibrostenosis)
	}
	if f.FoodImpactionERVisit {
		if f.Under18 {
			points += 4
		} else {
			points += 2
		}
	}
	if f.HospitalizationDueToEoE {
		points += 4
	}
	for _, on := range [...]bool{f.EsophagealPerforation, f.Malnutrition, f.PersistentInflammation} {
		if on {
			points += 15
		}
	}
	switch {
	case f.PeakEosinophilsPerHPF > 60:
		points += 2
	case f.PeakEosinophilsPerHPF >= 15:
		points++
	}
	if f.FibrostenoticHistology {
		points += 2
	}
	return points, nil
}

// ISEECategory is an I-SEE severity category with the label of Dellon 2022, verbatim.
type ISEECategory string

const (
	ISEEInactive       ISEECategory = "Inactive EoE"        // total < 1
	ISEEMildActive     ISEECategory = "Mild Active EoE"     // 1–6
	ISEEModerateActive ISEECategory = "Moderate Active EoE" // 7–14
	ISEESevereActive   ISEECategory = "Severe Active EoE"   // ≥ 15
)

// ISEESeverity returns the I-SEE severity category for a total: "Inactive EoE" (< 1), "Mild Active EoE"
// (1–6), "Moderate Active EoE" (7–14), "Severe Active EoE" (≥ 15).
// Errors: ErrOutOfRange for a negative total.
//
// ID:      MED-GASTRO-007
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Dellon ES, Khoury P, Muir AB, et al. Gastroenterology. 2022;163(1):59-76. doi:10.1053/j.gastro.2022.03.025. PMID 35606197. PMC9233087, Table 1: "Total Score: <1: Inactive EoE; 1–6: Mild Active EoE; 7–14: Moderate Active EoE; ≥15: Severe Active EoE"
// Tags:    gastro, eosinophilic-esophagitis, eoe, i-see, severity, category
// Units:   total points (ISEEScore); result category label
// Range:   total ≥ 0
// Example: 0 → "Inactive EoE"; 6 → "Mild Active EoE"; 7 → "Moderate Active EoE"; 15 → "Severe Active EoE"
func ISEESeverity(total int) (ISEECategory, error) {
	switch {
	case total < 0:
		return "", fmt.Errorf("%w: total = %d, want ≥ 0", ErrOutOfRange, total)
	case total >= 15:
		return ISEESevereActive, nil
	case total >= 7:
		return ISEEModerateActive, nil
	case total >= 1:
		return ISEEMildActive, nil
	default:
		return ISEEInactive, nil
	}
}
