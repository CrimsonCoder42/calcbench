package business

// DegreeOfOperatingLeverage returns contribution margin divided by net operating income, as a decimal. It is
// not rounded. Either input may be negative; the sign of the result follows the inputs.
// Errors: ErrOutOfRange when an input is outside ±2⁵³; ErrZeroDenominator when operatingIncomeCents is 0.
//
// ID:      FIN-BIZ-038
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Franklin, Graybeal and Cooper, Principles of Accounting, Vol. 2: Managerial Accounting, OpenStax (Rice University), 2019, §3.5 "Calculate and Interpret a Company's Margin of Safety and Operating Leverage", https://openstax.org/books/principles-managerial-accounting/pages/3-5-calculate-and-interpret-a-companys-margin-of-safety-and-operating-leverage (read 2026-11-12)
// Tags:    ratio, leverage, operating-leverage, cost-volume-profit
// Units:   contributionMarginCents and operatingIncomeCents in cents; the result is a dimensionless decimal
// Example: Company A $148,000 ÷ $85,000 = 1.741…, printed 1.74 (Principles of Accounting Vol. 2 §3.5)
func DegreeOfOperatingLeverage(contributionMarginCents, operatingIncomeCents int64) (float64, error) {
	if err := checkCents("contributionMarginCents", contributionMarginCents); err != nil {
		return 0, err
	}
	if err := checkCents("operatingIncomeCents", operatingIncomeCents); err != nil {
		return 0, err
	}
	return ratioCents(contributionMarginCents, "operatingIncomeCents", operatingIncomeCents)
}
