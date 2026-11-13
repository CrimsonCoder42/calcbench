package fitness

import "fmt"

// VolumeLoad returns the volume load of one set: repetitions × external load, in load-unit repetitions
// (kg·rep or lb·rep). A session or a week is the sum of its sets' volume loads; this is the per-set term.
// Errors: ErrNotFinite for NaN or ±Inf; ErrOutOfRange when reps < 0 or loadKgOrLb < 0; ErrDomain when the
// product is beyond float64.
//
// ID:      LIF-FIT-008
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  McBride JM, McCaulley GO, Cormie P, Nuzzo JL, Cavill MJ, Triplett NT. Comparison of methods to quantify volume during resistance exercise. J Strength Cond Res. 2009;23(1):106-110. doi:10.1519/jsc.0b013e31818efdfe. PMID 19130641. Abstract: "volume load (VL) (repetitions [no.] x external load [kg])"
// Tags:    fitness, resistance-training, volume, volume-load, tonnage
// Units:   reps repetitions; loadKgOrLb kg or lb (external load); result kg·rep or lb·rep, the load's unit
// Range:   reps ≥ 0; loadKgOrLb ≥ 0 (0 = no external load)
// Example: 10 reps × 135 lb → 1350 lb·rep (HAND-DERIVED)
func VolumeLoad(reps int, loadKgOrLb float64) (float64, error) {
	if err := checkFinite("loadKgOrLb", loadKgOrLb); err != nil {
		return 0, err
	}
	if reps < 0 {
		return 0, fmt.Errorf("%w: reps = %d, want ≥ 0", ErrOutOfRange, reps)
	}
	if loadKgOrLb < 0 {
		return 0, fmt.Errorf("%w: loadKgOrLb = %v, want ≥ 0", ErrOutOfRange, loadKgOrLb)
	}
	v := float64(reps) * loadKgOrLb
	if err := checkResult(v); err != nil {
		return 0, err
	}
	return v, nil
}
