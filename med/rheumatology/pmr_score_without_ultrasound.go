package rheumatology

// PMRScoreWithoutUltrasound returns the 2012 EULAR/ACR polymyalgia rheumatica scoring-algorithm score without
// ultrasound, 0–6 (Table 6): morning stiffness > 45 min, 2; hip pain or limited range of motion, 1; absence of
// RF or ACPA, 2; absence of other joint involvement, 1. The algorithm applies to people aged ≥ 50 with
// bilateral shoulder aching and abnormal CRP and/or ESR (the caller's to check); "A score of 4 or more is
// categorised as PMR in the algorithm without US".
//
// ID:      MED-RHEUM-002
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Dasgupta B, Cimmino MA, Maradit-Kremers H, et al. 2012 provisional classification criteria for polymyalgia rheumatica: a European League Against Rheumatism/American College of Rheumatology collaborative initiative. Ann Rheum Dis. 2012;71(4):484-492. doi:10.1136/annrheumdis-2011-200329, PMID 22388996, Table 6
// Tags:    rheumatology, polymyalgia rheumatica, pmr, classification criteria, score
// Units:   each item present or absent; result points 0–6
// Range:   any combination of the four items
// Example: stiffness and absent RF/ACPA → 4; all four → 6; none → 0 (Table 6 weights)
func PMRScoreWithoutUltrasound(morningStiffnessOver45Min, hipPainOrLimitedROM, rfAndACPAAbsent, otherJointInvolvementAbsent bool) int {
	return pts(morningStiffnessOver45Min, 2) + pts(hipPainOrLimitedROM, 1) + pts(rfAndACPAAbsent, 2) +
		pts(otherJointInvolvementAbsent, 1)
}
