package respiratory

// DRIPFindings are the risk factors of the Drug Resistance in Pneumonia score (Webb 2016, Table 3), worded as
// in the table.
type DRIPFindings struct {
	// Major risk factors, 2 points each.
	AntibioticsWithin60Days  bool // "Antibiotic use within previous 60 days"
	LongTermCareResidence    bool // "Residence in a long-term-care facility"
	TubeFeeding              bool // "Tube feeding"
	PriorDRPInfectionWithin1 bool // "Prior infection with a DRP (1 yr)"
	// Minor risk factors, 1 point each.
	HospitalizationWithin60Days bool // "Hospitalization within previous 60 days"
	ChronicPulmonaryDisease     bool // "Chronic pulmonary disease"
	PoorFunctionalStatus        bool // "Poor functional status"
	GastricAcidSuppression      bool // "Gastric acid suppression"
	WoundCare                   bool // "Wound care"
	MRSAColonizationWithin1     bool // "MRSA colonization (1 yr)"
}

// DRIPScore returns the Drug Resistance in Pneumonia (DRIP) score (0–14): 2 points for each major risk factor
// present and 1 point for each minor one. The paper's threshold: "At a threshold of ≥4 points, the DRIP score
// differentiated a high risk versus a low risk of pneumonia due to drug-resistant pathogens"; that comparison is
// the caller's.
//
// ID:      MED-RESP-032
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Webb BJ, Dascomb K, Stenehjem E, et al. Derivation and multicenter validation of the drug resistance in pneumonia clinical prediction score. Antimicrob Agents Chemother. 2016;60(5):2652-2663. doi:10.1128/AAC.03071-15. PMID 26856838, PMC4862530, Table 3 (major 2 points, minor 1 point, total 14)
// Tags:    respiratory, pneumonia, drug-resistance, drip, score
// Units:   every item bool; result points 0–14
// Range:   any combination of items
// Example: antibiotics within 60 days, chronic pulmonary disease, gastric acid suppression → 2 + 1 + 1 = 4 (hand-derived from Table 3)
func DRIPScore(f DRIPFindings) int {
	s := 0
	for _, major := range [4]bool{f.AntibioticsWithin60Days, f.LongTermCareResidence, f.TubeFeeding, f.PriorDRPInfectionWithin1} {
		if major {
			s += 2
		}
	}
	for _, minor := range [6]bool{f.HospitalizationWithin60Days, f.ChronicPulmonaryDisease, f.PoorFunctionalStatus,
		f.GastricAcidSuppression, f.WoundCare, f.MRSAColonizationWithin1} {
		if minor {
			s++
		}
	}
	return s
}
