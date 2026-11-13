package fitness

import "fmt"

// The RIR-based RPE scale's top: "RPE-10 = 0-RIR" (Zourdos et al. 2016, abstract).
const rpeScaleTop = 10

// The lowest RPE with a single RIR number: Helms et al. 2016 state the scale uses "RIR descriptors for scores
// of 5–10" and that "scores of 5–6 were grouped as 4–6 RIR", so one RIR value exists for 7–10 only.
const rpeMinSingleRIR = 7

// RepsInReserve returns the repetitions in reserve for a whole-number rating of perceived exertion on the
// resistance-training RIR-based RPE scale: 10 − RPE ("RPE-10 = 0-RIR, RPE-9 = 1-RIR, and so forth").
// Errors: ErrOutOfRange when rpe is outside 7–10 (RPE 5–6 map to the range 4–6 RIR and 1–4 to no RIR).
//
// ID:      LIF-FIT-006
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Zourdos MC, Klemp A, Dolan C, et al. Novel resistance training-specific rating of perceived exertion scale measuring repetitions in reserve. J Strength Cond Res. 2016;30(1):267-275. doi:10.1519/JSC.0000000000001049. PMID 26049792. Abstract: "RPE-10 = 0-RIR, RPE-9 = 1-RIR, and so forth"
// Source:  Helms ER, Cronin J, Storey A, Zourdos MC. Application of the repetitions in reserve-based rating of perceived exertion scale for resistance training. Strength Cond J. 2016;38(4):42-49. doi:10.1519/SSC.0000000000000218. PMC4961270: "RIR descriptors for scores of 5–10"; "scores of 5–6 were grouped as 4–6 RIR"
// Tags:    fitness, resistance-training, rpe, rir, reps-in-reserve, effort
// Units:   rpe scale points (whole); result repetitions
// Range:   rpe 7–10, whole points (population: Zourdos 2016, 29 individuals, 24.0 ± 3.4 years, experienced and novice squatters, squat)
// Example: RPE 8 → 2 RIR (Helms 2016: "an RIR-based RPE of 8–10 (RIR 0–2)")
func RepsInReserve(rpe int) (int, error) {
	if rpe < rpeMinSingleRIR || rpe > rpeScaleTop {
		return 0, fmt.Errorf("%w: rpe = %d, want %d–%d", ErrOutOfRange, rpe, rpeMinSingleRIR, rpeScaleTop)
	}
	return rpeScaleTop - rpe, nil
}
