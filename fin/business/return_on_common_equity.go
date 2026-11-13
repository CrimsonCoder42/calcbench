package business

// ReturnOnCommonEquity returns (net income − preferred dividends) divided by average common stockholders'
// equity, as a decimal. The subtraction is exact in int64; when its magnitude exceeds 2⁵³ the conversion to
// float64 adds one rounding before the division. A loss or negative equity is allowed; the sign of the
// result follows the inputs. Principles of Finance labels this computation ROE; FIN-BIZ-031 is the
// version without the preferred-dividend deduction.
// Errors: ErrOutOfRange when an input is outside ±2⁵³; ErrNegative when preferredDividendsCents is negative;
// ErrZeroDenominator when averageCommonEquityCents is 0.
//
// ID:      FIN-BIZ-032
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Dahlquist and Knight, Principles of Finance, 1st ed., OpenStax (Rice University), 2022, §6.6 "Profitability Ratios and the DuPont Method" (worked ROE deducts preferred dividends), https://openstax.org/books/principles-finance/pages/6-6-profitability-ratios-and-the-dupont-method (read 2026-11-12)
// Tags:    ratio, profitability, return, roe, return-on-common-equity, rcse
// Units:   every input in cents; the result is a decimal per period
// Example: Clear Lake Sporting Goods ($35,000 − $5,000) ÷ $95,000 = 0.3158, printed 0.32 (Principles of Finance §6.6)
func ReturnOnCommonEquity(netIncomeCents, preferredDividendsCents, averageCommonEquityCents int64) (float64, error) {
	if err := checkCents("netIncomeCents", netIncomeCents); err != nil {
		return 0, err
	}
	if err := checkNonNegativeCents("preferredDividendsCents", preferredDividendsCents); err != nil {
		return 0, err
	}
	if err := checkCents("averageCommonEquityCents", averageCommonEquityCents); err != nil {
		return 0, err
	}
	return ratioCents(netIncomeCents-preferredDividendsCents, "averageCommonEquityCents", averageCommonEquityCents)
}
