package endocrine

// IADPSG 75-g OGTT thresholds, IADPSG 2010 Table 1 (PMC2827530), as printed in each unit system.
const (
	iadpsgFastingMgDL  = 92
	iadpsgOneHourMgDL  = 180
	iadpsgTwoHourMgDL  = 153
	iadpsgFastingMmolL = 5.1
	iadpsgOneHourMmolL = 10.0
	iadpsgTwoHourMmolL = 8.5
)

// checkOGTT returns ErrNotFinite or ErrOutOfRange (in that order over the three values) unless every value is
// finite and > 0.
func checkOGTT(fasting, oneHour, twoHour float64) error {
	if err := checkFinite("fasting", fasting); err != nil {
		return err
	}
	if err := checkFinite("oneHour", oneHour); err != nil {
		return err
	}
	if err := checkFinite("twoHour", twoHour); err != nil {
		return err
	}
	if err := checkPositive("fasting", fasting); err != nil {
		return err
	}
	if err := checkPositive("oneHour", oneHour); err != nil {
		return err
	}
	return checkPositive("twoHour", twoHour)
}

// GDMThresholdIADPSGMgDL reports whether one or more of the three 75-g oral glucose tolerance test plasma
// glucose values equals or exceeds its IADPSG threshold (fasting 92, 1 h 180, 2 h 153 mg/dL), the source's
// rule: "One or more of these values from a 75-g OGTT must be equaled or exceeded for the diagnosis of GDM."
// Errors: ErrNotFinite for NaN or ±Inf; ErrOutOfRange outside Range.
//
// ID:      MED-ENDO-005
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  International Association of Diabetes and Pregnancy Study Groups Consensus Panel; Metzger BE, Gabbe SG, Persson B, et al. International Association of Diabetes and Pregnancy Study Groups recommendations on the diagnosis and classification of hyperglycemia in pregnancy. Diabetes Care. 2010;33(3):676-682. doi:10.2337/dc09-1848. PMID 20190296. PMC2827530, Table 1
// Tags:    endocrine, diabetes, pregnancy, gdm, ogtt, glucose, iadpsg
// Units:   fastingMgDL, oneHourMgDL, twoHourMgDL plasma glucose mg/dL; result true when a threshold is met
// Range:   every value > 0 (all three OGTT values are required; the source states no other limits)
// Example: 92, 150, 140 → true (fasting equals 92); 91, 179, 152 → false (Table 1 thresholds)
func GDMThresholdIADPSGMgDL(fastingMgDL, oneHourMgDL, twoHourMgDL float64) (bool, error) {
	if err := checkOGTT(fastingMgDL, oneHourMgDL, twoHourMgDL); err != nil {
		return false, err
	}
	return fastingMgDL >= iadpsgFastingMgDL || oneHourMgDL >= iadpsgOneHourMgDL || twoHourMgDL >= iadpsgTwoHourMgDL, nil
}

// GDMThresholdIADPSGMmolL reports whether one or more of the three 75-g oral glucose tolerance test plasma
// glucose values equals or exceeds its IADPSG threshold in mmol/L (fasting 5.1, 1 h 10.0, 2 h 8.5), the
// thresholds the source prints in its own mmol/L column (no conversion is applied).
// Errors: ErrNotFinite for NaN or ±Inf; ErrOutOfRange outside Range.
//
// ID:      MED-ENDO-006
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Metzger BE, Gabbe SG, Persson B, et al. (IADPSG Consensus Panel). Diabetes Care. 2010;33(3):676-682. doi:10.2337/dc09-1848. PMID 20190296. PMC2827530, Table 1 (mmol/L column)
// Tags:    endocrine, diabetes, pregnancy, gdm, ogtt, glucose, iadpsg, si
// Units:   fastingMmolL, oneHourMmolL, twoHourMmolL plasma glucose mmol/L; result true when a threshold is met
// Range:   every value > 0 (all three OGTT values are required; the source states no other limits)
// Example: 5.1, 8.0, 7.0 → true; 5.0, 9.9, 8.4 → false (Table 1 thresholds)
func GDMThresholdIADPSGMmolL(fastingMmolL, oneHourMmolL, twoHourMmolL float64) (bool, error) {
	if err := checkOGTT(fastingMmolL, oneHourMmolL, twoHourMmolL); err != nil {
		return false, err
	}
	return fastingMmolL >= iadpsgFastingMmolL || oneHourMmolL >= iadpsgOneHourMmolL || twoHourMmolL >= iadpsgTwoHourMmolL, nil
}
