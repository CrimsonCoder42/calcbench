package gastro

// RansonFindings holds Ranson's eleven early objective signs for acute pancreatitis (non-gallstone criteria),
// five at admission and six during the initial 48 h, each true when the criterion is met. Each field names
// the source's threshold.
type RansonFindings struct {
	AgeOver55Years                bool // at admission: age > 55 years
	WBCOver16000PerMm3            bool // at admission: white blood count > 16,000/mm³
	GlucoseOver200MgDL            bool // at admission: blood glucose > 200 mg/dL
	LDHOver350UL                  bool // at admission: serum LDH > 350 U/L
	ASTOver250UL                  bool // at admission: serum AST > 250 U/L
	HematocritFallOver10Points    bool // within 48 h: hematocrit fall > 10 percentage points
	BUNRiseOver5MgDL              bool // within 48 h: blood urea nitrogen rise > 5 mg/dL
	PaO2Below60MmHg               bool // within 48 h: arterial PO₂ < 60 mmHg
	CalciumBelow8MgDL             bool // within 48 h: serum calcium < 8 mg/dL
	BaseDeficitOver4MEqL          bool // within 48 h: base deficit > 4 mEq/L
	FluidSequestrationOver6Liters bool // within 48 h: estimated fluid sequestration > 6 L
}

// RansonScore returns Ranson's score for acute pancreatitis, 0–11: one point for each of the eleven early
// objective signs (non-gallstone criteria) that is met. No band is returned.
// Errors: none; every combination of the booleans is valid.
//
// ID:      MED-GASTRO-009
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Ranson JH, Rifkind KM, Roses DF, Fink SD, Eng K, Spencer FC. Prognostic signs and the role of operative management in acute pancreatitis. Surg Gynecol Obstet. 1974;139(1):69-81. PMID 4834279 (full text not opened; abstract states no criteria)
// Source:  Lee DW, Cho CM. Predicting severity of acute pancreatitis. Medicina (Kaunas). 2022. PMID 35744050. PMC9227091, Table 4 "Ranson score in acute pancreatitis (non-biliary and biliary)", non-biliary column; "Each factor 1 point (total 0–11 points)"
// Source:  Abu-Eshy SA, Abolfotouh MA, Nawar E, Abu Sabib AR. Ranson's criteria for acute pancreatitis in high altitude: do they need to be modified? Saudi J Gastroenterol. 2008. PMID 19568489. PMC2702891, Table 1 "Ranson's criteria: The 11 early objective signs used to classify the severity of acute pancreatitis" [ref 6, Ranson 1974]
// Verified: secondary (Lee 2022 PMC9227091; Abu-Eshy 2008 PMC2702891)
// Tags:    gastro, pancreas, acute-pancreatitis, ranson, severity, score
// Units:   eleven booleans; result points 0–11
// Range:   adults with acute pancreatitis, non-gallstone criteria (the gallstone thresholds of Lee 2022 Table 4 are not applied)
// Example: age > 55, WBC > 16,000, glucose > 200, calcium < 8 → 4 (hand-derived)
func RansonScore(f RansonFindings) int {
	points := 0
	for _, met := range [...]bool{
		f.AgeOver55Years, f.WBCOver16000PerMm3, f.GlucoseOver200MgDL, f.LDHOver350UL, f.ASTOver250UL,
		f.HematocritFallOver10Points, f.BUNRiseOver5MgDL, f.PaO2Below60MmHg, f.CalciumBelow8MgDL,
		f.BaseDeficitOver4MEqL, f.FluidSequestrationOver6Liters,
	} {
		if met {
			points++
		}
	}
	return points
}
