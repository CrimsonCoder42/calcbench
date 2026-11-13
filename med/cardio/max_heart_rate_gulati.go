package cardio

// MaxHeartRateGulati returns the peak heart rate on exercise testing predicted for women, in beats/min, by
// Gulati's equation 206 − 0.88 × age.
// Errors: ErrNotFinite for NaN or ±Inf; ErrOutOfRange outside Range.
//
// ID:      MED-CARDIO-029
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Gulati M, Shaw LJ, Thisted RA, Black HR, Bairey Merz CN, Arnsdorf MF. Heart rate response to exercise stress testing in asymptomatic women: the St. James Women Take Heart Project. Circulation. 2010;122(2):130-137. doi:10.1161/CIRCULATIONAHA.110.939249. PMID 20585008. Abstract: "mean peak HR=206-0.88(age)"
// Tags:    cardio, heart-rate, maximal-heart-rate, exercise, women, gulati
// Units:   ageYears years; result beats/min
// Range:   ageYears 18–100 (5,437 asymptomatic women, mean age 52 ± 11; the abstract states no age limits)
// Example: age 40 → 206 − 35.2 = 170.8 (HAND-DERIVED)
func MaxHeartRateGulati(ageYears float64) (float64, error) {
	if err := checkFinite("ageYears", ageYears); err != nil {
		return 0, err
	}
	if err := checkRange("ageYears", ageYears, 18, 100); err != nil {
		return 0, err
	}
	return 206 - 0.88*ageYears, nil
}
