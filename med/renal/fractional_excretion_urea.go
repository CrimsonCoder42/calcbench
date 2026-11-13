package renal

// FractionalExcretionUrea returns the fractional excretion of urea (FEUrea), in percent:
// (UUrea × PCr) / (PUrea × UCr) × 100. The urea pair may be in any one unit (mg/dL of urea nitrogen or
// mmol/L), the same for both, and the creatinine pair likewise: each ratio cancels its unit. No category is
// returned.
// Errors: ErrNotFinite for NaN or ±Inf; ErrOutOfRange outside Range; ErrDomain when the result is beyond float64.
//
// ID:      MED-RENAL-010
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Carvounis CP, Nisar S, Guro-Razuman S. Significance of the fractional excretion of urea in the differential diagnosis of acute renal failure. Kidney Int. 2002;62(6):2223-2229. doi:10.1046/j.1523-1755.2002.00683.x. PMID 12427149
// Tags:    renal, feurea, urea, creatinine, urine
// Units:   urineUrea, plasmaUrea any one unit, the same for both; urineCreatinine, plasmaCreatinine any one unit, the same for both; result percent
// Range:   every input > 0 (no input limits stated by the source)
// Example: UUrea 300, PUrea 30, UCr 100, PCr 2 → (300 × 2) / (30 × 100) × 100 = 20 % (hand-derived)
func FractionalExcretionUrea(urineUrea, plasmaUrea, urineCreatinine, plasmaCreatinine float64) (float64, error) {
	if err := checkFractionalExcretion("urineUrea", urineUrea, "plasmaUrea", plasmaUrea,
		urineCreatinine, plasmaCreatinine); err != nil {
		return 0, err
	}
	return fractionalExcretion(urineUrea, plasmaUrea, urineCreatinine, plasmaCreatinine)
}
