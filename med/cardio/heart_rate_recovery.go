package cardio

import "fmt"

// HeartRateRecovery1Min returns the heart-rate recovery at one minute, in beats/min: the decrease from the
// heart rate at peak exercise to the heart rate one minute after exercise stops. Cole 1999 calls a value of
// 12 beats/min or less abnormal; that label is not returned.
// Errors: ErrNotFinite for NaN or ±Inf; ErrOutOfRange when a rate is ≤ 0; ErrDomain when the one-minute rate
// is above the peak rate (no decrease).
//
// ID:      MED-CARDIO-028
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Cole CR, Blackstone EH, Pashkow FJ, Snader CE, Lauer MS. Heart-rate recovery immediately after exercise as a predictor of mortality. N Engl J Med. 1999;341(18):1351-1357. doi:10.1056/NEJM199910283411804. PMID 10536127. Abstract: "the decrease in the heart rate from peak exercise to one minute after the cessation of exercise"; abnormal "a reduction of 12 beats per minute or less"
// Tags:    cardio, heart-rate, exercise, recovery, cole
// Units:   peakBpm, oneMinuteBpm beats/min; result beats/min
// Range:   rates > 0; oneMinuteBpm ≤ peakBpm (2,428 adults, mean age 57, symptom-limited exercise testing)
// Example: peak 150, one minute 135 → 15 (HAND-DERIVED)
func HeartRateRecovery1Min(peakBpm, oneMinuteBpm float64) (float64, error) {
	if err := checkFinite("peakBpm", peakBpm); err != nil {
		return 0, err
	}
	if err := checkFinite("oneMinuteBpm", oneMinuteBpm); err != nil {
		return 0, err
	}
	if err := checkPositive("peakBpm", peakBpm); err != nil {
		return 0, err
	}
	if err := checkPositive("oneMinuteBpm", oneMinuteBpm); err != nil {
		return 0, err
	}
	if oneMinuteBpm > peakBpm {
		return 0, fmt.Errorf("%w: oneMinuteBpm = %v above peakBpm = %v", ErrDomain, oneMinuteBpm, peakBpm)
	}
	return peakBpm - oneMinuteBpm, nil
}
