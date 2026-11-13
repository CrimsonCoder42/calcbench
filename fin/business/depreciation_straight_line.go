package business

import "fmt"

// maxLifePeriods bounds a depreciation schedule's length so that a huge lifePeriods cannot exhaust memory.
const maxLifePeriods = 1_000_000

// checkDepreciationInputs validates the inputs shared by the depreciation schedules.
func checkDepreciationInputs(costCents, salvageCents int64, lifePeriods int) error {
	if err := checkNonNegativeCents("costCents", costCents); err != nil {
		return err
	}
	if err := checkNonNegativeCents("salvageCents", salvageCents); err != nil {
		return err
	}
	if salvageCents > costCents {
		return fmt.Errorf("%w: salvageCents = %d > costCents = %d", ErrDomain, salvageCents, costCents)
	}
	if lifePeriods < 1 || lifePeriods > maxLifePeriods {
		return fmt.Errorf("%w: lifePeriods = %d, want 1 ≤ lifePeriods ≤ %d", ErrPeriods, lifePeriods, maxLifePeriods)
	}
	return nil
}

// DepreciationStraightLine returns the per-period depreciation expense of an asset by the straight-line
// method, (cost − salvage) ÷ life, as a schedule of lifePeriods amounts in cents. Each period gets
// ⌊(cost − salvage) ÷ lifePeriods⌋ (exact integer division); the last period also takes the remainder, so
// the schedule sums to cost − salvage exactly. The result is a new slice; the book value after period t is
// cost minus the sum of the first t amounts.
// Errors: ErrOutOfRange when a money input is outside ±2⁵³; ErrNegative when a money input is negative;
// ErrDomain when salvageCents > costCents; ErrPeriods when lifePeriods < 1 or > 1 000 000.
//
// ID:      FIN-BIZ-046
// Time:    O(n)   n = lifePeriods
// Space:   O(n) for the result (one allocation), O(1) extra
// Source:  Franklin, Graybeal and Cooper, Principles of Accounting, Vol. 1: Financial Accounting, OpenStax (Rice University), 2019, §11.3 "Explain and Apply Depreciation Methods to Allocate Capitalized Costs" (straight-line: (cost − salvage value) ÷ useful life), https://openstax.org/books/principles-financial-accounting/pages/11-3-explain-and-apply-depreciation-methods-to-allocate-capitalized-costs (read 2026-11-12)
// Tags:    depreciation, straight-line, schedule, fixed-assets, accounting
// Units:   costCents, salvageCents and every schedule amount in cents; lifePeriods a count of periods
// Example: Kenzie Company: cost $58,000, salvage $10,000, 5 years → $9,600 each year, total $48,000 (Principles of Accounting Vol. 1 §11.3)
func DepreciationStraightLine(costCents, salvageCents int64, lifePeriods int) ([]int64, error) {
	if err := checkDepreciationInputs(costCents, salvageCents, lifePeriods); err != nil {
		return nil, err
	}
	depreciable := costCents - salvageCents
	n := int64(lifePeriods)
	per := depreciable / n
	schedule := make([]int64, lifePeriods)
	for i := range schedule {
		schedule[i] = per
	}
	schedule[lifePeriods-1] = depreciable - per*(n-1)
	return schedule, nil
}
