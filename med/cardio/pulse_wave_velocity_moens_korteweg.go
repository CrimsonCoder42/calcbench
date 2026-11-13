package cardio

import (
	"fmt"
	"math"
)

// PulseWaveVelocityMoensKorteweg returns the propagation speed of a pressure wave, in m/s, in a thin-walled
// elastic tube filled with an incompressible liquid, by Korteweg's eq. (10), β = √(a₁E₁ / (2ρR₁)), with a₁
// the wall thickness, E₁ the elastic modulus of the wall, ρ the liquid's density and R₁ the tube radius.
// Korteweg states the conditions (pp. 530–531): the wavelength is very long compared with the tube
// diameter and the kinetic energy of the transverse motion is negligible against the longitudinal one.
// Errors: ErrNotFinite for NaN or ±Inf; ErrOutOfRange for a non-positive input; ErrDomain when the result
// overflows or underflows float64.
//
// ID:      MED-CARDIO-003
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Korteweg DJ. Ueber die Fortpflanzungsgeschwindigkeit des Schalles in elastischen Röhren. Ann Phys Chem. 1878;241(12):525-542. doi:10.1002/andp.18782411206. Eq. (10), p. 530; symbols p. 526 (scan read 2026-11-12, zenodo.org/record/1423762)
// Source:  Moens AI. Die Pulscurve. Leiden: E. J. Brill; 1878 (cited by Korteweg p. 530, footnote 5, as the experimental counterpart)
// Tags:    cardio, pulse-wave-velocity, arterial-stiffness, moens-korteweg, hemodynamics
// Units:   SI throughout: wallThicknessM m; elasticModulusPa Pa; densityKgM3 kg/m³; radiusM m; result m/s
// Range:   every input > 0; the physics has no population (Korteweg's conditions above apply)
// Example: a₁ 0.001 m, E₁ 4×10⁵ Pa, ρ 1050 kg/m³, R₁ 0.01 m → √(400/21) = 4.36435… m/s (hand-derived)
func PulseWaveVelocityMoensKorteweg(wallThicknessM, elasticModulusPa, densityKgM3, radiusM float64) (float64, error) {
	if err := checkFinite("wallThicknessM", wallThicknessM); err != nil {
		return 0, err
	}
	if err := checkFinite("elasticModulusPa", elasticModulusPa); err != nil {
		return 0, err
	}
	if err := checkFinite("densityKgM3", densityKgM3); err != nil {
		return 0, err
	}
	if err := checkFinite("radiusM", radiusM); err != nil {
		return 0, err
	}
	if err := checkPositive("wallThicknessM", wallThicknessM); err != nil {
		return 0, err
	}
	if err := checkPositive("elasticModulusPa", elasticModulusPa); err != nil {
		return 0, err
	}
	if err := checkPositive("densityKgM3", densityKgM3); err != nil {
		return 0, err
	}
	if err := checkPositive("radiusM", radiusM); err != nil {
		return 0, err
	}
	v := math.Sqrt(wallThicknessM * elasticModulusPa / (2 * densityKgM3 * radiusM))
	if math.IsInf(v, 0) || math.IsNaN(v) || v == 0 {
		return 0, fmt.Errorf("%w: result %v is outside float64 for these inputs", ErrDomain, v)
	}
	return v, nil
}
