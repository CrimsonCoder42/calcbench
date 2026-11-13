package cardio

// DAPTFactors are the DAPT score items besides age (Yeh 2016).
type DAPTFactors struct {
	VeinGraftPCI          bool // 2
	CigaretteSmoker       bool // "current cigarette smoker or within past year": 1
	Diabetes              bool // 1
	MIAtPresentation      bool // 1
	StentDiameterBelow3mm bool // "stent diameter: <3mm": 1
	CHFOrLVEFBelow30      bool // "history of CHF or LVEF<30%": 2
	PriorPCIOrMI          bool // 1
	PaclitaxelStent       bool // "paclitaxel-eluting stent": 1
}

// DAPTScore returns the DAPT score, −2 to 10: age < 65 0, 65 to < 75 −1, ≥ 75 −2; vein graft PCI 2; smoker 1;
// diabetes 1; MI at presentation 1; stent < 3 mm 1; CHF or LVEF < 30% 2; prior PCI or MI 1; paclitaxel-eluting
// stent 1. The paper calls a score ≥ 2 high and < 2 low.
// Errors: ErrNotFinite for NaN or ±Inf; ErrOutOfRange when ageYears is outside 18–120.
//
// ID:      MED-CARDIO-027
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Yeh RW, Secemsky EA, Kereiakes DJ, et al. Development and validation of a prediction rule for benefit and harm of dual antiplatelet therapy beyond 1 year after percutaneous coronary intervention. JAMA. 2016;315(16):1735-1749. doi:10.1001/jama.2016.3775. PMID 27022822. PMC5408574: "0 for age <65, −1 for age 65–<75, −2 for age≥75, 2 for vein graft PCI, 1 for current cigarette smoker or within past year, 1 for diabetes mellitus, 1 for MI at presentation, 1 for stent diameter: <3mm, 2 for history of CHF or LVEF<30%, 1 for prior PCI or prior MI, and 1 for paclitaxel-eluting stent"
// Tags:    cardio, pci, antiplatelet, dapt, risk-score
// Units:   ageYears years; result points −2 to 10
// Range:   ageYears 18–120 (11,648 randomized adults, mean age 61.3; no age limits stated)
// Example: age 70, diabetes, MI at presentation → −1 + 1 + 1 = 1 (HAND-DERIVED)
func DAPTScore(ageYears float64, f DAPTFactors) (int, error) {
	if err := checkFinite("ageYears", ageYears); err != nil {
		return 0, err
	}
	if err := checkRange("ageYears", ageYears, 18, 120); err != nil {
		return 0, err
	}
	s := 0
	switch {
	case ageYears >= 75:
		s = -2
	case ageYears >= 65:
		s = -1
	}
	add := func(present bool, pts int) {
		if present {
			s += pts
		}
	}
	add(f.VeinGraftPCI, 2)
	add(f.CigaretteSmoker, 1)
	add(f.Diabetes, 1)
	add(f.MIAtPresentation, 1)
	add(f.StentDiameterBelow3mm, 1)
	add(f.CHFOrLVEFBelow30, 2)
	add(f.PriorPCIOrMI, 1)
	add(f.PaclitaxelStent, 1)
	return s, nil
}
