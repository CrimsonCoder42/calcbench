package invest

import "fmt"

// CapitalizationRate returns a property's capitalization rate (cap rate, overall rate): the annual net
// operating income divided by the property's price or value, as a decimal per year (0.08 = 8 %). Net
// operating income is annual gross income less operating expenses, before interest, principal and
// income taxes; a negative NOI gives a negative rate. The ratio is not rounded.
// Errors: ErrOutOfRange if an input is outside ±2⁵³ cents; ErrNegative if priceCents is negative;
// ErrZeroDenominator if priceCents is 0.
//
// ID:      FIN-INVEST-053
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Office of the Comptroller of the Currency, Comptroller's Handbook: Commercial Real Estate Lending, Version 2.0, March 2022, Glossary p. 138 "Capitalization rate: The ratio between a property's stabilized NOI and the property's sales price"; "Value Analysis" p. 43 (value = NOI ÷ cap rate); "Net operating income (NOI)" p. 140, https://www.occ.treas.gov/publications-and-resources/publications/comptrollers-handbook/files/commercial-real-estate-lending/pub-ch-commercial-real-estate.pdf (read 2026-11-13)
// Tags:    invest, real-estate, cap-rate, capitalization-rate, noi, rental-property, valuation
// Units:   annualNetOperatingIncomeCents and priceCents in cents; result decimal per year
// Example: hand-checked: NOI $16,100 a year on a $200,000 price → 0.0805
func CapitalizationRate(annualNetOperatingIncomeCents, priceCents int64) (float64, error) {
	if err := checkCents("annualNetOperatingIncomeCents", annualNetOperatingIncomeCents); err != nil {
		return 0, err
	}
	if err := checkNonNegativeCents("priceCents", priceCents); err != nil {
		return 0, err
	}
	if priceCents == 0 {
		return 0, fmt.Errorf("%w: priceCents = 0", ErrZeroDenominator)
	}
	return float64(annualNetOperatingIncomeCents) / float64(priceCents), nil
}
