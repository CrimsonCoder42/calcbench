package business

// ContributionMargin returns sales minus variable costs, in cents, exactly. Totals give the total margin;
// per-unit amounts give the margin per unit. Variable costs above sales give a negative margin.
// Errors: ErrOutOfRange when an input is outside ±2⁵³; ErrNegative when an input is negative. (Both
// inputs are in [0, 2⁵³], so the difference cannot overflow.)
//
// ID:      FIN-BIZ-040
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Franklin, Graybeal and Cooper, Principles of Accounting, Vol. 2: Managerial Accounting, OpenStax (Rice University), 2019, §3.1 "Explain Contribution Margin and Calculate Contribution Margin per Unit, Contribution Margin Ratio, and Total Contribution Margin", https://openstax.org/books/principles-managerial-accounting/pages/3-1-explain-contribution-margin-and-calculate-contribution-margin-per-unit-contribution-margin-ratio-and-total-contribution-margin (read 2026-11-12)
// Tags:    cost-volume-profit, contribution-margin, managerial
// Units:   salesCents, variableCostsCents and the result in cents (or all in cents per unit)
// Example: Hicks Manufacturing, April: $50,000 − $10,000 = $40,000; per unit $100 − $20 = $80 (Principles of Accounting Vol. 2 §3.1)
func ContributionMargin(salesCents, variableCostsCents int64) (int64, error) {
	if err := checkNonNegativeCents("salesCents", salesCents); err != nil {
		return 0, err
	}
	if err := checkNonNegativeCents("variableCostsCents", variableCostsCents); err != nil {
		return 0, err
	}
	return salesCents - variableCostsCents, nil
}
