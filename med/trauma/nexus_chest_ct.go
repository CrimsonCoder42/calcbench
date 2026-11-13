package trauma

// NEXUSChestCTFindings are the seven criteria of the NEXUS Chest CT decision instruments (Rodriguez 2015).
type NEXUSChestCTFindings struct {
	AbnormalCXR             bool // "(1) abnormal CXR"
	RapidDeceleration       bool // "(2) rapid deceleration mechanism" (Chest CT-All only)
	DistractingInjury       bool // "(3) distracting injury"
	ChestWallTenderness     bool // "(4) chest wall tenderness"
	SternalTenderness       bool // "(5) sternal tenderness"
	ThoracicSpineTenderness bool // "(6) thoracic spine tenderness"
	ScapularTenderness      bool // "(7) scapular tenderness"
}

// NEXUSChestCTResult reports, for each of the two instruments, whether any of its criteria is present.
type NEXUSChestCTResult struct {
	AnyChestCTAll   bool // any of the seven Chest CT-All criteria
	AnyChestCTMajor bool // any of the six Chest CT-Major criteria ("the same criteria without rapid deceleration mechanism")
}

// NEXUSChestCT applies the two NEXUS Chest CT decision instruments: Chest CT-All (seven criteria) and Chest
// CT-Major (the same without rapid deceleration mechanism). Each instrument's negative result is "not ... any
// DI criteria" present (false here).
//
// ID:      MED-TRAUMA-004
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Rodriguez RM, Langdorf MI, Nishijima D, et al. Derivation and validation of two decision instruments for selective chest CT in blunt trauma: a multicenter prospective observational study (NEXUS Chest CT). PLoS Med. 2015;12(10):e1001883. doi:10.1371/journal.pmed.1001883. PMID 26440607, PMC4595216 (abstract read from the full-text XML): "(1) abnormal CXR, (2) rapid deceleration mechanism, (3) distracting injury, (4) chest wall tenderness, (5) sternal tenderness, (6) thoracic spine tenderness, and (7) scapular tenderness"; Chest CT-Major: "the same criteria without rapid deceleration mechanism"
// Tags:    trauma, chest, ct, nexus, blunt trauma, decision rule
// Units:   result booleans
// Range:   any combination of findings (blunt trauma, the instruments' population)
// Example: rapid deceleration only → AnyChestCTAll true, AnyChestCTMajor false (from the instruments' definitions)
func NEXUSChestCT(f NEXUSChestCTFindings) NEXUSChestCTResult {
	major := f.AbnormalCXR || f.DistractingInjury || f.ChestWallTenderness || f.SternalTenderness ||
		f.ThoracicSpineTenderness || f.ScapularTenderness
	return NEXUSChestCTResult{AnyChestCTAll: major || f.RapidDeceleration, AnyChestCTMajor: major}
}
