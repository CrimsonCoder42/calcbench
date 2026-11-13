package loan

// PITIPayment returns the total monthly mortgage payment: principal and interest, plus one twelfth of the
// annual property tax and of the annual homeowners insurance (the escrow), plus the monthly mortgage
// insurance. The escrow twelfth is rounded half-even to the cent in exact integer arithmetic; the other
// parts are added unchanged.
// Errors: ErrOutOfRange if any input is outside ±2⁵³ cents; ErrNegative if any input is negative;
// ErrOverflow if the total exceeds 2⁵³ cents.
//
// ID:      FIN-LOAN-007
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Consumer Financial Protection Bureau, "On a mortgage, what's the difference between my principal and interest payment and my total monthly payment?" ("Principal + interest + mortgage insurance (if applicable) + escrow (homeowners insurance and tax) = total monthly payment"), https://www.consumerfinance.gov/ask-cfpb/what-is-a-piti-en-1941/ (read 2026-11-12)
// Source:  U.S. Department of Housing and Urban Development, Handbook 4000.1 Glossary and Acronyms (last revised 8/12/2026), "PITI - Principal, Interest, Taxes, and Insurance", https://www.hud.gov/sites/default/files/Housing/documents/40001-hsgh-Update-18-Glossary.pdf (read 2026-11-12)
// Tags:    loan, mortgage, piti, payment, escrow, property-tax, insurance
// Units:   every input and the result in cents; the tax and insurance are annual amounts, the others monthly
// Example: hand-checked: P&I $1,500, tax $3,600/yr, insurance $1,200/yr, mortgage insurance $50 → $1,950.00
func PITIPayment(principalAndInterestCents, annualPropertyTaxCents, annualInsuranceCents, monthlyMortgageInsuranceCents int64) (int64, error) {
	if err := checkNonNegativeCents("principalAndInterestCents", principalAndInterestCents); err != nil {
		return 0, err
	}
	if err := checkNonNegativeCents("annualPropertyTaxCents", annualPropertyTaxCents); err != nil {
		return 0, err
	}
	if err := checkNonNegativeCents("annualInsuranceCents", annualInsuranceCents); err != nil {
		return 0, err
	}
	if err := checkNonNegativeCents("monthlyMortgageInsuranceCents", monthlyMortgageInsuranceCents); err != nil {
		return 0, err
	}
	annual := annualPropertyTaxCents + annualInsuranceCents // ≤ 2⁵⁴, fits int64
	escrow, rem := annual/12, annual%12
	if rem > 6 || (rem == 6 && escrow%2 == 1) {
		escrow++
	}
	total := principalAndInterestCents + monthlyMortgageInsuranceCents + escrow // ≤ 3·2⁵³, fits int64
	if err := checkResultCents("total", total); err != nil {
		return 0, err
	}
	return total, nil
}
