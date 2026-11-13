package cardio

// MaxHeartRateNes returns the age-predicted maximal heart rate, in beats/min, by the HUNT Fitness Study
// equation 211 − 0.64 × age.
// Errors: ErrNotFinite for NaN or ±Inf; ErrOutOfRange outside Range.
//
// ID:      MED-CARDIO-030
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Nes BM, Janszky I, Wisløff U, Støylen A, Karlsen T. Age-predicted maximal heart rate in healthy subjects: The HUNT Fitness Study. Scand J Med Sci Sports. 2013;23(6):697-704. doi:10.1111/j.1600-0838.2012.01445.x. PMID 22376273. Abstract: "HRmax was univariately explained by the formula 211 - 0.64·age"
// Tags:    cardio, heart-rate, maximal-heart-rate, exercise, hunt, nes
// Units:   ageYears years; result beats/min
// Range:   ageYears 18–100 (3,320 healthy Norwegian adults, "wide age range"; exact limits not read)
// Example: age 40 → 211 − 25.6 = 185.4 (HAND-DERIVED)
func MaxHeartRateNes(ageYears float64) (float64, error) {
	if err := checkFinite("ageYears", ageYears); err != nil {
		return 0, err
	}
	if err := checkRange("ageYears", ageYears, 18, 100); err != nil {
		return 0, err
	}
	return 211 - 0.64*ageYears, nil
}
