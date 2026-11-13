package pediatrics

import (
	"errors"
	"fmt"

	"github.com/CrimsonCoder42/calcbench/algo/reduce"
)

var (
	// ErrNotFinite is returned when a float input is NaN or ±Inf.
	ErrNotFinite = errors.New("pediatrics: value is NaN or infinite")
	// ErrOutOfRange is returned when an input is outside its physical range or the range the source validated.
	ErrOutOfRange = errors.New("pediatrics: input outside the range the source validated")
)

// checkFloat returns ErrNotFinite for NaN or ±Inf and ErrOutOfRange outside lo–hi, wrapped with name.
func checkFloat(name string, v, lo, hi float64) error {
	if !reduce.AllFinite(v) {
		return fmt.Errorf("%w: %s = %v", ErrNotFinite, name, v)
	}
	if v < lo || v > hi {
		return fmt.Errorf("%w: %s = %v, want %v–%v", ErrOutOfRange, name, v, lo, hi)
	}
	return nil
}

// KocherCriteriaCount returns the number of Kocher's four predictors present, 0–4: "history of fever,
// non-weight-bearing, erythrocyte sedimentation rate of at least forty millimeters per hour, and serum white
// blood-cell count of more than 12,000 cells per cubic millimeter". The paper's predicted probabilities by
// count are its population's and are not returned.
// Errors: ErrNotFinite for NaN or ±Inf; ErrOutOfRange outside Range.
//
// ID:      MED-PEDS-001
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Kocher MS, Zurakowski D, Kasser JR. Differentiating between septic arthritis and transient synovitis of the hip in children: an evidence-based clinical prediction algorithm. J Bone Joint Surg Am. 1999;81(12):1662-1670. doi:10.2106/00004623-199912000-00002, PMID 10608376, abstract
// Tags:    pediatrics, orthopedics, septic arthritis, hip, kocher, score
// Units:   esrMmPerHour mm/h; wbcPerMm3 cells/mm³; findings present/absent; result count 0–4
// Range:   0 ≤ esrMmPerHour ≤ 200; 0 ≤ wbcPerMm3 ≤ 500000 (physical limits)
// Example: fever, non-weight-bearing, ESR 40, WBC 12,001 → 4; ESR 39.9, WBC 12,000, neither finding → 0 (abstract thresholds)
func KocherCriteriaCount(historyOfFever, nonWeightBearing bool, esrMmPerHour, wbcPerMm3 float64) (int, error) {
	if err := checkFloat("esrMmPerHour", esrMmPerHour, 0, 200); err != nil {
		return 0, err
	}
	if err := checkFloat("wbcPerMm3", wbcPerMm3, 0, 500000); err != nil {
		return 0, err
	}
	n := 0
	for _, v := range [...]bool{historyOfFever, nonWeightBearing, esrMmPerHour >= 40, wbcPerMm3 > 12000} {
		if v {
			n++
		}
	}
	return n, nil
}
