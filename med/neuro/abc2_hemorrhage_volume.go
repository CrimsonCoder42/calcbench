package neuro

import (
	"errors"
	"fmt"
	"math"

	"github.com/CrimsonCoder42/calcbench/algo/reduce"
)

var (
	// ErrNotFinite is returned when a float input is NaN or ±Inf.
	ErrNotFinite = errors.New("neuro: value is NaN or infinite")
	// ErrOutOfRange is returned when an input is outside the range the source validated or allows.
	ErrOutOfRange = errors.New("neuro: input outside the range the source validated")
)

// checkRange returns ErrNotFinite for NaN or ±Inf and ErrOutOfRange outside [lo, hi], naming the parameter.
func checkRange(name string, v, lo, hi float64) error {
	if !reduce.AllFinite(v) {
		return fmt.Errorf("%w: %s = %v", ErrNotFinite, name, v)
	}
	if v < lo || v > hi {
		return fmt.Errorf("%w: %s = %v, want %v ≤ %s ≤ %v", ErrOutOfRange, name, v, lo, name, hi)
	}
	return nil
}

// checkInt returns ErrOutOfRange when an integer item is outside [lo, hi], naming the parameter.
func checkInt(name string, v, lo, hi int) error {
	if v < lo || v > hi {
		return fmt.Errorf("%w: %s = %d, want %d–%d", ErrOutOfRange, name, v, lo, hi)
	}
	return nil
}

// ABC2HemorrhageVolumeML returns the intracerebral hemorrhage volume, in mL (cm³), by the ABC/2 method:
// A × B × C / 2, where A is the greatest hemorrhage diameter on CT, B the diameter at 90° to A, and C the
// number of CT slices with hemorrhage multiplied by the slice thickness.
// Errors: ErrNotFinite for NaN or ±Inf; ErrOutOfRange when a dimension is not > 0.
//
// ID:      MED-NEURO-001
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Kothari RU, Brott T, Broderick JP, et al. The ABCs of measuring intracerebral hemorrhage volumes. Stroke. 1996;27(8):1304-1305. doi:10.1161/01.str.27.8.1304, PMID 8711791, abstract ("The formula ABC/2 was used, where A is the greatest hemorrhage diameter by CT, B is the diameter 90 degrees to A, and C is the approximate number of CT slices with hemorrhage multiplied by the slice thickness")
// Tags:    neuro, stroke, intracerebral hemorrhage, ich, volume, ct
// Units:   aCm, bCm, cCm centimetres (cCm = slices with hemorrhage × slice thickness in cm, the caller's count); result mL (= cm³)
// Range:   aCm, bCm, cCm > 0 (the abstract states no limits)
// Example: A 4 cm, B 3 cm, C 2.5 cm → 4 × 3 × 2.5 / 2 = 15 mL (hand-derived from the abstract's formula)
func ABC2HemorrhageVolumeML(aCm, bCm, cCm float64) (float64, error) {
	for _, p := range [3]struct {
		name string
		v    float64
	}{{"aCm", aCm}, {"bCm", bCm}, {"cCm", cCm}} {
		if err := checkRange(p.name, p.v, 0, math.MaxFloat64); err != nil {
			return 0, err
		}
		if p.v == 0 {
			return 0, fmt.Errorf("%w: %s = 0, want > 0", ErrOutOfRange, p.name)
		}
	}
	v := aCm * bCm * cCm / 2
	if math.IsInf(v, 0) {
		return 0, fmt.Errorf("%w: product A × B × C overflows", ErrOutOfRange)
	}
	return v, nil
}
