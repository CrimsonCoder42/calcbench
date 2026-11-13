package fitness

import (
	"errors"
	"fmt"

	"github.com/CrimsonCoder42/calcbench/algo/reduce"
)

var (
	// ErrNotFinite is returned when an input is NaN or ±Inf.
	ErrNotFinite = errors.New("fitness: value is NaN or infinite")
	// ErrOutOfRange is returned when an input is outside the range the source defines or validated.
	ErrOutOfRange = errors.New("fitness: input outside the range the source validated")
	// ErrDomain is returned when the inputs are valid one by one but the result is beyond float64.
	ErrDomain = errors.New("fitness: inputs valid one by one but the result is undefined for them together")
)

// checkFinite returns ErrNotFinite, naming the parameter, when v is NaN or ±Inf.
func checkFinite(name string, v float64) error {
	if !reduce.AllFinite(v) {
		return fmt.Errorf("%w: %s = %v", ErrNotFinite, name, v)
	}
	return nil
}

// checkResult returns ErrDomain when a computed result overflowed to ±Inf or became NaN.
func checkResult(v float64) error {
	if !reduce.AllFinite(v) {
		return fmt.Errorf("%w: result = %v", ErrDomain, v)
	}
	return nil
}

// The ACSM position stand's "2-10% increase in load", as multipliers.
const (
	acsmIncreaseLow  = 1.02
	acsmIncreaseHigh = 1.10
)

// LoadIncreaseRange is the ACSM 2009 load range for the next session and whether the stand's condition
// ("one to two repetitions over the desired number") is met. LowLoad and HighLoad are in the unit of the load
// that was passed in.
type LoadIncreaseRange struct {
	LowLoad      float64 // current load × 1.02 (the 2% end)
	HighLoad     float64 // current load × 1.10 (the 10% end)
	ConditionMet bool    // repsOverTarget ≥ 1
}

// LoadProgressionACSM returns the 2–10% load-increase range of the ACSM 2009 position stand for a current
// load, and whether the stand's condition is met: the current workload was performed for at least one
// repetition over the desired number. The stand says "one to two"; a person who did more than two reps over
// also did one to two, so ConditionMet is repsOverTarget ≥ 1. The range is returned whether or not the
// condition is met. No rounding or plate step is applied: the stand states none.
// Errors: ErrNotFinite for NaN or ±Inf; ErrOutOfRange when currentLoadKgOrLb ≤ 0 or repsOverTarget < 0;
// ErrDomain when the high end is beyond float64.
//
// ID:      LIF-FIT-001
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  American College of Sports Medicine. American College of Sports Medicine position stand. Progression models in resistance training for healthy adults. Med Sci Sports Exerc. 2009;41(3):687-708. doi:10.1249/MSS.0b013e3181915670. PMID 19204579. Abstract: "2-10% increase in load be applied when the individual can perform the current workload for one to two repetitions over the desired number"
// Tags:    fitness, resistance-training, progression, load, acsm
// Units:   currentLoadKgOrLb kg or lb (the result is in the same unit); repsOverTarget repetitions; LowLoad, HighLoad same unit as the input
// Range:   currentLoadKgOrLb > 0; repsOverTarget ≥ 0 (population: healthy adults, the stand's title; no ages stated)
// Example: 225 lb, 1 rep over → 229.5 / 247.5 lb, condition met (HAND-DERIVED: 225 × 1.02, 225 × 1.10)
func LoadProgressionACSM(currentLoadKgOrLb float64, repsOverTarget int) (LoadIncreaseRange, error) {
	if err := checkFinite("currentLoadKgOrLb", currentLoadKgOrLb); err != nil {
		return LoadIncreaseRange{}, err
	}
	if currentLoadKgOrLb <= 0 {
		return LoadIncreaseRange{}, fmt.Errorf("%w: currentLoadKgOrLb = %v, want > 0", ErrOutOfRange, currentLoadKgOrLb)
	}
	if repsOverTarget < 0 {
		return LoadIncreaseRange{}, fmt.Errorf("%w: repsOverTarget = %d, want ≥ 0", ErrOutOfRange, repsOverTarget)
	}
	r := LoadIncreaseRange{
		LowLoad:      currentLoadKgOrLb * acsmIncreaseLow,
		HighLoad:     currentLoadKgOrLb * acsmIncreaseHigh,
		ConditionMet: repsOverTarget >= 1,
	}
	if err := checkResult(r.HighLoad); err != nil {
		return LoadIncreaseRange{}, err
	}
	return r, nil
}
