package respiratory

import (
	"fmt"
	"math"
)

// RapidShallowBreathingIndex returns the rapid shallow breathing index f/VT, in breaths/min/L: "the ratio of
// respiratory frequency to tidal volume". No threshold is applied (none was seen in the source's abstract).
// Errors: ErrNotFinite for NaN or ±Inf; ErrOutOfRange for a non-positive input; ErrDomain when the ratio overflows.
//
// ID:      MED-RESP-017
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Yang KL, Tobin MJ. A prospective study of indexes predicting the outcome of trials of weaning from mechanical ventilation. N Engl J Med. 1991;324(21):1445-1450. doi:10.1056/NEJM199105233242101. PMID 2023603, abstract: "the ratio of respiratory frequency to tidal volume (f/VT)"
// Tags:    respiratory, rsbi, f-vt, yang-tobin, weaning, ventilation, index
// Units:   respiratoryRatePerMin breaths/min; tidalVolumeL litres; result breaths/min/L
// Range:   respiratoryRatePerMin > 0; tidalVolumeL > 0
// Example: 30 breaths/min, VT 0.3 L → 100 breaths/min/L (hand-derived)
func RapidShallowBreathingIndex(respiratoryRatePerMin, tidalVolumeL float64) (float64, error) {
	if err := checkPositiveFinite("respiratoryRatePerMin", respiratoryRatePerMin); err != nil {
		return 0, err
	}
	if err := checkPositiveFinite("tidalVolumeL", tidalVolumeL); err != nil {
		return 0, err
	}
	r := respiratoryRatePerMin / tidalVolumeL
	if math.IsInf(r, 0) {
		return 0, fmt.Errorf("%w: respiratoryRatePerMin / tidalVolumeL overflows", ErrDomain)
	}
	return r, nil
}
