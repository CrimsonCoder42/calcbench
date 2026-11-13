package invest

import "fmt"

// SharpeRatio returns the reward-to-variability ratio (Rp − Rf)/σp: the portfolio's return above the
// risk-free rate per unit of the standard deviation of its return, all over the same period.
// Errors: ErrRate if portfolioReturn or riskFreeRate is NaN, ±Inf or ≤ −1; ErrNotFinite if
// returnStdDev is NaN or ±Inf; ErrNegative if returnStdDev < 0; ErrZeroDenominator if it is 0.
//
// ID:      FIN-INVEST-025
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Dahlquist, Knight et al., Principles of Finance, 1st ed., OpenStax (Rice University), 2022, §15.4 "Applications in Performance Measurement", eq. 15.14, https://openstax.org/books/principles-finance/pages/15-4-applications-in-performance-measurement
// Source:  Sharpe, "Mutual Fund Performance", Journal of Business 39(S1):119–138, 1966, doi:10.1086/294846; Sharpe, "The Sharpe Ratio", Journal of Portfolio Management 21(1):49–58, 1994, doi:10.3905/jpm.1994.409501 (the originals; DOIs confirmed on Crossref)
// Tags:    invest, portfolio, performance, sharpe-ratio, risk-adjusted-return, reward-to-variability
// Units:   returns and returnStdDev decimal per period; result dimensionless
// Example: OpenStax §15.4 McKinley: 20 % return, 3 % risk-free, σ 26 % → 17 %/26 % = 0.65
func SharpeRatio(portfolioReturn, riskFreeRate, returnStdDev float64) (float64, error) {
	if err := checkRate("portfolioReturn", portfolioReturn); err != nil {
		return 0, err
	}
	if err := checkRate("riskFreeRate", riskFreeRate); err != nil {
		return 0, err
	}
	if err := checkFinite("returnStdDev", returnStdDev); err != nil {
		return 0, err
	}
	if returnStdDev < 0 {
		return 0, fmt.Errorf("%w: returnStdDev = %v", ErrNegative, returnStdDev)
	}
	if returnStdDev == 0 {
		return 0, fmt.Errorf("%w: returnStdDev = 0", ErrZeroDenominator)
	}
	return finiteResult((portfolioReturn - riskFreeRate) / returnStdDev)
}
