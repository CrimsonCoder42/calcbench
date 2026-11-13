package business

import (
	"fmt"
	"math"

	"github.com/CrimsonCoder42/calcbench/algo/reduce"
)

// checkTaxRate returns ErrNotFinite for NaN or ±Inf and ErrRate outside [0, 1].
func checkTaxRate(taxRate float64) error {
	if !reduce.AllFinite(taxRate) {
		return fmt.Errorf("%w: taxRate = %v", ErrNotFinite, taxRate)
	}
	if taxRate < 0 || taxRate > 1 {
		return fmt.Errorf("%w: taxRate = %v, want 0 ≤ taxRate ≤ 1", ErrRate, taxRate)
	}
	return nil
}

// FreeCashFlowToFirm returns EBIT × (1 − taxRate) + depreciation − capital expenditures − change in working
// capital (fixed-capital and working-capital investment), in cents, rounded half-even once on the total.
// It is computed as EBIT − EBIT × taxRate (no rounding of 1 − taxRate); the tax product's rounding error is
// recovered with math.FMA, so the result is the correctly rounded exact value (taxRate taken as the
// float64 it holds).
// Negative EBIT, capital expenditures (net disposals) or working-capital change are allowed.
// Errors: ErrOutOfRange when a money input is outside ±2⁵³; ErrNegative when depreciationCents is
// negative; ErrNotFinite when taxRate is NaN or ±Inf; ErrRate when taxRate is outside [0, 1]; ErrOverflow
// when the result is outside ±2⁵³ cents.
//
// ID:      FIN-BIZ-035
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  CFA Institute, "Free Cash Flow Valuation", 2026 Curriculum CFA Program Level II Equity Valuation, refresher reading, Summary (FCFF = EBIT(1 – Tax rate) + Dep – FCInv – WCInv), https://www.cfainstitute.org/insights/professional-learning/refresher-readings/2025/free-cash-flow-valuation (read 2026-11-12)
// Source:  Dahlquist and Knight, Principles of Finance, 1st ed., OpenStax (Rice University), 2022, §5.6 "Operating Cash Flow and Free Cash Flow to the Firm (FCFF)", https://openstax.org/books/principles-finance/pages/5-6-operating-cash-flow-and-free-cash-flow-to-the-firm-fcff (read 2026-11-12)
// Tags:    cash-flow, free-cash-flow, fcff, valuation
// Units:   every money input and the result in cents; taxRate a decimal (0.25 = 25 %)
// Example: EBIT $100,000, tax 0.25, depreciation $10,000, capex $20,000, working capital +$5,000 → $60,000 (hand-checked; the source page gives no numbers)
func FreeCashFlowToFirm(ebitCents int64, taxRate float64, depreciationCents, capitalExpendituresCents, changeInWorkingCapitalCents int64) (int64, error) {
	if err := checkCents("ebitCents", ebitCents); err != nil {
		return 0, err
	}
	if err := checkTaxRate(taxRate); err != nil {
		return 0, err
	}
	if err := checkNonNegativeCents("depreciationCents", depreciationCents); err != nil {
		return 0, err
	}
	if err := checkCents("capitalExpendituresCents", capitalExpendituresCents); err != nil {
		return 0, err
	}
	if err := checkCents("changeInWorkingCapitalCents", changeInWorkingCapitalCents); err != nil {
		return 0, err
	}
	tax := float64(float64(ebitCents) * taxRate) // explicit conversion: rounded, never fused (Go spec, Floating-point operators); |tax| ≤ 2⁵³
	taxErr := math.FMA(float64(ebitCents), taxRate, -tax)
	fcff := roundHalfEvenExact(-tax, -taxErr, ebitCents+depreciationCents-capitalExpendituresCents-changeInWorkingCapitalCents)
	if err := checkResultCents("free cash flow to firm", fcff); err != nil {
		return 0, err
	}
	return fcff, nil
}
