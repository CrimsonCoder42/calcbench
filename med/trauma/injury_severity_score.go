package trauma

import "fmt"

// InjurySeverityScore returns the Injury Severity Score, 0–75, from the highest Abbreviated Injury Scale (AIS)
// severity in each of the six ISS body regions (head and neck, face, chest, abdomen and pelvic contents,
// extremities and pelvic bones, external/general): the sum of the squares of the three largest region values,
// and 75 when any region has an AIS of 6.
// Errors: ErrOutOfRange for a region value outside 0–6.
//
// ID:      MED-TRAUMA-016
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Baker SP, O'Neill B, Haddon W Jr, Long WB. The injury severity score: a method for describing patients with multiple injuries and evaluating emergency care. J Trauma. 1974;14(3):187-196. PMID 4814394 (primary; not opened)
// Source:  Clark DE, Black AW, Skavdahl DH, et al. Open-access programs for injury categorization using ICD-9 or ICD-10. Inj Epidemiol. 2018;5(1):11. doi:10.1186/s40621-018-0149-8. PMID 29629480. PMC5890002 ("sum of the squares of the three largest" region maxima; six regions; "assign RISS = 75 when any injury is assigned a severity of 6"; cites Baker 1974)
// Source:  Deng Q, Tang B, Xue C, et al. Comparison of the Ability to Predict Mortality between the Injury Severity Score and the New Injury Severity Score: A Meta-Analysis. Int J Environ Res Public Health. 2016;13(8):825. doi:10.3390/ijerph13080825. PMID 27537902. PMC4997511 (ISS definition, AIS 6 → 75, range 0–75; cites Baker 1974)
// Verified: secondary (Clark 2018; Deng 2016)
// Tags:    trauma, injury severity, iss, ais, score
// Units:   each argument the region's highest AIS severity, 0 (no injury) to 6; result points 0–75
// Range:   0 ≤ each region ≤ 6
// Example: head and neck 4, chest 3, extremities 2 → 16 + 9 + 4 = 29; any region 6 → 75 (HAND-DERIVED)
func InjurySeverityScore(headNeck, face, chest, abdomen, extremities, external int) (int, error) {
	regions := [...]int{headNeck, face, chest, abdomen, extremities, external}
	names := [...]string{"headNeck", "face", "chest", "abdomen", "extremities", "external"}
	var a, b, c int // the three largest, a ≥ b ≥ c
	for i, v := range regions {
		if v < 0 || v > 6 {
			return 0, fmt.Errorf("%w: %s = %d, want 0–6", ErrOutOfRange, names[i], v)
		}
	}
	for _, v := range regions {
		if v == 6 {
			return 75, nil
		}
		switch {
		case v > a:
			a, b, c = v, a, b
		case v > b:
			b, c = v, b
		case v > c:
			c = v
		}
	}
	return a*a + b*b + c*c, nil
}
