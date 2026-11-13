package respiratory

// FeverPAINScore returns the FeverPAIN score (0–5): "1 point each for: Fever during the last 24 h, Purulent
// tonsils, Attend rapidly (3 days or less), very Inflamed throat and No cough or coryza".
//
// ID:      MED-RESP-033
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Little P, Moore M, Hobbs FDR, et al. PRImary care Streptococcal Management (PRISM) study: identifying clinical variables associated with Lancefield group A β-haemolytic streptococci and Lancefield non-Group A streptococcal throat infections from two cohorts of patients presenting with an acute sore throat. BMJ Open. 2013;3(10):e003943. doi:10.1136/bmjopen-2013-003943. PMID 24163209, PMC3808825, Table 2 (FeverPAIN, 1 point each)
// Source:  Little P, Hobbs FDR, Moore M, et al. Clinical score and rapid antigen detection test to guide antibiotic use for sore throats: randomised controlled trial of PRISM (primary care streptococcal management). BMJ. 2013;347:f5806. doi:10.1136/bmj.f5806. PMID 24114306, PMC3805475, Methods (the five features of "score 2")
// Tags:    respiratory, pharyngitis, sore-throat, streptococcus, feverpain, score
// Units:   every item bool; result points 0–5
// Range:   any combination of items
// Example: fever, purulent tonsils, attended within 3 days → 3 (hand-derived: 1 point each)
func FeverPAINScore(feverLast24h, purulentTonsils, attendWithin3Days, veryInflamedThroat, noCoughOrCoryza bool) int {
	s := 0
	for _, present := range [5]bool{feverLast24h, purulentTonsils, attendWithin3Days, veryInflamedThroat, noCoughOrCoryza} {
		if present {
			s++
		}
	}
	return s
}
