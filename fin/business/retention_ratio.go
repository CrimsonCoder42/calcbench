package business

// RetentionRatio returns (net income attributable to common shares − common share dividends) ÷ net income
// attributable to common shares, as a decimal (= 1 − DividendPayoutRatio; the plowback ratio). It is not
// rounded. The subtraction is exact int64 and the difference is within ±2⁵³, so the result is one correctly
// rounded division; it is negative when dividends exceed earnings. A loss is rejected.
// Errors: ErrOutOfRange when an input is outside ±2⁵³; ErrNegative when dividendsCents or netIncomeCents is
// negative; ErrZeroDenominator when netIncomeCents is 0.
//
// ID:      FIN-BIZ-029
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  CFA Institute, "Financial Ratio List", CFA Program Level II, ratio 36 "Retention rate = (Net income attributable to common shares – Common share dividends) ÷ Net income attributable to common shares = 1 – Payout ratio", https://www.cfainstitute.org/sites/default/files/-/media/documents/support/programs/cfa/cfa_program_level_ii_financial_ratio_list.pdf (read 2026-11-12)
// Tags:    ratio, dividend, retention, retention-ratio, plowback, sustainable-growth
// Units:   dividendsCents and netIncomeCents in cents; the result is a dimensionless decimal
// Example: dividends $250,000 on net income $1,000,000 → 750,000 ÷ 1,000,000 = 0.75 (hand-checked; the source gives no numbers)
func RetentionRatio(dividendsCents, netIncomeCents int64) (float64, error) {
	if err := checkNonNegativeCents("dividendsCents", dividendsCents); err != nil {
		return 0, err
	}
	if err := checkNetIncomeForPayout(netIncomeCents); err != nil {
		return 0, err
	}
	// 0 ≤ dividends, net income ≤ 2⁵³: the difference is within ±2⁵³ and exact in float64.
	return float64(netIncomeCents-dividendsCents) / float64(netIncomeCents), nil
}
