package loan

import "fmt"

// DebtServiceCoverageRatio returns the debt-service coverage ratio, netOperatingIncomeCents divided by
// debtServiceCents (principal and interest due) over the same period, as an unrounded decimal (1.25 =
// income covers the debt service 1.25 times). A negative net operating income is allowed and gives a
// negative ratio.
// Errors: ErrOutOfRange if either input is outside ±2⁵³ cents; ErrNegative if debtServiceCents < 0;
// ErrZeroDenominator if debtServiceCents is 0.
//
// ID:      FIN-LOAN-008
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Office of the Comptroller of the Currency, Comptroller's Handbook: Commercial Real Estate Lending, Version 2.0, March 2022, "Debt-Service Coverage Ratio", p. 43 ("calculated by dividing the NOI by the annual debt service requirements") and Glossary p. 137, https://www.occ.treas.gov/publications-and-resources/publications/comptrollers-handbook/files/commercial-real-estate-lending/pub-ch-commercial-real-estate.pdf (read 2026-11-12)
// Tags:    loan, commercial-real-estate, ratio, dscr, debt-service, coverage, underwriting
// Units:   both inputs in cents over the same period; result decimal (times covered)
// Example: hand-checked: NOI $125,000, annual debt service $100,000 → 1.25
func DebtServiceCoverageRatio(netOperatingIncomeCents, debtServiceCents int64) (float64, error) {
	if err := checkCents("netOperatingIncomeCents", netOperatingIncomeCents); err != nil {
		return 0, err
	}
	if err := checkNonNegativeCents("debtServiceCents", debtServiceCents); err != nil {
		return 0, err
	}
	if debtServiceCents == 0 {
		return 0, fmt.Errorf("%w: debtServiceCents = 0", ErrZeroDenominator)
	}
	return float64(netOperatingIncomeCents) / float64(debtServiceCents), nil
}
