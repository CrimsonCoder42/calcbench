package business

import "fmt"

// checkNetIncomeForPayout returns ErrOutOfRange when |netIncomeCents| > 2⁵³, ErrZeroDenominator when it is 0
// and ErrNegative when it is negative: a payout or retention share of a loss has no meaning.
func checkNetIncomeForPayout(netIncomeCents int64) error {
	if err := checkCents("netIncomeCents", netIncomeCents); err != nil {
		return err
	}
	if netIncomeCents == 0 {
		return fmt.Errorf("%w: netIncomeCents = 0", ErrZeroDenominator)
	}
	if netIncomeCents < 0 {
		return fmt.Errorf("%w: netIncomeCents = %d", ErrNegative, netIncomeCents)
	}
	return nil
}

// DividendPayoutRatio returns common share dividends divided by net income attributable to common shares,
// as a decimal (0.25 = a quarter of earnings paid out). It is not rounded and may exceed 1 when dividends
// exceed earnings. netIncomeCents is the net income attributable to common shares (after preferred
// dividends); a loss is rejected.
// Errors: ErrOutOfRange when an input is outside ±2⁵³; ErrNegative when dividendsCents or netIncomeCents is
// negative; ErrZeroDenominator when netIncomeCents is 0.
//
// ID:      FIN-BIZ-028
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  CFA Institute, "Financial Ratio List", CFA Program Level II, ratio 35 "Dividend payout ratio = Common share dividends ÷ Net income attributable to common shares", https://www.cfainstitute.org/sites/default/files/-/media/documents/support/programs/cfa/cfa_program_level_ii_financial_ratio_list.pdf (read 2026-11-12)
// Tags:    ratio, dividend, payout, dividend-payout-ratio, per-share
// Units:   dividendsCents and netIncomeCents in cents; the result is a dimensionless decimal
// Example: dividends $250,000 on net income $1,000,000 → 0.25 (hand-checked; the source gives no numbers)
func DividendPayoutRatio(dividendsCents, netIncomeCents int64) (float64, error) {
	if err := checkNonNegativeCents("dividendsCents", dividendsCents); err != nil {
		return 0, err
	}
	if err := checkNetIncomeForPayout(netIncomeCents); err != nil {
		return 0, err
	}
	return float64(dividendsCents) / float64(netIncomeCents), nil
}
