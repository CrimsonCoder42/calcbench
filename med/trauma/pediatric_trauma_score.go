package trauma

import "fmt"

// PTSGrade is the grade of a categorical Pediatric Trauma Score item: +2, +1 or −1, the item's printed points.
type PTSGrade int

// Pediatric Trauma Score grades (Kim 2021 Table 1; Misirlioglu 2024).
const (
	PTSMajor PTSGrade = -1 // airway unmaintainable; comatose or decerebrate; open or multiple fractures; major, penetrating or burn wounds
	PTSMinor PTSGrade = 1  // airway maintainable; obtunded or loss of consciousness; closed or suspected fracture; minor wounds
	PTSNone  PTSGrade = 2  // airway normal (patent); awake; no fracture; no wound
)

// PediatricTraumaScore returns the Pediatric Trauma Score, −6 to +12: weight > 20 kg +2, 10–20 kg +1, < 10 kg −1;
// systolic blood pressure > 90 mmHg +2, 50–90 +1, < 50 −1; and the airway, central nervous system, skeletal and
// cutaneous items graded +2, +1 or −1. The secondaries differ at the top edges (Kim "> 20" and "> 90";
// Misirlioglu "≥ 20" and "≥ 90"); Kim's form is used, so 20 kg and 90 mmHg score +1 (CONFIRM).
// Errors: ErrNotFinite for NaN or ±Inf; ErrOutOfRange for a negative weight or pressure or an unknown grade.
//
// ID:      MED-TRAUMA-018
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Tepas JJ 3rd, Mollitt DL, Talbert JL, Bryant M. The pediatric trauma score as a predictor of injury severity in the injured child. J Pediatr Surg. 1987;22(1):14-18. doi:10.1016/s0022-3468(87)80006-4. PMID 3102714 (primary; not opened)
// Source:  Kim HJ, Eun S, Yoon SH, et al. Paediatric Trauma Score as a non-imaging tool for predicting intracranial haemorrhage in patients with traumatic brain injury. Sci Rep. 2021;11(1):20911. doi:10.1038/s41598-021-00419-y. PMID 34686729. PMC8536669, Table 1 (cites Tepas 1987)
// Source:  Misirlioglu M, Alakaya M, Arslankoylu AE, et al. Evaluation of pediatric trauma score and pediatric age-adjusted shock index in pediatric patients admitted to the hospital after an earthquake. Ulus Travma Acil Cerrahi Derg. 2024;30(4):254-262. doi:10.14744/tjtes.2024.47835. PMID 38634847. PMC11065969 (PTS table; cites Tepas 1987)
// Verified: secondary (Kim 2021 Table 1; Misirlioglu 2024)
// Tags:    trauma, pediatrics, pts, triage, score
// Units:   weightKg kg; systolicMmHg mmHg; grades points; result points −6 to 12
// Range:   weightKg ≥ 0; systolicMmHg ≥ 0 (physical limits; injured children, the source's population)
// Example: 15 kg (+1), SBP 95 (+2), airway normal (+2), awake (+2), closed fracture (+1), minor wound (+1) → 9 (HAND-DERIVED)
func PediatricTraumaScore(weightKg, systolicMmHg float64, airway, cns, skeletal, cutaneous PTSGrade) (int, error) {
	for _, e := range [...]error{
		finite("weightKg", weightKg), finite("systolicMmHg", systolicMmHg),
		atLeastZero("weightKg", weightKg), atLeastZero("systolicMmHg", systolicMmHg),
	} {
		if e != nil {
			return 0, e
		}
	}
	names := [...]string{"airway", "cns", "skeletal", "cutaneous"}
	grades := [...]PTSGrade{airway, cns, skeletal, cutaneous}
	s := 0
	for i, g := range grades {
		if g != PTSMajor && g != PTSMinor && g != PTSNone {
			return 0, fmt.Errorf("%w: %s = %d, want −1, 1 or 2", ErrOutOfRange, names[i], g)
		}
		s += int(g)
	}
	switch {
	case weightKg > 20:
		s += 2
	case weightKg >= 10:
		s++
	default:
		s--
	}
	switch {
	case systolicMmHg > 90:
		s += 2
	case systolicMmHg >= 50:
		s++
	default:
		s--
	}
	return s, nil
}
