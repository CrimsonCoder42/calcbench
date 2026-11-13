package critical

// LRINECScore returns the Laboratory Risk Indicator for Necrotizing Fasciitis score, 0–13: C-reactive protein
// > 150 mg/L → 4; white cell count 15–25 ×10³/mm³ → 1, > 25 → 2; hemoglobin 11–13.5 g/dL → 1, < 11 → 2; sodium
// < 135 mmol/L → 2; creatinine > 1.6 mg/dL → 2; glucose > 180 mg/dL → 1. The two secondaries print the CRP row as
// "≤150 = 0, >150 = 4" and "<150 = 0, >150 = 4"; a CRP of exactly 150 scores 0 here (the first form).
// Errors: ErrNotFinite for NaN or ±Inf; ErrOutOfRange for a negative value.
//
// ID:      MED-CRIT-019
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Wong CH, Khin LW, Heng KS, Tan KC, Low CO. The LRINEC (Laboratory Risk Indicator for Necrotizing Fasciitis) score: a tool for distinguishing necrotizing fasciitis from other soft tissue infections. Crit Care Med. 2004;32(7):1535-1541. doi:10.1097/01.ccm.0000129486.35458.7d. PMID 15241098 (primary; not opened)
// Source:  Hoesl V, Kempa S, Prantl L, et al. The LRINEC Score-An Indicator for the Course and Prognosis of Necrotizing Fasciitis? J Clin Med. 2022;11(13):3583. doi:10.3390/jcm11133583. PMID 35806870. PMC9267597, Table 1 (after Wong 2004)
// Source:  El-Menyar A, Asim M, Mudali IN, et al. The laboratory risk indicator for necrotizing fasciitis (LRINEC) scoring: the diagnostic and potential prognostic role. Scand J Trauma Resusc Emerg Med. 2017;25(1):28. doi:10.1186/s13049-017-0359-z. PMID 28270166. PMC5341454, LRINEC table (after Wong 2004)
// Verified: secondary (Hoesl 2022 Table 1; El-Menyar 2017)
// Tags:    critical, infection, necrotizing fasciitis, lrinec, score
// Units:   crpMgL mg/L; wbcThousandPerMm3 ×10³/mm³; hemoglobinGDL g/dL; sodiumMmolL mmol/L; creatinineMgDL mg/dL; glucoseMgDL mg/dL; result points 0–13
// Range:   every input ≥ 0 (physical limits; adults with soft-tissue infection, the source's population)
// Example: CRP 200, WBC 20, Hb 12, Na 132, creatinine 1.8, glucose 150 → 4+1+1+2+2+0 = 10 (HAND-DERIVED from Hoesl Table 1)
func LRINECScore(crpMgL, wbcThousandPerMm3, hemoglobinGDL, sodiumMmolL, creatinineMgDL, glucoseMgDL float64) (int, error) {
	for _, e := range [...]error{
		finite("crpMgL", crpMgL), finite("wbcThousandPerMm3", wbcThousandPerMm3), finite("hemoglobinGDL", hemoglobinGDL),
		finite("sodiumMmolL", sodiumMmolL), finite("creatinineMgDL", creatinineMgDL), finite("glucoseMgDL", glucoseMgDL),
		atLeastZero("crpMgL", crpMgL), atLeastZero("wbcThousandPerMm3", wbcThousandPerMm3),
		atLeastZero("hemoglobinGDL", hemoglobinGDL), atLeastZero("sodiumMmolL", sodiumMmolL),
		atLeastZero("creatinineMgDL", creatinineMgDL), atLeastZero("glucoseMgDL", glucoseMgDL),
	} {
		if e != nil {
			return 0, e
		}
	}
	s := 0
	if crpMgL > 150 {
		s += 4
	}
	switch {
	case wbcThousandPerMm3 > 25:
		s += 2
	case wbcThousandPerMm3 >= 15:
		s++
	}
	switch {
	case hemoglobinGDL < 11:
		s += 2
	case hemoglobinGDL <= 13.5:
		s++
	}
	if sodiumMmolL < 135 {
		s += 2
	}
	if creatinineMgDL > 1.6 {
		s += 2
	}
	if glucoseMgDL > 180 {
		s++
	}
	return s, nil
}
