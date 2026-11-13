package trauma

// NEXUSChestFindings are the NEXUS Chest criteria besides age (Rodriguez 2013 abstract).
type NEXUSChestFindings struct {
	RapidDeceleration         bool // "rapid deceleration mechanism"
	ChestPain                 bool // "chest pain"
	Intoxication              bool // "intoxication"
	AbnormalAlertness         bool // "abnormal alertness/mental status"
	DistractingPainfulInjury  bool // "distracting painful injury"
	ChestWallTendernessOnPalp bool // "tenderness to chest wall palpation"
}

// NEXUSChest returns the number of NEXUS Chest decision-instrument criteria present (0–7): age above 60 years,
// rapid deceleration mechanism, chest pain, intoxication, abnormal alertness/mental status, distracting painful
// injury and tenderness to chest wall palpation. A count of 0 is the instrument's negative result (none of the
// seven criteria present).
// Errors: ErrNotFinite for NaN or ±Inf; ErrOutOfRange when ageYears is outside 0–130.
//
// ID:      MED-TRAUMA-003
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Rodriguez RM, Anglin D, Langdorf MI, et al. NEXUS chest: validation of a decision instrument for selective chest imaging in blunt trauma. JAMA Surg. 2013;148(10):940-946. doi:10.1001/jamasurg.2013.2757. PMID 23925583. Abstract: "the NEXUS Chest 7 clinical criteria (age >60 years, rapid deceleration mechanism, chest pain, intoxication, abnormal alertness/mental status, distracting painful injury, and tenderness to chest wall palpation)"
// Source:  Rodriguez RM, Hendey GW, Mower W, et al. Derivation of a decision instrument for selective chest radiography in blunt trauma. J Trauma. 2011;71(3):549-553. doi:10.1097/TA.0b013e3181f2ac9d. PMID 21045745 (derivation; same seven items)
// Tags:    trauma, chest, nexus, blunt trauma, decision rule
// Units:   ageYears years; result count 0–7
// Range:   0 ≤ ageYears ≤ 130 (physical limits; blunt trauma, the instrument's population)
// Example: age 61, nothing else → 1; age 60, none → 0 (HAND-DERIVED from "age >60 years")
func NEXUSChest(ageYears float64, f NEXUSChestFindings) (int, error) {
	if err := checkAge(ageYears); err != nil {
		return 0, err
	}
	return count(ageYears > 60, f.RapidDeceleration, f.ChestPain, f.Intoxication, f.AbnormalAlertness,
		f.DistractingPainfulInjury, f.ChestWallTendernessOnPalp), nil
}
