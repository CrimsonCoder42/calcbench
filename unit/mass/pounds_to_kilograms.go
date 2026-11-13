package mass

// PoundsToKilograms returns the mass in kilograms of a mass in pounds (avoirdupois):
// m/kg = 0.453 592 37 · m/lb.
// Errors: ErrNotFinite for NaN or ±Inf; ErrNegative below 0 lb.
//
// ID:      UNT-MASS-002
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Thompson A, Taylor BN. Guide for the Use of the International System of Units (SI). NIST Special Publication 811, 2008 ed., Appendix B.8 "Factors for units listed alphabetically", row "pound (avoirdupois) (lb) to kilogram (kg) 4.535 924 E-01" and footnote 22: "The exact conversion factor is 4.535 923 7 E-01"
// Tags:    unit, mass, weight, pound, kilogram, avoirdupois, conversion
// Units:   pounds lb; result kg
// Range:   pounds ≥ 0
// Example: 1 lb → 0.45359237 kg; 150 lb → 68.0388555 kg; 2.2046226218487757 lb → 1 kg
func PoundsToKilograms(pounds float64) (float64, error) {
	if err := checkMass("pounds", pounds); err != nil {
		return 0, err
	}
	// The factor is below 1, so a finite input gives a finite result.
	return pounds * kilogramsPerPound, nil
}
