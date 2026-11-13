package invest

import "fmt"

// NetAssetValuePerShare returns a fund's net asset value per share, (total assets − total
// liabilities)/shares outstanding, in cents rounded half-even. The division is exact integer
// arithmetic (quotient and remainder), so no float rounding enters. A negative NAV (liabilities above
// assets) gives a negative result.
// Errors: ErrOutOfRange if totalAssetsCents or totalLiabilitiesCents is above 2⁵³ cents; ErrNegative
// if an input is negative; ErrZeroDenominator if sharesOutstanding is 0.
//
// ID:      FIN-INVEST-028
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  US Securities and Exchange Commission, Fast Answers "Net Asset Value", http://www.sec.gov/answers/nav.htm (now investor.gov glossary; read 2026-11-12 from the copy archived 2013-02-14 at https://webarchive.library.unt.edu/web/20130214105729mp_/http://www.sec.gov/answers/nav.htm)
// Tags:    invest, fund, mutual-fund, nav, net-asset-value, per-share
// Units:   totalAssetsCents, totalLiabilitiesCents and the result in cents; sharesOutstanding shares
// Example: SEC: NAV $100 million over 10,000,000 shares → $10; assets $100 million, liabilities $10 million → NAV $90 million
func NetAssetValuePerShare(totalAssetsCents, totalLiabilitiesCents, sharesOutstanding int64) (int64, error) {
	if err := checkNonNegativeCents("totalAssetsCents", totalAssetsCents); err != nil {
		return 0, err
	}
	if err := checkNonNegativeCents("totalLiabilitiesCents", totalLiabilitiesCents); err != nil {
		return 0, err
	}
	if sharesOutstanding < 0 {
		return 0, fmt.Errorf("%w: sharesOutstanding = %d", ErrNegative, sharesOutstanding)
	}
	if sharesOutstanding == 0 {
		return 0, fmt.Errorf("%w: sharesOutstanding = 0", ErrZeroDenominator)
	}
	nav := totalAssetsCents - totalLiabilitiesCents // |nav| ≤ 2⁵³
	q, r := nav/sharesOutstanding, nav%sharesOutstanding
	if r < 0 { // floor division so the half-even test below sees 0 ≤ r < shares
		q, r = q-1, r+sharesOutstanding
	}
	// Round half-even: compare 2r with shares without overflow (r < shares ≤ 2⁶³−1).
	if r > sharesOutstanding-r || (r == sharesOutstanding-r && q%2 != 0) {
		q++
	}
	return q, nil
}
