package anthropometry

import (
	"fmt"
	"math"
)

// BodyRoundnessIndex returns the body roundness index (BRI), dimensionless, from waist circumference and
// height in cm: the body is modelled as an ellipse with semi-major axis height/2 and semi-minor axis
// waist/(2π), its eccentricity is ε = √(1 − (waist/(2π))² / (height/2)²), and BRI = 364.2 − 365.5 × ε.
// Not rounded; no band is returned.
// Errors: ErrNotFinite for NaN or ±Inf; ErrOutOfRange outside Range; ErrDomain when waist/(2π) ≥ height/2
// (the ellipse's minor axis is not shorter than its major axis, so the eccentricity is undefined).
//
// ID:      MED-ANTHRO-014
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Thomas DM, Bredlau C, Bosy-Westphal A, et al. Relationships between body roundness with body fat and visceral adipose tissue emerging from a new geometrical model. Obesity (Silver Spring). 2013;21(11):2264-2271. doi:10.1002/oby.20408. PMID 23519954. PMC3692604, Methods: "Body Roundness Index = 364.2 − 365.5 × Eccentricity"; the ellipse's axes from waist circumference ("½WC/π") and height, cm
// Source:  Zhang X, Ma N, Lin Q, et al. Body roundness index and all-cause mortality among US adults. JAMA Netw Open. 2024;7(6):e2415051. doi:10.1001/jamanetworkopen.2024.15051. PMID 38837158. PMC11154161, Methods: "BRI was calculated as 364.2 − 365.5 × √(1 − [waist circumference in centimeters / 2π]² / [0.5 × height in centimeters]²)" [Thomas 2013]
// Tags:    anthropometry, body-shape, bri, body-roundness-index, waist, height, eccentricity
// Units:   waistCm cm; heightCm cm; result dimensionless
// Range:   every input > 0 and waistCm < π × heightCm (adults, NHANES 1999–2004 derivation; input limits not stated in the source)
// Example: waist 90 cm, height 175 cm → ε = 0.98651…, BRI = 3.6307 (hand-derived from the Zhang 2024 form)
func BodyRoundnessIndex(waistCm, heightCm float64) (float64, error) {
	if err := checkFinite("waistCm", waistCm); err != nil {
		return 0, err
	}
	if err := checkFinite("heightCm", heightCm); err != nil {
		return 0, err
	}
	if err := checkPositive("waistCm", waistCm); err != nil {
		return 0, err
	}
	if err := checkPositive("heightCm", heightCm); err != nil {
		return 0, err
	}
	return bodyRoundnessIndex(waistCm, heightCm)
}

// BodyRoundnessIndexImperial returns the body roundness index (BRI), dimensionless, from waist circumference
// and height in inches. The eccentricity depends only on the ratio waist/height, which has no unit, so the
// inches enter the formula of MED-ANTHRO-014 as they are; no conversion factor is applied and the result is
// identical to BodyRoundnessIndex on the same body in cm.
// Errors: ErrNotFinite for NaN or ±Inf; ErrOutOfRange outside Range; ErrDomain when waistIn ≥ π × heightIn.
//
// ID:      MED-ANTHRO-015
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Thomas DM, Bredlau C, Bosy-Westphal A, et al. Relationships between body roundness with body fat and visceral adipose tissue emerging from a new geometrical model. Obesity (Silver Spring). 2013;21(11):2264-2271. doi:10.1002/oby.20408. PMID 23519954. PMC3692604, Methods: "Body Roundness Index = 364.2 − 365.5 × Eccentricity"
// Source:  Zhang X, Ma N, Lin Q, et al. Body roundness index and all-cause mortality among US adults. JAMA Netw Open. 2024;7(6):e2415051. doi:10.1001/jamanetworkopen.2024.15051. PMID 38837158. PMC11154161, Methods (the full formula, cm)
// Tags:    anthropometry, body-shape, bri, body-roundness-index, waist, height, imperial
// Units:   waistIn in; heightIn in; result dimensionless
// Range:   every input > 0 and waistIn < π × heightIn (adults; input limits not stated in the source)
// Example: waist 24 in, height 70 in → BRI = 0.8831 (hand-derived; the same as 60.96 cm and 177.8 cm)
func BodyRoundnessIndexImperial(waistIn, heightIn float64) (float64, error) {
	if err := checkFinite("waistIn", waistIn); err != nil {
		return 0, err
	}
	if err := checkFinite("heightIn", heightIn); err != nil {
		return 0, err
	}
	if err := checkPositive("waistIn", waistIn); err != nil {
		return 0, err
	}
	if err := checkPositive("heightIn", heightIn); err != nil {
		return 0, err
	}
	return bodyRoundnessIndex(waistIn, heightIn)
}

// bodyRoundnessIndex applies BRI = 364.2 − 365.5 √(1 − r²), r = (waist/(2π)) / (height/2) = waist/(π·height),
// to positive finite waist and height in the same unit.
func bodyRoundnessIndex(waist, height float64) (float64, error) {
	r := waist / (math.Pi * height) // π·height cannot overflow to 0; an overflow to +Inf gives r = 0
	if !(r < 1) {
		return 0, fmt.Errorf("%w: waist = %v ≥ π × height = π × %v, eccentricity undefined", ErrDomain, waist, height)
	}
	return 364.2 - 365.5*math.Sqrt(1-r*r), nil
}
