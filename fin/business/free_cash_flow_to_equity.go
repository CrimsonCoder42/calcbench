package business

import "math"

// FreeCashFlowToEquity returns FCFF − interest expense × (1 − taxRate) + net borrowing, in cents, rounded
// half-even once on the total. It is computed as FCFF − interest + interest × taxRate + net borrowing; the
// tax-shield product's rounding error is recovered with math.FMA, so the result is the correctly rounded
// exact value (taxRate taken as the float64 it holds).
// fcffCents is typically FIN-BIZ-035's result; negative FCFF or net repayment (negative net borrowing) are
// allowed.
// Errors: ErrOutOfRange when a money input is outside ±2⁵³; ErrNegative when interestExpenseCents is
// negative; ErrNotFinite when taxRate is NaN or ±Inf; ErrRate when taxRate is outside [0, 1]; ErrOverflow
// when the result is outside ±2⁵³ cents.
//
// ID:      FIN-BIZ-036
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  CFA Institute, "Free Cash Flow Valuation", 2026 Curriculum CFA Program Level II Equity Valuation, refresher reading, Summary (FCFE = FCFF – Int(1 – Tax rate) + Net borrowing), https://www.cfainstitute.org/insights/professional-learning/refresher-readings/2025/free-cash-flow-valuation (read 2026-11-12)
// Tags:    cash-flow, free-cash-flow, fcfe, valuation
// Units:   every money input and the result in cents; taxRate a decimal (0.25 = 25 %)
// Example: FCFF $60,000, interest $4,000, tax 0.25, net borrowing $2,000 → $59,000 (hand-checked; the source page gives no numbers)
func FreeCashFlowToEquity(fcffCents, interestExpenseCents int64, taxRate float64, netBorrowingCents int64) (int64, error) {
	if err := checkCents("fcffCents", fcffCents); err != nil {
		return 0, err
	}
	if err := checkNonNegativeCents("interestExpenseCents", interestExpenseCents); err != nil {
		return 0, err
	}
	if err := checkTaxRate(taxRate); err != nil {
		return 0, err
	}
	if err := checkCents("netBorrowingCents", netBorrowingCents); err != nil {
		return 0, err
	}
	shield := float64(float64(interestExpenseCents) * taxRate) // explicit conversion: rounded, never fused (Go spec, Floating-point operators); |shield| ≤ 2⁵³
	shieldErr := math.FMA(float64(interestExpenseCents), taxRate, -shield)
	fcfe := roundHalfEvenExact(shield, shieldErr, fcffCents-interestExpenseCents+netBorrowingCents)
	if err := checkResultCents("free cash flow to equity", fcfe); err != nil {
		return 0, err
	}
	return fcfe, nil
}
