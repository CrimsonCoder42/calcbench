package trig

import (
	"fmt"
	"math"

	"github.com/CrimsonCoder42/calcbench/algo/reduce"
)

// UnwrapAngles returns the sequence anglesRad with its period jumps removed (phase unwrapping): the first
// angle is kept, and each later one is the previous unwrapped angle plus the step between the two inputs
// reduced to the interval [−P/2, P/2] by a multiple of P, P = periodRad (2π for radians in a circle). A
// step already within (−P/2, P/2) is kept exactly; a step of exactly ±P/2 keeps its sign. The result
// is correct when the true signal changes by less than P/2 between samples (Itoh's condition).
// The result is written into dst when cap(dst) ≥ len(anglesRad), else a new slice is allocated; dst may
// be anglesRad itself (in place). anglesRad is otherwise read, never modified. Empty in, empty out.
// Errors: ErrNotFinite when an angle or periodRad is NaN or ±Inf; ErrDomain when periodRad ≤ 0;
// ErrOverflow when an unwrapped angle is beyond the float64 range.
//
// ID:      NUM-TRIG-032
// Time:    O(n)   n = len(anglesRad); one pass
// Space:   O(1) extra; O(n) for the result unless dst has the capacity (then no allocation)
// Source:  Itoh, "Analysis of the phase unwrapping algorithm", Applied Optics, 1982, 21(14):2470, doi:10.1364/AO.21.002470 (eq. 4: unwrapped = first sample + Σ of wrapped differences)
// Tags:    trig, angle, phase, unwrap, signal, sequence
// Units:   anglesRad and periodRad in one angular unit (radians: periodRad = 2π); result in that unit
// Example: 0, 3π/4, −3π/4 (= 5π/4 wrapped), −π/4 (= 7π/4) with P = 2π → 0, 3π/4, 5π/4, 7π/4 (hand-checked)
func UnwrapAngles(dst, anglesRad []float64, periodRad float64) ([]float64, error) {
	if !reduce.AllFinite(periodRad) {
		return nil, fmt.Errorf("%w: periodRad", ErrNotFinite)
	}
	if !(periodRad > 0) {
		return nil, fmt.Errorf("%w: periodRad = %v, want > 0", ErrDomain, periodRad)
	}
	n := len(anglesRad)
	for i, a := range anglesRad {
		if !reduce.AllFinite(a) {
			return nil, fmt.Errorf("%w: anglesRad[%d] = %v", ErrNotFinite, i, a)
		}
	}
	if cap(dst) < n {
		dst = make([]float64, n)
	}
	dst = dst[:n]
	if n == 0 {
		return dst, nil
	}
	half := periodRad / 2
	prevIn := anglesRad[0]
	acc := prevIn
	dst[0] = acc
	for i := 1; i < n; i++ {
		cur := anglesRad[i] // read before writing dst[i]: dst may alias anglesRad
		step := cur - prevIn
		if math.IsInf(step, 0) {
			step = cur/2 - prevIn/2                            // the true step is 2·this; halves cannot overflow
			step = wrapStep(step, half) + wrapStep(step, half) // ≡ 2·step mod P, in [−P, P]
		}
		acc += wrapStep(step, half)
		if math.IsInf(acc, 0) {
			return nil, fmt.Errorf("%w: unwrapped anglesRad[%d]", ErrOverflow, i)
		}
		dst[i] = acc
		prevIn = cur
	}
	return dst, nil
}

// wrapStep returns step reduced by a multiple of 2·half into [−half, half]: unchanged when already
// strictly inside, and a reduced value of exactly −half becomes +half when step > 0.
func wrapStep(step, half float64) float64 {
	if math.Abs(step) < half {
		return step
	}
	r := math.Remainder(step, 2*half) // in [−half, half]
	if r == -half && step > 0 {
		r = half
	} else if r == half && step < 0 {
		r = -half
	}
	return r
}
