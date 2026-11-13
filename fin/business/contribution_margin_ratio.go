package business

// ContributionMarginRatio returns (sales − variable costs) divided by sales, as a decimal (0.80 = 80 %).
// Totals or per-unit amounts give the same ratio. The subtraction is exact; the result has one rounding.
// Variable costs above sales give a negative ratio.
// Errors: ErrOutOfRange when an input is outside ±2⁵³; ErrNegative when an input is negative;
// ErrZeroDenominator when salesCents is 0.
//
// ID:      FIN-BIZ-041
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Franklin, Graybeal and Cooper, Principles of Accounting, Vol. 2: Managerial Accounting, OpenStax (Rice University), 2019, §3.1 "Explain Contribution Margin and Calculate Contribution Margin per Unit, Contribution Margin Ratio, and Total Contribution Margin", https://openstax.org/books/principles-managerial-accounting/pages/3-1-explain-contribution-margin-and-calculate-contribution-margin-per-unit-contribution-margin-ratio-and-total-contribution-margin (read 2026-11-12)
// Tags:    ratio, cost-volume-profit, contribution-margin, managerial
// Units:   salesCents and variableCostsCents in cents; the result is a dimensionless decimal
// Example: Hicks Manufacturing Blue Jay: ($100 − $20) ÷ $100 = 0.80 (Principles of Accounting Vol. 2 §3.1)
func ContributionMarginRatio(salesCents, variableCostsCents int64) (float64, error) {
	if err := checkNonNegativeCents("salesCents", salesCents); err != nil {
		return 0, err
	}
	if err := checkNonNegativeCents("variableCostsCents", variableCostsCents); err != nil {
		return 0, err
	}
	return ratioCents(salesCents-variableCostsCents, "salesCents", salesCents)
}
