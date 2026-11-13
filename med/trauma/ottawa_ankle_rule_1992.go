package trauma

// OttawaAnkle1992Findings are the findings of the 1992 derivation of the Ottawa ankle rules (Stiell 1992 abstract).
type OttawaAnkle1992Findings struct {
	PainNearMalleoli             bool // "pain near the malleoli"
	MalleolusPosteriorTenderness bool // "localized bone tenderness of the posterior edge or tip of either malleolus"
	UnableToBearWeight           bool // "unable to bear weight both immediately after the injury and in the ED"
	PainInMidfoot                bool // "pain in the midfoot"
	FifthMetatarsalBaseTender    bool // "bone tenderness at the base of the fifth metatarsal"
	CuboidTender                 bool // "the cuboid"
	NavicularTender              bool // "or the navicular"
}

// OttawaAnkle1992Result reports the two finding sets of the 1992 abstract.
type OttawaAnkle1992Result struct {
	MalleolarSet bool // pain near the malleoli and (age ≥ 55, posterior-edge/tip tenderness, or unable to bear weight)
	MidfootSet   bool // pain in the midfoot and bone tenderness at the base of the fifth metatarsal, the cuboid or the navicular
}

// OttawaAnkleRule1992 applies the finding sets of the 1992 derivation study of the Ottawa ankle rules: all
// significant malleolar fractures were among people with pain near the malleoli who were 55 years or older,
// had localized bone tenderness of the posterior edge or tip of either malleolus, or were unable to bear weight
// both immediately and in the ED; all significant midfoot fractures were among those with midfoot pain and
// bone tenderness at the base of the fifth metatarsal, the cuboid or the navicular. This is the 1992 derived
// form, not the 1993 refined rule.
// Errors: ErrNotFinite for NaN or ±Inf; ErrOutOfRange when ageYears is outside 0–130.
//
// ID:      MED-TRAUMA-007
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Stiell IG, Greenberg GH, McKnight RD, Nair RC, McDowell I, Worthington JR. A study to develop clinical decision rules for the use of radiography in acute ankle injuries. Ann Emerg Med. 1992;21(4):384-390. doi:10.1016/s0196-0644(05)82656-3. PMID 1554175. Abstract: "had pain near the malleoli and were age 55 years or more, had localized bone tenderness of the posterior edge or tip of either malleolus, or were unable to bear weight both immediately after the injury and in the ED"; "pain in the midfoot and bone tenderness at the base of the fifth metatarsal, the cuboid, or the navicular"
// Tags:    trauma, ankle, foot, ottawa ankle rules, decision rule, radiography
// Units:   ageYears years; result booleans
// Range:   0 ≤ ageYears ≤ 130 (physical limits; acute ankle injuries in adults, the study's population, which the caller establishes)
// Example: malleolar pain, age 56 → MalleolarSet; midfoot pain, cuboid tender → MidfootSet (from the abstract's sets)
func OttawaAnkleRule1992(ageYears float64, f OttawaAnkle1992Findings) (OttawaAnkle1992Result, error) {
	if err := checkAge(ageYears); err != nil {
		return OttawaAnkle1992Result{}, err
	}
	return OttawaAnkle1992Result{
		MalleolarSet: f.PainNearMalleoli && (ageYears >= 55 || f.MalleolusPosteriorTenderness || f.UnableToBearWeight),
		MidfootSet:   f.PainInMidfoot && (f.FifthMetatarsalBaseTender || f.CuboidTender || f.NavicularTender),
	}, nil
}
