package infectious

// MenzaScore returns the Menza score for HIV acquisition among men who have sex with men (simple model,
// Table 2), 0–19: diagnosis or history of bacterial sexually transmitted infection at baseline 4;
// methamphetamine or inhaled nitrites in the prior 6 months 11; ≥ 10 male sex partners in the prior year 3;
// non-concordant unprotected anal sex in the prior year 1.
//
// ID:      MED-INFECT-002
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Menza TW, Hughes JP, Celum CL, Golden MR. Prediction of HIV acquisition among men who have sex with men. Sex Transm Dis. 2009;36(9):547-555. doi:10.1097/OLQ.0b013e3181a9cc41, PMID 19707108, Table 2
// Tags:    infectious, hiv, msm, menza, score
// Units:   each item present or absent; result points 0–19
// Range:   any combination of the four items
// Example: STI and methamphetamine → 15; all four → 19; none → 0 (Table 2)
func MenzaScore(bacterialSTI, methamphetamineOrNitrites, tenOrMoreMalePartners, nonConcordantUnprotectedAnalSex bool) int {
	s := 0
	if bacterialSTI {
		s += 4
	}
	if methamphetamineOrNitrites {
		s += 11
	}
	if tenOrMoreMalePartners {
		s += 3
	}
	if nonConcordantUnprotectedAnalSex {
		s++
	}
	return s
}
