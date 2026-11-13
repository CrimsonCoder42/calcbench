package scales

import "fmt"

// checkApgarSign returns ErrOutOfRange, wrapped with name, unless v is 0, 1 or 2.
func checkApgarSign(name string, v int) error {
	if v < 0 || v > 2 {
		return fmt.Errorf("%w: %s = %d, want 0, 1 or 2", ErrOutOfRange, name, v)
	}
	return nil
}

// ApgarScore returns the Apgar score, 0–10: the sum of the five signs (heart rate, respiratory effort, reflex
// irritability, muscle tone, color), each rated 0, 1 or 2 (Apgar 1953: "A rating of zero, one or two, was
// given to each sign"; "A score of ten indicated a baby in the best possible condition"). The paper rated the
// signs at "sixty seconds after the complete birth of the baby"; the function sums whatever ratings it is given.
// Errors: ErrOutOfRange when a sign is outside 0–2, naming the sign.
//
// ID:      MED-SCALE-005
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Apgar V. A proposal for a new method of evaluation of the newborn infant. Curr Res Anesth Analg. 1953;32(4):260-267. PMID 13083014. pp. 261–262 (read in the NLM scan)
// Tags:    scales, apgar, newborn, neonatal, score
// Units:   each sign points 0–2; result points 0–10
// Range:   heartRate, respiratoryEffort, reflexIrritability, muscleTone, color each 0–2
// Example: all five rated 2 → 10; all 0 → 0; 2, 1, 1, 1, 0 → 5 (HAND-DERIVED from pp. 261–262)
func ApgarScore(heartRate, respiratoryEffort, reflexIrritability, muscleTone, color int) (int, error) {
	if err := checkApgarSign("heartRate", heartRate); err != nil {
		return 0, err
	}
	if err := checkApgarSign("respiratoryEffort", respiratoryEffort); err != nil {
		return 0, err
	}
	if err := checkApgarSign("reflexIrritability", reflexIrritability); err != nil {
		return 0, err
	}
	if err := checkApgarSign("muscleTone", muscleTone); err != nil {
		return 0, err
	}
	if err := checkApgarSign("color", color); err != nil {
		return 0, err
	}
	return heartRate + respiratoryEffort + reflexIrritability + muscleTone + color, nil
}
