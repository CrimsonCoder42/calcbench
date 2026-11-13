package business

// GoodwillFromAcquisition returns the purchase price minus the fair value of the net assets acquired, in
// cents, exactly. The source covers the positive case (price above net assets, recorded as goodwill); a
// negative result is returned as is and means the price was below the net assets — this function does not
// say how that case is accounted for.
// Errors: ErrOutOfRange when an input is outside ±2⁵³; ErrNegative when purchasePriceCents is negative;
// ErrOverflow when the result is outside ±2⁵³ cents (possible only with negative net assets).
//
// ID:      FIN-BIZ-048
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Franklin, Graybeal and Cooper, Principles of Accounting, Vol. 1: Financial Accounting, OpenStax (Rice University), 2019, §11.4 "Describe Accounting for Intangible Assets and Record Related Transactions" (goodwill = purchase price − net assets), https://openstax.org/books/principles-financial-accounting/pages/11-4-describe-accounting-for-intangible-assets-and-record-related-transactions (read 2026-11-12)
// Tags:    goodwill, acquisition, intangible-assets, accounting
// Units:   purchasePriceCents, fairValueOfNetAssetsCents and the result in cents
// Example: London Hoops: $10 million − $7 million = $3 million (Principles of Accounting Vol. 1 §11.4)
func GoodwillFromAcquisition(purchasePriceCents, fairValueOfNetAssetsCents int64) (int64, error) {
	if err := checkNonNegativeCents("purchasePriceCents", purchasePriceCents); err != nil {
		return 0, err
	}
	if err := checkCents("fairValueOfNetAssetsCents", fairValueOfNetAssetsCents); err != nil {
		return 0, err
	}
	goodwill := purchasePriceCents - fairValueOfNetAssetsCents
	if err := checkResultCents("goodwill", goodwill); err != nil {
		return 0, err
	}
	return goodwill, nil
}
