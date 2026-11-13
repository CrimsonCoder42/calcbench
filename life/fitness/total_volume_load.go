package fitness

import (
	"fmt"

	"github.com/CrimsonCoder42/calcbench/algo/reduce"
)

// LoadSet is one performed set: its repetitions and its external load in kg or lb (one unit for every set
// passed together; 0 = no external load).
type LoadSet struct {
	Reps       int
	LoadKgOrLb float64
}

// TotalVolumeLoad returns the volume load of a session or a week: the sum over the sets of each set's
// volume load, repetitions × external load (LIF-FIT-008, called per set). The sum is compensated
// (Neumaier). sets is read, never modified; no sets give 0.
// Errors: as LIF-FIT-008 for a set (ErrNotFinite, ErrOutOfRange for negative reps or load, wrapped with the
// set's index); ErrDomain when the total is beyond float64.
//
// ID:      LIF-FIT-011
// Time:    O(n)   n = len(sets)
// Space:   O(1), no allocation
// Source:  McBride JM, McCaulley GO, Cormie P, Nuzzo JL, Cavill MJ, Triplett NT. Comparison of methods to quantify volume during resistance exercise. J Strength Cond Res. 2009;23(1):106-110. doi:10.1519/jsc.0b013e31818efdfe. PMID 19130641. Abstract: "volume load (VL) (repetitions [no.] x external load [kg])"
// Tags:    fitness, resistance-training, volume, volume-load, tonnage, session
// Units:   Reps repetitions; LoadKgOrLb kg or lb; result kg·rep or lb·rep, the loads' unit
// Range:   Reps ≥ 0; LoadKgOrLb ≥ 0
// Example: 10 × 100 + 8 × 105 + 12 × 0 lb → 1000 + 840 + 0 = 1840 lb·rep (HAND-DERIVED)
func TotalVolumeLoad(sets []LoadSet) (float64, error) {
	var sum, comp float64
	for i, s := range sets {
		v, err := VolumeLoad(s.Reps, s.LoadKgOrLb)
		if err != nil {
			return 0, fmt.Errorf("sets[%d]: %w", i, err)
		}
		sum, comp = reduce.NeumaierAdd(sum, comp, v)
	}
	total := sum + comp
	if err := checkResult(total); err != nil {
		return 0, err
	}
	if err := checkResult(sum); err != nil {
		return 0, err
	}
	return total, nil
}
