package invest

// RiskPremium returns the extra return earned for taking on risk, expectedReturn − riskFreeRate, as
// a decimal per period (the market risk premium when expectedReturn is the market's return).
// Errors: ErrRate if a rate is NaN, ±Inf or ≤ −1.
//
// ID:      FIN-INVEST-017
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Dahlquist, Knight et al., Principles of Finance, 1st ed., OpenStax (Rice University), 2022, §15.3 "The Capital Asset Pricing Model (CAPM)" (Risk Premium = S&P 500 Avg Return − US T-Bill Avg Return), https://openstax.org/books/principles-finance/pages/15-3-the-capital-asset-pricing-model-capm
// Tags:    invest, risk-premium, market-risk-premium, excess-return, capm
// Units:   rates decimal per period; result decimal per period
// Example: OpenStax §15.3: 11.64 % − 3.36 % = 8.28 %
func RiskPremium(expectedReturn, riskFreeRate float64) (float64, error) {
	if err := checkRate("expectedReturn", expectedReturn); err != nil {
		return 0, err
	}
	if err := checkRate("riskFreeRate", riskFreeRate); err != nil {
		return 0, err
	}
	return expectedReturn - riskFreeRate, nil
}
