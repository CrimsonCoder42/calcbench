package psychiatry

import "fmt"

// COWSItems are the points of the eleven Clinical Opiate Withdrawal Scale items. Each item takes only the
// values the scale prints for it (Wesson & Ling 2003, Appendix 1):
//
//	RestingPulse 0, 1, 2, 4            Sweating 0, 1, 2, 3, 4          Restlessness 0, 1, 3, 5
//	PupilSize 0, 1, 2, 5               BoneJointAches 0, 1, 2, 4       RunnyNoseTearing 0, 1, 2, 4
//	GIUpset 0, 1, 2, 3, 5              Tremor 0, 1, 2, 4               Yawning 0, 1, 2, 4
//	AnxietyIrritability 0, 1, 2, 4     GoosefleshSkin 0, 3, 5
type COWSItems struct {
	RestingPulse, Sweating, Restlessness, PupilSize, BoneJointAches, RunnyNoseTearing int
	GIUpset, Tremor, Yawning, AnxietyIrritability, GoosefleshSkin                     int
}

// COWSScore returns the Clinical Opiate Withdrawal Scale total, 0–48: "the sum of all 11 items".
// Errors: ErrOutOfRange when an item takes a value the scale does not print for it, naming the item.
//
// ID:      MED-PSYCH-016
// Time:    O(1)   11 items, fixed
// Space:   O(1), no allocation
// Source:  Wesson DR, Ling W. The Clinical Opiate Withdrawal Scale (COWS). J Psychoactive Drugs. 2003;35(2):253-259. doi:10.1080/02791072.2003.10400007, PMID 12924748, Appendix 1 (item points; "The total score is the sum of all 11 items"), read in the reprint NIDA distributes: https://nida.nih.gov/sites/default/files/ClinicalOpiateWithdrawalScale.pdf (read 2026-11-12, SHA-256 4c04fc126ae8252932d55caf2a9ef0bf89df3ee5cbaa5b56bd266f260a25ebf9)
// Tags:    psychiatry, opioid, withdrawal, cows, score
// Units:   item points as printed; result points 0–48
// Range:   each item one of its printed values (COWSItems)
// Example: pulse 2, sweating 1, restlessness 1, pupils 1, aches 1, nose 2, GI 2, tremor 1, yawning 1, anxiety 1, gooseflesh 0 → 13 (HAND-DERIVED sum)
func COWSScore(it COWSItems) (int, error) {
	const (
		p0124  = 1<<0 | 1<<1 | 1<<2 | 1<<4
		p01234 = 1<<0 | 1<<1 | 1<<2 | 1<<3 | 1<<4
		p0135  = 1<<0 | 1<<1 | 1<<3 | 1<<5
		p0125  = 1<<0 | 1<<1 | 1<<2 | 1<<5
		p01235 = 1<<0 | 1<<1 | 1<<2 | 1<<3 | 1<<5
		p035   = 1<<0 | 1<<3 | 1<<5
	)
	items := [...]struct {
		name    string
		v       int
		allowed uint8
	}{
		{"RestingPulse", it.RestingPulse, p0124}, {"Sweating", it.Sweating, p01234},
		{"Restlessness", it.Restlessness, p0135}, {"PupilSize", it.PupilSize, p0125},
		{"BoneJointAches", it.BoneJointAches, p0124}, {"RunnyNoseTearing", it.RunnyNoseTearing, p0124},
		{"GIUpset", it.GIUpset, p01235}, {"Tremor", it.Tremor, p0124}, {"Yawning", it.Yawning, p0124},
		{"AnxietyIrritability", it.AnxietyIrritability, p0124}, {"GoosefleshSkin", it.GoosefleshSkin, p035},
	}
	total := 0
	for _, x := range items {
		if x.v < 0 || x.v > 5 || x.allowed&(1<<uint(x.v)) == 0 {
			return 0, fmt.Errorf("%w: %s = %d, not a value the scale prints for this item", ErrOutOfRange, x.name, x.v)
		}
		total += x.v
	}
	return total, nil
}

// COWSSeverity is a COWS withdrawal band; its value is the source's words, verbatim.
type COWSSeverity string

// Wesson & Ling 2003, Appendix 1: "Score: 5-12 = mild; 13-24 = moderate; 25-36 = moderately severe; more than 36
// = severe withdrawal". The scale prints no band for 0–4.
const (
	COWSNoBand           COWSSeverity = ""                  // 0–4: no band printed
	COWSMild             COWSSeverity = "mild"              // 5–12
	COWSModerate         COWSSeverity = "moderate"          // 13–24
	COWSModeratelySevere COWSSeverity = "moderately severe" // 25–36
	COWSSevere           COWSSeverity = "severe withdrawal" // 37–48
)

// COWSSeverityOf returns the COWS band for a total (0–48): 5–12, 13–24, 25–36, more than 36; COWSNoBand ("") for
// 0–4, where the scale prints none.
// Errors: ErrOutOfRange when total is outside 0–48.
//
// ID:      MED-PSYCH-017
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Wesson DR, Ling W. The Clinical Opiate Withdrawal Scale (COWS). J Psychoactive Drugs. 2003;35(2):253-259. doi:10.1080/02791072.2003.10400007, PMID 12924748, Appendix 1 (the bands, quoted), read in the NIDA reprint (see COWSScore)
// Tags:    psychiatry, opioid, withdrawal, cows, category
// Units:   total points 0–48; result the source's words
// Range:   total 0–48
// Example: 4 → ""; 5 → "mild"; 13 → "moderate"; 25 → "moderately severe"; 37 → "severe withdrawal"
func COWSSeverityOf(total int) (COWSSeverity, error) {
	if err := checkPoints("total", total, 0, 48); err != nil {
		return "", err
	}
	switch {
	case total <= 4:
		return COWSNoBand, nil
	case total <= 12:
		return COWSMild, nil
	case total <= 24:
		return COWSModerate, nil
	case total <= 36:
		return COWSModeratelySevere, nil
	}
	return COWSSevere, nil
}
