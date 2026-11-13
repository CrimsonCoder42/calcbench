package endocrine

// BeAMValue returns the BeAM value, in mg/dL: the bedtime glucose minus the next morning's pre-breakfast
// glucose. Negative when the morning value is the higher one. Not rounded.
// Errors: ErrNotFinite for NaN or ±Inf; ErrOutOfRange outside Range.
//
// ID:      MED-ENDO-003
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Zisman A, Morales F, Stewart J, Stuhr A, Vlajnic A, Zhou R. BeAM value: an indicator of the need to initiate and intensify prandial therapy in patients with type 2 diabetes mellitus receiving basal insulin. BMJ Open Diabetes Res Care. 2016;4(1):e000171. doi:10.1136/bmjdrc-2015-000171. PMID 27110368. PMC4838666, Box 1: "BeAM value = Bedtime (Be) glucose value minus morning (AM) glucose value"
// Tags:    endocrine, diabetes, beam, glucose, basal-insulin
// Units:   bedtimeGlucoseMgDL mg/dL; morningGlucoseMgDL mg/dL (pre-breakfast); result mg/dL
// Range:   both > 0 (population: adults with type 2 diabetes on basal insulin; input limits not stated in the source)
// Example: bedtime 200, morning 130 → 70 (HAND-DERIVED from Box 1)
func BeAMValue(bedtimeGlucoseMgDL, morningGlucoseMgDL float64) (float64, error) {
	if err := checkFinite("bedtimeGlucoseMgDL", bedtimeGlucoseMgDL); err != nil {
		return 0, err
	}
	if err := checkFinite("morningGlucoseMgDL", morningGlucoseMgDL); err != nil {
		return 0, err
	}
	if err := checkPositive("bedtimeGlucoseMgDL", bedtimeGlucoseMgDL); err != nil {
		return 0, err
	}
	if err := checkPositive("morningGlucoseMgDL", morningGlucoseMgDL); err != nil {
		return 0, err
	}
	// The difference of two positive finite values is finite.
	return bedtimeGlucoseMgDL - morningGlucoseMgDL, nil
}
