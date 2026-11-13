package physics

import (
	"fmt"
	"math"
)

// PowerFromWork returns the average power of a constant force moving its point of application a distance
// along its line in a time, P = W / t = F · d / t, in watts.
// Errors: ErrNotFinite; ErrNegative for any negative input; ErrZero for durationSeconds = 0; ErrOverflow.
//
// ID:      SCI-PHYS-007
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Ling SJ, Sanny J, Moebs W. University Physics Volume 1. OpenStax, 2016, https://openstax.org/details/books/university-physics-volume-1, §7.1 "Work" (W = F·d for a constant force along the displacement) and §7.4 "Power" (average power P = W/t)
// Tags:    power, work, force, horsepower, mechanics, physics
// Units:   forceNewtons N; distanceMeters m; durationSeconds s; result W
// Example: 10 N over 3 m in 2 s → 15 W
func PowerFromWork(forceNewtons, distanceMeters, durationSeconds float64) (float64, error) {
	w, err := product(forceNewtons, "forceNewtons", distanceMeters, "distanceMeters")
	if err != nil {
		return 0, err
	}
	p, err := quotient(w, "work", durationSeconds, "durationSeconds")
	if err != nil {
		return 0, fmt.Errorf("PowerFromWork: %w", err)
	}
	if math.IsInf(p, 0) {
		return 0, fmt.Errorf("%w: power", ErrOverflow)
	}
	return p, nil
}
