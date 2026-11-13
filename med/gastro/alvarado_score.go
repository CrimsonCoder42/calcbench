package gastro

// AlvaradoFindings holds the eight Alvarado (MANTRELS) items as the caller has judged them. The definitions
// are those Alvarado restates in 2016: elevated temperature "oral temperature of 37.3 C or more";
// leukocytosis "more than 10,000/ml"; shift to the left "stabs or bands (more than 5 %) or increased segmented
// neutrophils (more than 75 %)".
type AlvaradoFindings struct {
	MigrationOfPain     bool // migration of pain to the right lower quadrant: 1
	Anorexia            bool // anorexia: 1
	NauseaOrVomiting    bool // nausea or vomiting: 1
	TendernessRLQ       bool // tenderness in the right lower quadrant: 2
	ReboundPain         bool // rebound pain: 1
	ElevatedTemperature bool // oral temperature ≥ 37.3 °C: 1
	Leukocytosis        bool // white cells > 10,000/mL: 2
	ShiftToLeft         bool // bands > 5 % or segmented neutrophils > 75 %: 1
}

// AlvaradoScore returns the Alvarado score, 0–10: "2 points to each of the two of the more important factors
// (tenderness and leukocytosis), and a value of 1 to each one of the others, for a possible total score of 10".
// Categories are not built here (they were seen only in the author's 2016 restatement).
//
// ID:      MED-GASTRO-011
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Alvarado A. A practical score for the early diagnosis of acute appendicitis. Ann Emerg Med. 1986;15(5):557-564. doi:10.1016/s0196-0644(86)80993-3, PMID 3963537 (primary; not opened)
// Source:  Prabhu BL Y, Mehra B, Ghoshal S, Dubhashi SP. Cureus. 2024, PMC11087879, Table 1 (the eight items and points, 0–10; ref 5 Alvarado 1986)
// Source:  Jose T, Rajesh PS. Appendicitis Inflammatory Response Score in Comparison to Alvarado Score in Acute Appendicitis. Surg J (N Y). 2021;7(3):e127-e131. doi:10.1055/s-0041-1731446, PMID 34295969, PMC8289675, Methods (item weights)
// Verified: secondary (Jose & Rajesh 2021 PMC8289675, ref 7 Alvarado 1986; Prabhu 2024 PMC11087879, ref 5 Alvarado 1986)
// Tags:    gastro, appendicitis, alvarado, mantrels, score
// Units:   every item bool; result points 0–10
// Range:   any combination of items
// Example: tenderness, leukocytosis, migration → 2 + 2 + 1 = 5 (hand-derived from the 2016 weights)
func AlvaradoScore(f AlvaradoFindings) int {
	s := 0
	for _, one := range [...]bool{f.MigrationOfPain, f.Anorexia, f.NauseaOrVomiting, f.ReboundPain, f.ElevatedTemperature, f.ShiftToLeft} {
		if one {
			s++
		}
	}
	if f.TendernessRLQ {
		s += 2
	}
	if f.Leukocytosis {
		s += 2
	}
	return s
}
