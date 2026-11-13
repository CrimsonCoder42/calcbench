package trauma

import "fmt"

// TASHScore returns the Trauma Associated Severe Hemorrhage (TASH) score, 0–28, as weighted in the 2006
// abstract: systolic blood pressure < 100 mmHg 4, < 120 mmHg 1; hemoglobin < 7 g/dL 8, < 9 6, < 10 4, < 11 3,
// < 12 2; intra-abdominal fluid 3; complex long bone and/or pelvic fractures AIS 3/4 3, AIS 5 6; heart rate
// > 120 2; base excess < −10 mmol/L 4, < −6 3, < −2 1; male 1. Each graded item scores its highest band met.
// longBoneOrPelvicFractureAIS is the Abbreviated Injury Scale severity of the worst complex long bone and/or
// pelvic fracture, 0 when there is none; values 1–2 score 0.
// Errors: ErrNotFinite for NaN or ±Inf; ErrOutOfRange for a negative hemoglobin, pressure or heart rate, or an AIS
// outside 0–5.
//
// ID:      MED-TRAUMA-013
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Yücel N, Lefering R, Maegele M, et al. Trauma Associated Severe Hemorrhage (TASH)-Score: probability of mass transfusion as surrogate for life threatening hemorrhage after multiple trauma. J Trauma. 2006;60(6):1228-1236. doi:10.1097/01.ta.0000220386.84012.bf. PMID 16766965. Abstract: "systolic blood pressure (<100 mm Hg=4 pts, <120 mm Hg=1 pt), hemoglobin (<7 g/dL=8 pts, <9 g/dL=6 pts, <10 g/dL=4 pts, <11 g/dL=3 pts, and <12 g/dL=2 pts), intra-abdominal fluid (3 pts), complex long bone and/or pelvic fractures (AIS 3/4=3 pts and AIS 5=6 pts), heart rate (>120=2 pts), base excess (<-10 mmol/L=4 pts, <-6 mmol/L=3 pts, and <-2 mmol/L=1 pt), and gender (male=1 pt)"
// Tags:    trauma, hemorrhage, massive transfusion, tash, score
// Units:   hemoglobinGDL g/dL; baseExcessMmolL mmol/L; systolicMmHg mmHg; heartRatePerMin beats/min; longBoneOrPelvicFractureAIS AIS severity 0–5; result points 0–28
// Range:   hemoglobinGDL ≥ 0; systolicMmHg ≥ 0; heartRatePerMin ≥ 0; baseExcessMmolL any finite value; 0 ≤ longBoneOrPelvicFractureAIS ≤ 5 (physical limits; multiple trauma, the study's population)
// Example: male, Hb 6.5, BE −11, SBP 95, HR 130, fluid, AIS 5 → 1+8+4+4+2+3+6 = 28 (HAND-DERIVED from the abstract's weights)
func TASHScore(male bool, hemoglobinGDL, baseExcessMmolL, systolicMmHg, heartRatePerMin float64, intraAbdominalFluid bool, longBoneOrPelvicFractureAIS int) (int, error) {
	for _, e := range [...]error{
		finite("hemoglobinGDL", hemoglobinGDL), finite("baseExcessMmolL", baseExcessMmolL),
		finite("systolicMmHg", systolicMmHg), finite("heartRatePerMin", heartRatePerMin),
		atLeastZero("hemoglobinGDL", hemoglobinGDL), atLeastZero("systolicMmHg", systolicMmHg),
		atLeastZero("heartRatePerMin", heartRatePerMin),
	} {
		if e != nil {
			return 0, e
		}
	}
	if longBoneOrPelvicFractureAIS < 0 || longBoneOrPelvicFractureAIS > 5 {
		return 0, fmt.Errorf("%w: longBoneOrPelvicFractureAIS = %d, want 0–5", ErrOutOfRange, longBoneOrPelvicFractureAIS)
	}
	s := 0
	switch {
	case systolicMmHg < 100:
		s += 4
	case systolicMmHg < 120:
		s++
	}
	switch {
	case hemoglobinGDL < 7:
		s += 8
	case hemoglobinGDL < 9:
		s += 6
	case hemoglobinGDL < 10:
		s += 4
	case hemoglobinGDL < 11:
		s += 3
	case hemoglobinGDL < 12:
		s += 2
	}
	if intraAbdominalFluid {
		s += 3
	}
	switch longBoneOrPelvicFractureAIS {
	case 3, 4:
		s += 3
	case 5:
		s += 6
	}
	if heartRatePerMin > 120 {
		s += 2
	}
	switch {
	case baseExcessMmolL < -10:
		s += 4
	case baseExcessMmolL < -6:
		s += 3
	case baseExcessMmolL < -2:
		s++
	}
	if male {
		s++
	}
	return s, nil
}
