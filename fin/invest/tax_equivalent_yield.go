package invest

import "fmt"

// TaxEquivalentYield returns the taxable yield that equals a tax-exempt yield after tax,
// taxFreeYield/(1 − marginalTaxRate), as a decimal.
// Errors: ErrRate if taxFreeYield is NaN, ±Inf or ≤ −1; ErrNotFinite if marginalTaxRate is NaN or
// ±Inf; ErrDomain if marginalTaxRate is outside [0, 1); ErrOverflow if the result exceeds float64.
//
// ID:      FIN-INVEST-027
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Municipal Securities Rulemaking Board (MSRB), "Understanding Taxable Municipal Bonds" (taxable equivalent yield = tax-exempt yield / (1 − tax bracket)), https://msrb.org/node/13519 (read 2026-11-12)
// Tags:    invest, bond, municipal-bond, tax-equivalent-yield, taxable-equivalent-yield, yield, tax
// Units:   taxFreeYield decimal per year; marginalTaxRate decimal; result decimal per year
// Example: MSRB: 3.00 % tax-exempt, 37 % bracket → 3.00 %/(1 − 0.37) = 4.76 %
func TaxEquivalentYield(taxFreeYield, marginalTaxRate float64) (float64, error) {
	if err := checkRate("taxFreeYield", taxFreeYield); err != nil {
		return 0, err
	}
	if err := checkFinite("marginalTaxRate", marginalTaxRate); err != nil {
		return 0, err
	}
	if marginalTaxRate < 0 || marginalTaxRate >= 1 {
		return 0, fmt.Errorf("%w: marginalTaxRate = %v, want 0 ≤ T < 1", ErrDomain, marginalTaxRate)
	}
	return finiteResult(taxFreeYield / (1 - marginalTaxRate))
}
