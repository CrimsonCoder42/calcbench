package rheumatology

// PMRScoreWithUltrasound returns the 2012 EULAR/ACR polymyalgia rheumatica scoring-algorithm score with
// ultrasound, 0–8 (Table 6): the four clinical items as in PMRScoreWithoutUltrasound (2, 1, 2, 1) plus 1 for
// at least one shoulder with subdeltoid bursitis and/or biceps tenosynovitis and/or glenohumeral synovitis and
// at least one hip with synovitis and/or trochanteric bursitis, and 1 for both shoulders with subdeltoid
// bursitis, biceps tenosynovitis or glenohumeral synovitis. "A score of 5 or more is categorised as PMR in the
// algorithm with US".
//
// ID:      MED-RHEUM-003
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Dasgupta B, Cimmino MA, Maradit-Kremers H, et al. 2012 provisional classification criteria for polymyalgia rheumatica. Ann Rheum Dis. 2012;71(4):484-492. doi:10.1136/annrheumdis-2011-200329, PMID 22388996, Table 6
// Tags:    rheumatology, polymyalgia rheumatica, pmr, ultrasound, classification criteria, score
// Units:   each item present or absent; result points 0–8
// Range:   any combination of the six items
// Example: all six → 8; stiffness, absent RF/ACPA and both-shoulder US → 5 (Table 6 weights)
func PMRScoreWithUltrasound(morningStiffnessOver45Min, hipPainOrLimitedROM, rfAndACPAAbsent, otherJointInvolvementAbsent, usShoulderAndHip, usBothShoulders bool) int {
	return PMRScoreWithoutUltrasound(morningStiffnessOver45Min, hipPainOrLimitedROM, rfAndACPAAbsent, otherJointInvolvementAbsent) +
		pts(usShoulderAndHip, 1) + pts(usBothShoulders, 1)
}
