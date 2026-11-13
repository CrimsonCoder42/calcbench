package renal

import "fmt"

// ThakarSurgeryType is the type of cardiac surgery in the Thakar (Cleveland Clinic) score.
type ThakarSurgeryType int

const (
	ThakarCABGOnly      ThakarSurgeryType = iota // coronary artery bypass grafting only, 0 points
	ThakarValveOnly                              // valve repair/replacement only, 1 point
	ThakarCABGPlusValve                          // CABG + valve, 2 points
	ThakarOtherCardiac                           // any other cardiac surgery, 2 points
)

// ThakarInput holds the preoperative items of the Thakar (Cleveland Clinic) acute renal failure score.
type ThakarInput struct {
	Female                     bool              // 1 point
	CongestiveHeartFailure     bool              // 1 point
	LVEFPercent                float64           // left ventricular ejection fraction, %; < 35 % scores 1 point
	PreoperativeIABP           bool              // preoperative intra-aortic balloon pump, 2 points
	COPD                       bool              // chronic obstructive pulmonary disease, 1 point
	InsulinRequiringDiabetes   bool              // 1 point
	PreviousCardiacSurgery     bool              // 1 point
	EmergencySurgery           bool              // 2 points
	Surgery                    ThakarSurgeryType // 0, 1, 2 or 2 points
	PreoperativeCreatinineMgDL float64           // mg/dL; 1.2 to < 2.1 scores 2, ≥ 2.1 scores 5
}

// ThakarScore returns the Thakar (Cleveland Clinic) score for acute renal failure after cardiac surgery, 0–17:
// female 1; congestive heart failure 1; LVEF < 35 % 1; preoperative IABP 2; COPD 1; insulin-requiring diabetes 1;
// previous cardiac surgery 1; emergency surgery 2; surgery type CABG only 0, valve only 1, CABG + valve 2, other
// cardiac surgery 2; preoperative creatinine < 1.2 mg/dL 0, 1.2 to < 2.1 mg/dL 2, ≥ 2.1 mg/dL 5.
// Errors: ErrNotFinite for NaN or ±Inf; ErrOutOfRange for LVEFPercent outside 0–100,
// PreoperativeCreatinineMgDL ≤ 0 or an unknown Surgery.
//
// ID:      MED-RENAL-027
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Thakar CV, Arrigain S, Worley S, Yared JP, Paganini EP. A clinical score to predict acute renal failure after cardiac surgery. J Am Soc Nephrol. 2005;16(1):162-168. doi:10.1681/ASN.2004040331. PMID 15563569, abstract ("The score ranges between 0 and 17 points"); point table not in the abstract
// Source:  Crosina et al. Kidney Int Rep. 2017;2(2):172-179. PMC5678656, Table 1 (every item and its points; creatinine "107 to <186 μM (1.2 to <2.1 mg/dl)" 2, "≥186 μM (≥2.1 mg/dl)" 5; maximum 17); cites Thakar 2005
// Source:  Kozioł et al. PeerJ. 2026;14:e20533. PMC12812272, Table 1 (the same items and points; creatinine "<1.2 mg/dl" 0, "1.2 to <2.1 mg/dl" 2, "≥2.1 mg/dl" 5; "minimum score 0, maximum 17"); cites Thakar 2005
// Verified: secondary (Crosina 2017 PMC5678656 Table 1; Kozioł 2026 PMC12812272 Table 1)
// Tags:    renal, aki, acute renal failure, cardiac surgery, thakar, cleveland clinic, score
// Units:   LVEFPercent %; PreoperativeCreatinineMgDL mg/dL; result points 0–17
// Range:   0 ≤ LVEFPercent ≤ 100; PreoperativeCreatinineMgDL > 0 (the abstract reports no other limits)
// Example: female, LVEF 30 %, valve only, creatinine 1.5 mg/dL → 1 + 1 + 1 + 2 = 5 (HAND-DERIVED from Table 1)
func ThakarScore(in ThakarInput) (int, error) {
	if err := checkFinite("LVEFPercent", in.LVEFPercent); err != nil {
		return 0, err
	}
	if err := checkFinite("PreoperativeCreatinineMgDL", in.PreoperativeCreatinineMgDL); err != nil {
		return 0, err
	}
	if err := checkBetween("LVEFPercent", in.LVEFPercent, 0, 100); err != nil {
		return 0, err
	}
	if err := checkPositive("PreoperativeCreatinineMgDL", in.PreoperativeCreatinineMgDL); err != nil {
		return 0, err
	}
	s := 0
	switch in.Surgery {
	case ThakarCABGOnly:
	case ThakarValveOnly:
		s = 1
	case ThakarCABGPlusValve, ThakarOtherCardiac:
		s = 2
	default:
		return 0, fmt.Errorf("%w: Surgery = %d, want ThakarCABGOnly…ThakarOtherCardiac", ErrOutOfRange, in.Surgery)
	}
	for _, f := range [...]bool{in.Female, in.CongestiveHeartFailure, in.LVEFPercent < 35, in.COPD, in.InsulinRequiringDiabetes, in.PreviousCardiacSurgery} {
		if f {
			s++
		}
	}
	if in.PreoperativeIABP {
		s += 2
	}
	if in.EmergencySurgery {
		s += 2
	}
	switch c := in.PreoperativeCreatinineMgDL; {
	case c >= 2.1:
		s += 5
	case c >= 1.2:
		s += 2
	}
	return s, nil
}
