package fitness

import "fmt"

// Pace returns the pace, in seconds per unit of distance: the elapsed time divided by the distance
// covered, the reciprocal of the average speed. The distance may be in kilometres or miles; the result is
// seconds per kilometre or seconds per mile accordingly.
// Errors: ErrNotFinite for NaN or ±Inf; ErrOutOfRange when durationSeconds ≤ 0 or distanceKmOrMi ≤ 0;
// ErrDomain when the result is beyond float64.
//
// ID:      LIF-FIT-009
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  BIPM, The International System of Units (SI), 9th ed., 2019, §2.3.4 and Table 5 (speed: the derived quantity length per time, m s⁻¹); pace is its reciprocal, time per length
// Tags:    fitness, running, pace, speed, endurance
// Units:   durationSeconds s; distanceKmOrMi km or mi; result s/km or s/mi, the distance's unit
// Range:   durationSeconds > 0; distanceKmOrMi > 0
// Example: 1860 s over 3.1 mi → 600 s/mi (10:00 per mile) (HAND-DERIVED)
func Pace(durationSeconds, distanceKmOrMi float64) (float64, error) {
	if err := checkFinite("durationSeconds", durationSeconds); err != nil {
		return 0, err
	}
	if err := checkFinite("distanceKmOrMi", distanceKmOrMi); err != nil {
		return 0, err
	}
	if durationSeconds <= 0 {
		return 0, fmt.Errorf("%w: durationSeconds = %v, want > 0", ErrOutOfRange, durationSeconds)
	}
	if distanceKmOrMi <= 0 {
		return 0, fmt.Errorf("%w: distanceKmOrMi = %v, want > 0", ErrOutOfRange, distanceKmOrMi)
	}
	v := durationSeconds / distanceKmOrMi
	if err := checkResult(v); err != nil {
		return 0, err
	}
	return v, nil
}
