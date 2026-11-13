package electrical

import (
	"fmt"
	"math"
)

// WireResistance returns the resistance of a round wire of uniform section, R = L / (σ · A) with
// A = π d² / 4 (R = ρL/A, ρ = 1/σ).
// Errors: ErrNotFinite; ErrNegative for any negative input; ErrZero for a zero diameter or conductivity;
// ErrOverflow.
//
// ID:      ENG-ELEC-011
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Ling SJ, Sanny J, Moebs W. University Physics Volume 2. OpenStax, 2016, https://openstax.org/details/books/university-physics-volume-2, §9.3 "Resistivity and Resistance" (R = ρL/A; σ = 1/ρ)
// Tags:    wire, resistance, resistivity, conductivity, conductor, electrical
// Units:   lengthMeters m; diameterMeters m; conductivitySiemensPerMeter S/m; result Ω
// Example: 30.48 m (100 ft) of 1.27 mm (0.05 in) wire at 5.96e7 S/m (annealed copper) → 0.4037117406548 Ω
func WireResistance(lengthMeters, diameterMeters, conductivitySiemensPerMeter float64) (float64, error) {
	for _, p := range [3]struct {
		name string
		v    float64
	}{{"lengthMeters", lengthMeters}, {"diameterMeters", diameterMeters}, {"conductivitySiemensPerMeter", conductivitySiemensPerMeter}} {
		if math.IsNaN(p.v) || math.IsInf(p.v, 0) || p.v < 0 {
			return 0, checkInput(p.name, p.v)
		}
	}
	if diameterMeters == 0 {
		return 0, fmt.Errorf("%w: diameterMeters = 0", ErrZero)
	}
	if conductivitySiemensPerMeter == 0 {
		return 0, fmt.Errorf("%w: conductivitySiemensPerMeter = 0", ErrZero)
	}
	area := math.Pi * diameterMeters * diameterMeters / 4
	den := conductivitySiemensPerMeter * area
	if den == 0 {
		return 0, fmt.Errorf("%w: conductivity × area underflows to 0", ErrOverflow)
	}
	r := lengthMeters / den
	if math.IsInf(r, 0) {
		return 0, fmt.Errorf("%w: wire resistance", ErrOverflow)
	}
	return r, nil
}
