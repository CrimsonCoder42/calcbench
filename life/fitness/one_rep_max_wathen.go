package fitness

import "math"

// Wathen's constants: 1RM = 100 × load / (48.8 + 53.8 × e^(−0.075 × reps)).
const (
	wathenIntercept = 48.8
	wathenScale     = 53.8
	wathenRate      = 0.075
)

// OneRepMaxWathen returns the one-repetition maximum estimated from a load lifted for reps repetitions to
// fatigue by the Wathen equation, 1RM = 100 × load / (48.8 + 53.8 × e^(−0.075 × reps)), in the unit of the
// load. The denominator lies in (48.8, 102.6) for every reps ≥ 1, so it is never zero.
// Errors: ErrNotFinite for NaN or ±Inf; ErrOutOfRange when loadKgOrLb ≤ 0 or reps < 1; ErrDomain when the
// result is beyond float64.
//
// ID:      LIF-FIT-014
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Wathen D. Load assignment. In: Baechle TR, ed. Essentials of Strength Training and Conditioning. Champaign, IL: Human Kinetics; 1994:435-444 (pages per Beia ref. 17; Lin ref. 51 gives 435-439) (primary; book chapter, not opened)
// Source:  Lin TY, Hsieh SS, Chueh TY, Huang CJ, Hung TM. The effects of barbell resistance exercise on information processing speed and conflict-related ERP in older adults: a crossover randomized controlled trial. Sci Rep. 2021;11:9137. doi:10.1038/s41598-021-88634-5. PMID 33911153, PMC8080630, Table 1 "Wathan" [sic]: 1RM = 100 × weight ÷ (48.8 + 53.8e^(−0.075 × reps)); worked example 80 kg × 7 reps = 99 kg
// Source:  Beia R, Wassermann A, Raps S, Mayhew J, Uder M, Kemmler W. Developing Accurate Repetition Prediction Equations for Trained Older Adults with Osteopenia. Sports (Basel). 2024;12(9):233. doi:10.3390/sports12090233. PMID 39330710, PMC11435939, Table 1 "Wathen [17]": 1RM = 100 × rep wt/(48.8 + 53.8 × exp [−0.075 × reps]
// Verified: secondary (Lin 2021 Sci Rep Table 1, ref. 51 = Wathen 1994; Beia 2024 Sports Table 1, ref. 17 = Wathen 1994)
// Tags:    fitness, strength, one-rep-max, 1rm, wathen, reps-to-fatigue
// Units:   loadKgOrLb kg or lb (the result is in the same unit, D-063); reps repetitions to fatigue; result kg or lb
// Range:   loadKgOrLb > 0; reps ≥ 1 (no validated rep range is stated in the sources read)
// Example: 80 kg × 7 reps → 8000 / 80.6257 = 99.22 kg; Lin 2021 prints "= 99 kg"
func OneRepMaxWathen(loadKgOrLb float64, reps int) (float64, error) {
	if err := checkRepsToFatigue(loadKgOrLb, reps); err != nil {
		return 0, err
	}
	// For very large reps the exponential underflows to 0 and the denominator tends to 48.8; no special case.
	den := wathenIntercept + wathenScale*math.Exp(-wathenRate*float64(reps))
	v := 100 * loadKgOrLb / den
	if err := checkResult(v); err != nil {
		return 0, err
	}
	return v, nil
}
