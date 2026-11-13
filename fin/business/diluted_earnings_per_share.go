package business

import "fmt"

// DilutedEarningsPerShare returns (net income − preferred dividends + after-tax interest on convertible
// debt) ÷ (weighted average ordinary shares + dilutive potential ordinary shares), in cents per share,
// rounded half-even. The division is exact integer arithmetic (no float64), so the rounding is exact.
// preferredDividendsCents is the dividend on preferred shares that stay outstanding: dividends on
// convertible preferred shares assumed converted are left out of it (IAS 33 ¶33(a) adds them back), and their
// conversion shares go in dilutiveShares. convertibleInterestAfterTaxCents is the period's interest on
// convertible debt net of tax (¶33(b)). dilutiveShares is the weighted average number of shares issued on
// conversion of the dilutive potential ordinary shares only (¶36); deciding which instruments are dilutive
// (¶41–44) and any other income or expense changes from conversion (¶33(c)) are the caller's. A loss gives a
// negative result.
// Errors: ErrOutOfRange when a money input is outside ±2⁵³ or a share count > 2⁵³; ErrNegative when
// preferredDividendsCents, convertibleInterestAfterTaxCents, weightedAverageShares or dilutiveShares is
// negative; ErrZeroDenominator when weightedAverageShares is 0; ErrOverflow when the result is outside ±2⁵³
// cents.
//
// ID:      FIN-BIZ-022
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  International Accounting Standard 33 "Earnings per Share", ¶12, ¶31–33 and ¶36, as adopted in the Annex to Commission Regulation (EU) 2023/1803 of 13 August 2023, OJ L 237, 26.9.2023, https://eur-lex.europa.eu/legal-content/EN/TXT/HTML/?uri=CELEX:32023R1803 (read 2026-11-12)
// Source:  CFA Institute, "CFA Program Level I Sample Questions", Q5 (Sublyme Corporation, diluted EPS with convertible preferred shares), https://www.cfainstitute.org/programs/cfa-program/cfa-program-level-i-sample-questions (read 2026-11-12)
// Tags:    per-share, market-value, earnings-per-share, diluted-eps, eps, if-converted
// Units:   netIncomeCents, preferredDividendsCents and convertibleInterestAfterTaxCents in cents; weightedAverageShares and dilutiveShares share counts; the result in cents per share
// Example: Sublyme USD200,000,000 ÷ (50,000,000 + 2,000,000 × 2) = USD3.70 (CFA Level I sample Q5) → 370 cents
func DilutedEarningsPerShare(netIncomeCents, preferredDividendsCents, convertibleInterestAfterTaxCents, weightedAverageShares, dilutiveShares int64) (int64, error) {
	if err := checkCents("netIncomeCents", netIncomeCents); err != nil {
		return 0, err
	}
	if err := checkNonNegativeCents("preferredDividendsCents", preferredDividendsCents); err != nil {
		return 0, err
	}
	if err := checkNonNegativeCents("convertibleInterestAfterTaxCents", convertibleInterestAfterTaxCents); err != nil {
		return 0, err
	}
	if err := checkShares("weightedAverageShares", weightedAverageShares); err != nil {
		return 0, err
	}
	if dilutiveShares > maxExactCents {
		return 0, fmt.Errorf("%w: dilutiveShares = %d, want ≤ 2^53", ErrOutOfRange, dilutiveShares)
	}
	if dilutiveShares < 0 {
		return 0, fmt.Errorf("%w: dilutiveShares = %d", ErrNegative, dilutiveShares)
	}
	// |numerator| ≤ 3·2⁵³ and 1 ≤ denominator ≤ 2⁵⁴: exact in int64, and divRoundHalfEven's 2·remainder
	// stays below 2⁵⁵.
	eps := divRoundHalfEven(netIncomeCents-preferredDividendsCents+convertibleInterestAfterTaxCents, weightedAverageShares+dilutiveShares)
	if err := checkResultCents("diluted earnings per share", eps); err != nil {
		return 0, err
	}
	return eps, nil
}
