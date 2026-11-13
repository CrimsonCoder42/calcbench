package cardio

import "fmt"

// IMPEDEDexamethasone is the dexamethasone item of IMPEDE VTE.
type IMPEDEDexamethasone int

// Dexamethasone values (Sanfilippo 2019 score table): high dose is "total monthly dosing was greater than
// 160 mg", low dose is ≤ 160 mg per month.
const (
	IMPEDEDexNone IMPEDEDexamethasone = iota // 0
	IMPEDEDexLow                             // "Dexamethasone (Low-Dose)": 2
	IMPEDEDexHigh                            // "Dexamethasone (High-Dose)": 4
)

// IMPEDEThromboprophylaxis is the anticoagulant/antiplatelet item of IMPEDE VTE.
type IMPEDEThromboprophylaxis int

// Thromboprophylaxis values (score table).
const (
	IMPEDEProphylaxisNone         IMPEDEThromboprophylaxis = iota // 0
	IMPEDEProphylaxisProphylactic                                 // "Prophylactic LMWH or Aspirin": −3
	IMPEDEProphylaxisTherapeutic                                  // "Therapeutic LMWH or Warfarin": −4
)

// IMPEDEFactors are the yes/no items of IMPEDE VTE (score table labels in the comments).
type IMPEDEFactors struct {
	ImmunomodulatoryDrug      bool // 4
	PelvicHipFemurFracture    bool // "Pelvic, Hip or Femur Fracture": 4
	ErythropoiesisStimulating bool // 1
	Doxorubicin               bool // 3
	AsianPacificIslander      bool // "Ethnicity/Race = Asian/Pacific Islander": −3
	PriorVTE                  bool // "History of Venous Thromboembolism before MM": 5
	TunneledLineOrCVC         bool // "Tunneled Line/Central Venous Catheter": 2
}

// IMPEDEVTEScore returns the IMPEDE VTE score for multiple myeloma, −7 to 24: the item points plus 1 for body
// mass index ≥ 25 kg/m². The paper's tiers are ≤ 3, 4–7 and ≥ 8 (6-month VTE 3.3%, 8.3%, 15.2% in derivation);
// their names are not printed, so no category is returned.
// Errors: ErrNotFinite for NaN or ±Inf; ErrOutOfRange when bmiKgM2 is negative or an enumerated item is not
// one of its constants.
//
// ID:      MED-CARDIO-044
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Sanfilippo KM, Luo S, Wang TF, et al. Predicting venous thromboembolism in multiple myeloma. Am J Hematol. 2019;94(11):1176-1184. doi:10.1002/ajh.25603. PMID 31379000. PMC7058359, IMPEDE VTE score table (points 4, 1, 4, 1, 3, 4/2, −3, 5, 2, −4, −3)
// Tags:    hematology, multiple-myeloma, vte, impede, risk-score
// Units:   bmiKgM2 kg/m²; result points
// Range:   bmiKgM2 ≥ 0 (4,446 US veterans with multiple myeloma starting chemotherapy; SEER-Medicare validation)
// Example: IMiD, high-dose dexamethasone, BMI 27, prophylactic aspirin → 4 + 4 + 1 − 3 = 6 (HAND-DERIVED)
func IMPEDEVTEScore(bmiKgM2 float64, dex IMPEDEDexamethasone, prophylaxis IMPEDEThromboprophylaxis, f IMPEDEFactors) (int, error) {
	if err := checkFinite("bmiKgM2", bmiKgM2); err != nil {
		return 0, err
	}
	if err := checkRange("bmiKgM2", bmiKgM2, 0, maxFinite); err != nil {
		return 0, err
	}
	s := 0
	switch dex {
	case IMPEDEDexNone:
	case IMPEDEDexLow:
		s += 2
	case IMPEDEDexHigh:
		s += 4
	default:
		return 0, fmt.Errorf("%w: dex = %d, want 0, 1 or 2", ErrOutOfRange, dex)
	}
	switch prophylaxis {
	case IMPEDEProphylaxisNone:
	case IMPEDEProphylaxisProphylactic:
		s -= 3
	case IMPEDEProphylaxisTherapeutic:
		s -= 4
	default:
		return 0, fmt.Errorf("%w: prophylaxis = %d, want 0, 1 or 2", ErrOutOfRange, prophylaxis)
	}
	add := func(present bool, pts int) {
		if present {
			s += pts
		}
	}
	add(bmiKgM2 >= 25, 1)
	add(f.ImmunomodulatoryDrug, 4)
	add(f.PelvicHipFemurFracture, 4)
	add(f.ErythropoiesisStimulating, 1)
	add(f.Doxorubicin, 3)
	add(f.AsianPacificIslander, -3)
	add(f.PriorVTE, 5)
	add(f.TunneledLineOrCVC, 2)
	return s, nil
}
