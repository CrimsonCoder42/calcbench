package business

// FixedChargeCoverageRatio returns (EBIT + lease payments) ÷ (interest payments + lease payments), as a
// decimal ("times"). fixedChargesCents is the fixed charge other than interest: the lease payments of the
// CFA definition (for an SEC "ratio of earnings to fixed charges", the estimate of the interest within rental
// expense). It is not rounded. A negative EBIT is allowed and gives a ratio below 1 (or negative). The sums
// are exact int64; when a sum exceeds 2⁵³ in magnitude its float64 conversion adds one rounding, so the
// result is within 2 ulp of the exact quotient (exactly rounded when both sums are within ±2⁵³).
// The SEC definition is historical: Item 503(d) and its instructions were removed from Regulation S-K by
// Release 33-10532 (2018); it is cited for what "fixed charges" meant, not as a current requirement.
// Errors: ErrOutOfRange when an input is outside ±2⁵³; ErrNegative when fixedChargesCents or
// interestExpenseCents is negative; ErrZeroDenominator when interestExpenseCents + fixedChargesCents is 0.
//
// ID:      FIN-BIZ-008
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  CFA Institute, "Financial Ratio List", CFA Program Level II, ratio 34 "Fixed charge coverage ratio = (Earnings before interest and taxes + Lease payments) ÷ (Interest payments + Lease payments)", https://www.cfainstitute.org/sites/default/files/-/media/documents/support/programs/cfa/cfa_program_level_ii_financial_ratio_list.pdf (read 2026-11-12)
// Source:  U.S. SEC, Regulation S-K, 17 CFR 229.503(d) "Ratio of earnings to fixed charges" and its Instruction 1 (earnings = pre-tax income from continuing operations plus fixed charges; fixed charges = interest expensed and capitalized plus the estimated interest within rental expense, among others), as in force 2018-01-01, https://www.ecfr.gov/api/renderer/v1/content/enhanced/2018-01-01/title-17?part=229&section=229.503 (read 2026-11-12)
// Source:  U.S. SEC, "Disclosure Update and Simplification", Release No. 33-10532, August 17, 2018, amendatory instruction to § 229.503 ("Removing paragraph (d) and the instructions to paragraph (d)"), https://www.sec.gov/rules/final/2018/33-10532.pdf (read 2026-11-12)
// Tags:    ratio, solvency, coverage, fixed-charge-coverage, lease
// Units:   ebitCents, fixedChargesCents and interestExpenseCents in cents; the result is a dimensionless count of times
// Example: EBIT $50,000, lease payments $10,000, interest $20,000 → 60,000 ÷ 30,000 = 2 (hand-checked; the sources give no numbers)
func FixedChargeCoverageRatio(ebitCents, fixedChargesCents, interestExpenseCents int64) (float64, error) {
	if err := checkCents("ebitCents", ebitCents); err != nil {
		return 0, err
	}
	if err := checkNonNegativeCents("fixedChargesCents", fixedChargesCents); err != nil {
		return 0, err
	}
	if err := checkNonNegativeCents("interestExpenseCents", interestExpenseCents); err != nil {
		return 0, err
	}
	// |sums| ≤ 2⁵⁴: exact in int64.
	return ratioCents(ebitCents+fixedChargesCents, "interestExpenseCents + fixedChargesCents", interestExpenseCents+fixedChargesCents)
}
