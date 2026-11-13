package anthropometry

import "math"

// BodySurfaceAreaHaycock returns the body surface area, in m², by the Haycock height-weight formula:
// weightKg^0.5378 × heightCm^0.3964 × 0.024265. Every finite positive input gives a finite positive result
// (the exponents sum to less than 1, so the result never overflows; the smallest inputs give about 2e-304).
// Errors: ErrNotFinite for NaN or ±Inf; ErrOutOfRange outside Range.
//
// ID:      MED-ANTHRO-004
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Haycock GB, Schwartz GJ, Wisotsky DH. Geometric method for measuring body surface area: a height-weight formula validated in infants, children, and adults. J Pediatr. 1978;93(1):62-66. doi:10.1016/s0022-3476(78)80601-5. PMID 650346. Abstract: "SA (m2) = weight (kg)0.5378 X height (cm)0.3964 X 0.024265"
// Tags:    anthropometry, bsa, body-surface-area, haycock, pediatric
// Units:   weightKg kg; heightCm cm; result m²
// Range:   weightKg > 0; heightCm > 0 (no hard limits stated; the 81 subjects ran from premature infants to adults, measured SA from less than 0.2 m² to greater than 2.0 m²)
// Example: 70 kg, 175 cm → 70^0.5378 × 175^0.3964 × 0.024265 (HAND-DERIVED; the test writes the formula out)
func BodySurfaceAreaHaycock(weightKg, heightCm float64) (float64, error) {
	if err := checkFinite("weightKg", weightKg); err != nil {
		return 0, err
	}
	if err := checkFinite("heightCm", heightCm); err != nil {
		return 0, err
	}
	if err := checkPositive("weightKg", weightKg); err != nil {
		return 0, err
	}
	if err := checkPositive("heightCm", heightCm); err != nil {
		return 0, err
	}
	return math.Pow(weightKg, 0.5378) * math.Pow(heightCm, 0.3964) * 0.024265, nil
}
