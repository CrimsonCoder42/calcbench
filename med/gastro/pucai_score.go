package gastro

import "fmt"

// PUCAIFindings holds the six PUCAI items (Moon 2019 Table 3; Fell 2016 Table 1). Level fields are the
// row index of the item's table; each comment gives the cell text and its points.
type PUCAIFindings struct {
	AbdominalPain    int  // 0 "No pain" 0; 1 "Pain can be ignored" 5; 2 "Pain cannot be ignored" 10
	RectalBleeding   int  // 0 "None" 0; 1 "Small amount only, in <50% of stools" 10; 2 "Small amount with most stools" 20; 3 "Large amount (>50% of stools)" 30
	StoolConsistency int  // consistency of most stools: 0 "Formed" 0; 1 "Partially formed" 5; 2 "Completely unformed" 10
	StoolsPer24h     int  // number of stools per 24 h: 0–2: 0; 3–5: 5; 6–8: 10; > 8: 15
	NocturnalStools  bool // "Nocturnal stools (any episode causing awakening)": yes 10
	ActivityLevel    int  // 0 "No limitation" 0; 1 "Occasional limitation" 5; 2 "Severe restricted activity" 10
}

// PUCAIScore returns the Pediatric Ulcerative Colitis Activity Index, 0–85, in steps of 5.
// Errors: ErrOutOfRange for a level outside its row count or a negative stool count.
//
// ID:      MED-GASTRO-014
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Turner D, Otley AR, Mack D, et al. Development, validation, and evaluation of a pediatric ulcerative colitis activity index: a prospective multicenter study. Gastroenterology. 2007;133(2):423-432. doi:10.1053/j.gastro.2007.05.029, PMID 17681163 (primary; not opened)
// Source:  Moon JS. Clinical aspects and treatments for pediatric inflammatory bowel diseases. Intest Res. 2019;17(1):17-23. doi:10.5217/ir.2018.00139, PMID 30625263, PMC6361015, Table 3 (PUCAI)
// Source:  Fell JM, et al. Arch Dis Child. 2016, PMC4853583, Table 1 "Paediatric Ulcerative Colitis Activity Index (PUCAI)" (all items and points; ref 8 Turner 2007)
// Verified: secondary (Moon 2019 PMC6361015 Table 3, ref 29 Turner 2007; Fell 2016 PMC4853583 Table 1, ref 8 Turner 2007)
// Tags:    gastro, pediatric, ulcerative colitis, ibd, pucai, score
// Units:   levels are row indexes; StoolsPer24h count; result points 0–85
// Range:   children and adolescents with ulcerative colitis (the source population); any valid levels
// Example: pain ignorable, small blood <50 %, partially formed, 4 stools, no nocturnal, no limitation → 5+10+5+5 = 25 (hand-derived from Table 3)
func PUCAIScore(f PUCAIFindings) (int, error) {
	for _, in := range [...]struct {
		name   string
		v, max int
	}{{"AbdominalPain", f.AbdominalPain, 2}, {"RectalBleeding", f.RectalBleeding, 3}, {"StoolConsistency", f.StoolConsistency, 2}, {"ActivityLevel", f.ActivityLevel, 2}} {
		if in.v < 0 || in.v > in.max {
			return 0, fmt.Errorf("%w: %s = %d, want 0–%d", ErrOutOfRange, in.name, in.v, in.max)
		}
	}
	if f.StoolsPer24h < 0 {
		return 0, fmt.Errorf("%w: StoolsPer24h = %d, want ≥ 0", ErrOutOfRange, f.StoolsPer24h)
	}
	s := 5*f.AbdominalPain + 10*f.RectalBleeding + 5*f.StoolConsistency + 5*f.ActivityLevel
	switch n := f.StoolsPer24h; {
	case n > 8:
		s += 15
	case n >= 6:
		s += 10
	case n >= 3:
		s += 5
	}
	if f.NocturnalStools {
		s += 10
	}
	return s, nil
}
