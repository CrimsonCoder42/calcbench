package business

import (
	"fmt"
	"math"
)

// DepreciationDecliningBalance returns the per-period depreciation expense of an asset by the declining-
// balance method, as a schedule of lifePeriods amounts in cents. The rate per period is
// rateMultiple ÷ lifePeriods (rateMultiple = 2 is double-declining balance). Salvage is not subtracted
// first: each period's expense is book value × rateMultiple ÷ lifePeriods, rounded half-even to the cent,
// but never more than book value − salvage, so the book value never falls below salvage. Once it reaches
// salvage the remaining periods are 0.
// The schedule can total LESS than cost − salvage: a declining balance only approaches salvage (with
// salvage 0 it never reaches it), and the source gives no rule for the balance left at the end of the life
// (no switch to straight line, no final-period adjustment), so none is applied. The caller who needs the
// full depreciable amount adds cost − salvage − Σschedule to the last period.
// Accuracy: each amount is computed as float64(book) × rateMultiple ÷ lifePeriods (the product is exact for
// small multiples such as 2, then one correctly rounded division) and rounded half-even; it equals the
// correctly rounded exact value except when that value lies within about 2·ε·|x| of a half cent, and is
// within 1 cent of it for amounts ≤ 2⁴³ cents. The result is a new slice.
// Errors: ErrOutOfRange when a money input is outside ±2⁵³; ErrNegative when a money input is negative;
// ErrDomain when salvageCents > costCents or rateMultiple ≤ 0; ErrNotFinite when rateMultiple is NaN or
// ±Inf; ErrPeriods when lifePeriods < 1 or > 1 000 000.
//
// ID:      FIN-BIZ-047
// Time:    O(n)   n = lifePeriods
// Space:   O(n) for the result (one allocation), O(1) extra
// Source:  Franklin, Graybeal and Cooper, Principles of Accounting, Vol. 1: Financial Accounting, OpenStax (Rice University), 2019, §11.3 "Explain and Apply Depreciation Methods to Allocate Capitalized Costs" (double-declining-balance: book value × straight-line percentage × 2; not below salvage), https://openstax.org/books/principles-financial-accounting/pages/11-3-explain-and-apply-depreciation-methods-to-allocate-capitalized-costs (read 2026-11-12)
// Tags:    depreciation, declining-balance, double-declining-balance, schedule, fixed-assets, accounting
// Units:   costCents, salvageCents and every schedule amount in cents; lifePeriods a count of periods; rateMultiple dimensionless
// Example: Kenzie Company: cost $58,000, salvage $10,000, 5 years, 40 %: $23,200, $13,920, $8,352, $2,528, $0 (Principles of Accounting Vol. 1 §11.3)
func DepreciationDecliningBalance(costCents, salvageCents int64, lifePeriods int, rateMultiple float64) ([]int64, error) {
	if err := checkDepreciationInputs(costCents, salvageCents, lifePeriods); err != nil {
		return nil, err
	}
	if math.IsNaN(rateMultiple) || math.IsInf(rateMultiple, 0) {
		return nil, fmt.Errorf("%w: rateMultiple = %v", ErrNotFinite, rateMultiple)
	}
	if rateMultiple <= 0 {
		return nil, fmt.Errorf("%w: rateMultiple = %v, want > 0", ErrDomain, rateMultiple)
	}
	life := float64(lifePeriods)
	schedule := make([]int64, lifePeriods)
	book := costCents
	for i := range schedule {
		remaining := book - salvageCents
		if remaining == 0 {
			break // the rest of the schedule stays 0
		}
		expense := float64(book) * rateMultiple / life
		var e int64
		if expense >= float64(remaining) {
			e = remaining
		} else {
			e = min(roundHalfEvenPlus(expense, 0), remaining)
		}
		schedule[i] = e
		book -= e
	}
	return schedule, nil
}
