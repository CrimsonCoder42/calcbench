package fitness

import "fmt"

// Brzycki's constants as printed in the secondaries: 1RM = load / (1.0278 − 0.0278 × reps).
const (
	brzyckiIntercept = 1.0278
	brzyckiSlope     = 0.0278
	// brzyckiMaxReps is the largest whole rep count with a positive denominator
	// (1.0278 − 0.0278 × 36 = 0.0270; at 37 it is −0.0008).
	brzyckiMaxReps = 36
)

// checkRepsToFatigue checks the inputs every repetitions-to-fatigue 1 RM estimate shares: a finite load
// above zero and at least one repetition.
func checkRepsToFatigue(loadKgOrLb float64, reps int) error {
	if err := checkFinite("loadKgOrLb", loadKgOrLb); err != nil {
		return err
	}
	if loadKgOrLb <= 0 {
		return fmt.Errorf("%w: loadKgOrLb = %v, want > 0", ErrOutOfRange, loadKgOrLb)
	}
	if reps < 1 {
		return fmt.Errorf("%w: reps = %d, want ≥ 1", ErrOutOfRange, reps)
	}
	return nil
}

// OneRepMaxBrzycki returns the one-repetition maximum estimated from a load lifted for reps repetitions to
// fatigue by the Brzycki equation, 1RM = load / (1.0278 − 0.0278 × reps), in the unit of the load.
// At one repetition the estimate is the load itself (the denominator is exactly 1.0000).
// Errors: ErrNotFinite for NaN or ±Inf; ErrOutOfRange when loadKgOrLb ≤ 0, reps < 1 or reps > 36 (the
// denominator is not positive from 37 on); ErrDomain when the result is beyond float64.
//
// ID:      LIF-FIT-012
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Brzycki M. Strength testing—predicting a one-rep max from reps-to-fatigue. J Phys Educ Recreat Dance. 1993;64(1):88-90. doi:10.1080/07303084.1993.10606684 (primary; not opened: publisher 403, not in PubMed)
// Source:  Lin TY, Hsieh SS, Chueh TY, Huang CJ, Hung TM. The effects of barbell resistance exercise on information processing speed and conflict-related ERP in older adults: a crossover randomized controlled trial. Sci Rep. 2021;11:9137. doi:10.1038/s41598-021-88634-5. PMID 33911153, PMC8080630, Table 1 "Brzycki": 1RM = 100 × weight ÷ (102.78 − 2.78 × reps)
// Source:  Beia R, Wassermann A, Raps S, Mayhew J, Uder M, Kemmler W. Developing Accurate Repetition Prediction Equations for Trained Older Adults with Osteopenia. Sports (Basel). 2024;12(9):233. doi:10.3390/sports12090233. PMID 39330710, PMC11435939, Table 1 "Brzycki [12]": 1RM = rep weight/(1.0278 − (0.0278 × reps))
// Verified: secondary (Lin 2021 Sci Rep Table 1; Beia 2024 Sports Table 1, ref. 12 = Brzycki 1993)
// Tags:    fitness, strength, one-rep-max, 1rm, brzycki, reps-to-fatigue
// Units:   loadKgOrLb kg or lb (the result is in the same unit, D-063); reps repetitions to fatigue; result kg or lb
// Range:   loadKgOrLb > 0; reps 1–36 (no validated rep range is stated in the sources read; 36 is the formula's own limit)
// Example: 100 lb × 10 reps → 100 / 0.7498 = 133.369 lb (HAND-DERIVED from the equation)
func OneRepMaxBrzycki(loadKgOrLb float64, reps int) (float64, error) {
	if err := checkRepsToFatigue(loadKgOrLb, reps); err != nil {
		return 0, err
	}
	if reps > brzyckiMaxReps {
		return 0, fmt.Errorf("%w: reps = %d, want ≤ %d (denominator 1.0278 − 0.0278 × reps must be > 0)",
			ErrOutOfRange, reps, brzyckiMaxReps)
	}
	v := loadKgOrLb / (brzyckiIntercept - brzyckiSlope*float64(reps))
	if err := checkResult(v); err != nil {
		return 0, err
	}
	return v, nil
}
