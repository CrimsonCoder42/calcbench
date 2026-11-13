package fitness

import "fmt"

// LoadAtPercentOneRepMax returns the training load that corresponds to a percent of the one-repetition
// maximum: oneRepMax × (percent / 100), in the unit of the one-repetition maximum. This is the meaning of the
// position stand's "loads corresponding to … % of 1 RM"; no rounding or plate step is applied.
// Errors: ErrNotFinite for NaN or ±Inf; ErrOutOfRange when oneRepMaxKgOrLb ≤ 0 or percentOneRepMax is
// outside (0, 100].
//
// ID:      LIF-FIT-007
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  American College of Sports Medicine. American College of Sports Medicine position stand. Progression models in resistance training for healthy adults. Med Sci Sports Exerc. 2009;41(3):687-708. doi:10.1249/MSS.0b013e3181915670. PMID 19204579. Abstract: loads stated as a percent of 1 RM ("60-70% of 1 RM")
// Source:  McBride JM, McCaulley GO, Cormie P, Nuzzo JL, Cavill MJ, Triplett NT. Comparison of methods to quantify volume during resistance exercise. J Strength Cond Res. 2009;23(1):106-110. doi:10.1519/jsc.0b013e31818efdfe. PMID 19130641. Abstract: protocols "at 75% of a 1-repetition maximum (1RM)"
// Tags:    fitness, strength, one-rep-max, 1rm, percent, training-load
// Units:   oneRepMaxKgOrLb kg or lb (the result is in the same unit); percentOneRepMax percent; result kg or lb
// Range:   oneRepMaxKgOrLb > 0; 0 < percentOneRepMax ≤ 100 (a load above 1 RM is not a percent of it)
// Example: 1 RM 200 lb at 75 % → 150 lb (HAND-DERIVED)
func LoadAtPercentOneRepMax(oneRepMaxKgOrLb, percentOneRepMax float64) (float64, error) {
	if err := checkFinite("oneRepMaxKgOrLb", oneRepMaxKgOrLb); err != nil {
		return 0, err
	}
	if err := checkFinite("percentOneRepMax", percentOneRepMax); err != nil {
		return 0, err
	}
	if oneRepMaxKgOrLb <= 0 {
		return 0, fmt.Errorf("%w: oneRepMaxKgOrLb = %v, want > 0", ErrOutOfRange, oneRepMaxKgOrLb)
	}
	if percentOneRepMax <= 0 || percentOneRepMax > 100 {
		return 0, fmt.Errorf("%w: percentOneRepMax = %v, want 0 < p ≤ 100", ErrOutOfRange, percentOneRepMax)
	}
	// percent / 100 is at most 1, so the product is at most oneRepMaxKgOrLb and cannot overflow.
	return oneRepMaxKgOrLb * (percentOneRepMax / 100), nil
}
