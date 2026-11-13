package gastro

import "fmt"

// RockallComorbidity is the comorbidity row of the Rockall numerical risk scoring system (Rockall 1996,
// Table III); each constant's comment quotes the table's cell.
type RockallComorbidity int

const (
	RockallComorbidityNone                             RockallComorbidity = iota // score 0: "No major comorbidity"
	RockallComorbidityMajor                                                      // score 2: "Cardiac failure, ischaemic heart disease, any major comorbidity"
	RockallComorbidityRenalLiverDisseminatedMalignancy                           // score 3: "Renal failure, liver failure, disseminated malignancy"
)

// RockallDiagnosis is the endoscopic diagnosis row of Rockall 1996, Table III.
type RockallDiagnosis int

const (
	RockallDiagnosisMalloryWeissOrNoLesion RockallDiagnosis = iota // score 0: "Mallory-Weiss tear, no lesion identified and no SRH"
	RockallDiagnosisOther                                          // score 1: "All other diagnoses"
	RockallDiagnosisUpperGIMalignancy                              // score 2: "Malignancy of upper GI tract"
)

// RockallFindings holds the Rockall items. Diagnosis and MajorStigmata are read only by RockallScoreComplete.
type RockallFindings struct {
	AgeYears      float64            // completed age, years: < 60: 0, 60–79: 1, ≥ 80: 2
	SystolicMmHg  float64            // systolic blood pressure, mmHg
	PulseBpm      float64            // pulse, beats/min
	Comorbidity   RockallComorbidity // 0, 2 or 3 points
	Diagnosis     RockallDiagnosis   // 0, 1 or 2 points
	MajorStigmata bool               // major SRH, "Blood in upper GI tract, adherent clot, visible or spurting vessel": 2 points; else "None or dark spot only": 0
}

// rockallPre checks f and returns the age, shock and comorbidity points (Table III).
func rockallPre(f RockallFindings) (int, error) {
	if err := checkFinite("AgeYears", f.AgeYears); err != nil {
		return 0, err
	}
	if err := checkFinite("SystolicMmHg", f.SystolicMmHg); err != nil {
		return 0, err
	}
	if err := checkFinite("PulseBpm", f.PulseBpm); err != nil {
		return 0, err
	}
	if !(f.AgeYears >= 0) {
		return 0, fmt.Errorf("%w: AgeYears = %v, want ≥ 0", ErrOutOfRange, f.AgeYears)
	}
	if !(f.SystolicMmHg >= 0) {
		return 0, fmt.Errorf("%w: SystolicMmHg = %v, want ≥ 0", ErrOutOfRange, f.SystolicMmHg)
	}
	if !(f.PulseBpm >= 0) {
		return 0, fmt.Errorf("%w: PulseBpm = %v, want ≥ 0", ErrOutOfRange, f.PulseBpm)
	}
	points := 0
	switch f.Comorbidity {
	case RockallComorbidityNone:
	case RockallComorbidityMajor:
		points += 2
	case RockallComorbidityRenalLiverDisseminatedMalignancy:
		points += 3
	default:
		return 0, fmt.Errorf("%w: Comorbidity = %d, want 0–2", ErrOutOfRange, f.Comorbidity)
	}
	switch {
	case f.AgeYears >= 80:
		points += 2
	case f.AgeYears >= 60:
		points++
	}
	switch {
	case f.SystolicMmHg < 100: // "Hypotension", systolic BP < 100
		points += 2
	case f.PulseBpm >= 100: // "Tachycardia", systolic BP ≥ 100, pulse ≥ 100
		points++
	}
	return points, nil
}

// RockallScorePreEndoscopy returns the Rockall score before endoscopy, 0–7: age (< 60: 0, 60–79: 1, ≥ 80: 2),
// shock ("No shock" systolic BP ≥ 100 and pulse < 100: 0; "Tachycardia" systolic BP ≥ 100 and pulse ≥ 100: 1;
// "Hypotension" systolic BP < 100: 2) and comorbidity (0, 2 or 3), from Table III. Diagnosis and
// MajorStigmata are ignored.
// Errors: ErrNotFinite for NaN or ±Inf; ErrOutOfRange outside Range or for an unknown Comorbidity.
//
// ID:      MED-GASTRO-002
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Rockall TA, Logan RF, Devlin HB, Northfield TC. Risk assessment after acute upper gastrointestinal haemorrhage. Gut. 1996;38(3):316-321. doi:10.1136/gut.38.3.316. PMID 8675081. PMC1383057, p. 318, Table III "Numerical risk scoring system": "Maximum additive score prior to diagnosis=7"
// Tags:    gastro, gi-bleeding, upper-gi, rockall, score
// Units:   AgeYears years; SystolicMmHg mmHg; PulseBpm beats/min; result points 0–7
// Range:   AgeYears, SystolicMmHg, PulseBpm ≥ 0 (population: 4185 cases of acute upper GI haemorrhage, UK national audit; input limits not stated in the source)
// Example: age 72, systolic 95, pulse 110, cardiac failure → 1 + 2 + 2 = 5 (HAND-DERIVED from Table III)
func RockallScorePreEndoscopy(f RockallFindings) (int, error) {
	return rockallPre(f)
}

// RockallScoreComplete returns the complete Rockall score, 0–11: the pre-endoscopy items plus diagnosis
// (0, 1 or 2) and major stigmata of recent haemorrhage (0 or 2), from Table III.
// Errors: ErrNotFinite for NaN or ±Inf; ErrOutOfRange outside Range or for an unknown Comorbidity or Diagnosis.
//
// ID:      MED-GASTRO-003
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Rockall TA, Logan RF, Devlin HB, Northfield TC. Gut. 1996;38(3):316-321. doi:10.1136/gut.38.3.316. PMID 8675081. PMC1383057, p. 318, Table III: "Maximum additive score following diagnosis=11"
// Tags:    gastro, gi-bleeding, upper-gi, rockall, endoscopy, score
// Units:   AgeYears years; SystolicMmHg mmHg; PulseBpm beats/min; result points 0–11
// Range:   AgeYears, SystolicMmHg, PulseBpm ≥ 0 (input limits not stated in the source)
// Example: age 85, systolic 90, renal failure, upper GI malignancy, visible vessel → 2 + 2 + 3 + 2 + 2 = 11 (HAND-DERIVED from Table III)
func RockallScoreComplete(f RockallFindings) (int, error) {
	points, err := rockallPre(f)
	if err != nil {
		return 0, err
	}
	switch f.Diagnosis {
	case RockallDiagnosisMalloryWeissOrNoLesion:
	case RockallDiagnosisOther:
		points++
	case RockallDiagnosisUpperGIMalignancy:
		points += 2
	default:
		return 0, fmt.Errorf("%w: Diagnosis = %d, want 0–2", ErrOutOfRange, f.Diagnosis)
	}
	if f.MajorStigmata {
		points += 2
	}
	return points, nil
}
