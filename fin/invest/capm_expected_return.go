package invest

// CAPMExpectedReturn returns the capital asset pricing model's expected (required) return on a
// security, Rf + β·(Rm − Rf), as a decimal per period.
// Errors: ErrRate if riskFreeRate or marketReturn is NaN, ±Inf or ≤ −1; ErrNotFinite if beta is NaN
// or ±Inf; ErrOverflow if the result exceeds float64.
//
// ID:      FIN-INVEST-016
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Dahlquist, Knight et al., Principles of Finance, 1st ed., OpenStax (Rice University), 2022, §15.3 "The Capital Asset Pricing Model (CAPM)" (Re = Rf + Beta × (Rm − Rf)), https://openstax.org/books/principles-finance/pages/15-3-the-capital-asset-pricing-model-capm
// Source:  Sharpe, "Capital Asset Prices: A Theory of Market Equilibrium under Conditions of Risk", Journal of Finance 19(3):425–442, 1964, doi:10.1111/j.1540-6261.1964.tb02865.x (the original; DOI confirmed on Crossref)
// Tags:    invest, capm, expected-return, required-return, beta, cost-of-equity, risk
// Units:   rates decimal per period; beta dimensionless; result decimal per period
// Example: OpenStax §15.3 Delta (DAL): Rf 3.36 %, Rm 11.64 %, β 1.39 → 14.87 %
func CAPMExpectedReturn(riskFreeRate, beta, marketReturn float64) (float64, error) {
	if err := checkRate("riskFreeRate", riskFreeRate); err != nil {
		return 0, err
	}
	if err := checkFinite("beta", beta); err != nil {
		return 0, err
	}
	if err := checkRate("marketReturn", marketReturn); err != nil {
		return 0, err
	}
	return finiteResult(riskFreeRate + beta*(marketReturn-riskFreeRate))
}
