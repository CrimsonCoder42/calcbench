package respiratory

import "fmt"

// LIPSHighRiskSurgery is the high-risk surgery item of the Lung Injury Prediction Score (one at most).
type LIPSHighRiskSurgery int

// The high-risk surgery rows of Gajic 2011, Table 3, with their points.
const (
	LIPSSurgeryNone           LIPSHighRiskSurgery = iota // no high-risk surgery, 0
	LIPSSurgeryOrthoSpine                                // "Orthopedic spine" 1
	LIPSSurgeryAcuteAbdomen                              // "Acute abdomen" 2
	LIPSSurgeryCardiac                                   // "Cardiac" 2.5
	LIPSSurgeryAorticVascular                            // "Aortic vascular" 3.5
)

// LIPSFindings are the predisposing conditions and risk modifiers of the Lung Injury Prediction Score (Gajic
// 2011, Table 3), worded as in the table.
type LIPSFindings struct {
	Shock              bool                // "Shock" 2
	Aspiration         bool                // "Aspiration" 2
	Sepsis             bool                // "Sepsis" 1
	Pneumonia          bool                // "Pneumonia" 1.5
	HighRiskSurgery    LIPSHighRiskSurgery // "High-risk surgery"
	EmergencySurgery   bool                // "Emergency surgery" +1.5
	TraumaticBrainInj  bool                // "Traumatic brain injury" 2
	SmokeInhalation    bool                // "Smoke inhalation" 2
	NearDrowning       bool                // "Near drowning" 2
	LungContusion      bool                // "Lung contusion" 1.5
	MultipleFractures  bool                // "Multiple fractures" 1.5
	AlcoholAbuse       bool                // "Alcohol abuse" 1
	ObesityBMIOver30   bool                // "Obesity (BMI > 30)" 1
	Hypoalbuminemia    bool                // "Hypoalbuminemia" 1
	Chemotherapy       bool                // "Chemotherapy" 1
	FiO2Over035        bool                // "FiO2 > 0.35 (> 4 L/min)" 2
	TachypneaRROver30  bool                // "Tachypnea (RR > 30)" 1.5
	SpO2Below95        bool                // "SpO2 < 95%" 1
	AcidosisPHBelow735 bool                // "Acidosis (pH < 7.35)" 1.5
	DiabetesMellitus   bool                // "Diabetes mellitus" −1, "only if sepsis"
}

// LIPSScore returns the Lung Injury Prediction Score of Gajic 2011, Table 3, in points (multiples of 0.5): shock
// 2, aspiration 2, sepsis 1, pneumonia 1.5; high-risk surgery orthopedic spine 1, acute abdomen 2, cardiac 2.5,
// aortic vascular 3.5, plus 1.5 for emergency surgery; high-risk trauma: traumatic brain injury 2, smoke
// inhalation 2, near drowning 2, lung contusion 1.5, multiple fractures 1.5; risk modifiers: alcohol abuse 1,
// obesity 1, hypoalbuminemia 1, chemotherapy 1, FiO₂ > 0.35 2, tachypnea 1.5, SpO₂ < 95 % 1, acidosis 1.5,
// diabetes mellitus −1 only if sepsis is present. The paper's cut-off, LIPS ≥ 4, is the caller's comparison.
// Errors: ErrOutOfRange for a HighRiskSurgery value that is not one of the constants.
//
// ID:      MED-RESP-043
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Gajic O, Dabbagh O, Park PK, et al. Early identification of patients at risk of acute lung injury: evaluation of lung injury prediction score in a multicenter cohort study. Am J Respir Crit Care Med. 2011;183(4):462-470. doi:10.1164/rccm.201004-0549OC. PMID 20802164, PMC3056224, Table 3 (LIPS points)
// Tags:    respiratory, ards, acute-lung-injury, lips, score
// Units:   every item bool or a LIPSHighRiskSurgery; result points (exact multiples of 0.5)
// Range:   any combination of items; adults at risk of acute lung injury at hospital admission
// Example: sepsis, pneumonia, tachypnea, SpO₂ < 95 %, diabetes → 1 + 1.5 + 1.5 + 1 − 1 = 4 (hand-derived from Table 3)
func LIPSScore(f LIPSFindings) (float64, error) {
	var surgeryHalves int
	switch f.HighRiskSurgery {
	case LIPSSurgeryNone:
	case LIPSSurgeryOrthoSpine:
		surgeryHalves = 2
	case LIPSSurgeryAcuteAbdomen:
		surgeryHalves = 4
	case LIPSSurgeryCardiac:
		surgeryHalves = 5
	case LIPSSurgeryAorticVascular:
		surgeryHalves = 7
	default:
		return 0, fmt.Errorf("%w: HighRiskSurgery = %d, want 0–4", ErrOutOfRange, f.HighRiskSurgery)
	}
	halves := surgeryHalves // the score in half points
	for _, t := range [19]struct {
		on     bool
		halves int
	}{
		{f.Shock, 4}, {f.Aspiration, 4}, {f.Sepsis, 2}, {f.Pneumonia, 3}, {f.EmergencySurgery, 3},
		{f.TraumaticBrainInj, 4}, {f.SmokeInhalation, 4}, {f.NearDrowning, 4}, {f.LungContusion, 3}, {f.MultipleFractures, 3},
		{f.AlcoholAbuse, 2}, {f.ObesityBMIOver30, 2}, {f.Hypoalbuminemia, 2}, {f.Chemotherapy, 2}, {f.FiO2Over035, 4},
		{f.TachypneaRROver30, 3}, {f.SpO2Below95, 2}, {f.AcidosisPHBelow735, 3}, {f.DiabetesMellitus && f.Sepsis, -2},
	} {
		if t.on {
			halves += t.halves
		}
	}
	return float64(halves) / 2, nil
}
