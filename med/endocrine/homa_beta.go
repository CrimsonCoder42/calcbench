package endocrine

import "fmt"

// HOMABeta returns the homeostasis model assessment of beta-cell function (HOMA-%B, the HOMA1 approximation),
// in percent: 20 × fasting insulin (µU/mL) / (fasting glucose (mmol/L) − 3.5). The source's units are mmol/L
// for glucose; a glucose in mg/dL is converted by the caller before the call. Not rounded.
// Errors: ErrNotFinite for NaN or ±Inf; ErrOutOfRange outside Range; ErrDomain when glucose ≤ 3.5 mmol/L (the
// denominator is zero or negative) or the result is beyond float64.
//
// ID:      MED-ENDO-025
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Matthews DR, Hosker JP, Rudenski AS, Naylor BA, Treacher DF, Turner RC. Homeostasis model assessment: insulin resistance and beta-cell function from fasting plasma glucose and insulin concentrations in man. Diabetologia. 1985;28(7):412-419. doi:10.1007/BF00280883. PMID 3899825 (full text not opened; abstract states no formula)
// Source:  Singh B, Saxena A. Surrogate markers of insulin resistance: a review. World J Diabetes. 2010;1(2):36-47. doi:10.4239/wjd.v1.i2.36. PMC3083884, HOMA section: "HOMA1 - %B = (20 × FPI)/(FPG - 3.5)", FPI µU/mL, FPG mmol/L, cites Matthews 1985
// Source:  Lee SH, Kwon HS, Park YM, et al. Predicting the development of diabetes using the product of triglycerides and glucose: the Chungju Metabolic Disease Cohort (CMC) study. PLoS One. 2014;9(2):e90430. doi:10.1371/journal.pone.0090430. PMID 24587359. PMC3938726, Methods: "20×fasting insulin (mIU/L)/[FPG (mmol/L)–3.5]" [ref 19, Matthews 1985]
// Verified: secondary (Singh 2010 PMC3083884; Lee 2014 PMC3938726)
// Tags:    endocrine, diabetes, beta-cell, homa, homa-b, glucose, insulin
// Units:   fastingInsulinMicroUnitsPerML µU/mL (= mIU/L); fastingGlucoseMmolL mmol/L; result percent
// Range:   insulin > 0; glucose > 3.5 mmol/L (the formula's own domain; the sources state no further limits)
// Example: insulin 10 µU/mL, glucose 5 mmol/L → 200 / 1.5 = 133.33… % (hand-derived)
func HOMABeta(fastingInsulinMicroUnitsPerML, fastingGlucoseMmolL float64) (float64, error) {
	if err := checkFinite("fastingInsulinMicroUnitsPerML", fastingInsulinMicroUnitsPerML); err != nil {
		return 0, err
	}
	if err := checkFinite("fastingGlucoseMmolL", fastingGlucoseMmolL); err != nil {
		return 0, err
	}
	if err := checkPositive("fastingInsulinMicroUnitsPerML", fastingInsulinMicroUnitsPerML); err != nil {
		return 0, err
	}
	if err := checkPositive("fastingGlucoseMmolL", fastingGlucoseMmolL); err != nil {
		return 0, err
	}
	d := fastingGlucoseMmolL - 3.5
	if !(d > 0) {
		return 0, fmt.Errorf("%w: fastingGlucoseMmolL = %v, want > 3.5 (denominator glucose − 3.5)", ErrDomain, fastingGlucoseMmolL)
	}
	r := 20 * fastingInsulinMicroUnitsPerML / d
	if err := checkResult(r); err != nil {
		return 0, err
	}
	return r, nil
}
